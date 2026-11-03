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
func (dashboardHandler *DashboardHandler) GetModeTimeline(writer http.ResponseWriter, request *http.Request) {
	serveDashboard(writer, request, dashboardHandler.service.GetModeTimeline)
}

// GetOEE handles GET /api/v1/productivity/oee?machineIds=1,2&from=&to=&shiftIds=&program=&schedule=with:
// one OEE row per day, averaged over the machines.
func (dashboardHandler *DashboardHandler) GetOEE(writer http.ResponseWriter, request *http.Request) {
	serveDashboard(writer, request, dashboardHandler.service.GetOEE)
}

// GetCycleCounts handles GET /api/v1/productivity/cycle-counts?machineIds=1,2&from=&to=&shiftIds=:
// cycles started, finished and unfinished per day.
func (dashboardHandler *DashboardHandler) GetCycleCounts(writer http.ResponseWriter, request *http.Request) {
	serveDashboard(writer, request, dashboardHandler.service.GetCycleCounts)
}

// GetFinishCycles handles GET /api/v1/productivity/finish-cycles?machineId=&date=&shiftIds=&program=:
// the machine's completed cycles of the day with the day's totals.
func (dashboardHandler *DashboardHandler) GetFinishCycles(writer http.ResponseWriter, request *http.Request) {
	serveDashboard(writer, request, dashboardHandler.service.GetFinishCycles)
}

// GetMachineTimeByDay handles GET /api/v1/machine-health/machine-time-by-day?machineId=&from=&to=:
// hours per machine mode for each day.
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
