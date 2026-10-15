package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

func newLicenceKeyRouter(authenticate middleware.Middleware) http.Handler {
	return routes.New(routes.Handlers{
		LicenceKey: handler.NewLicenceKeyHandler(service.NewLicenceKeyService(
			repository.NewLicenceKeyRepository(testDBPool), repository.NewUserRepository(testDBPool))),
	}, authenticate)
}

// generateLicenceKeys calls POST /licence-keys/generate and returns the keys.
func generateLicenceKeys(t *testing.T, router http.Handler, count int) []string {
	t.Helper()
	result := sendRequest(t, router, http.MethodPost, apiV1+"/licence-keys/generate", map[string]any{"count": count})
	expectStatus(t, result, http.StatusOK, "success")
	var generated model.GeneratedLicenceKeys
	decodeData(t, result, &generated)
	return generated.LicenceKeys
}

// licenceKeyItems pairs each key with duration for the POST /licence-keys body.
func licenceKeyItems(duration string, keys ...string) []map[string]string {
	items := make([]map[string]string, len(keys))
	for index, key := range keys {
		items[index] = map[string]string{"licenceKey": key, "duration": duration}
	}
	return items
}

func TestLicenceKeyCRUD(t *testing.T) {
	router := newLicenceKeyRouter(asTestUser) // superadmin
	suffix := uniqueSuffix()
	customerID := createAuthTestUser(t, "licence_customer_"+suffix, model.RoleCustomer)
	otherCustomerID := createAuthTestUser(t, "licence_other_"+suffix, model.RoleCustomer)
	adminID := createAuthTestUser(t, "licence_admin_"+suffix, model.RoleAdmin)
	licenceKeyPath := func(id int64) string { return fmt.Sprintf("%s/licence-keys/%d", apiV1, id) }

	var keys []string
	t.Run("generate", func(t *testing.T) {
		keys = generateLicenceKeys(t, router, 3)
		if len(keys) != 3 || keys[0] == keys[1] || keys[1] == keys[2] {
			t.Fatalf("keys = %v; want 3 different keys", keys)
		}
		for _, key := range keys {
			if !dto.ValidLicenceKey(key) {
				t.Errorf("key %q does not have the licence key format", key)
			}
		}
	})
	if len(keys) != 3 {
		t.FailNow()
	}

	t.Run("generate validation", func(t *testing.T) {
		for _, count := range []int{0, dto.LicenceKeyMaxCount + 1} {
			result := sendRequest(t, router, http.MethodPost, apiV1+"/licence-keys/generate", map[string]any{"count": count})
			expectStatus(t, result, http.StatusBadRequest, "Please enter valid number of keys (1 to 100)")
		}
	})

	var created []model.LicenceKey
	t.Run("create", func(t *testing.T) {
		// Lowercase and padded keys are normalized
		items := append(licenceKeyItems("m3", " "+strings.ToLower(keys[0])+" "), licenceKeyItems("LIFETIME", keys[1])...)
		body := map[string]any{"userId": customerID, "licenceKeys": items}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/licence-keys", body)
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		if len(created) != 2 || created[0].LicenceKey != keys[0] || created[0].UserID != customerID ||
			created[0].Status != "active" || created[0].CustomerName == "" {
			t.Fatalf("created = %+v", created)
		}
		wantExpiry := time.Now().AddDate(0, 3, -1).Format(time.DateOnly)
		switch {
		case created[0].Duration != "M3" || created[0].ExpiresOn == nil || *created[0].ExpiresOn != wantExpiry || created[0].IsExpired:
			t.Errorf("3 month key = %s %v %v; want M3 %s not expired", created[0].Duration, created[0].ExpiresOn, created[0].IsExpired, wantExpiry)
		case created[1].Duration != "LIFETIME" || created[1].ExpiresOn != nil:
			t.Errorf("lifetime key = %s %v; want no end date", created[1].Duration, created[1].ExpiresOn)
		}
	})
	if len(created) != 2 {
		t.FailNow()
	}

	t.Run("duplicate key is refused, nothing is saved", func(t *testing.T) {
		body := map[string]any{"userId": customerID, "licenceKeys": licenceKeyItems("TRIAL", keys[2], keys[0])}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/licence-keys", body)
		expectStatus(t, result, http.StatusBadRequest, "Licence key already exists, please generate again")

		result = sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/licence-keys?userId=%d", apiV1, customerID), nil)
		var saved []model.LicenceKey
		decodeData(t, result, &saved)
		if len(saved) != 2 {
			t.Errorf("saved %d keys; want 2 (the batch is all or none)", len(saved))
		}
	})

	t.Run("create validation", func(t *testing.T) {
		cases := map[string]struct {
			body    map[string]any
			message string
		}{
			"no customer":  {map[string]any{"licenceKeys": licenceKeyItems("TRIAL", keys[2])}, "Please select customer"},
			"no keys":      {map[string]any{"userId": customerID, "licenceKeys": []string{}}, "Please generate licence keys"},
			"bad format":   {map[string]any{"userId": customerID, "licenceKeys": licenceKeyItems("TRIAL", "ABC-123")}, "Please enter valid licence key"},
			"repeated key": {map[string]any{"userId": customerID, "licenceKeys": licenceKeyItems("TRIAL", keys[2], keys[2])}, "Licence key is repeated"},
			"bad duration": {map[string]any{"userId": customerID, "licenceKeys": licenceKeyItems("M13", keys[2])}, "Please select valid subscription duration"},
			"not customer": {map[string]any{"userId": adminID, "licenceKeys": licenceKeyItems("TRIAL", keys[2])}, "Please select valid customer"},
		}
		for name, testCase := range cases {
			t.Run(name, func(t *testing.T) {
				result := sendRequest(t, router, http.MethodPost, apiV1+"/licence-keys", testCase.body)
				expectStatus(t, result, http.StatusBadRequest, testCase.message)
			})
		}
	})

	t.Run("update", func(t *testing.T) {
		body := map[string]any{"userId": otherCustomerID, "duration": "TRIAL", "status": "inactive"}
		result := sendRequest(t, router, http.MethodPut, licenceKeyPath(created[0].ID), body)
		expectStatus(t, result, http.StatusOK, "success")
		var updated model.LicenceKey
		decodeData(t, result, &updated)
		wantExpiry := created[0].CreatedAt.AddDate(0, 0, dto.LicenceKeyTrialDays-1).Format(time.DateOnly)
		if updated.UserID != otherCustomerID || updated.Status != "inactive" || updated.LicenceKey != keys[0] ||
			updated.Duration != "TRIAL" || updated.ExpiresOn == nil || *updated.ExpiresOn != wantExpiry {
			t.Errorf("updated = %+v; want trial ending %s", updated, wantExpiry)
		}
	})

	t.Run("company scope", func(t *testing.T) {
		customer := newLicenceKeyRouter(actAs(customerID, model.RoleCustomer))
		result := sendRequest(t, customer, http.MethodGet, licenceKeyPath(created[1].ID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		result = sendRequest(t, customer, http.MethodGet, licenceKeyPath(created[0].ID), nil) // moved to the other customer
		expectStatus(t, result, http.StatusNotFound, "Licence key not found")

		// Only the superadmin generates, saves, edits and deletes
		result = sendRequest(t, customer, http.MethodPost, apiV1+"/licence-keys/generate", map[string]any{"count": 1})
		expectStatus(t, result, http.StatusForbidden, "You do not have permission to do this")
		result = sendRequest(t, customer, http.MethodDelete, licenceKeyPath(created[1].ID), nil)
		expectStatus(t, result, http.StatusForbidden, "You do not have permission to do this")
	})

	t.Run("delete", func(t *testing.T) {
		for _, licenceKey := range created {
			result := sendRequest(t, router, http.MethodDelete, licenceKeyPath(licenceKey.ID), nil)
			expectStatus(t, result, http.StatusOK, "Licence key deleted")
		}
		result := sendRequest(t, router, http.MethodGet, licenceKeyPath(created[0].ID), nil)
		expectStatus(t, result, http.StatusNotFound, "Licence key not found")
	})
}
