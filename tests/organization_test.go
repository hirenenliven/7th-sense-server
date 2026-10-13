package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// orgRecord holds the fields the organization levels return; each level fills
// its own parent fields (a city has userStateId and stateName, ...).
type orgRecord struct {
	ID          int64  `json:"id"`
	UserID      int64  `json:"userId"`
	UserStateID int64  `json:"userStateId"`
	StateName   string `json:"stateName"`
	Name        string `json:"name"`
}

// newOrganizationRouter builds the real router with the country, state, city,
// branch and location handlers wired to testDBPool. Requests act as the dev
// user (middleware.Authenticate).
func newOrganizationRouter() http.Handler {
	userCountryService := service.NewUserCountryService(repository.NewUserCountryRepository(testDBPool))
	userStateService := service.NewUserStateService(repository.NewUserStateRepository(testDBPool), userCountryService)
	userCityService := service.NewUserCityService(repository.NewUserCityRepository(testDBPool), userStateService)
	userBranchService := service.NewUserBranchService(repository.NewUserBranchRepository(testDBPool), userCityService)
	userLocationService := service.NewUserLocationService(repository.NewUserLocationRepository(testDBPool), userBranchService)
	return routes.New(routes.Handlers{
		UserCountry:  handler.NewUserCountryHandler(userCountryService),
		UserState:    handler.NewUserStateHandler(userStateService),
		UserCity:     handler.NewUserCityHandler(userCityService),
		UserBranch:   handler.NewUserBranchHandler(userBranchService),
		UserLocation: handler.NewUserLocationHandler(userLocationService),
	})
}

