package tests

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

const (
	reasonsPath           = apiV1 + "/machine-mode-reasons"
	reasonOutsideMessage  = "Reason time must be inside the period"
	reasonOverlapMessage  = "This time overlaps another reason of this period"
	machineNotFoundText   = "Machine not found"
	invalidDowntimeReason = "Please select valid downtime reason"
	reasonDeletedMessage  = "Machine mode reason deleted"
)

// newMachineModeReasonRouter serves /machine-mode-reasons, acting as authenticate's user.
func newMachineModeReasonRouter(authenticate middleware.Middleware) http.Handler {
	return routes.New(routes.Handlers{MachineModeReason: handler.NewMachineModeReasonHandler(
		service.NewMachineModeReasonService(testDBPool, newMachineService()))}, authenticate)
}

// insertModePeriod adds a machine_mode_log row and removes it (with its reasons) when the test
// ends. A nil end is a running period.
func insertModePeriod(t *testing.T, machineID int64, modeID int, startedAt time.Time, endedAt *time.Time) int64 {
	t.Helper()
	var id int64
	err := testDBPool.QueryRow(context.Background(), `INSERT INTO machine_mode_log (machine_id, mode_id, started_at, ended_at)
		VALUES ($1, $2, $3, $4) RETURNING id`, machineID, modeID, startedAt, endedAt).Scan(&id)
	if err != nil {
		t.Fatalf("insert mode period: %v", err)
	}
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `DELETE FROM machine_mode_log WHERE id = $1`, id)
	})
	return id
}

// insertReason adds a master reason to table for customerID and returns its id.
func insertReason(t *testing.T, table string, customerID int64, reason, status string) int64 {
	t.Helper()
	return insertTestRow(t, table, `INSERT INTO `+table+` (user_id, reason, status) VALUES ($1, $2, $3) RETURNING id`, customerID, reason, status)
}

func reasonBody(periodID, reasonID int64, startedAt, endedAt time.Time) dto.CreateMachineModeReasonRequest {
	return dto.CreateMachineModeReasonRequest{PeriodID: periodID, ReasonID: reasonID, StartedAt: startedAt, EndedAt: endedAt}
}

