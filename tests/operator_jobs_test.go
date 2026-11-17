package tests

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"7th-sense-server/internal/model"
)

// TestOperatorJobs covers GET /production-schedules/operator-jobs (Operator Wise Production
// Schedule): one row per job with the names the screen shows, every filter, the order, and the
// company scoping of the schedule list.
func TestOperatorJobs(t *testing.T) {
	router := newProductionScheduleRouter(asTestUser) // superadmin
	partRouter := newPartRouter(asTestUser, t.TempDir())
	suffix := uniqueSuffix()
	fixture := newMachineFixture(t, suffix)
	customer := fixture.customerID
	machine := createTestMachine(t, fixture, "OPJ-"+suffix)
	supervisor := createCompanyMember(t, customer, "opj_sup_"+suffix, model.RoleSupervisor)
	firstOperator := createCompanyMember(t, customer, "opj_op1_"+suffix, model.RoleOperator)
	secondOperator := createCompanyMember(t, customer, "opj_op2_"+suffix, model.RoleOperator)

	var drawing model.PartDrawing
	decodeData(t, uploadDrawing(t, partRouter, "d.pdf", testPDF), &drawing)
	partResult := sendRequest(t, partRouter, http.MethodPost, apiV1+"/parts", map[string]any{
		"userId": customer, "partNo": "OPJ-" + suffix, "partName": "Flange", "partCount": 1,
		"processes": []map[string]any{partProcessBody("P1", drawing, []int64{machine})},
	})
	expectStatus(t, partResult, http.StatusCreated, "created")
	var part model.Part
	decodeData(t, partResult, &part)

	shiftResult := sendRequest(t, newShiftRouter(asTestUser), http.MethodPost, apiV1+"/shifts", shiftBody(customer, "Night "+suffix))
	expectStatus(t, shiftResult, http.StatusCreated, "created")
	var shift model.Shift
	decodeData(t, shiftResult, &shift)

	item := func(jobID string, operatorID int64, start, end string) map[string]any {
		return map[string]any{
			"jobId": jobID, "supervisorId": supervisor, "operatorId": operatorID, "shiftId": shift.ID,
			"partId": part.ID, "processId": part.Processes[0].ID, "quantity": 15, "startTime": start, "endTime": end,
		}
	}
	createSchedule := func(day string, items ...map[string]any) {
		t.Helper()
		result := sendRequest(t, router, http.MethodPost, apiV1+"/production-schedules", map[string]any{
			"userId": customer, "machineId": machine, "scheduleType": "single", "startDate": day, "items": items,
		})
		expectStatus(t, result, http.StatusCreated, "created")
		var created model.ProductionSchedule
		decodeData(t, result, &created)
		t.Cleanup(func() {
			_, _ = testDBPool.Exec(context.Background(), `UPDATE production_schedules SET is_deleted = true WHERE id = $1`, created.ID)
		})
	}
	createSchedule("2030-05-10", item("J1", firstOperator, "22:00", "02:00"), item("J2", secondOperator, "02:00", "06:00"))
	createSchedule("2030-05-11", item("J3", firstOperator, "22:00", "23:30"))

	jobsPath := apiV1 + "/production-schedules/operator-jobs"
	getJobs := func(t *testing.T, query string) []model.OperatorJob {
		t.Helper()
		result := sendRequest(t, router, http.MethodGet, jobsPath+"?"+query, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var jobs []model.OperatorJob
		decodeData(t, result, &jobs)
		return jobs
	}
	customerQuery := fmt.Sprintf("userId=%d", customer)

	t.Run("every job of the customer, latest day first", func(t *testing.T) {
		jobs := getJobs(t, customerQuery)
		if len(jobs) != 3 || jobs[0].JobID != "J3" || jobs[1].JobID != "J1" || jobs[2].JobID != "J2" {
			t.Fatalf("jobs = %+v; want J3, J1, J2", jobs)
		}
		first := jobs[1]
		switch {
		case first.ScheduleDate != "2030-05-10" || first.StartTime != "22:00" || first.EndTime != "02:00" || first.Quantity != 15:
			t.Errorf("job times = %+v", first)
		case first.OperatorID != firstOperator || first.OperatorName == "" || first.ShiftName != "Night "+suffix:
			t.Errorf("operator / shift = %+v", first)
		case first.MachineID != machine || first.MachineName != "OPJ-"+suffix || first.PartNo != "OPJ-"+suffix || first.PartName != "Flange":
			t.Errorf("machine / part = %+v", first)
		case first.ProcessName != "P1" || first.ProgramName != "PRG_P1" || first.UserID != customer || first.CustomerName == "":
			t.Errorf("process / customer = %+v", first)
		}
	})

	t.Run("filters", func(t *testing.T) {
		cases := []struct {
			query string
			want  int
		}{
			{fmt.Sprintf("%s&operatorId=%d", customerQuery, firstOperator), 2},
			{fmt.Sprintf("%s&operatorId=%d", customerQuery, secondOperator), 1},
			{customerQuery + "&date=2030-05-10", 2},
			{fmt.Sprintf("%s&date=2030-05-10&operatorId=%d", customerQuery, secondOperator), 1},
			{fmt.Sprintf("%s&machineId=%d&partId=%d&shiftId=%d", customerQuery, machine, part.ID, shift.ID), 3},
			{customerQuery + "&date=2030-05-12", 0},
		}
		for _, testCase := range cases {
			if jobs := getJobs(t, testCase.query); len(jobs) != testCase.want {
				t.Errorf("%s → %d jobs; want %d", testCase.query, len(jobs), testCase.want)
			}
		}
	})

	t.Run("validation", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodGet, jobsPath+"?date=10-05-2030", nil), http.StatusBadRequest, "Please enter valid date")
		expectStatus(t, sendRequest(t, router, http.MethodGet, jobsPath+"?machineId=abc", nil), http.StatusBadRequest, "Please enter valid machineId")
	})

	t.Run("company scoping", func(t *testing.T) {
		member := newProductionScheduleRouter(actAsMember(firstOperator, model.RoleOperator, customer))
		result := sendRequest(t, member, http.MethodGet, jobsPath, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var jobs []model.OperatorJob
		decodeData(t, result, &jobs)
		if len(jobs) != 3 {
			t.Fatalf("company member sees %d jobs; want the company's 3", len(jobs))
		}

		other := newMachineFixture(t, suffix+"9")
		otherAdmin := createCompanyMember(t, other.customerID, "opj_adm_"+suffix, model.RoleAdmin)
		outsider := newProductionScheduleRouter(actAsMember(otherAdmin, model.RoleAdmin, other.customerID))
		result = sendRequest(t, outsider, http.MethodGet, jobsPath, nil)
		expectStatus(t, result, http.StatusOK, "success")
		decodeData(t, result, &jobs)
		if len(jobs) != 0 {
			t.Fatalf("another company sees %d jobs; want none", len(jobs))
		}
		expectStatus(t, sendRequest(t, outsider, http.MethodGet, jobsPath+"?"+customerQuery, nil), http.StatusForbidden, "You do not have permission to do this")
	})
}
