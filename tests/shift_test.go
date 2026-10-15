package tests

import (
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

// newShiftRouter builds the real router with the shift, break and holiday handlers wired
// to testDBPool; authenticate decides who the requests act as.
func newShiftRouter(authenticate middleware.Middleware) http.Handler {
	userRepository := repository.NewUserRepository(testDBPool)
	return routes.New(routes.Handlers{
		Shift:   handler.NewShiftHandler(service.NewShiftService(testDBPool, userRepository)),
		Break:   handler.NewBreakHandler(service.NewBreakService(testDBPool, userRepository)),
		Holiday: handler.NewHolidayHandler(service.NewHolidayService(repository.NewHolidayRepository(testDBPool), userRepository)),
	}, authenticate)
}

// createBreak adds a break of customerID through the API and returns it.
func createBreak(t *testing.T, router http.Handler, customerID int64, name, start, end string) model.Break {
	t.Helper()
	result := sendRequest(t, router, http.MethodPost, apiV1+"/breaks",
		map[string]any{"userId": customerID, "name": name, "startTime": start, "endTime": end})
	expectStatus(t, result, http.StatusCreated, "created")
	var created model.Break
	decodeData(t, result, &created)
	return created
}

// shiftBody is a valid night shift (past midnight) taking breakIDs; change fields per case.
func shiftBody(customerID int64, name string, breakIDs ...int64) map[string]any {
	if breakIDs == nil {
		breakIDs = []int64{}
	}
	return map[string]any{
		"userId":      customerID,
		"name":        name,
		"shiftType":   "Night",
		"startTime":   "22:00",
		"endTime":     "06:00",
		"workingDays": []int{7, 1, 2, 3, 1},
		"breakIds":    breakIDs,
	}
}

func TestBreakAndShiftCRUD(t *testing.T) {
	router := newShiftRouter(asTestUser) // superadmin
	suffix := uniqueSuffix()
	customerID := createAuthTestUser(t, "shift_customer_"+suffix, model.RoleCustomer)
	otherCustomerID := createAuthTestUser(t, "shift_other_"+suffix, model.RoleCustomer)
	shiftPath := func(id int64) string { return fmt.Sprintf("%s/shifts/%d", apiV1, id) }
	breakPath := func(id int64) string { return fmt.Sprintf("%s/breaks/%d", apiV1, id) }
	getBreak := func(t *testing.T, id int64) model.Break {
		t.Helper()
		result := sendRequest(t, router, http.MethodGet, breakPath(id), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var found model.Break
		decodeData(t, result, &found)
		return found
	}

	// The customer's breaks, not on a shift yet
	dinner := createBreak(t, router, customerID, "Dinner", "23:30", "00:15") // past midnight
	tea := createBreak(t, router, customerID, "N-Tea", "02:00", "02:10")
	lunch := createBreak(t, router, customerID, "Lunch", "13:00", "13:30") // a day break
	overlapping := createBreak(t, router, customerID, "Snack", "02:05", "02:20")
	otherBreak := createBreak(t, router, otherCustomerID, "Other", "01:00", "01:10")
	if dinner.ShiftID != nil {
		t.Fatalf("new break shift = %v; want none", *dinner.ShiftID)
	}

	t.Run("break rules", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/breaks",
			map[string]any{"userId": customerID, "name": "dinner", "startTime": "20:00", "endTime": "20:30"})
		expectStatus(t, result, http.StatusBadRequest, "This customer already has a break with this name")
		result = sendRequest(t, router, http.MethodPost, apiV1+"/breaks",
			map[string]any{"userId": customerID, "name": "Bad", "startTime": "10:00", "endTime": "10:00"})
		expectStatus(t, result, http.StatusBadRequest, "Break end time must differ from the start time")
		result = sendRequest(t, router, http.MethodPost, apiV1+"/breaks",
			map[string]any{"userId": customerID, "name": "Bad", "startTime": "9:00", "endTime": "10:00"})
		expectStatus(t, result, http.StatusBadRequest, "Please enter valid break start time")
	})

	var created model.Shift
	t.Run("create night shift taking breaks", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/shifts", shiftBody(customerID, "A-"+suffix, tea.ID, dinner.ID, tea.ID))
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		switch {
		case created.ID == 0 || created.UserID != customerID || created.ShiftType != "night" || created.Status != "active":
			t.Fatalf("created = %+v", created)
		case fmt.Sprint(created.WorkingDays) != "[1 2 3 7]":
			t.Errorf("working days = %v; want sorted, without duplicates", created.WorkingDays)
		case created.DurationMinutes != 480 || created.BreakMinutes != 55:
			t.Errorf("minutes = %d / break %d; want 480 / 55", created.DurationMinutes, created.BreakMinutes)
		// From the shift start: Dinner (23:30) before N-Tea (02:00)
		case len(created.Breaks) != 2 || created.Breaks[0].ID != dinner.ID || created.Breaks[1].ID != tea.ID:
			t.Errorf("breaks = %+v; want Dinner then N-Tea", created.Breaks)
		}
		if found := getBreak(t, tea.ID); found.ShiftID == nil || *found.ShiftID != created.ID || found.ShiftName != "A-"+suffix {
			t.Errorf("N-Tea shift = %v %q; want shift %d", found.ShiftID, found.ShiftName, created.ID)
		}
	})
	if created.ID == 0 {
		t.FailNow()
	}

	t.Run("shift rules", func(t *testing.T) {
		cases := map[string]struct {
			body    map[string]any
			message string
		}{
			"duplicate name":         {shiftBody(customerID, "A-"+suffix), "This customer already has a shift with this name"},
			"break outside shift":    {shiftBody(customerID, "V1-"+suffix, lunch.ID), "Break Lunch (13:00 to 13:30) is not within the shift time"},
			"break of another shift": {shiftBody(customerID, "V3-"+suffix, tea.ID), "Break N-Tea is already taken by shift A-" + suffix},
			"other customer break":   {shiftBody(customerID, "V4-"+suffix, otherBreak.ID), "Please select breaks of the selected customer"},
			"bad type":               {with(shiftBody(customerID, "V5-"+suffix), "shiftType", "evening"), "Please select shift type"},
			"no working days":        {with(shiftBody(customerID, "V6-"+suffix), "workingDays", []int{}), "Please select valid working days"},
			"same start and end":     {with(shiftBody(customerID, "V7-"+suffix), "endTime", "22:00"), "Shift end time must differ from the start time"},
		}
		for name, testCase := range cases {
			t.Run(name, func(t *testing.T) {
				expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/shifts", testCase.body), http.StatusBadRequest, testCase.message)
			})
		}
		// Breaks of one shift must not overlap
		body := shiftBody(customerID, "A-"+suffix, tea.ID, dinner.ID, overlapping.ID)
		expectStatus(t, sendRequest(t, router, http.MethodPut, shiftPath(created.ID), body), http.StatusBadRequest, "Breaks N-Tea and Snack overlap")
	})

	t.Run("break page keeps its shift's rules", func(t *testing.T) {
		body := map[string]any{"userId": customerID, "shiftId": created.ID, "name": "N-Tea", "startTime": "05:55", "endTime": "06:20"}
		result := sendRequest(t, router, http.MethodPut, breakPath(tea.ID), body)
		expectStatus(t, result, http.StatusBadRequest, "Shift A-"+suffix+": Break N-Tea (05:55 to 06:20) is not within the shift time")

		body["startTime"], body["endTime"] = "03:00", "03:15" // still inside: allowed
		expectStatus(t, sendRequest(t, router, http.MethodPut, breakPath(tea.ID), body), http.StatusOK, "success")

		// Adding a break straight onto the shift from the Breaks page
		result = sendRequest(t, router, http.MethodPost, apiV1+"/breaks",
			map[string]any{"userId": customerID, "shiftId": created.ID, "name": "Late tea", "startTime": "04:00", "endTime": "04:10"})
		expectStatus(t, result, http.StatusCreated, "created")
		var late model.Break
		decodeData(t, result, &late)
		if late.ShiftID == nil || *late.ShiftID != created.ID {
			t.Errorf("Late tea shift = %v; want %d", late.ShiftID, created.ID)
		}

		result = sendRequest(t, router, http.MethodPost, apiV1+"/breaks",
			map[string]any{"userId": otherCustomerID, "shiftId": created.ID, "name": "X", "startTime": "04:00", "endTime": "04:10"})
		expectStatus(t, result, http.StatusBadRequest, "Please select a shift of the selected customer")

		// Deleting a break takes it off its shift
		expectStatus(t, sendRequest(t, router, http.MethodDelete, breakPath(late.ID), nil), http.StatusOK, "Break deleted")
	})

	t.Run("update releases the breaks left out", func(t *testing.T) {
		body := shiftBody(customerID, "A-"+suffix, lunch.ID)
		body["shiftType"], body["startTime"], body["endTime"], body["status"] = "day", "09:00", "17:30", "inactive"
		result := sendRequest(t, router, http.MethodPut, shiftPath(created.ID), body)
		expectStatus(t, result, http.StatusOK, "success")
		var updated model.Shift
		decodeData(t, result, &updated)
		if updated.ShiftType != "day" || updated.Status != "inactive" || len(updated.Breaks) != 1 || updated.Breaks[0].ID != lunch.ID {
			t.Errorf("updated = %+v", updated)
		}
		if found := getBreak(t, tea.ID); found.ShiftID != nil || found.ShiftName != "" {
			t.Errorf("N-Tea shift = %v %q; want none", found.ShiftID, found.ShiftName)
		}
	})

	t.Run("company scope", func(t *testing.T) {
		own := newShiftRouter(actAs(customerID, model.RoleCustomer))
		other := newShiftRouter(actAs(otherCustomerID, model.RoleCustomer))

		expectStatus(t, sendRequest(t, own, http.MethodGet, shiftPath(created.ID), nil), http.StatusOK, "success")
		expectStatus(t, sendRequest(t, other, http.MethodGet, shiftPath(created.ID), nil), http.StatusNotFound, "Shift not found")
		expectStatus(t, sendRequest(t, other, http.MethodGet, breakPath(lunch.ID), nil), http.StatusNotFound, "Break not found")
		expectStatus(t, sendRequest(t, other, http.MethodDelete, breakPath(lunch.ID), nil), http.StatusNotFound, "Break not found")

		operator := newShiftRouter(actAs(customerID, model.RoleOperator))
		expectStatus(t, sendRequest(t, operator, http.MethodDelete, shiftPath(created.ID), nil), http.StatusForbidden, "You do not have permission to do this")
	})

	t.Run("delete releases the shift's breaks", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, shiftPath(created.ID), nil), http.StatusOK, "Shift deleted")
		expectStatus(t, sendRequest(t, router, http.MethodGet, shiftPath(created.ID), nil), http.StatusNotFound, "Shift not found")
		if found := getBreak(t, lunch.ID); found.ShiftID != nil {
			t.Errorf("Lunch shift = %v; want none after the shift was deleted", *found.ShiftID)
		}
	})
}

