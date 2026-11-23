package tests

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"testing"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/service"
)

// locationRecord holds the organization path fields of a machine response.
type locationRecord struct {
	ID               int64    `json:"id"`
	UserCountryID    *int64   `json:"userCountryId"`
	UserStateID      *int64   `json:"userStateId"`
	UserCityID       *int64   `json:"userCityId"`
	BranchID         *int64   `json:"branchId"`
	UserLocationID   *int64   `json:"userLocationId"`
	UserShopID       *int64   `json:"userShopId"`
	LineID           *int64   `json:"lineId"`
	UserCellID       *int64   `json:"userCellId"`
	UserSupervisorID *int64   `json:"userSupervisorId"`
	UserStateName    string   `json:"userStateName"`
	UserCellName     string   `json:"userCellName"`
	CustomerFilters  []string `json:"customerFilters"`
}

// organizationFixture is a machine fixture plus every level of its customer's organization:
// country … line (newMachineFixture) and a cell with a supervisor under the line.
type organizationFixture struct {
	machineFixture
	countryID, stateID, cityID, locationID, shopID, cellID, supervisorID int64
}

func newOrganizationFixture(t *testing.T, suffix string) organizationFixture {
	t.Helper()
	fixture := organizationFixture{machineFixture: newMachineFixture(t, suffix)}
	err := testDBPool.QueryRow(context.Background(),
		`SELECT user_country_id, user_state_id, user_city_id, user_location_id, user_shop_id
		FROM organization_paths WHERE user_line_id = $1`, fixture.lineID,
	).Scan(&fixture.countryID, &fixture.stateID, &fixture.cityID, &fixture.locationID, &fixture.shopID)
	if err != nil {
		t.Fatalf("read organization path: %v", err)
	}
	customer := fixture.customerID
	fixture.cellID = insertTestRow(t, "user_cell", `INSERT INTO user_cell (user_id, user_line_id, name) VALUES ($1, $2, $3) RETURNING id`, customer, fixture.lineID, "Cell "+suffix)
	fixture.supervisorID = insertTestRow(t, "user_supervisor", `INSERT INTO user_supervisor (user_id, user_cell_id, name) VALUES ($1, $2, $3) RETURNING id`, customer, fixture.cellID, "Supervisor "+suffix)
	return fixture
}

// setCustomerFilters saves filters as the customer's filters (user_profiles.filters).
func setCustomerFilters(t *testing.T, customerID int64, filters ...string) {
	t.Helper()
	_, err := testDBPool.Exec(context.Background(),
		`INSERT INTO user_profiles (user_id, filters) VALUES ($1, $2)
		ON CONFLICT (user_id) DO UPDATE SET filters = EXCLUDED.filters`, customerID, filters)
	if err != nil {
		t.Fatalf("save customer filters: %v", err)
	}
}

// locationBody is a valid machine body for fixture placed at location (no lineId).
func locationBody(fixture organizationFixture, name string, location map[string]any) map[string]any {
	return with(with(machineBody(fixture.machineFixture, name), "lineId", nil), "location", location)
}

