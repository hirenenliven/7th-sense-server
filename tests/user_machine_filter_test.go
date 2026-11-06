package tests

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"testing"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
)

// organizationPathOf returns the country and state of line lineID (organization_paths).
func organizationPathOf(t *testing.T, lineID int64) (countryID, stateID int64) {
	t.Helper()
	err := testDBPool.QueryRow(context.Background(),
		`SELECT user_country_id, user_state_id FROM organization_paths WHERE user_line_id = $1`, lineID).Scan(&countryID, &stateID)
	if err != nil {
		t.Fatalf("find organization path of line %d: %v", lineID, err)
	}
	return countryID, stateID
}

// addTestMachine adds a machine of company on line lineID (as the superadmin) and
// soft-deletes it when the test ends.
func addTestMachine(t *testing.T, company testCompany, name string, lineID int64) int64 {
	t.Helper()
	result := sendRequest(t, newUserAccessRouter(asTestUser), http.MethodPost, apiV1+"/machines",
		with(machineBody(company.machineFixture, name), "lineId", lineID))
	expectStatus(t, result, http.StatusCreated, "created")
	var machine machineRecord
	decodeData(t, result, &machine)
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `UPDATE machines SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, machine.ID)
	})
	return machine.ID
}

// TestMachineFilterAccess checks the user form's machine filter: with all machines on, a
// member may access every company machine under the organization levels named, machines
// added later included; the levels must be one chain of the company's organization.
func TestMachineFilterAccess(t *testing.T) {
	suffix := uniqueSuffix()
	company := newTestCompany(t, "F"+suffix, 2) // two machines on the fixture's line
	other := newTestCompany(t, "G"+suffix, 0)
	customer := newUserAccessRouter(actAs(company.customerID, model.RoleCustomer))
	userPath := func(id int64) string { return fmt.Sprintf("%s/users/%d", apiV1, id) }
	countryID, stateID := organizationPathOf(t, company.lineID)
	otherCountryID, _ := organizationPathOf(t, other.lineID)

	// A second state of the same country, down to a line with one machine
	owner := company.customerID
	state2 := insertTestRow(t, "user_state", `INSERT INTO user_state (user_id, user_country_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, countryID, "State2 "+suffix)
	city2 := insertTestRow(t, "user_city", `INSERT INTO user_city (user_id, user_state_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, state2, "City2 "+suffix)
	branch2 := insertTestRow(t, "branches", `INSERT INTO branches (code, name, user_city_id) VALUES ($1, $2, $3) RETURNING id`, "F2"+suffix, "Branch2 "+suffix, city2)
	location2 := insertTestRow(t, "user_location", `INSERT INTO user_location (user_id, branch_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, branch2, "Location2 "+suffix)
	shop2 := insertTestRow(t, "user_shop", `INSERT INTO user_shop (user_id, user_location_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, location2, "Shop2 "+suffix)
	line2 := insertTestRow(t, "user_line", `INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3) RETURNING id`, owner, shop2, "Line2 "+suffix)
	state2Machine := addTestMachine(t, company, "UA-F2-"+suffix, line2)

	body := userBody(t, company, "flt"+suffix, model.RoleSupervisor, []int64{state2Machine})
	body.MachineFilter = &dto.MachineFilter{UserCountryID: &countryID, UserStateID: &stateID, AllMachines: true}
	supervisor := createUser(t, customer, body)
	supervisorRouter := newUserAccessRouter(actAsMember(supervisor.ID, model.RoleSupervisor, company.customerID))
	seenMachines := func() []int64 {
		return listIDs(t, sendRequest(t, supervisorRouter, http.MethodGet, apiV1+"/machines", nil))
	}
	saveFilter := func(machineIDs []int64, filter *dto.MachineFilter) {
		update := updateBody(body)
		update.MachineIDs = machineIDs
		update.MachineFilter = filter
		expectStatus(t, sendRequest(t, customer, http.MethodPut, userPath(supervisor.ID), update), http.StatusOK, "success")
	}

	t.Run("all machines of the levels named, none stored one by one", func(t *testing.T) {
		if len(supervisor.MachineIDs) != 0 {
			t.Fatalf("supervisor machineIds = %v; want none stored with all machines", supervisor.MachineIDs)
		}
		if got, want := seenMachines(), slices.Sorted(slices.Values(company.machineIDs)); !slices.Equal(got, want) {
			t.Fatalf("supervisor sees %v; want the state's machines %v", got, want)
		}
		expectStatus(t, sendRequest(t, supervisorRouter, http.MethodGet, fmt.Sprintf("%s/machines/%d", apiV1, state2Machine), nil),
			http.StatusNotFound, "Machine not found")
	})

	t.Run("a machine added under the filter is accessible at once", func(t *testing.T) {
		added := addTestMachine(t, company, "UA-F3-"+suffix, company.lineID)
		if got := seenMachines(); !slices.Contains(got, added) {
			t.Fatalf("supervisor sees %v; want the new machine %d too", got, added)
		}
	})

	t.Run("the filter is returned with the user", func(t *testing.T) {
		result := sendRequest(t, customer, http.MethodGet, userPath(supervisor.ID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var user struct {
			MachineFilter *model.UserMachineFilter `json:"machineFilter"`
		}
		decodeData(t, result, &user)
		filter := user.MachineFilter
		if filter == nil || !filter.AllMachines || filter.UserStateID == nil || *filter.UserStateID != stateID || filter.UserLineID != nil {
			t.Fatalf("machineFilter = %+v; want state %d with all machines", filter, stateID)
		}
	})

	t.Run("without all machines the filter only narrows the form", func(t *testing.T) {
		saveFilter([]int64{state2Machine}, &dto.MachineFilter{UserCountryID: &countryID, UserStateID: &state2})
		if got := seenMachines(); !slices.Equal(got, []int64{state2Machine}) {
			t.Fatalf("supervisor sees %v; want only the assigned machine %d", got, state2Machine)
		}
	})

	t.Run("a line narrows to its machines", func(t *testing.T) {
		saveFilter(nil, &dto.MachineFilter{UserShopID: &shop2, UserLineID: &line2, AllMachines: true})
		if got := seenMachines(); !slices.Equal(got, []int64{state2Machine}) {
			t.Fatalf("supervisor sees %v; want the line's machine %d", got, state2Machine)
		}
	})

	t.Run("all machines without levels is every company machine", func(t *testing.T) {
		saveFilter(nil, &dto.MachineFilter{AllMachines: true})
		want := slices.Sorted(slices.Values(append([]int64{state2Machine}, company.machineIDs...)))
		if got := seenMachines(); !slices.Equal(got, want) {
			t.Fatalf("supervisor sees %v; want every company machine %v", got, want)
		}
	})

	t.Run("levels must be one chain of the company's organization", func(t *testing.T) {
		for name, filter := range map[string]dto.MachineFilter{
			"line of another state":   {UserStateID: &stateID, UserLineID: &line2, AllMachines: true},
			"another company's level": {UserCountryID: &otherCountryID, AllMachines: true},
		} {
			t.Run(name, func(t *testing.T) {
				invalid := userBody(t, company, "bad"+suffix, model.RoleOperator, nil)
				invalid.MachineFilter = &filter
				expectStatus(t, sendRequest(t, customer, http.MethodPost, apiV1+"/users", invalid), http.StatusBadRequest, "Please select valid machine filter")
			})
		}
	})

	t.Run("deleting the user removes its filter", func(t *testing.T) {
		expectStatus(t, sendRequest(t, customer, http.MethodDelete, userPath(supervisor.ID), nil), http.StatusOK, "User deleted")
		var count int
		if err := testDBPool.QueryRow(context.Background(), `SELECT COUNT(*) FROM user_machine_filters WHERE user_id = $1`, supervisor.ID).Scan(&count); err != nil {
			t.Fatalf("count machine filters: %v", err)
		}
		if count != 0 {
			t.Fatalf("deleted supervisor still has %d machine filter rows", count)
		}
	})
}
