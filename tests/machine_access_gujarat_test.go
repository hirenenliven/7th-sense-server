package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"testing"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
)

// gujaratCustomerID is the customer of migrations/seed_machine_access_gujarat.sql (cust_15).
const gujaratCustomerID int64 = 49

// gujaratSeed is the part of the Gujarat sample data the scenario uses.
type gujaratSeed struct {
	countryID, stateID, branchID, lineCID int64
	machines                              map[string]int64 // GJ-VMC-01, …, Poway
}

// loadGujaratSeed reads the sample data, skipping the test when the seed has not been run.
func loadGujaratSeed(t *testing.T) gujaratSeed {
	t.Helper()
	ctx := context.Background()
	seed := gujaratSeed{machines: map[string]int64{}}
	err := testDBPool.QueryRow(ctx, `SELECT co.id, st.id FROM user_country co
		JOIN user_state st ON st.user_country_id = co.id AND st.name = 'Gujarat' AND st.is_deleted = false
		WHERE co.user_id = $1 AND co.name = 'India' AND co.is_deleted = false`, gujaratCustomerID).Scan(&seed.countryID, &seed.stateID)
	if errors.Is(err, pgx.ErrNoRows) {
		t.Skip("run migrations/seed_machine_access_gujarat.sql first")
	}
	if err != nil {
		t.Fatalf("find India › Gujarat: %v", err)
	}
	if err := testDBPool.QueryRow(ctx, `SELECT id FROM branches WHERE code = 'GJRJK01' AND is_deleted = false`).Scan(&seed.branchID); err != nil {
		t.Fatalf("find branch Rajkot Main: %v", err)
	}
	if err := testDBPool.QueryRow(ctx, `SELECT id FROM user_line WHERE user_id = $1 AND name = 'GJ Line C' AND is_deleted = false`, gujaratCustomerID).Scan(&seed.lineCID); err != nil {
		t.Fatalf("find GJ Line C: %v", err)
	}
	rows, err := testDBPool.Query(ctx, `SELECT name, id FROM machines WHERE user_id = $1 AND is_deleted = false`, gujaratCustomerID)
	if err != nil {
		t.Fatalf("find machines: %v", err)
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		var id int64
		if err := rows.Scan(&name, &id); err != nil {
			t.Fatalf("read machine: %v", err)
		}
		seed.machines[name] = id
	}
	return seed
}

// logResponse writes result as indented JSON to the test log (go test -v shows it).
func logResponse(t *testing.T, title string, result apiResponse) {
	t.Helper()
	raw, _ := json.Marshal(result)
	var pretty bytes.Buffer
	_ = json.Indent(&pretty, raw, "", "  ")
	t.Logf("%s:\n%s", title, pretty.String())
}

// machineNames returns the names of the machines of a GET /machines response, sorted.
func machineNames(t *testing.T, result apiResponse) []string {
	t.Helper()
	var machines []struct {
		Name string `json:"name"`
	}
	if len(result.Data) > 0 {
		decodeData(t, result, &machines)
	}
	names := make([]string, 0, len(machines))
	for _, machine := range machines {
		names = append(names, machine.Name)
	}
	slices.Sort(names)
	return names
}