// createLocatedMachine adds a machine through router and soft-deletes it when the test ends.
func createLocatedMachine(t *testing.T, router http.Handler, body map[string]any) locationRecord {
	t.Helper()
	result := sendRequest(t, router, http.MethodPost, apiV1+"/machines", body)
	expectStatus(t, result, http.StatusCreated, "created")
	var created locationRecord
	decodeData(t, result, &created)
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `UPDATE machines SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, created.ID)
	})
	return created
}

// savedLevel returns the machine's saved level columns that are set, as "column=id".
func savedLevel(t *testing.T, machineID int64) []string {
	t.Helper()
	var saved []string
	err := testDBPool.QueryRow(context.Background(),
		`SELECT array_remove(ARRAY[
			'user_country_id=' || user_country_id, 'user_state_id=' || user_state_id,
			'user_city_id=' || user_city_id, 'branch_id=' || branch_id,
			'user_location_id=' || user_location_id, 'user_shop_id=' || user_shop_id,
			'line_id=' || line_id, 'user_cell_id=' || user_cell_id,
			'user_supervisor_id=' || user_supervisor_id], NULL)
		FROM machines WHERE id = $1`, machineID).Scan(&saved)
	if err != nil {
		t.Fatalf("read saved level: %v", err)
	}
	return saved
}

func equalID(id *int64, want int64) bool { return id != nil && *id == want }

func TestMachineLocation(t *testing.T) {
	router := newMachineRouter(asTestUser) // superadmin
	suffix := uniqueSuffix()
	fixture := newOrganizationFixture(t, suffix)
	setCustomerFilters(t, fixture.customerID, "location", "country", "state")

	t.Run("location filters are the customer's, in hierarchy order", func(t *testing.T) {
		var filters struct {
			Filters []string `json:"filters"`
		}
		decodeData(t, sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/machines/location-filters?userId=%d", apiV1, fixture.customerID), nil), &filters)
		if !slices.Equal(filters.Filters, []string{"country", "state", "location"}) {
			t.Fatalf("filters = %v; want [country state location]", filters.Filters)
		}
		expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines/location-filters", nil), http.StatusBadRequest, "Please select customer")

		customerRouter := newMachineRouter(actAs(fixture.customerID, model.RoleCustomer))
		decodeData(t, sendRequest(t, customerRouter, http.MethodGet, apiV1+"/machines/location-filters", nil), &filters)
		if len(filters.Filters) != 3 {
			t.Fatalf("customer's own filters = %v", filters.Filters)
		}
		expectStatus(t, sendRequest(t, customerRouter, http.MethodGet, fmt.Sprintf("%s/machines/location-filters?userId=%d", apiV1, testUserID), nil), http.StatusForbidden, "You do not have permission to do this")
	})

	t.Run("every level is optional", func(t *testing.T) {
		onlyState := createLocatedMachine(t, router, locationBody(fixture, "State-"+suffix, map[string]any{"userStateId": fixture.stateID}))
		if !equalID(onlyState.UserCountryID, fixture.countryID) || !equalID(onlyState.UserStateID, fixture.stateID) || onlyState.UserCityID != nil {
			t.Errorf("machine at a state = %+v; want the state and its country only", onlyState)
		}
		none := createLocatedMachine(t, router, locationBody(fixture, "Empty-"+suffix, map[string]any{}))
		if none.UserCountryID != nil {
			t.Errorf("machine without location = %+v", none)
		}
	})

	t.Run("the levels must be one chain of the customer's organization", func(t *testing.T) {
		other := newOrganizationFixture(t, suffix+"o")
		body := locationBody(fixture, "Chain-"+suffix, map[string]any{
			"userCountryId": fixture.countryID, "userStateId": fixture.stateID, "userLocationId": other.locationID,
		})
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", body), http.StatusBadRequest, "Please select valid location")
		body = locationBody(fixture, "Bad-"+suffix, map[string]any{"userCountryId": fixture.countryID, "userStateId": -1, "userLocationId": fixture.locationID})
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", body), http.StatusBadRequest, "Please select valid location")
	})

	t.Run("a machine is saved at the lowest level and returns its whole path", func(t *testing.T) {
		created := createLocatedMachine(t, router, locationBody(fixture, "Loc-"+suffix, map[string]any{
			"userCountryId": fixture.countryID, "userStateId": fixture.stateID, "userLocationId": fixture.locationID,
		}))
		switch {
		case !equalID(created.UserCountryID, fixture.countryID) || !equalID(created.UserStateID, fixture.stateID) ||
			!equalID(created.UserCityID, fixture.cityID) || !equalID(created.BranchID, fixture.branchID) ||
			!equalID(created.UserLocationID, fixture.locationID):
			t.Errorf("path = %+v; want country … location of the fixture", created)
		case created.UserShopID != nil || created.LineID != nil || created.UserCellID != nil:
			t.Errorf("levels below the location = %v / %v / %v; want none", created.UserShopID, created.LineID, created.UserCellID)
		case created.UserStateName != "State "+suffix || !slices.Equal(created.CustomerFilters, []string{"country", "state", "location"}):
			t.Errorf("state name %q, customer filters %v", created.UserStateName, created.CustomerFilters)
		}
		if saved := savedLevel(t, created.ID); !slices.Equal(saved, []string{fmt.Sprintf("user_location_id=%d", fixture.locationID)}) {
			t.Errorf("saved levels = %v; want only the location", saved)
		}
	})

	t.Run("a machine saved at a supervisor is reached through its line", func(t *testing.T) {
		setCustomerFilters(t, fixture.customerID, "line", "supervisor")
		created := createLocatedMachine(t, router, locationBody(fixture, "Sup-"+suffix, map[string]any{
			"userLineId": fixture.lineID, "userSupervisorId": fixture.supervisorID,
		}))
		if !equalID(created.LineID, fixture.lineID) || !equalID(created.UserCellID, fixture.cellID) || created.UserCellName != "Cell "+suffix {
			t.Fatalf("path = %+v; want line %d and cell %d", created, fixture.lineID, fixture.cellID)
		}
		if saved := savedLevel(t, created.ID); !slices.Equal(saved, []string{fmt.Sprintf("user_supervisor_id=%d", fixture.supervisorID)}) {
			t.Errorf("saved levels = %v; want only the supervisor", saved)
		}

		// An operator whose machine filter takes every machine of the line sees it
		operatorID := createAuthTestUser(t, "machine_location_operator_"+suffix, model.RoleOperator)
		if _, err := testDBPool.Exec(context.Background(),
			`INSERT INTO user_machine_filters (user_id, user_line_id, all_machines) VALUES ($1, $2, true)`, operatorID, fixture.lineID); err != nil {
			t.Fatalf("save machine filter: %v", err)
		}
		t.Cleanup(func() {
			_, _ = testDBPool.Exec(context.Background(), `DELETE FROM user_machine_filters WHERE user_id = $1`, operatorID)
		})
		var visible []machineRecord
		decodeData(t, sendRequest(t, newMachineRouter(actAs(operatorID, model.RoleOperator)), http.MethodGet, apiV1+"/machines", nil), &visible)
		if !containsMachineID(visible, created.ID) {
			t.Errorf("operator's machines = %+v; want machine %d (saved under their line)", visible, created.ID)
		}

		// The supervisor it is saved at cannot be deleted
		supervisorService := service.NewUserSupervisorService(repository.NewUserSupervisorRepository(testDBPool), nil)
		err := supervisorService.Delete(context.Background(), fixture.customerID, fixture.supervisorID)
		if err == nil || err.Error() != "This supervisor has machines. Move or delete them first" {
			t.Errorf("delete supervisor = %v; want it refused", err)
		}
	})

	t.Run("removing filters from the customer keeps its machines' places", func(t *testing.T) {
		setCustomerFilters(t, fixture.customerID, "country", "state", "location")
		created := createLocatedMachine(t, router, locationBody(fixture, "Kept-"+suffix, map[string]any{
			"userCountryId": fixture.countryID, "userStateId": fixture.stateID, "userLocationId": fixture.locationID,
		}))
		setCustomerFilters(t, fixture.customerID, "country")

		var read locationRecord
		decodeData(t, sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/machines/%d", apiV1, created.ID), nil), &read)
		if !equalID(read.UserLocationID, fixture.locationID) || !slices.Equal(read.CustomerFilters, []string{"country"}) {
			t.Fatalf("after removing filters = %+v; want the location kept and filters [country]", read)
		}
		// Saving again with the old, hidden levels still names one chain: accepted, place kept
		body := locationBody(fixture, "Kept-"+suffix, map[string]any{
			"userCountryId": fixture.countryID, "userStateId": fixture.stateID, "userLocationId": fixture.locationID,
		})
		expectStatus(t, sendRequest(t, router, http.MethodPut, fmt.Sprintf("%s/machines/%d", apiV1, created.ID), body), http.StatusOK, "success")
		// Saving with only the remaining level moves it up to the country
		body = locationBody(fixture, "Kept-"+suffix, map[string]any{"userCountryId": fixture.countryID})
		expectStatus(t, sendRequest(t, router, http.MethodPut, fmt.Sprintf("%s/machines/%d", apiV1, created.ID), body), http.StatusOK, "success")
		if saved := savedLevel(t, created.ID); !slices.Equal(saved, []string{fmt.Sprintf("user_country_id=%d", fixture.countryID)}) {
			t.Errorf("saved levels = %v; want only the country", saved)
		}
	})
}
