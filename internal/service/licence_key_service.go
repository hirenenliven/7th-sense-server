package service

import (
	"context"
	"crypto/rand"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// LicenceKeyService holds the business rules for licence keys. The superadmin generates
// keys and saves them for a customer, edits and deletes them (superAdminOnly routes);
// a customer and its admins read their company's (managersOnly, see companyOwner).
type LicenceKeyService struct {
	licenceKeyRepository *repository.LicenceKeyRepository
	userRepository       *repository.UserRepository
}

// NewLicenceKeyService creates a LicenceKeyService that stores keys with
// licenceKeyRepository and checks customers with userRepository.
func NewLicenceKeyService(licenceKeyRepository *repository.LicenceKeyRepository, userRepository *repository.UserRepository) *LicenceKeyService {
	return &LicenceKeyService{licenceKeyRepository: licenceKeyRepository, userRepository: userRepository}
}

// GetAll returns the licence keys actor may read. A non-nil customerID keeps only that
// customer's.
func (licenceKeyService *LicenceKeyService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.LicenceKey, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	return licenceKeyService.licenceKeyRepository.FindAll(ctx, ownerID)
}

// GetByID returns the licence key when actor may read it, else "Licence key not found".
func (licenceKeyService *LicenceKeyService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.LicenceKey, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return model.LicenceKey{}, err
	}
	licenceKey, err := licenceKeyService.licenceKeyRepository.FindByID(ctx, ownerID, id)
	return licenceKey, toLicenceKeyResponseError(err)
}

// Generate returns count new random keys. They are not saved: the superadmin reviews them
// and saves them with Create. A collision with a saved key (about 1 in 2^100) is refused
// by Create.
func (licenceKeyService *LicenceKeyService) Generate(generateRequest dto.GenerateLicenceKeysRequest) (model.GeneratedLicenceKeys, error) {
	licenceKeys := make([]string, 0, generateRequest.Count)
	seen := make(map[string]bool, generateRequest.Count)
	for len(licenceKeys) < generateRequest.Count {
		licenceKey, err := newLicenceKey()
		if err != nil {
			return model.GeneratedLicenceKeys{}, err
		}
		if !seen[licenceKey] {
			seen[licenceKey] = true
			licenceKeys = append(licenceKeys, licenceKey)
		}
	}
	return model.GeneratedLicenceKeys{LicenceKeys: licenceKeys}, nil
}

// Create saves the keys for the customer named, all or none, and returns them. Each
// licence runs for its duration from today.
func (licenceKeyService *LicenceKeyService) Create(ctx context.Context, createRequest dto.CreateLicenceKeysRequest) ([]model.LicenceKey, error) {
	if err := licenceKeyService.checkCustomer(ctx, createRequest.UserID); err != nil {
		return nil, err
	}
	today := time.Now()
	licenceKeys := make([]model.LicenceKey, len(createRequest.LicenceKeys))
	for index, newKey := range createRequest.LicenceKeys {
		licenceKeys[index] = model.LicenceKey{
			UserID:     createRequest.UserID,
			LicenceKey: newKey.LicenceKey,
			Duration:   newKey.Duration,
			ExpiresOn:  licenceExpiresOn(newKey.Duration, today),
			Status:     dto.LicenceKeyStatusActive,
		}
	}
	ids, err := licenceKeyService.licenceKeyRepository.CreateMany(ctx, licenceKeys)
	if err != nil {
		return nil, toLicenceKeyResponseError(err)
	}
	return licenceKeyService.licenceKeyRepository.FindByIDs(ctx, ids)
}

// Update changes the customer, duration and status of the key (a blank status keeps the
// saved one). The end date is counted again from the day the key was saved.
func (licenceKeyService *LicenceKeyService) Update(ctx context.Context, id int64, updateRequest dto.UpdateLicenceKeyRequest) (model.LicenceKey, error) {
	saved, err := licenceKeyService.licenceKeyRepository.FindByID(ctx, nil, id)
	if err != nil {
		return model.LicenceKey{}, toLicenceKeyResponseError(err)
	}
	if err := licenceKeyService.checkCustomer(ctx, updateRequest.UserID); err != nil {
		return model.LicenceKey{}, err
	}
	licenceKey := model.LicenceKey{
		ID:        id,
		UserID:    updateRequest.UserID,
		Duration:  updateRequest.Duration,
		ExpiresOn: licenceExpiresOn(updateRequest.Duration, saved.CreatedAt),
		Status:    updateRequest.Status,
	}
	if err := licenceKeyService.licenceKeyRepository.Update(ctx, licenceKey); err != nil {
		return model.LicenceKey{}, toLicenceKeyResponseError(err)
	}
	updated, err := licenceKeyService.licenceKeyRepository.FindByID(ctx, nil, id)
	return updated, toLicenceKeyResponseError(err)
}

// Delete soft-deletes the licence key.
func (licenceKeyService *LicenceKeyService) Delete(ctx context.Context, id int64) error {
	return toLicenceKeyResponseError(licenceKeyService.licenceKeyRepository.Delete(ctx, id))
}

// checkCustomer returns "Please select valid customer" unless userID is a customer.
func (licenceKeyService *LicenceKeyService) checkCustomer(ctx context.Context, userID int64) error {
	isCustomer, err := licenceKeyService.userRepository.IsCustomer(ctx, userID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	return nil
}

// licenceExpiresOn returns the last valid day ("YYYY-MM-DD") of a licence of duration
// starting on start, or nil when it never ends (LIFETIME). duration is a validated code.
func licenceExpiresOn(duration string, start time.Time) *string {
	plan := dto.LicenceKeyDurations[duration]
	if plan.Lifetime {
		return nil
	}
	// A licence of 1 month started on 15 Jan is valid up to and including 14 Feb
	lastDay := start.AddDate(0, plan.Months, plan.Days-1).Format(time.DateOnly)
	return &lastDay
}

// newLicenceKey returns a random key in the dto licence key format ("K7M2Q-XR9TB-4HWZP-N8CDF").
// The alphabet has 32 characters, so the low 5 bits of a random byte pick one without bias.
func newLicenceKey() (string, error) {
	random := make([]byte, dto.LicenceKeyGroups*dto.LicenceKeyGroupLength)
	if _, err := rand.Read(random); err != nil {
		return "", fmt.Errorf("generate licence key: %w", err)
	}

	var licenceKey strings.Builder
	for index, value := range random {
		if index > 0 && index%dto.LicenceKeyGroupLength == 0 {
			licenceKey.WriteByte('-')
		}
		licenceKey.WriteByte(dto.LicenceKeyAlphabet[value&31])
	}
	return licenceKey.String(), nil
}

// toLicenceKeyResponseError maps repository errors to licence key responses. Other errors
// pass through and are sent as ServerError.
func toLicenceKeyResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Licence key not found")
	case errors.Is(err, repository.ErrDuplicateLicenceKey):
		return response.NewError(http.StatusBadRequest, "Licence key already exists, please generate again")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	default:
		return err
	}
}
