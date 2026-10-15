package tests

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
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
	UserCityID  int64  `json:"userCityId"`
	CityName    string `json:"cityName"`
	BranchID    int64  `json:"branchId"`
	BranchName  string `json:"branchName"`
	Name        string `json:"name"`
}

// newOrganizationRouter builds the real router with the country, state, city,
// branch and location handlers wired to testDBPool. Requests act as
// testUserID (asTestUser).
func newOrganizationRouter() http.Handler {
	userCountryService := service.NewUserCountryService(repository.NewUserCountryRepository(testDBPool))
	userStateService := service.NewUserStateService(repository.NewUserStateRepository(testDBPool), userCountryService)
	userCityService := service.NewUserCityService(repository.NewUserCityRepository(testDBPool), userStateService)
	branchService := service.NewBranchService(repository.NewBranchRepository(testDBPool), userCityService)
	userLocationService := service.NewUserLocationService(repository.NewUserLocationRepository(testDBPool), branchService)
	return routes.New(routes.Handlers{
		UserCountry:  handler.NewUserCountryHandler(userCountryService),
		UserState:    handler.NewUserStateHandler(userStateService),
		UserCity:     handler.NewUserCityHandler(userCityService),
		Branch:       handler.NewBranchHandler(branchService),
		UserLocation: handler.NewUserLocationHandler(userLocationService),
	}, asTestUser)
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
	branchCode := "T" + suffix
	branch := create(t, "/branches", map[string]any{"code": branchCode, "name": "Branch " + suffix, "userCityId": city.ID})
	if branch.UserCityID != city.ID || branch.CityName != city.Name {
		t.Fatalf("branch = %+v; want city %d %q", branch, city.ID, city.Name)
	}
	location := create(t, "/user-location", map[string]any{"name": "Location " + suffix, "branchId": branch.ID})
	if location.BranchID != branch.ID || location.BranchName != branch.Name {
		t.Fatalf("location = %+v; want branch %d %q", location, branch.ID, branch.Name)
	}

	t.Run("validation", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-city", map[string]any{"name": "  ", "userStateId": state.ID}), http.StatusBadRequest, "Please enter valid name")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-city", map[string]any{"name": "No parent"}), http.StatusBadRequest, "Please enter valid state")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-city", map[string]any{"name": "Unknown parent", "userStateId": unknownID}), http.StatusBadRequest, "Please enter valid state")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-location", map[string]any{"name": "Unknown parent", "branchId": unknownID}), http.StatusBadRequest, "Please enter valid branch")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/branches", map[string]any{"code": branchCode + "X", "name": "Unknown city", "userCityId": unknownID}), http.StatusBadRequest, "Please enter valid city")
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
		path := fmt.Sprintf("%s/user-location/%d", apiV1, location.ID)
		result := sendRequest(t, router, http.MethodPut, path, map[string]any{"name": "Location renamed " + suffix, "branchId": branch.ID})
		expectStatus(t, result, http.StatusOK, "success")
		var updated orgRecord
		decodeData(t, result, &updated)
		if updated.Name != "Location renamed "+suffix {
			t.Fatalf("updated = %+v", updated)
		}
		expectStatus(t, sendRequest(t, router, http.MethodPut, path, map[string]any{"name": "x", "branchId": unknownID}), http.StatusBadRequest, "Please enter valid branch")
		expectStatus(t, sendRequest(t, router, http.MethodPut, fmt.Sprintf("%s/user-location/%d", apiV1, unknownID), map[string]any{"name": "x", "branchId": branch.ID}), http.StatusNotFound, "Location not found")
	})

	t.Run("another user's rows are invisible", func(t *testing.T) {
		otherCountryID := insertOtherUserCountry(t, "Other "+suffix)
		expectStatus(t, sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/user-country/%d", apiV1, otherCountryID), nil), http.StatusNotFound, "Country not found")
		expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/user-country/%d", apiV1, otherCountryID), nil), http.StatusNotFound, "Country not found")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/user-state", map[string]any{"name": "Sneaky", "userCountryId": otherCountryID}), http.StatusBadRequest, "Please enter valid country")
	})

	t.Run("delete is blocked while children exist", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/user-city/%d", apiV1, city.ID), nil), http.StatusBadRequest, "This city has branches. Delete its branches first")
		expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/branches/%d", apiV1, branch.ID), nil), http.StatusBadRequest, "This branch has locations. Delete its locations first")
	})

	t.Run("delete bottom-up and reuse a deleted name", func(t *testing.T) {
		for _, step := range []struct {
			path  string
			id    int64
			label string
		}{{"/user-location", location.ID, "Location"}, {"/branches", branch.ID, "Branch"}, {"/user-city", city.ID, "City"}} {
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

// TestOrganizationForAnotherUser checks ?userId=: an admin reads and changes
// another user's rows, a non-admin may only name themself.
func TestOrganizationForAnotherUser(t *testing.T) {
	adminRouter := newOrganizationRouter()
	suffix := uniqueSuffix()
	otherUserID := createAuthTestUser(t, "org_owner_"+suffix, model.RoleCustomer)
	otherPath := func(path string) string { return fmt.Sprintf("%s%s?userId=%d", apiV1, path, otherUserID) }

	result := sendRequest(t, adminRouter, http.MethodPost, otherPath("/user-country"), map[string]any{"name": "Owned " + suffix})
	expectStatus(t, result, http.StatusCreated, "created")
	var country orgRecord
	decodeData(t, result, &country)
	if country.UserID != otherUserID {
		t.Fatalf("country = %+v; want user %d", country, otherUserID)
	}

	t.Run("admin lists the named user's rows, or every user's without userId", func(t *testing.T) {
		var others, all []orgRecord
		decodeData(t, sendRequest(t, adminRouter, http.MethodGet, otherPath("/user-country"), nil), &others)
		decodeData(t, sendRequest(t, adminRouter, http.MethodGet, apiV1+"/user-country", nil), &all)
		if !containsOrgID(others, country.ID) || !containsOrgID(all, country.ID) {
			t.Fatalf("country %d: in other user's list = %v, in unfiltered list = %v; want true, true", country.ID, containsOrgID(others, country.ID), containsOrgID(all, country.ID))
		}
		for _, record := range others {
			if record.UserID != otherUserID {
				t.Fatalf("named user's list has %+v; want only user %d", record, otherUserID)
			}
		}
	})

	t.Run("parents are checked against the named user", func(t *testing.T) {
		expectStatus(t, sendRequest(t, adminRouter, http.MethodPost, otherPath("/user-state"), map[string]any{"name": "State " + suffix, "userCountryId": country.ID}), http.StatusCreated, "created")
		expectStatus(t, sendRequest(t, adminRouter, http.MethodPost, apiV1+"/user-state", map[string]any{"name": "State " + suffix, "userCountryId": country.ID}), http.StatusBadRequest, "Please enter valid country")
	})

	t.Run("invalid or unknown user", func(t *testing.T) {
		expectStatus(t, sendRequest(t, adminRouter, http.MethodGet, apiV1+"/user-country?userId=abc", nil), http.StatusBadRequest, "Please enter valid userId")
		expectStatus(t, sendRequest(t, adminRouter, http.MethodPost, fmt.Sprintf("%s/user-country?userId=%d", apiV1, unknownID), map[string]any{"name": "Nobody " + suffix}), http.StatusBadRequest, "Please enter valid user")
	})

	t.Run("non-admin may only name themself", func(t *testing.T) {
		asOtherUser := func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
				user := model.AuthUser{ID: otherUserID, Role: model.RoleCustomer}
				next.ServeHTTP(writer, request.WithContext(middleware.WithAuthUser(request.Context(), user)))
			})
		}
		userCountryService := service.NewUserCountryService(repository.NewUserCountryRepository(testDBPool))
		customerRouter := routes.New(routes.Handlers{UserCountry: handler.NewUserCountryHandler(userCountryService)}, asOtherUser)

		expectStatus(t, sendRequest(t, customerRouter, http.MethodGet, otherPath("/user-country"), nil), http.StatusOK, "success")
		expectStatus(t, sendRequest(t, customerRouter, http.MethodGet, fmt.Sprintf("%s/user-country?userId=%d", apiV1, testUserID), nil), http.StatusForbidden, "You do not have permission to do this")
	})
}

func containsOrgID(records []orgRecord, id int64) bool {
	for _, record := range records {
		if record.ID == id {
			return true
		}
	}
	return false
}
