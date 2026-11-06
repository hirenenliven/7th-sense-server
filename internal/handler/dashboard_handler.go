package handler

import (
	"net/http"
	"strconv"
	"strings"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// DashboardHandler serves the charts of the Home (/dashboard/home), Productivity (/productivity/...)
// and Machine Health (/machine-health/...) pages. Every logged-in user reads the machines they may access.
type DashboardHandler struct {
	service *service.DashboardService
}

// NewDashboardHandler creates a DashboardHandler that uses dashboardService.
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: dashboardService}
}

// GetHomeDashboard handles GET /api/v1/dashboard/home?date=YYYY-MM-DD[&userId=<customer id>]
// (date defaults to today; userId narrows the superadmin's machines to one customer): the active
// machines together, the fleet now and the day's OEE, cycles, parts vs plan, time per mode group,
// parts per hour, per-machine rows and the most raised alarms.
//
// @Summary      Get home dashboard
// @Description  GetHomeDashboard handles GET /api/v1/dashboard/home?date=YYYY-MM-DD[&userId=<customer id>] (date defaults to today; userId narrows the superadmin's machines to one customer): the active machines together, the fleet now and the day's OEE, cycles, parts vs plan, time per mode group, parts per hour, per-machine rows and the most raised alarms.
// @Description
// @Description  Used by pages: Home, Operations Center, Plant Intelligence
// @Tags         Dashboard
// @ID           getHomeDashboard
// @Produce      json
// @Security     cookieAuth
// @Param        date query string false "Day YYYY-MM-DD (default today, not after today)"
// @Param        userId query int false "Customer id (superadmin; others only their own company)"
// @Success      200 {object} response.Envelope{data=dto.HomeDashboardResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Another company"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/dashboard/home [get]
func (dashboardHandler *DashboardHandler) GetHomeDashboard(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	customerID, err := queryID(request, "userId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	homeRequest := dto.HomeDashboardRequest{Date: request.URL.Query().Get("date"), CustomerID: customerID}
	if err := homeRequest.Validate(); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	home, err := dashboardHandler.service.GetHomeDashboard(request.Context(), actor, homeRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, home)
}

// GetModeTimeline handles GET /api/v1/productivity/mode-timeline?machineId=&date=YYYY-MM-DD
// (date defaults to today) or ?machineId=&minutes=30 (the last minutes up to now): the machine's
// mode periods in that window with its current program and operation mode.
//
// @Summary      Get mode timeline
// @Description  GetModeTimeline handles GET /api/v1/productivity/mode-timeline?machineId=&date=YYYY-MM-DD (date defaults to today) or ?machineId=&minutes=30 (the last minutes up to now): the machine's mode periods in that window with its current program and operation mode.
// @Description
// @Description  Used by pages: Productivity
// @Tags         Productivity
// @ID           getModeTimeline
// @Produce      json
// @Security     cookieAuth
// @Param        machineId query int true "Machine id"
// @Param        date query string false "Day YYYY-MM-DD (default today)"
// @Param        minutes query int false "Last N minutes up to now (1-1440); overrides date"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/productivity/mode-timeline [get]
func (dashboardHandler *DashboardHandler) GetModeTimeline(writer http.ResponseWriter, request *http.Request) {
	actor, dashboardRequest, err := readDashboardRequest(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	timeline, err := dashboardHandler.service.GetModeTimeline(request.Context(), actor, dashboardRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, timeline)
}

// GetOEE handles GET /api/v1/productivity/oee?machineIds=1,2&from=&to=&shiftIds=&program=&schedule=with:
// one OEE row per day over the machines together, weighted by their planned / run time and closed cycles.
//
// @Summary      Get oee
// @Description  GetOEE handles GET /api/v1/productivity/oee?machineIds=1,2&from=&to=&shiftIds=&program=&schedule=with: one OEE row per day over the machines together, weighted by their planned / run time and closed cycles.
// @Description
// @Description  Used by pages: Productivity, Operations Center
// @Tags         Productivity
// @ID           getOee
// @Produce      json
// @Security     cookieAuth
// @Param        machineIds query string true "Machine ids, comma separated"
// @Param        from query string false "First day YYYY-MM-DD"
// @Param        to query string false "Last day YYYY-MM-DD (default today)"
// @Param        shiftIds query string false "Shift ids, comma separated (none = whole day)"
// @Param        program query string false "Program name (empty = all)"
// @Param        schedule query string false "with = only the scheduled shift time"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/productivity/oee [get]
func (dashboardHandler *DashboardHandler) GetOEE(writer http.ResponseWriter, request *http.Request) {
	actor, dashboardRequest, err := readDashboardRequest(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	oee, err := dashboardHandler.service.GetOEE(request.Context(), actor, dashboardRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, oee)
}

// GetCycleCounts handles GET /api/v1/productivity/cycle-counts?machineIds=1,2&from=&to=&shiftIds=:
// cycles started, finished and unfinished per day.
//
// @Summary      Get cycle counts
// @Description  GetCycleCounts handles GET /api/v1/productivity/cycle-counts?machineIds=1,2&from=&to=&shiftIds=: cycles started, finished and unfinished per day.
// @Description
// @Description  Used by pages: Productivity, Operations Center
// @Tags         Productivity
// @ID           getCycleCounts
// @Produce      json
// @Security     cookieAuth
// @Param        machineIds query string true "Machine ids, comma separated"
// @Param        from query string false "First day YYYY-MM-DD"
// @Param        to query string false "Last day YYYY-MM-DD (default today)"
// @Param        shiftIds query string false "Shift ids, comma separated"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/productivity/cycle-counts [get]
func (dashboardHandler *DashboardHandler) GetCycleCounts(writer http.ResponseWriter, request *http.Request) {
	actor, dashboardRequest, err := readDashboardRequest(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	rows, err := dashboardHandler.service.GetCycleCounts(request.Context(), actor, dashboardRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, rows)
}

// GetMachineSummaries handles GET /api/v1/productivity/machine-summaries?date=YYYY-MM-DD (default
// today): OEE and cycle counts of the day for every machine the user may read.
//
// @Summary      Get machine summaries
// @Description  GetMachineSummaries handles GET /api/v1/productivity/machine-summaries?date=YYYY-MM-DD (default today): OEE (availability, performance, quality) and cycle counts (total, completed, unfinished) of the whole day for every machine the user may read. Availability is null when the day has no planned time yet; performance, quality and oee are null when no cycle closed.
// @Description
// @Description  Used by pages: Monitoring (machine list)
// @Tags         Productivity
// @ID           getMachineSummaries
// @Produce      json
// @Security     cookieAuth
// @Param        date query string false "Day YYYY-MM-DD (default today)"
// @Success      200 {object} response.Envelope{data=[]dto.MachineSummary}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/productivity/machine-summaries [get]
func (dashboardHandler *DashboardHandler) GetMachineSummaries(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	date := request.URL.Query().Get("date")
	// Only the day is checked: the machines are every one the user may read
	dateRequest := dto.MachineCycleListRequest{MachineID: 1, From: date, To: date}
	if err := dateRequest.Validate(); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	summaries, err := dashboardHandler.service.GetMachineSummaries(request.Context(), actor, dto.DashboardRequest{From: date, To: date})
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, summaries)
}

// GetFinishCycles handles GET /api/v1/productivity/finish-cycles?machineId=&date=&shiftIds=&program=:
// the machine's completed cycles of the day with the day's totals.
//
// @Summary      Get finish cycles
// @Description  GetFinishCycles handles GET /api/v1/productivity/finish-cycles?machineId=&date=&shiftIds=&program=: the machine's completed cycles of the day with the day's totals.
// @Description
// @Description  Used by pages: Productivity
// @Tags         Productivity
// @ID           getFinishCycles
// @Produce      json
// @Security     cookieAuth
// @Param        machineId query int true "Machine id"
// @Param        date query string false "Day YYYY-MM-DD (default today)"
// @Param        shiftIds query string false "Shift ids, comma separated"
// @Param        program query string false "Program name (empty = all)"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/productivity/finish-cycles [get]
func (dashboardHandler *DashboardHandler) GetFinishCycles(writer http.ResponseWriter, request *http.Request) {
	actor, dashboardRequest, err := readDashboardRequest(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	finishCycles, err := dashboardHandler.service.GetFinishCycles(request.Context(), actor, dashboardRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, finishCycles)
}

// GetModeTotals handles GET /api/v1/productivity/mode-totals?machineId=&from=&to= (or &date=
// for one day, or &all=true for every recorded period): the time the machine spent in each mode.
//
// @Summary      Get mode totals
// @Description  GetModeTotals handles GET /api/v1/productivity/mode-totals?machineId=&from=&to= (or &date= for one day, or &all=true for every recorded period): the time the machine spent in each mode.
// @Description
// @Description  Used by pages: Productivity
// @Tags         Productivity
// @ID           getModeTotals
// @Produce      json
// @Security     cookieAuth
// @Param        machineId query int true "Machine id"
// @Param        date query string false "One day YYYY-MM-DD (default today)"
// @Param        from query string false "First day YYYY-MM-DD"
// @Param        to query string false "Last day YYYY-MM-DD (default today); at most 31 days"
// @Param        all query bool false "Every recorded period (dates ignored)"
// @Success      200 {object} response.Envelope{data=dto.ModeTotalsResponse}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/productivity/mode-totals [get]
func (dashboardHandler *DashboardHandler) GetModeTotals(writer http.ResponseWriter, request *http.Request) {
	actor, dashboardRequest, err := readDashboardRequest(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	totals, err := dashboardHandler.service.GetModeTotals(request.Context(), actor, dashboardRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, totals)
}

// GetMachineTimeByDay handles GET /api/v1/machine-health/machine-time-by-day?machineId=&from=&to=:
// hours per machine mode for each day.
//
// @Summary      Get machine time by day
// @Description  GetMachineTimeByDay handles GET /api/v1/machine-health/machine-time-by-day?machineId=&from=&to=: hours per machine mode for each day.
// @Description
// @Description  Used by pages: Machine health
// @Tags         Machine health
// @ID           getMachineTimeByDay
// @Produce      json
// @Security     cookieAuth
// @Param        machineId query int true "Machine id"
// @Param        from query string false "First day YYYY-MM-DD"
// @Param        to query string false "Last day YYYY-MM-DD (default today)"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/machine-health/machine-time-by-day [get]
func (dashboardHandler *DashboardHandler) GetMachineTimeByDay(writer http.ResponseWriter, request *http.Request) {
	actor, dashboardRequest, err := readDashboardRequest(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	rows, err := dashboardHandler.service.GetMachineTimeByDay(request.Context(), actor, dashboardRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, rows)
}

// readDashboardRequest reads the logged-in user and the chart query, and validates the query.
func readDashboardRequest(request *http.Request) (model.AuthUser, dto.DashboardRequest, error) {
	actor, err := currentUser(request)
	if err != nil {
		return model.AuthUser{}, dto.DashboardRequest{}, err
	}
	dashboardRequest, err := dashboardQuery(request)
	if err != nil {
		return model.AuthUser{}, dto.DashboardRequest{}, err
	}
	if err := dashboardRequest.Validate(); err != nil {
		return model.AuthUser{}, dto.DashboardRequest{}, err
	}
	return actor, dashboardRequest, nil
}

// dashboardQuery reads the chart query: machineIds (or machineId), from / to (or date for one
// day), shiftIds, program, schedule, minutes and all (every recorded period).
func dashboardQuery(request *http.Request) (dto.DashboardRequest, error) {
	query := request.URL.Query()
	machineIDs, err := queryIDList(query.Get("machineIds"), "machineIds")
	if err != nil {
		return dto.DashboardRequest{}, err
	}
	if machineID, err := queryID(request, "machineId"); err != nil {
		return dto.DashboardRequest{}, err
	} else if machineID != nil {
		machineIDs = append([]int64{*machineID}, machineIDs...)
	}
	shiftIDs, err := queryIDList(query.Get("shiftIds"), "shiftIds")
	if err != nil {
		return dto.DashboardRequest{}, err
	}

	dashboardRequest := dto.DashboardRequest{
		MachineIDs: machineIDs, From: query.Get("from"), To: query.Get("to"), ShiftIDs: shiftIDs,
		Program: strings.TrimSpace(query.Get("program")), WithSchedule: query.Get("schedule") == "with",
		AllTime: query.Get("all") == "true",
	}
	if minutes := query.Get("minutes"); minutes != "" {
		value, err := strconv.Atoi(minutes)
		if err != nil {
			return dto.DashboardRequest{}, response.NewError(http.StatusBadRequest, "Please enter valid minutes")
		}
		dashboardRequest.Minutes = value
	}
	if date := query.Get("date"); date != "" {
		dashboardRequest.From, dashboardRequest.To = date, date
	}
	return dashboardRequest, nil
}

// queryIDList parses comma-separated ids ("1,2,3"); empty gives none.
func queryIDList(value, name string) ([]int64, error) {
	ids := []int64{}
	for part := range strings.SplitSeq(value, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		id, err := strconv.ParseInt(part, 10, 64)
		if err != nil || id <= 0 {
			return nil, response.NewError(http.StatusBadRequest, "Please enter valid "+name)
		}
		ids = append(ids, id)
	}
	return ids, nil
}
