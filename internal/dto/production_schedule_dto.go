package dto

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"7th-sense-server/internal/response"
)

// Production schedule types: one day, or a range of days.
const (
	ScheduleTypeSingle   = "single"
	ScheduleTypeMultiple = "multiple"
)

// Production schedule statuses accepted by UpdateProductionScheduleRequest.
const (
	ScheduleStatusActive   = "active"
	ScheduleStatusInactive = "inactive"
)

// Limits of the production_schedules and production_schedule_items columns (migration 026).
const (
	scheduleJobIDMaxLength = 50
	scheduleMaxItems       = 100
	scheduleMaxQuantity    = 1000000
	scheduleMaxDays        = 366
)

// ProductionScheduleItemFields is one job row of the Add / Edit Production Schedule form.
// ID names a saved item of the schedule when editing; nil adds a new one.
type ProductionScheduleItemFields struct {
	ID           *int64 `json:"id"`
	ScheduleDate string `json:"scheduleDate"` // a day of the schedule; single day: the schedule date
	JobID        string `json:"jobId"`
	SupervisorID int64  `json:"supervisorId"` // a supervisor of the customer
	OperatorID   int64  `json:"operatorId"`   // an operator of the customer
	ShiftID      int64  `json:"shiftId"`      // a shift of the customer
	PartID       int64  `json:"partId"`       // a part of the customer
	ProcessID    int64  `json:"processId"`    // a process of the part that runs on the machine
	Quantity     int    `json:"quantity"`
	StartTime    string `json:"startTime"` // "HH:MM", inside the shift
	EndTime      string `json:"endTime"`   // "HH:MM"; not after the start runs past midnight
}

// ProductionScheduleFields are the form fields shared by create and update.
type ProductionScheduleFields struct {
	UserID       int64                          `json:"userId"`    // customer that owns the schedule
	MachineID    int64                          `json:"machineId"` // a machine of the customer
	ScheduleType string                         `json:"scheduleType"`
	StartDate    string                         `json:"startDate"`
	EndDate      string                         `json:"endDate"` // ignored for a single day
	Items        []ProductionScheduleItemFields `json:"items"`
}

// CreateProductionScheduleRequest is the body for POST /production-schedules.
type CreateProductionScheduleRequest struct {
	ProductionScheduleFields
}

// Validate trims the fields and checks they are usable.
func (request *CreateProductionScheduleRequest) Validate() error {
	return request.ProductionScheduleFields.validate()
}

// UpdateProductionScheduleRequest is the body for PUT /production-schedules/{id}.
// A blank status keeps the saved one.
type UpdateProductionScheduleRequest struct {
	ProductionScheduleFields
	Status string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateProductionScheduleRequest) Validate() error {
	if err := request.ProductionScheduleFields.validate(); err != nil {
		return err
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != ScheduleStatusActive && request.Status != ScheduleStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

func (fields *ProductionScheduleFields) validate() error {
	fields.ScheduleType = strings.ToLower(strings.TrimSpace(fields.ScheduleType))
	fields.StartDate = strings.TrimSpace(fields.StartDate)
	fields.EndDate = strings.TrimSpace(fields.EndDate)
	if fields.ScheduleType == ScheduleTypeSingle {
		fields.EndDate = fields.StartDate
	}

	switch {
	case fields.UserID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select customer")
	case fields.MachineID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select machine")
	case fields.ScheduleType != ScheduleTypeSingle && fields.ScheduleType != ScheduleTypeMultiple:
		return response.NewError(http.StatusBadRequest, "Please select schedule type")
	case !validDate(fields.StartDate):
		return response.NewError(http.StatusBadRequest, "Please enter valid schedule date")
	case !validDate(fields.EndDate):
		return response.NewError(http.StatusBadRequest, "Please enter valid schedule end date")
	case fields.EndDate < fields.StartDate: // YYYY-MM-DD sorts as text
		return response.NewError(http.StatusBadRequest, "Schedule end date cannot be before the start date")
	case ScheduleDay(fields.EndDate)-ScheduleDay(fields.StartDate) >= scheduleMaxDays:
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("A schedule can cover up to %d days", scheduleMaxDays))
	case len(fields.Items) == 0:
		return response.NewError(http.StatusBadRequest, "Please add at least one job")
	case len(fields.Items) > scheduleMaxItems:
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("A schedule can have up to %d jobs", scheduleMaxItems))
	}

	seenIDs := make(map[int64]bool, len(fields.Items))
	for index := range fields.Items {
		item := &fields.Items[index]
		if fields.ScheduleType == ScheduleTypeSingle {
			item.ScheduleDate = fields.StartDate
		}
		if err := item.validate(index+1, fields.StartDate, fields.EndDate); err != nil {
			return err
		}
		if item.ID != nil {
			if seenIDs[*item.ID] {
				return response.NewError(http.StatusBadRequest, "Please enter valid jobs")
			}
			seenIDs[*item.ID] = true
		}
	}
	return nil
}

// validate checks job number position (1-based, used in the messages) of a schedule
// from startDate to endDate.
func (item *ProductionScheduleItemFields) validate(position int, startDate, endDate string) error {
	item.ScheduleDate = strings.TrimSpace(item.ScheduleDate)
	item.JobID = strings.TrimSpace(item.JobID)
	item.StartTime = strings.TrimSpace(item.StartTime)
	item.EndTime = strings.TrimSpace(item.EndTime)
	invalid := func(message string) error {
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("%s (job %d)", message, position))
	}

	switch {
	case item.ID != nil && *item.ID <= 0:
		return invalid("Please enter valid job")
	case !validDate(item.ScheduleDate) || item.ScheduleDate < startDate || item.ScheduleDate > endDate:
		return invalid("Please select a date of the schedule")
	case !validText(item.JobID, scheduleJobIDMaxLength):
		return invalid("Please enter valid job ID")
	case item.SupervisorID <= 0:
		return invalid("Please select supervisor / programmer")
	case item.OperatorID <= 0:
		return invalid("Please select operator")
	case item.ShiftID <= 0:
		return invalid("Please select shift")
	case item.PartID <= 0:
		return invalid("Please select part")
	case item.ProcessID <= 0:
		return invalid("Please select part process")
	case item.Quantity < 1 || item.Quantity > scheduleMaxQuantity:
		return invalid("Please enter valid quantity")
	case !clockPattern.MatchString(item.StartTime):
		return invalid("Please enter valid start time")
	case !clockPattern.MatchString(item.EndTime):
		return invalid("Please enter valid end time")
	case item.StartTime == item.EndTime:
		return invalid("End time cannot be the same as the start time")
	}
	return nil
}

// ScheduleDay returns the day number of a valid "YYYY-MM-DD" date.
func ScheduleDay(date string) int {
	day, _ := time.Parse(dateLayout, date)
	return int(day.Unix() / 86400)
}
