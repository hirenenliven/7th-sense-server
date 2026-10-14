package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"golang.org/x/crypto/bcrypt"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// UserService holds the business rules for users. A user is saved together
// with its customer details (company profile, user profile and address) in
// one transaction, so a failure never leaves a half-created customer.
type UserService struct {
	dbPool        *pgxpool.Pool
	branchService *BranchService
	// services runs on the pool; writes use services created on a transaction.
	services userServices
}

// userServices are the services that store a user and its customer details.
type userServices struct {
	userRepository        *repository.UserRepository
	userProfileService    *UserProfileService
	companyProfileService *CompanyProfileService
	addressService        *AddressService
}

// newUserServices creates the user services on database (a pool or a transaction).
func newUserServices(database db.Querier) userServices {
	return userServices{
		userRepository:        repository.NewUserRepository(database),
		userProfileService:    NewUserProfileService(repository.NewUserProfileRepository(database)),
		companyProfileService: NewCompanyProfileService(repository.NewCompanyProfileRepository(database)),
		addressService:        NewAddressService(repository.NewAddressRepository(database)),
	}
}

// NewUserService creates a UserService that stores users in dbPool and checks
// branches with branchService.
func NewUserService(dbPool *pgxpool.Pool, branchService *BranchService) *UserService {
	return &UserService{
		dbPool:        dbPool,
		branchService: branchService,
		services:      newUserServices(dbPool),
	}
}

func (userService *UserService) GetAll(ctx context.Context) ([]model.User, error) {
	return userService.services.userRepository.FindAll(ctx)
}

func (userService *UserService) GetByID(ctx context.Context, id int64) (dto.UserResponse, error) {
	return userService.services.details(ctx, id)
}

// Create adds a customer user with its customer details. A request without a
// role makes the customer an admin, so they can manage their own users. actor
// is the logged-in user; only the superadmin may give the superadmin role.
func (userService *UserService) Create(ctx context.Context, actor model.AuthUser, userRequest dto.CreateUserRequest) (dto.UserResponse, error) {
	if userRequest.RoleID == 0 {
		adminRoleID, err := userService.roleID(ctx, model.RoleAdmin)
		if err != nil {
			return dto.UserResponse{}, err
		}
		userRequest.RoleID = adminRoleID
	}
	if err := userService.checkRoleAllowed(ctx, actor, userRequest.RoleID); err != nil {
		return dto.UserResponse{}, err
	}
	if err := userService.checkBranch(ctx, *userRequest.BranchID); err != nil {
		return dto.UserResponse{}, err
	}
	dayStart, err := dto.ParseDayStart(userRequest.DayStartFrom)
	if err != nil {
		return dto.UserResponse{}, response.NewError(http.StatusBadRequest, "Please enter valid day start time")
	}
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userRequest.Password), bcrypt.DefaultCost)
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("hash password: %w", err)
	}

	var created dto.UserResponse
	err = db.WithTransaction(ctx, userService.dbPool, func(tx pgx.Tx) error {
		services := newUserServices(tx)

		user, err := services.userRepository.Create(ctx, model.User{
			Username:  userRequest.Username,
			Email:     userRequest.Email,
			Password:  string(hashedPassword),
			FirstName: userRequest.FirstName,
			LastName:  userRequest.LastName,
			Mobile:    userRequest.PrimaryPhone,
			RoleID:    userRequest.RoleID,
			BranchID:  userRequest.BranchID,
			CompanyID: userRequest.CompanyID,
		})
		if err != nil {
			return toUserResponseError(err)
		}

		if _, err := services.companyProfileService.Create(ctx, companyProfileFrom(user.ID, userRequest.CustomerDetails)); err != nil {
			return err
		}
		profile := model.UserProfile{UserID: user.ID, Theme: themeOrDefault(userRequest.Theme), Language: defaultLanguage}
		setCustomerSettings(&profile, userRequest.CustomerDetails, dayStart)
		profile.Filters = userRequest.Filters
		if _, err := services.userProfileService.Create(ctx, profile); err != nil {
			return err
		}
		address := model.Address{UserID: user.ID, AddressType: dto.AddressTypeWork, IsDefault: true}
		setCustomerAddress(&address, userRequest.CustomerDetails)
		if _, err := services.addressService.Create(ctx, address); err != nil {
			return err
		}

		created, err = services.details(ctx, user.ID)
		return err
	})
	return created, err
}

