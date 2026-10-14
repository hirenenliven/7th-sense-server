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
// with its customer details (company profile, user profile and address) and
// machine access in one transaction, so a failure never leaves a half-created user.
//
// Customers are companies; the users they and their admins create are members of
// the company (users.company_id). Every method takes the logged-in user (actor):
// the superadmin manages every user, a customer or admin only their company's
// members, and only with the roles of model.AssignableRoles.
type UserService struct {
	dbPool         *pgxpool.Pool
	branchService  *BranchService
	roleService    *RoleService
	machineService *MachineService
	// services runs on the pool; writes use services created on a transaction.
	services userServices
}

// userServices are the services that store a user and its customer details.
type userServices struct {
	userRepository              *repository.UserRepository
	userProfileService          *UserProfileService
	companyProfileService       *CompanyProfileService
	addressService              *AddressService
	userMachineAccessRepository *repository.UserMachineAccessRepository
}

// newUserServices creates the user services on database (a pool or a transaction).
func newUserServices(database db.Querier) userServices {
	return userServices{
		userRepository:              repository.NewUserRepository(database),
		userProfileService:          NewUserProfileService(repository.NewUserProfileRepository(database)),
		companyProfileService:       NewCompanyProfileService(repository.NewCompanyProfileRepository(database)),
		addressService:              NewAddressService(repository.NewAddressRepository(database)),
		userMachineAccessRepository: repository.NewUserMachineAccessRepository(database),
	}
}

// NewUserService creates a UserService that stores users in dbPool and checks
// branches, roles and machines with the other services.
func NewUserService(dbPool *pgxpool.Pool, branchService *BranchService, roleService *RoleService, machineService *MachineService) *UserService {
	return &UserService{
		dbPool:         dbPool,
		branchService:  branchService,
		roleService:    roleService,
		machineService: machineService,
		services:       newUserServices(dbPool),
	}
}

// GetAll returns the users actor may see, narrowed by filter: every user for
// the superadmin, the company's members for a customer or admin. A customer or
// admin asking for another company gets Forbidden.
func (userService *UserService) GetAll(ctx context.Context, actor model.AuthUser, filter dto.UserListFilter) ([]model.User, error) {
	if actor.Role != model.RoleSuperAdmin {
		company := actor.Company()
		if company == nil || (filter.CompanyID != nil && *filter.CompanyID != *company) {
			return nil, response.Forbidden
		}
		filter.CompanyID = company
	}
	return userService.services.userRepository.FindAll(ctx, filter.Role, filter.CompanyID)
}