// TestGujaratMachineAccessScenario runs the user form flow on the Gujarat sample data, as the
// customer cust_15: one user gets "India › Gujarat + All machines", another two machines picked
// one by one. A machine added later in Gujarat goes to the first user only. Everything the test
// adds is soft-deleted at the end. Run with -v to see the API responses.
func TestGujaratMachineAccessScenario(t *testing.T) {
	seed := loadGujaratSeed(t)
	suffix := uniqueSuffix()
	customer := newUserAccessRouter(actAs(gujaratCustomerID, model.RoleCustomer))
	company := testCompany{machineFixture: machineFixture{customerID: gujaratCustomerID, branchID: seed.branchID}}
	seenBy := func(user userRecord, role string) []string {
		router := newUserAccessRouter(actAsMember(user.ID, role, gujaratCustomerID))
		return machineNames(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines", nil))
	}

	// User 1: Country India, State Gujarat, Machines = "All machines of this selection"
	allBody := userBody(t, company, "gj_all"+suffix, model.RoleSupervisor, nil)
	allBody.MachineFilter = &dto.MachineFilter{UserCountryID: &seed.countryID, UserStateID: &seed.stateID, AllMachines: true}
	allResult := sendRequest(t, customer, http.MethodPost, apiV1+"/users", allBody)
	logResponse(t, "POST /users — supervisor, India › Gujarat, All machines", allResult)
	expectStatus(t, allResult, http.StatusCreated, "created")
	var allUser userRecord
	decodeData(t, allResult, &allUser)
	t.Cleanup(func() { softDeleteUser(allUser.ID) })

	// User 2: Country India, State Gujarat, Machines = GJ-VMC-01 and GJ-LATHE-01 picked
	picked := []int64{seed.machines["GJ-VMC-01"], seed.machines["GJ-LATHE-01"]}
	pickedBody := userBody(t, company, "gj_pick"+suffix, model.RoleOperator, picked)
	pickedBody.MachineFilter = &dto.MachineFilter{UserCountryID: &seed.countryID, UserStateID: &seed.stateID}
	pickedResult := sendRequest(t, customer, http.MethodPost, apiV1+"/users", pickedBody)
	logResponse(t, "POST /users — operator, India › Gujarat, GJ-VMC-01 + GJ-LATHE-01", pickedResult)
	expectStatus(t, pickedResult, http.StatusCreated, "created")
	var pickedUser userRecord
	decodeData(t, pickedResult, &pickedUser)
	t.Cleanup(func() { softDeleteUser(pickedUser.ID) })

	gujarat := []string{"GJ-HMC-01", "GJ-LATHE-01", "GJ-VMC-01", "GJ-VMC-02"}
	t.Run("All machines: every Gujarat machine, not Poway (Andhra Pradesh)", func(t *testing.T) {
		if len(allUser.MachineIDs) != 0 {
			t.Errorf("machineIds = %v; want none stored with All machines", allUser.MachineIDs)
		}
		got := seenBy(allUser, model.RoleSupervisor)
		t.Logf("GET /machines as the supervisor → %v", got)
		if !slices.Equal(got, gujarat) {
			t.Fatalf("supervisor sees %v; want %v", got, gujarat)
		}
	})

	t.Run("Picked machines: only those two", func(t *testing.T) {
		got := seenBy(pickedUser, model.RoleOperator)
		t.Logf("GET /machines as the operator → %v", got)
		if want := []string{"GJ-LATHE-01", "GJ-VMC-01"}; !slices.Equal(got, want) {
			t.Fatalf("operator sees %v; want %v", got, want)
		}
	})

	t.Run("A machine added later in Gujarat goes to the All machines user only", func(t *testing.T) {
		// Inactive, so the running server's CNC collector does not try to connect to it
		var newID int64
		err := testDBPool.QueryRow(context.Background(), `INSERT INTO machines (user_id, name, serial_number, live_ip, port,
				username, password, local_ip, brand_id, controller_id, model_id, controller_version, accuracy_id, line_id,
				installation_date, hourly_cost, status)
			VALUES ($1, $2, $3, '10.0.0.99', 4840, 'operator', 'secret123', '192.168.10.99', 63, 151, 41, 'V1.0', 1, $4,
				DATE '2026-01-15', 750, 'inactive')
			RETURNING id`, gujaratCustomerID, "GJ-NEW-"+suffix, "SN-GJ-NEW-"+suffix, seed.lineCID).Scan(&newID)
		if err != nil {
			t.Fatalf("add machine on GJ Line C: %v", err)
		}
		t.Cleanup(func() {
			_, _ = testDBPool.Exec(context.Background(), `UPDATE machines SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, newID)
		})
		supervisor := seenBy(allUser, model.RoleSupervisor)
		operator := seenBy(pickedUser, model.RoleOperator)
		t.Logf("after adding GJ-NEW-%s: supervisor → %v, operator → %v", suffix, supervisor, operator)
		if !slices.Contains(supervisor, "GJ-NEW-"+suffix) {
			t.Errorf("supervisor does not get the new machine: %v", supervisor)
		}
		if slices.Contains(operator, "GJ-NEW-"+suffix) {
			t.Errorf("operator got the new machine without picking it: %v", operator)
		}
	})

	t.Run("The edit form loads each user's filter back", func(t *testing.T) {
		result := sendRequest(t, customer, http.MethodGet, fmt.Sprintf("%s/users/%d", apiV1, allUser.ID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var user struct {
			MachineIDs    []int64                  `json:"machineIds"`
			MachineFilter *model.UserMachineFilter `json:"machineFilter"`
		}
		decodeData(t, result, &user)
		t.Logf("GET /users/%d → machineIds %v, machineFilter %+v", allUser.ID, user.MachineIDs, user.MachineFilter)
		if user.MachineFilter == nil || !user.MachineFilter.AllMachines || *user.MachineFilter.UserStateID != seed.stateID {
			t.Fatalf("machineFilter = %+v; want Gujarat with all machines", user.MachineFilter)
		}
	})
}

// softDeleteUser soft-deletes a user a test created and removes its machine access and filter.
func softDeleteUser(id int64) {
	ctx := context.Background()
	_, _ = testDBPool.Exec(ctx, `UPDATE users SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, id)
	_, _ = testDBPool.Exec(ctx, `DELETE FROM user_machine_access WHERE user_id = $1`, id)
	_, _ = testDBPool.Exec(ctx, `DELETE FROM user_machine_filters WHERE user_id = $1`, id)
}