// TestOrganizationHierarchy creates country → state → city → branch for the dev
// user and checks ownership, duplicate, parent and delete rules. The rows are
// soft-deleted at the end and kept for inspection.
func TestOrganizationHierarchy(t *testing.T) {
	router := newOrganizationRouter()
	suffix := uniqueSuffix()

	create := func(t *testing.T, path string, body map[string]any) orgRecord {
		t.Helper()
		result := sendRequest(t, router, http.MethodPost, apiV1+path, body)
		expectStatus(t, result, http.StatusCreated, "created")
		var record orgRecord
		decodeData(t, result, &record)
		return record
	}

	country := create(t, "/user-country", map[string]any{"name": "  Testland " + suffix + "  "})
	if country.Name != "Testland "+suffix || country.UserID != testUserID {
		t.Fatalf("country = %+v; want trimmed name and user %d", country, testUserID)
	}
	state := create(t, "/user-state", map[string]any{"name": "State " + suffix, "userCountryId": country.ID})
	city := create(t, "/user-city", map[string]any{"name": "City " + suffix, "userStateId": state.ID})
	if city.UserStateID != state.ID || city.StateName != state.Name {
		t.Fatalf("city = %+v; want parent %d %q", city, state.ID, state.Name)
	}
	branch := create(t, "/user-branch", map[string]any{"name": "Branch " + suffix, "userCityId": city.ID})

	t.Run("validation", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-city", map[string]any{"name": "  ", "userStateId": state.ID}), http.StatusBadRequest, "Please enter valid name")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-city", map[string]any{"name": "No parent"}), http.StatusBadRequest, "Please enter valid state")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-city", map[string]any{"name": "Unknown parent", "userStateId": unknownID}), http.StatusBadRequest, "Please enter valid state")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-location", map[string]any{"name": "Unknown parent", "userBranchId": unknownID}), http.StatusBadRequest, "Please enter valid branch")
	})

	t.Run("duplicate name under the same parent is rejected, case-insensitively", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-city", map[string]any{"name": "CITY " + suffix, "userStateId": state.ID}), http.StatusBadRequest, "City already exists")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-country", map[string]any{"name": "testland " + suffix}), http.StatusBadRequest, "Country already exists")
	})

	t.Run("list filters by parent", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/user-city?userStateId=%d", apiV1, state.ID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var cities []orgRecord
		decodeData(t, result, &cities)
		if len(cities) != 1 || cities[0].ID != city.ID {
			t.Fatalf("cities of state %d = %+v; want only %d", state.ID, cities, city.ID)
		}
		expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/user-city?userStateId=abc", nil), http.StatusBadRequest, "Please enter valid userStateId")
	})

	t.Run("country and state cannot be updated", func(t *testing.T) {
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, httptest.NewRequest(http.MethodPut, fmt.Sprintf("%s/user-country/%d", apiV1, country.ID), nil))
		if recorder.Code != http.StatusMethodNotAllowed {
			t.Fatalf("PUT /user-country status = %d; want 405", recorder.Code)
		}
	})

	t.Run("update renames and checks the parent", func(t *testing.T) {
		path := fmt.Sprintf("%s/user-branch/%d", apiV1, branch.ID)
		result := sendRequest(t, router, http.MethodPut, path, map[string]any{"name": "Branch renamed " + suffix, "userCityId": city.ID})
		expectStatus(t, result, http.StatusOK, "success")
		var updated orgRecord
		decodeData(t, result, &updated)
		if updated.Name != "Branch renamed "+suffix {
			t.Fatalf("updated = %+v", updated)
		}
		expectStatus(t, sendRequest(t, router, http.MethodPut, path, map[string]any{"name": "x", "userCityId": unknownID}), http.StatusBadRequest, "Please enter valid city")
		expectStatus(t, sendRequest(t, router, http.MethodPut, fmt.Sprintf("%s/user-branch/%d", apiV1, unknownID), map[string]any{"name": "x", "userCityId": city.ID}), http.StatusNotFound, "Branch not found")
	})

	t.Run("another user's rows are invisible", func(t *testing.T) {
		otherCountryID := insertOtherUserCountry(t, "Other "+suffix)
		expectStatus(t, sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/user-country/%d", apiV1, otherCountryID), nil), http.StatusNotFound, "Country not found")
		expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/user-country/%d", apiV1, otherCountryID), nil), http.StatusNotFound, "Country not found")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-state", map[string]any{"name": "Sneaky", "userCountryId": otherCountryID}), http.StatusBadRequest, "Please enter valid country")
	})

	t.Run("delete is blocked while children exist", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/user-city/%d", apiV1, city.ID), nil), http.StatusBadRequest, "This city has branches. Delete its branches first")
	})

	t.Run("delete bottom-up and reuse a deleted name", func(t *testing.T) {
		for _, step := range []struct {
			path  string
			id    int64
			label string
		}{{"/user-branch", branch.ID, "Branch"}, {"/user-city", city.ID, "City"}} {
			expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s%s/%d", apiV1, step.path, step.id), nil), http.StatusOK, step.label+" deleted")
		}
		recreated := create(t, "/user-city", map[string]any{"name": "City " + suffix, "userStateId": state.ID})
		for _, step := range []struct {
			path  string
			id    int64
			label string
		}{{"/user-city", recreated.ID, "City"}, {"/user-state", state.ID, "State"}, {"/user-country", country.ID, "Country"}} {
			expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s%s/%d", apiV1, step.path, step.id), nil), http.StatusOK, step.label+" deleted")
		}
		expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/user-country/%d", apiV1, country.ID), nil), http.StatusNotFound, "Country not found")
	})
}

// insertOtherUserCountry inserts a country owned by a user other than
// testUserID and returns its id; it is soft-deleted when the test ends. The test is skipped when no other user exists.
func insertOtherUserCountry(t *testing.T, name string) int64 {
	t.Helper()
	ctx := context.Background()
	var otherUserID int64
	err := testDBPool.QueryRow(ctx, `SELECT id FROM users WHERE id <> $1 ORDER BY id LIMIT 1`, testUserID).Scan(&otherUserID)
	if err != nil {
		t.Skipf("no second user to own a foreign country: %v", err)
	}
	var id int64
	err = testDBPool.QueryRow(ctx,
		`INSERT INTO user_country (user_id, name) VALUES ($1, $2) RETURNING id`,
		otherUserID, name).Scan(&id)
	if err != nil {
		t.Fatalf("insert other user's country: %v", err)
	}
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `UPDATE user_country SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, id)
	})
	return id
}