// Update replaces the user and its customer details. Customer details the user
// does not have yet are created. User profile fields that are not part of the
// customer form (designation, birth date, ...) keep their values; an empty
// theme keeps the current theme. Only the superadmin may change the superadmin
// or give the superadmin role (actor is the logged-in user).
func (userService *UserService) Update(ctx context.Context, actor model.AuthUser, id int64, userRequest dto.UpdateUserRequest) (dto.UserResponse, error) {
	if err := userService.checkCanManage(ctx, actor, id); err != nil {
		return dto.UserResponse{}, err
	}
	if err := userService.checkRoleAllowed(ctx, actor, userRequest.RoleID); err != nil {
		return dto.UserResponse{}, err
	}
	if err := userService.checkBranch(ctx, *userRequest.BranchID); err != nil {
		return dto.UserResponse{}, err
	}
	dayStart, err := dto.ParseDayStart(userRequest.DayStartFrom)
	if err != nil {
		return dto.UserResponse{}, response.NewError(http.StatusBadRequest, "Please enter valid day start time")
	}

	var updated dto.UserResponse
	err = db.WithTransaction(ctx, userService.dbPool, func(tx pgx.Tx) error {
		services := newUserServices(tx)

		if _, err := services.userRepository.Update(ctx, model.User{
			ID:        id,
			Username:  userRequest.Username,
			Email:     userRequest.Email,
			FirstName: userRequest.FirstName,
			LastName:  userRequest.LastName,
			Mobile:    userRequest.PrimaryPhone,
			RoleID:    userRequest.RoleID,
			BranchID:  userRequest.BranchID,
			CompanyID: userRequest.CompanyID,
			Status:    userRequest.Status,
		}); err != nil {
			return toUserResponseError(err)
		}

		// Company profile: every column comes from the form.
		companyProfile := companyProfileFrom(id, userRequest.CustomerDetails)
		_, err := services.companyProfileService.Update(ctx, companyProfile)
		if err != nil && err.Error() == "Company profile not found" {
			_, err = services.companyProfileService.Create(ctx, companyProfile)
		}
		if err != nil {
			return err
		}

		// User profile: change only the form fields.
		profile, err := services.userProfileService.GetByUserID(ctx, id)
		switch {
		case err != nil && err.Error() == "User profile not found":
			profile = model.UserProfile{UserID: id, Theme: themeOrDefault(userRequest.Theme), Language: defaultLanguage}
			setCustomerSettings(&profile, userRequest.CustomerDetails, dayStart)
			profile.Filters = userRequest.Filters
			_, err = services.userProfileService.Create(ctx, profile)
		case err == nil:
			if userRequest.Theme != "" {
				profile.Theme = userRequest.Theme
			}
			setCustomerSettings(&profile, userRequest.CustomerDetails, dayStart)
			profile.Filters = userRequest.Filters
			_, err = services.userProfileService.Update(ctx, profile)
		}
		if err != nil {
			return err
		}

		// Address: change the primary address, keeping its type and default flag.
		addresses, err := services.addressService.GetByUserID(ctx, id)
		if err != nil {
			return err
		}
		if address := primaryAddress(addresses); address != nil {
			setCustomerAddress(address, userRequest.CustomerDetails)
			_, err = services.addressService.Update(ctx, *address)
		} else {
			newAddress := model.Address{UserID: id, AddressType: dto.AddressTypeWork, IsDefault: true}
			setCustomerAddress(&newAddress, userRequest.CustomerDetails)
			_, err = services.addressService.Create(ctx, newAddress)
		}
		if err != nil {
			return err
		}

		updated, err = services.details(ctx, id)
		return err
	})
	return updated, err
}

// Delete soft-deletes the user with its profiles and addresses. Only the
// superadmin may delete the superadmin (actor is the logged-in user).
func (userService *UserService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	if err := userService.checkCanManage(ctx, actor, id); err != nil {
		return err
	}
	return db.WithTransaction(ctx, userService.dbPool, func(tx pgx.Tx) error {
		services := newUserServices(tx)

		if err := services.userRepository.Delete(ctx, id); err != nil {
			return toUserResponseError(err)
		}
		if err := services.userProfileService.Delete(ctx, id); err != nil && err.Error() != "User profile not found" {
			return err
		}
		if err := services.companyProfileService.Delete(ctx, id); err != nil && err.Error() != "Company profile not found" {
			return err
		}
		addresses, err := services.addressService.GetByUserID(ctx, id)
		if err != nil {
			return err
		}
		for _, address := range addresses {
			if err := services.addressService.Delete(ctx, address.ID); err != nil {
				return err
			}
		}
		return nil
	})
}

// SaveTheme sets the UI theme of userID (the logged-in user).
func (userService *UserService) SaveTheme(ctx context.Context, userID int64, theme string) error {
	return userService.services.userProfileService.SaveTheme(ctx, userID, theme)
}

// Schema defaults of user_profiles.theme and language. The repository writes
// an empty value as NULL, so new profiles set them explicitly.
const (
	defaultTheme    = dto.ThemeDark
	defaultLanguage = "en"
)

// themeOrDefault returns theme, or defaultTheme when it is empty.
func themeOrDefault(theme string) string {
	if theme == "" {
		return defaultTheme
	}
	return theme
}

