package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// newProductionScheduleRouter builds the real router with the schedule handler wired to testDBPool.
func newProductionScheduleRouter(authenticate middleware.Middleware) http.Handler {
	scheduleService := service.NewProductionScheduleService(testDBPool, repository.NewUserRepository(testDBPool), newMachineService())
	return routes.New(routes.Handlers{
		ProductionSchedule: handler.NewProductionScheduleHandler(scheduleService),
	}, authenticate)
}

// createCompanyMember adds a user with role to customerID's company.
func createCompanyMember(t *testing.T, customerID int64, username, role string) int64 {
	t.Helper()
	id := createAuthTestUser(t, username, role)
	if _, err := testDBPool.Exec(context.Background(), `UPDATE users SET company_id = $1 WHERE id = $2`, customerID, id); err != nil {
		t.Fatalf("set company: %v", err)
	}
	return id
}

func TestProductionScheduleCRUD(t *testing.T) {
	router := newProductionScheduleRouter(asTestUser) // superadmin
	partRouter := newPartRouter(asTestUser, t.TempDir())
	suffix := uniqueSuffix()
	fixture := newMachineFixture(t, suffix)
	customer := fixture.customerID
	machineA := createTestMachine(t, fixture, "SCH-A-"+suffix)
	machineB := createTestMachine(t, fixture, "SCH-B-"+suffix)
	supervisor := createCompanyMember(t, customer, "sch_sup_"+suffix, model.RoleSupervisor)
	operator := createCompanyMember(t, customer, "sch_op_"+suffix, model.RoleOperator)
	schedulePath := func(id int64) string { return fmt.Sprintf("%s/production-schedules/%d", apiV1, id) }

	// A part with P1 on machine A and P2 on machine B only
	var drawing model.PartDrawing
	decodeData(t, uploadDrawing(t, partRouter, "d.pdf", testPDF), &drawing)
	partResult := sendRequest(t, partRouter, http.MethodPost, apiV1+"/parts", map[string]any{
		"userId": customer, "partNo": "SCH-" + suffix, "partName": "Gear", "partCount": 1,
		"processes": []map[string]any{
			partProcessBody("P1", drawing, []int64{machineA}),
			partProcessBody("P2", drawing, []int64{machineB}),
		},
	})
	expectStatus(t, partResult, http.StatusCreated, "created")
	var part model.Part
	decodeData(t, partResult, &part)

	// Night shift 22:00 – 06:00
	shiftResult := sendRequest(t, newShiftRouter(asTestUser), http.MethodPost, apiV1+"/shifts", shiftBody(customer, "Night "+suffix))
	expectStatus(t, shiftResult, http.StatusCreated, "created")
	var shift model.Shift
	decodeData(t, shiftResult, &shift)

	item := func(jobID, start, end string) map[string]any {
		return map[string]any{
			"jobId": jobID, "supervisorId": supervisor, "operatorId": operator, "shiftId": shift.ID,
			"partId": part.ID, "processId": part.Processes[0].ID, "quantity": 10, "startTime": start, "endTime": end,
		}
	}
	body := map[string]any{
		"userId": customer, "machineId": machineA, "scheduleType": "single", "startDate": "2030-01-10",
		"items": []map[string]any{item("J1", "22:00", "02:00"), item("J2", "02:00", "06:00")},
	}

	var created model.ProductionSchedule
	t.Run("create", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/production-schedules", body)
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		switch {
		case created.ID == 0 || created.EndDate != "2030-01-10" || created.Status != "active" || created.MachineName == "":
			t.Fatalf("created = %+v", created)
		case len(created.Items) != 2 || created.ItemCount != 2 || created.TotalQty != 20:
			t.Fatalf("items = %d (count %d, qty %d); want 2, 20", len(created.Items), created.ItemCount, created.TotalQty)
		case created.Items[0].TotalMinutes != 240 || created.Items[0].ScheduleDate != "2030-01-10" || created.Items[0].OperatorName == "":
			t.Errorf("item = %+v", created.Items[0])
		}
	})
	if created.ID == 0 {
		t.FailNow()
	}
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `UPDATE production_schedules SET is_deleted = true WHERE id = $1`, created.ID)
	})

	rejects := []struct {
		name    string
		body    map[string]any
		message string
	}{
		{"machine already scheduled that day", body, "This machine is already scheduled on the selected date(s)"},
	}
	other := with(body, "startDate", "2030-03-01") // a free day of machine A
	rejects = append(rejects, []struct {
		name    string
		body    map[string]any
		message string
	}{
		{"jobs overlap", with(other, "items", []map[string]any{item("J1", "22:00", "01:00"), item("J2", "00:30", "02:00")}),
			"Job 1 and job 2 overlap on the machine"},
		{"time outside the shift", with(other, "items", []map[string]any{item("J1", "20:00", "23:00")}),
			fmt.Sprintf("Start and end time must be inside the shift Night %s (22:00 – 06:00) (job 1)", suffix)},
		{"process does not run on the machine", with(other, "machineId", machineB),
			"Please select a process of the part that runs on the selected machine"},
		{"operator is not an operator", with(other, "items", []map[string]any{with(item("J1", "22:00", "23:00"), "operatorId", supervisor)}),
			"Please select an operator of the selected customer"},
		{"item date outside the schedule", with(with(with(other, "scheduleType", "multiple"), "endDate", "2030-03-05"), "items", []map[string]any{with(item("J1", "22:00", "23:00"), "scheduleDate", "2030-02-01")}),
			"Please select a date of the schedule (job 1)"},
	}...)
	for _, testCase := range rejects {
		t.Run(testCase.name, func(t *testing.T) {
			expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/production-schedules", testCase.body), http.StatusBadRequest, testCase.message)
		})
	}

	t.Run("update to multiple days keeps, drops and adds jobs", func(t *testing.T) {
		kept := item("J1-edited", "22:00", "23:00")
		kept["id"] = created.Items[0].ID
		kept["scheduleDate"] = "2030-01-10"
		added := item("J3", "22:00", "06:00")
		added["scheduleDate"] = "2030-01-12"
		update := with(with(with(body, "scheduleType", "multiple"), "endDate", "2030-01-12"), "status", "inactive")
		update["items"] = []map[string]any{kept, added}

		result := sendRequest(t, router, http.MethodPut, schedulePath(created.ID), update)
		expectStatus(t, result, http.StatusOK, "success")
		var updated model.ProductionSchedule
		decodeData(t, result, &updated)
		switch {
		case updated.ScheduleType != "multiple" || updated.EndDate != "2030-01-12" || updated.Status != "inactive":
			t.Errorf("updated = %+v", updated)
		case len(updated.Items) != 2 || updated.Items[0].ID != created.Items[0].ID || updated.Items[0].JobID != "J1-edited":
			t.Errorf("items = %+v", updated.Items)
		case updated.Items[1].ScheduleDate != "2030-01-12" || updated.Items[1].TotalMinutes != 480:
			t.Errorf("added = %+v", updated.Items[1])
		}
	})

	t.Run("another company's admin cannot see it", func(t *testing.T) {
		other := newMachineFixture(t, suffix+"7")
		otherAdmin := createCompanyMember(t, other.customerID, "sch_adm_"+suffix, model.RoleAdmin)
		otherRouter := newProductionScheduleRouter(actAs(otherAdmin, model.RoleAdmin))
		result := sendRequest(t, otherRouter, http.MethodGet, schedulePath(created.ID), nil)
		if result.StatusCode == http.StatusOK {
			t.Fatalf("status = %d; want the schedule hidden", result.StatusCode)
		}
	})

	t.Run("list and delete", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/production-schedules?userId=%d", apiV1, customer), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var schedules []model.ProductionSchedule
		decodeData(t, result, &schedules)
		if len(schedules) != 1 || schedules[0].ID != created.ID {
			t.Fatalf("schedules = %+v", schedules)
		}
		expectStatus(t, sendRequest(t, router, http.MethodDelete, schedulePath(created.ID), nil), http.StatusOK, "Production schedule deleted")
		expectStatus(t, sendRequest(t, router, http.MethodGet, schedulePath(created.ID), nil), http.StatusNotFound, "Production schedule not found")
	})
}
