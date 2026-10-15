package dto

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"7th-sense-server/internal/response"
)

// Licence key format: LicenceKeyGroups groups of LicenceKeyGroupLength characters of
// LicenceKeyAlphabet joined by "-" ("K7M2Q-XR9TB-4HWZP-N8CDF"). The alphabet leaves out
// I, O, 0 and 1 (easily misread) and has 32 characters, so each one carries 5 random bits.
const (
	LicenceKeyAlphabet    = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	LicenceKeyGroups      = 4
	LicenceKeyGroupLength = 5
)

// LicenceKeyMaxCount is the most keys generated or saved at once.
const LicenceKeyMaxCount = 100

// Licence key statuses accepted by UpdateLicenceKeyRequest.
const (
	LicenceKeyStatusActive   = "active"
	LicenceKeyStatusInactive = "inactive"
)

// LicenceKeyDuration is a subscription plan: the licence runs Months months plus Days days
// from the day the key is saved; Lifetime never ends.
type LicenceKeyDuration struct {
	Months   int
	Days     int
	Lifetime bool
}

// LicenceKeyTrialDays is the length of a "Considerable Trial" licence.
const LicenceKeyTrialDays = 30

// LicenceKeyDurations are the subscription plans by code (legacy sub_planduration):
// M1 … M12 months, Y1 … Y3 years, LIFETIME, TRIAL (LicenceKeyTrialDays) and TRIAL_1Y (1 year).
var LicenceKeyDurations = licenceKeyDurations()

func licenceKeyDurations() map[string]LicenceKeyDuration {
	durations := map[string]LicenceKeyDuration{
		"LIFETIME": {Lifetime: true},
		"TRIAL":    {Days: LicenceKeyTrialDays},
		"TRIAL_1Y": {Months: 12},
	}
	for months := 1; months <= 12; months++ {
		durations[fmt.Sprintf("M%d", months)] = LicenceKeyDuration{Months: months}
	}
	for years := 1; years <= 3; years++ {
		durations[fmt.Sprintf("Y%d", years)] = LicenceKeyDuration{Months: years * 12}
	}
	return durations
}

var licenceKeyPattern = regexp.MustCompile(`^[A-HJ-NP-Z2-9]{5}(-[A-HJ-NP-Z2-9]{5}){3}$`)

// GenerateLicenceKeysRequest is the body for POST /licence-keys/generate.
type GenerateLicenceKeysRequest struct {
	Count int `json:"count"` // number of keys, 1 … LicenceKeyMaxCount
}

// Validate checks the count is usable.
func (request *GenerateLicenceKeysRequest) Validate() error {
	if request.Count < 1 || request.Count > LicenceKeyMaxCount {
		return response.NewError(http.StatusBadRequest, "Please enter valid number of keys (1 to 100)")
	}
	return nil
}

// NewLicenceKey is one generated key to save with its subscription plan.
type NewLicenceKey struct {
	LicenceKey string `json:"licenceKey"`
	Duration   string `json:"duration"` // a LicenceKeyDurations code
}

// CreateLicenceKeysRequest is the body for POST /licence-keys: the generated keys to save
// for a customer, each with its subscription plan.
type CreateLicenceKeysRequest struct {
	UserID      int64           `json:"userId"` // customer the keys belong to
	LicenceKeys []NewLicenceKey `json:"licenceKeys"`
}

// Validate trims and uppercases the keys and checks they are well formed and not repeated.
func (request *CreateLicenceKeysRequest) Validate() error {
	if request.UserID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please select customer")
	}
	if len(request.LicenceKeys) == 0 {
		return response.NewError(http.StatusBadRequest, "Please generate licence keys")
	}
	if len(request.LicenceKeys) > LicenceKeyMaxCount {
		return response.NewError(http.StatusBadRequest, "Please save at most 100 licence keys at once")
	}

	seen := make(map[string]bool, len(request.LicenceKeys))
	for index := range request.LicenceKeys {
		newKey := &request.LicenceKeys[index]
		newKey.LicenceKey = strings.ToUpper(strings.TrimSpace(newKey.LicenceKey))
		if !ValidLicenceKey(newKey.LicenceKey) {
			return response.NewError(http.StatusBadRequest, "Please enter valid licence key")
		}
		if seen[newKey.LicenceKey] {
			return response.NewError(http.StatusBadRequest, "Licence key is repeated")
		}
		seen[newKey.LicenceKey] = true
		if err := validDuration(&newKey.Duration); err != nil {
			return err
		}
	}
	return nil
}

// UpdateLicenceKeyRequest is the body for PUT /licence-keys/{id}. The key itself never
// changes; a blank status keeps the saved one. A new Duration moves the end date (counted
// from the day the key was saved).
type UpdateLicenceKeyRequest struct {
	UserID   int64  `json:"userId"`
	Duration string `json:"duration"`
	Status   string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateLicenceKeyRequest) Validate() error {
	if request.UserID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please select customer")
	}
	if err := validDuration(&request.Duration); err != nil {
		return err
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != LicenceKeyStatusActive && request.Status != LicenceKeyStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

// validDuration uppercases *duration and checks it is a LicenceKeyDurations code.
func validDuration(duration *string) error {
	*duration = strings.ToUpper(strings.TrimSpace(*duration))
	if _, ok := LicenceKeyDurations[*duration]; !ok {
		return response.NewError(http.StatusBadRequest, "Please select valid subscription duration")
	}
	return nil
}

// ValidLicenceKey reports whether key has the licence key format ("K7M2Q-XR9TB-4HWZP-N8CDF").
func ValidLicenceKey(key string) bool {
	return licenceKeyPattern.MatchString(key)
}