// TestMachineModeReasons covers recording reasons against the Gantt chart's periods: which periods
// take one, the time and reason checks, edit, delete and who may do it.
func TestMachineModeReasons(t *testing.T) {
	suffix := uniqueSuffix()
	company := newTestCompany(t, "mr"+suffix, 1)
	other := newTestCompany(t, "mo"+suffix, 1)
	machineID := company.machineIDs[0]
	customerID := company.customerID

	// Microseconds on purpose: the stored bounds must be accepted exactly as the API returned them
	start := time.Date(2026, 3, 2, 10, 15, 20, 123456000, time.Local)
	end := start.Add(10 * time.Minute)
	downtimePeriod := insertModePeriod(t, machineID, model.ModeDowntime, start, &end)
	idlePeriod := insertModePeriod(t, machineID, model.ModeIdle, end, new(end.Add(5*time.Minute)))
	productionPeriod := insertModePeriod(t, machineID, model.ModeInProduction, end.Add(5*time.Minute), new(end.Add(6*time.Minute)))
	runningDowntime := insertModePeriod(t, machineID, model.ModeDowntime, end.Add(6*time.Minute), nil)

	downtimeReason := insertReason(t, "downtime_reasons", customerID, "Tool broke "+suffix, dto.DowntimeReasonStatusActive)
	secondReason := insertReason(t, "downtime_reasons", customerID, "No power "+suffix, dto.DowntimeReasonStatusActive)
	inactiveReason := insertReason(t, "downtime_reasons", customerID, "Old reason "+suffix, dto.DowntimeReasonStatusInactive)
	idleReason := insertReason(t, "idle_reasons", customerID, "No operator "+suffix, "active")
	otherCompanyReason := insertReason(t, "downtime_reasons", other.customerID, "Theirs "+suffix, dto.DowntimeReasonStatusActive)
	otherCompanyIdleReason := insertReason(t, "idle_reasons", other.customerID, "Theirs "+suffix, "active")

	router := newMachineModeReasonRouter(asTestUser)
	getPath := func(periodID int64) string { return fmt.Sprintf("%s?periodId=%d", reasonsPath, periodID) }

	t.Run("get returns the whole period and the list to pick from", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, getPath(downtimePeriod), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var loaded dto.MachineModeReasonsResponse
		decodeData(t, result, &loaded)
		if loaded.ReasonType != model.ReasonTypeDowntime || loaded.CustomerID == nil || *loaded.CustomerID != customerID {
			t.Errorf("got type %q customer %v", loaded.ReasonType, loaded.CustomerID)
		}
		if !loaded.Period.StartedAt.Equal(start) || !loaded.Period.EndedAt.Equal(end) || len(loaded.Reasons) != 0 {
			t.Errorf("got period %+v, %d reasons", loaded.Period, len(loaded.Reasons))
		}

		result = sendRequest(t, router, http.MethodGet, getPath(idlePeriod), nil)
		decodeData(t, result, &loaded)
		if loaded.ReasonType != model.ReasonTypeIdle {
			t.Errorf("idle period: got type %q", loaded.ReasonType)
		}
	})

	t.Run("periods that take no reason are refused with the message to show", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodGet, getPath(runningDowntime), nil), http.StatusBadRequest, "You can not add downtime reason in running mode.")
		expectStatus(t, sendRequest(t, router, http.MethodGet, getPath(productionPeriod), nil), http.StatusBadRequest, "You can not add a reason in In Production mode.")
		expectStatus(t, sendRequest(t, router, http.MethodGet, reasonsPath, nil), http.StatusBadRequest, "Please select machine mode period")
		expectStatus(t, sendRequest(t, router, http.MethodGet, getPath(1<<62), nil), http.StatusNotFound, "Machine mode reason not found")

		running := reasonBody(runningDowntime, downtimeReason, end.Add(7*time.Minute), end.Add(8*time.Minute))
		expectStatus(t, sendRequest(t, router, http.MethodPost, reasonsPath, running), http.StatusBadRequest, "You can not add downtime reason in running mode.")
	})

	t.Run("create checks the fields, the time and the reason", func(t *testing.T) {
		badCases := []struct {
			name    string
			body    dto.CreateMachineModeReasonRequest
			status  int
			message string
		}{
			{"no reason", reasonBody(downtimePeriod, 0, start, end), http.StatusBadRequest, "Please select reason"},
			{"no period", reasonBody(0, downtimeReason, start, end), http.StatusBadRequest, "Please select machine mode period"},
			{"no start", reasonBody(downtimePeriod, downtimeReason, time.Time{}, end), http.StatusBadRequest, "Please select start time"},
			{"no end", reasonBody(downtimePeriod, downtimeReason, start, time.Time{}), http.StatusBadRequest, "Please select end time"},
			{"end before start", reasonBody(downtimePeriod, downtimeReason, end, start), http.StatusBadRequest, "End time must be after start time"},
			{"empty span", reasonBody(downtimePeriod, downtimeReason, start, start), http.StatusBadRequest, "End time must be after start time"},
			{"before the period", reasonBody(downtimePeriod, downtimeReason, start.Add(-time.Microsecond), end), http.StatusBadRequest, reasonOutsideMessage},
			{"after the period", reasonBody(downtimePeriod, downtimeReason, start, end.Add(time.Microsecond)), http.StatusBadRequest, reasonOutsideMessage},
			{"inactive reason", reasonBody(downtimePeriod, inactiveReason, start, end), http.StatusBadRequest, invalidDowntimeReason},
			{"another company's reason", reasonBody(downtimePeriod, otherCompanyReason, start, end), http.StatusBadRequest, invalidDowntimeReason},
			{"another company's idle reason", reasonBody(idlePeriod, otherCompanyIdleReason, end, end.Add(time.Minute)), http.StatusBadRequest, "Please select valid idle reason"},
			{"reason id too big for its column", reasonBody(downtimePeriod, 1<<40, start, end), http.StatusBadRequest, "Please select reason"},
		}
		for _, badCase := range badCases {
			t.Run(badCase.name, func(t *testing.T) {
				expectStatus(t, sendRequest(t, router, http.MethodPost, reasonsPath, badCase.body), badCase.status, badCase.message)
			})
		}
	})

	t.Run("reasons split a period, edit and delete", func(t *testing.T) {
		middle := start.Add(4 * time.Minute)
		result := sendRequest(t, router, http.MethodPost, reasonsPath, reasonBody(downtimePeriod, downtimeReason, start, middle))
		expectStatus(t, result, http.StatusCreated, "created")
		var first model.MachineModeReason
		decodeData(t, result, &first)
		if first.ReasonType != model.ReasonTypeDowntime || first.ReasonID != downtimeReason || first.Reason != "Tool broke "+suffix {
			t.Errorf("got %+v", first)
		}

		// Touching the first is allowed, covering any of it is not
		result = sendRequest(t, router, http.MethodPost, reasonsPath, reasonBody(downtimePeriod, secondReason, middle, end))
		expectStatus(t, result, http.StatusCreated, "created")
		var second model.MachineModeReason
		decodeData(t, result, &second)
		expectStatus(t, sendRequest(t, router, http.MethodPost, reasonsPath, reasonBody(downtimePeriod, secondReason, middle.Add(-time.Second), middle.Add(time.Second))), http.StatusBadRequest, reasonOverlapMessage)

		var loaded dto.MachineModeReasonsResponse
		decodeData(t, sendRequest(t, router, http.MethodGet, getPath(downtimePeriod), nil), &loaded)
		if len(loaded.Reasons) != 2 || loaded.Reasons[0].ID != first.ID || loaded.Reasons[1].ID != second.ID {
			t.Fatalf("got reasons %+v, want first then second", loaded.Reasons)
		}

		// An edit overlaps other rows only, never itself
		secondPath := fmt.Sprintf("%s/%d", reasonsPath, second.ID)
		edit := dto.UpdateMachineModeReasonRequest{ReasonID: downtimeReason, StartedAt: middle.Add(time.Minute), EndedAt: end}
		result = sendRequest(t, router, http.MethodPut, secondPath, edit)
		expectStatus(t, result, http.StatusOK, "success")
		decodeData(t, result, &second)
		if second.ReasonID != downtimeReason || !second.StartedAt.Equal(middle.Add(time.Minute)) {
			t.Errorf("edit not saved: %+v", second)
		}
		edit.StartedAt = middle.Add(-time.Minute)
		expectStatus(t, sendRequest(t, router, http.MethodPut, secondPath, edit), http.StatusBadRequest, reasonOverlapMessage)
		edit.StartedAt, edit.EndedAt = middle, end.Add(time.Second)
		expectStatus(t, sendRequest(t, router, http.MethodPut, secondPath, edit), http.StatusBadRequest, reasonOutsideMessage)

		// A reason made inactive after it was picked stays on the row, but cannot be picked anew
		firstPath := fmt.Sprintf("%s/%d", reasonsPath, first.ID)
		if _, err := testDBPool.Exec(context.Background(), `UPDATE downtime_reasons SET status = 'inactive' WHERE id = $1`, downtimeReason); err != nil {
			t.Fatal(err)
		}
		keep := dto.UpdateMachineModeReasonRequest{ReasonID: downtimeReason, StartedAt: start, EndedAt: middle.Add(-time.Minute)}
		expectStatus(t, sendRequest(t, router, http.MethodPut, firstPath, keep), http.StatusOK, "success")
		keep.ReasonID = inactiveReason
		expectStatus(t, sendRequest(t, router, http.MethodPut, firstPath, keep), http.StatusBadRequest, invalidDowntimeReason)
		if _, err := testDBPool.Exec(context.Background(), `UPDATE downtime_reasons SET status = 'active' WHERE id = $1`, downtimeReason); err != nil {
			t.Fatal(err)
		}

		// A deleted row frees its time and is gone
		expectStatus(t, sendRequest(t, router, http.MethodDelete, secondPath, nil), http.StatusOK, reasonDeletedMessage)
		expectStatus(t, sendRequest(t, router, http.MethodDelete, secondPath, nil), http.StatusNotFound, "Machine mode reason not found")
		result = sendRequest(t, router, http.MethodPost, reasonsPath, reasonBody(downtimePeriod, secondReason, middle, end))
		expectStatus(t, result, http.StatusCreated, "created")
	})

	t.Run("an idle period uses the idle list", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, reasonsPath, reasonBody(idlePeriod, idleReason, end, end.Add(5*time.Minute)))
		expectStatus(t, result, http.StatusCreated, "created")
		var saved model.MachineModeReason
		decodeData(t, result, &saved)
		if saved.ReasonType != model.ReasonTypeIdle || saved.Reason != "No operator "+suffix {
			t.Errorf("got %+v", saved)
		}
	})

	t.Run("only one of two saves at once wins the same time", func(t *testing.T) {
		from := end.Add(time.Minute)
		concurrentPeriod := insertModePeriod(t, machineID, model.ModeDowntime, from, new(from.Add(time.Minute)))
		var group sync.WaitGroup
		statuses := make([]int, 8)
		for index := range statuses {
			group.Go(func() {
				statuses[index] = sendRequest(t, router, http.MethodPost, reasonsPath, reasonBody(concurrentPeriod, downtimeReason, from, from.Add(time.Minute))).StatusCode
			})
		}
		group.Wait()
		created := 0
		for _, status := range statuses {
			if status == http.StatusCreated {
				created++
			}
		}
		if created != 1 {
			t.Errorf("%d saves succeeded (%v), want exactly 1", created, statuses)
		}
	})

	t.Run("access follows the machine", func(t *testing.T) {
		accessPeriod := insertModePeriod(t, machineID, model.ModeDowntime, end.Add(time.Hour), new(end.Add(2*time.Hour)))
		from := end.Add(time.Hour)
		body := reasonBody(accessPeriod, downtimeReason, from, from.Add(time.Minute))

		// Another company's customer does not see the machine at all
		outsider := newMachineModeReasonRouter(actAs(other.customerID, model.RoleCustomer))
		expectStatus(t, sendRequest(t, outsider, http.MethodGet, getPath(accessPeriod), nil), http.StatusNotFound, machineNotFoundText)
		expectStatus(t, sendRequest(t, outsider, http.MethodPost, reasonsPath, body), http.StatusNotFound, machineNotFoundText)

		// The machine's own customer records one; the outsider can neither change nor delete it
		owner := newMachineModeReasonRouter(actAs(customerID, model.RoleCustomer))
		result := sendRequest(t, owner, http.MethodPost, reasonsPath, body)
		expectStatus(t, result, http.StatusCreated, "created")
		var saved model.MachineModeReason
		decodeData(t, result, &saved)
		savedPath := fmt.Sprintf("%s/%d", reasonsPath, saved.ID)
		update := dto.UpdateMachineModeReasonRequest{MachineModeReasonFields: body.MachineModeReasonFields}
		expectStatus(t, sendRequest(t, outsider, http.MethodPut, savedPath, update), http.StatusNotFound, machineNotFoundText)
		expectStatus(t, sendRequest(t, outsider, http.MethodDelete, savedPath, nil), http.StatusNotFound, machineNotFoundText)

		// An operator of the company without this machine assigned is refused; assigned they may
		operatorID := createAuthTestUser(t, "mr_operator_"+suffix, model.RoleOperator)
		operator := newMachineModeReasonRouter(actAsMember(operatorID, model.RoleOperator, customerID))
		expectStatus(t, sendRequest(t, operator, http.MethodGet, getPath(accessPeriod), nil), http.StatusNotFound, machineNotFoundText)
		if _, err := testDBPool.Exec(context.Background(), `INSERT INTO user_machine_access (user_id, machine_id) VALUES ($1, $2)`, operatorID, machineID); err != nil {
			t.Fatal(err)
		}
		t.Cleanup(func() {
			_, _ = testDBPool.Exec(context.Background(), `DELETE FROM user_machine_access WHERE user_id = $1`, operatorID)
		})
		expectStatus(t, sendRequest(t, operator, http.MethodGet, getPath(accessPeriod), nil), http.StatusOK, "success")
		expectStatus(t, sendRequest(t, operator, http.MethodDelete, savedPath, nil), http.StatusOK, reasonDeletedMessage)
	})
}