// GetByID returns user id with its details, or "User not found" when actor may
// not see it.
func (userService *UserService) GetByID(ctx context.Context, actor model.AuthUser, id int64) (dto.UserResponse, error) {
	userResponse, err := userService.services.details(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	if !canSee(actor, userResponse.User) {
		return dto.UserResponse{}, userNotFound
	}
	return userResponse, nil
}

// GetCurrent returns the logged-in user (userID) with its details.
func (userService *UserService) GetCurrent(ctx context.Context, userID int64) (dto.UserResponse, error) {
	return userService.services.details(ctx, userID)
}

// Create adds a user with its customer details and machine access. actor must
// be allowed to give the role (see accessFor).
func (userService *UserService) Create(ctx context.Context, actor model.AuthUser, userRequest dto.CreateUserRequest) (dto.UserResponse, error) {
	access, err := userService.accessFor(ctx, actor, userRequest.RoleID, userRequest.CompanyID, userRequest.MachineIDs)
	if err != nil {
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
			CompanyID: access.companyID,
		})
		if err != nil {
			return toUserResponseError(err)
		}
		if err := services.userMachineAccessRepository.Replace(ctx, user.ID, access.machineIDs); err != nil {
			return err
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
// theme keeps the current theme. The machine access is replaced with the
// request's. actor must be allowed to manage the user and give the new role.
func (userService *UserService) Update(ctx context.Context, actor model.AuthUser, id int64, userRequest dto.UpdateUserRequest) (dto.UserResponse, error) {
	if err := userService.checkCanManage(ctx, actor, id); err != nil {
		return dto.UserResponse{}, err
	}
	access, err := userService.accessFor(ctx, actor, userRequest.RoleID, userRequest.CompanyID, userRequest.MachineIDs)
	if err != nil {
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
			CompanyID: access.companyID,
			Status:    userRequest.Status,
		}); err != nil {
			return toUserResponseError(err)
		}
		if err := services.userMachineAccessRepository.Replace(ctx, id, access.machineIDs); err != nil {
			return err
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

// Delete soft-deletes the user with its profiles and addresses, and removes its
// machine access. actor must be allowed to manage the user.
func (userService *UserService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	if err := userService.checkCanManage(ctx, actor, id); err != nil {
		return err
	}
	return db.WithTransaction(ctx, userService.dbPool, func(tx pgx.Tx) error {
		services := newUserServices(tx)

		if err := services.userRepository.Delete(ctx, id); err != nil {
			return toUserResponseError(err)
		}
		if err := services.userMachineAccessRepository.Replace(ctx, id, nil); err != nil {
			return err
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

	machineIDs, err := services.userMachineAccessRepository.FindMachineIDs(ctx, id)
	if err != nil {
		return dto.UserResponse{}, err
	}
	userResponse.MachineIDs = append([]int64{}, machineIDs...) // [] rather than null in JSON
	return userResponse, nil
}

// userNotFound is the answer for a missing user and for one the actor may not
// see, so other companies' user ids are not revealed.
var userNotFound = response.NewError(http.StatusNotFound, "User not found")

// canSee reports whether actor may see user: the superadmin every user, anyone
// else only the members of their company.
func canSee(actor model.AuthUser, user model.User) bool {
	if actor.Role == model.RoleSuperAdmin {
		return true
	}
	company := actor.Company()
	return company != nil && user.CompanyID != nil && *user.CompanyID == *company
}

// checkCanManage returns "User not found" when actor may not see user id, and
// Forbidden when actor may not give the user's current role (so an admin cannot
// change another admin, and nobody changes the superadmin here).
func (userService *UserService) checkCanManage(ctx context.Context, actor model.AuthUser, id int64) error {
	user, err := userService.services.userRepository.FindByID(ctx, id)
	if err != nil {
		return toUserResponseError(err)
	}
	if !canSee(actor, user) {
		return userNotFound
	}
	if !model.CanAssignRole(actor.Role, user.Role) {
		return response.Forbidden
	}
	return nil
}

// userAccess is the company and machines a user is saved with.
type userAccess struct {
	companyID  *int64
	machineIDs []int64
}

// accessFor checks that actor may give roleID and returns the company and
// machines to save the user with:
//   - a customer is a company itself: no company, no machines;
//   - a member belongs to the actor's company; the superadmin names the company
//     (companyID), which must be a customer;
//   - an admin sees all of the company's machines, so only other members keep
//     machineIDs, which must all be machines of the company.
func (userService *UserService) accessFor(ctx context.Context, actor model.AuthUser, roleID int64, companyID *int64, machineIDs []int64) (userAccess, error) {
	role, err := userService.roleService.GetByID(ctx, roleID)
	if err != nil {
		return userAccess{}, asInvalidSelection(err, "Role not found", "Please enter valid role")
	}
	if !model.CanAssignRole(actor.Role, role.Name) {
		return userAccess{}, response.Forbidden
	}
	if role.Name == model.RoleCustomer {
		return userAccess{}, nil
	}

	company := actor.Company()
	if actor.Role == model.RoleSuperAdmin {
		if err := userService.checkCustomer(ctx, companyID); err != nil {
			return userAccess{}, err
		}
		company = companyID
	}
	if company == nil {
		return userAccess{}, response.Forbidden // an admin that belongs to no company
	}
	if model.SeesAllCompanyMachines(role.Name) {
		return userAccess{companyID: company}, nil
	}
	if err := userService.machineService.CheckOwnedBy(ctx, *company, machineIDs); err != nil {
		return userAccess{}, err
	}
	return userAccess{companyID: company, machineIDs: machineIDs}, nil
}

// checkCustomer returns "Please select valid customer" unless customerID is an
// existing customer.
func (userService *UserService) checkCustomer(ctx context.Context, customerID *int64) error {
	if customerID == nil {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	isCustomer, err := userService.services.userRepository.IsCustomer(ctx, *customerID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	return nil
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
