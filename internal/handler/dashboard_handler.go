package handler

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// DashboardHandler serves the charts of the Productivity (/productivity/...) and Machine Health
// (/machine-health/...) pages. Every logged-in user reads the machines they may access.
type DashboardHandler struct {
	service *service.DashboardService
}

// NewDashboardHandler creates a DashboardHandler that uses dashboardService.
func NewDashboardHandler(dashboardService *service.DashboardService) *DashboardHandler {
	return &DashboardHandler{service: dashboardService}
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
	serveDashboard(writer, request, dashboardHandler.service.GetModeTimeline)
}

// GetOEE handles GET /api/v1/productivity/oee?machineIds=1,2&from=&to=&shiftIds=&program=&schedule=with:
// one OEE row per day, averaged over the machines.
//
// @Summary      Get oee
// @Description  GetOEE handles GET /api/v1/productivity/oee?machineIds=1,2&from=&to=&shiftIds=&program=&schedule=with: one OEE row per day, averaged over the machines.
// @Description
// @Description  Used by pages: Productivity
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
	serveDashboard(writer, request, dashboardHandler.service.GetOEE)
}

// GetCycleCounts handles GET /api/v1/productivity/cycle-counts?machineIds=1,2&from=&to=&shiftIds=:
// cycles started, finished and unfinished per day.
//
// @Summary      Get cycle counts
// @Description  GetCycleCounts handles GET /api/v1/productivity/cycle-counts?machineIds=1,2&from=&to=&shiftIds=: cycles started, finished and unfinished per day.
// @Description
// @Description  Used by pages: Productivity
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
	serveDashboard(writer, request, dashboardHandler.service.GetCycleCounts)
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
	serveDashboard(writer, request, dashboardHandler.service.GetFinishCycles)
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
	serveDashboard(writer, request, dashboardHandler.service.GetMachineTimeByDay)
}

// serveDashboard reads the user and the chart query, then answers with what load returns.
func serveDashboard[T any](writer http.ResponseWriter, request *http.Request,
	load func(context.Context, model.AuthUser, dto.DashboardRequest) (T, error)) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	dashboardRequest, err := dashboardQuery(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	if err := dashboardRequest.Validate(); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	data, err := load(request.Context(), actor, dashboardRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, data)
}

// dashboardQuery reads the chart query: machineIds (or machineId), from / to (or date for one
// day), shiftIds, program, schedule and minutes.
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