func TestHolidayCRUD(t *testing.T) {
	router := newShiftRouter(asTestUser) // superadmin
	suffix := uniqueSuffix()
	customerID := createAuthTestUser(t, "holiday_customer_"+suffix, model.RoleCustomer)
	otherCustomerID := createAuthTestUser(t, "holiday_other_"+suffix, model.RoleCustomer)
	holidayPath := func(id int64) string { return fmt.Sprintf("%s/holidays/%d", apiV1, id) }

	var diwali model.Holiday
	t.Run("create holiday", func(t *testing.T) {
		body := map[string]any{"userId": customerID, "holidayType": "Holiday", "name": "Diwali",
			"startDate": "2026-11-08", "endDate": "2026-11-10", "weekDay": 3}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/holidays", body)
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &diwali)
		if diwali.HolidayType != "holiday" || diwali.StartDate == nil || *diwali.StartDate != "2026-11-08" || diwali.WeekDay != nil {
			t.Fatalf("created = %+v; want a holiday with dates and no weekday", diwali)
		}
	})

	var weekOff model.Holiday
	t.Run("create week off", func(t *testing.T) {
		body := map[string]any{"userId": customerID, "holidayType": "week_off", "weekDay": 3, "startDate": "2026-01-01"}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/holidays", body)
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &weekOff)
		if weekOff.WeekDay == nil || *weekOff.WeekDay != 3 || weekOff.StartDate != nil || weekOff.Name != "" {
			t.Fatalf("created = %+v; want Wednesday, no dates, no name", weekOff)
		}
	})

	t.Run("rules", func(t *testing.T) {
		cases := map[string]struct {
			body    map[string]any
			message string
		}{
			"week off day twice": {map[string]any{"userId": customerID, "holidayType": "week_off", "weekDay": 3},
				"This customer already has a week off on this day"},
			"same holiday twice": {map[string]any{"userId": customerID, "holidayType": "holiday", "name": "diwali",
				"startDate": "2026-11-08", "endDate": "2026-11-08"}, "This customer already has this holiday on this date"},
			"end before start": {map[string]any{"userId": customerID, "holidayType": "holiday", "name": "X",
				"startDate": "2026-11-08", "endDate": "2026-11-07"}, "End date cannot be before the start date"},
			"holiday without name": {map[string]any{"userId": customerID, "holidayType": "holiday",
				"startDate": "2026-11-08", "endDate": "2026-11-08"}, "Please enter valid holiday name"},
			"bad date": {map[string]any{"userId": customerID, "holidayType": "holiday", "name": "X",
				"startDate": "2026-02-30", "endDate": "2026-03-01"}, "Please enter valid start date"},
			"week off without day": {map[string]any{"userId": customerID, "holidayType": "week_off"}, "Please select week off day"},
			"unknown kind":         {map[string]any{"userId": customerID, "holidayType": "leave"}, "Please select holiday or week off"},
		}
		for name, testCase := range cases {
			t.Run(name, func(t *testing.T) {
				expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/holidays", testCase.body), http.StatusBadRequest, testCase.message)
			})
		}
		// Another customer may use the same weekday
		body := map[string]any{"userId": otherCustomerID, "holidayType": "week_off", "weekDay": 3}
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/holidays", body), http.StatusCreated, "created")
	})

	t.Run("update week off to a holiday", func(t *testing.T) {
		body := map[string]any{"userId": customerID, "holidayType": "holiday", "name": "Holi",
			"startDate": "2026-03-04", "endDate": "2026-03-04", "status": "inactive"}
		result := sendRequest(t, router, http.MethodPut, holidayPath(weekOff.ID), body)
		expectStatus(t, result, http.StatusOK, "success")
		var updated model.Holiday
		decodeData(t, result, &updated)
		if updated.HolidayType != "holiday" || updated.WeekDay != nil || updated.Name != "Holi" || updated.Status != "inactive" {
			t.Errorf("updated = %+v", updated)
		}
	})

	t.Run("company scope and delete", func(t *testing.T) {
		other := newShiftRouter(actAs(otherCustomerID, model.RoleCustomer))
		expectStatus(t, sendRequest(t, other, http.MethodGet, holidayPath(diwali.ID), nil), http.StatusNotFound, "Holiday not found")
		expectStatus(t, sendRequest(t, other, http.MethodDelete, holidayPath(diwali.ID), nil), http.StatusNotFound, "Holiday not found")

		expectStatus(t, sendRequest(t, router, http.MethodDelete, holidayPath(diwali.ID), nil), http.StatusOK, "Holiday deleted")
		expectStatus(t, sendRequest(t, router, http.MethodGet, holidayPath(diwali.ID), nil), http.StatusNotFound, "Holiday not found")
	})
}