// details reads the user with its customer details. Missing details are nil.
func (services userServices) details(ctx context.Context, id int64) (dto.UserResponse, error) {
	user, err := services.userRepository.FindByID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, toUserResponseError(err)
	}
	userResponse := dto.UserResponse{User: user}

	profile, err := services.userProfileService.GetByUserID(ctx, id)
	switch {
	case err == nil:
		userResponse.Profile = &profile
	case err != nil && err.Error() != "User profile not found":
		fmt.Printf("Error fetching profile for user %d: %v\n", id, err)
		return dto.UserResponse{}, err
	default:
		fmt.Printf("Profile not found for user %d: %v\n", id, err)
	}

	companyProfile, err := services.companyProfileService.GetByUserID(ctx, id)
	switch {
	case err == nil:
		userResponse.CompanyProfile = &companyProfile
	case err != nil && err.Error() != "Company profile not found":
		return dto.UserResponse{}, err
	}

	addresses, err := services.addressService.GetByUserID(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse.Address = primaryAddress(addresses)
	return userResponse, nil
}

// roleID returns the id of the role called name.
func (userService *UserService) roleID(ctx context.Context, name string) (int64, error) {
	roleID, err := userService.services.userRepository.RoleIDByName(ctx, name)
	if errors.Is(err, repository.ErrNotFound) {
		return 0, fmt.Errorf("role %q is missing from the roles table", name)
	}
	return roleID, err
}

// checkRoleAllowed returns Forbidden when actor, who is not the superadmin,
// gives roleID and it is the superadmin role. Admins may give every other role,
// admin included.
func (userService *UserService) checkRoleAllowed(ctx context.Context, actor model.AuthUser, roleID int64) error {
	if actor.Role == model.RoleSuperAdmin {
		return nil
	}
	superAdminRoleID, err := userService.services.userRepository.RoleIDByName(ctx, model.RoleSuperAdmin)
	if errors.Is(err, repository.ErrNotFound) {
		return nil // no superadmin role, so it cannot be given
	}
	if err != nil {
		return err
	}
	if roleID == superAdminRoleID {
		return response.Forbidden
	}
	return nil
}

// checkCanManage returns Forbidden when actor, who is not the superadmin, would
// change or delete user id and that user is the superadmin. A missing user is
// left to the caller, which answers "User not found".
func (userService *UserService) checkCanManage(ctx context.Context, actor model.AuthUser, id int64) error {
	if actor.Role == model.RoleSuperAdmin {
		return nil
	}
	user, err := userService.services.userRepository.FindByID(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return nil
	}
	if err != nil {
		return err
	}
	return userService.checkRoleAllowed(ctx, actor, user.RoleID)
}

// checkBranch returns ValidBranch when branchID is not an existing branch.
func (userService *UserService) checkBranch(ctx context.Context, branchID int64) error {
	_, err := userService.branchService.GetByID(ctx, branchID)
	if err != nil && err.Error() == "Branch not found" {
		return response.NewError(http.StatusBadRequest, "Please enter valid branch")
	}
	return err
}

// companyProfileFrom builds the company profile of userID from the form.
func companyProfileFrom(userID int64, details dto.CustomerDetails) model.CompanyProfile {
	return model.CompanyProfile{
		UserID:      userID,
		CompanyName: details.CompanyName,
		GSTNumber:   details.GSTNumber,
		PANNumber:   details.PANNumber,
		Website:     details.Website,
		FinixID:     details.FinixID,
	}
}

// setCustomerSettings copies the form's logo, contact and settings into profile.
func setCustomerSettings(profile *model.UserProfile, details dto.CustomerDetails, dayStart time.Time) {
	profile.ProfilePicture = details.CompanyLogo
	profile.ContactPerson = details.ContactPerson
	profile.ContactPhone = details.SecondaryPhone
	profile.Timezone = details.TimeZone
	profile.StartDayTime = &dayStart
}

// setCustomerAddress copies the form's address into address.
func setCustomerAddress(address *model.Address, details dto.CustomerDetails) {
	address.AddressLine1 = details.AddressLine1
	address.AddressLine2 = details.AddressLine2
	address.CountryName = details.Country
	address.StateName = details.State
	address.City = details.City
	address.PinCode = details.ZipCode
	address.ContactPerson = details.ContactPerson
	address.ContactPhone = details.PrimaryPhone
}

// primaryAddress returns the default address, else the first one, else nil.
func primaryAddress(addresses []model.Address) *model.Address {
	for index := range addresses {
		if addresses[index].IsDefault {
			return &addresses[index]
		}
	}
	if len(addresses) > 0 {
		return &addresses[0]
	}
	return nil
}

// toUserResponseError maps repository errors to user responses. Other errors
// pass through and are sent as ServerError.
func toUserResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "User not found")
	case errors.Is(err, repository.ErrDuplicateEmail):
		return response.NewError(http.StatusBadRequest, "Email already exists")
	case errors.Is(err, repository.ErrDuplicateUsername):
		return response.NewError(http.StatusBadRequest, "Username already exists")
	case errors.Is(err, repository.ErrInvalidRole):
		return response.NewError(http.StatusBadRequest, "Please enter valid role")
	default:
		return err
	}
}
