package dto

import (
	"fmt"
	"net/http"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"7th-sense-server/internal/response"
)

// Shift types (shifts.shift_type).
const (
	ShiftTypeDay   = "day"
	ShiftTypeNight = "night"
)

// Shift statuses accepted by UpdateShiftRequest.
const (
	ShiftStatusActive   = "active"
	ShiftStatusInactive = "inactive"
)

const (
	shiftNameMaxLength = 100
	shiftMaxBreaks     = 20
	minutesPerDay      = 24 * 60
)

// clockPattern is a 24-hour "HH:MM" time.
var clockPattern = regexp.MustCompile(`^([01][0-9]|2[0-3]):[0-5][0-9]$`)

// ShiftFields are the Add / Manage Shift form fields shared by create and update.
type ShiftFields struct {
	UserID      int64   `json:"userId"` // customer the shift belongs to
	Name        string  `json:"name"`
	ShiftType   string  `json:"shiftType"`   // ShiftTypeDay or ShiftTypeNight
	StartTime   string  `json:"startTime"`   // "HH:MM"
	EndTime     string  `json:"endTime"`     // "HH:MM"; not after start = past midnight
	WorkingDays []int32 `json:"workingDays"` // ISO weekdays 1 (Mon) … 7 (Sun)
	BreakIDs    []int64 `json:"breakIds"`    // the customer's breaks (POST /breaks)
}

// CreateShiftRequest is the body for POST /shifts.
type CreateShiftRequest struct {
	ShiftFields
}

// Validate trims the fields and checks they are usable.
func (request *CreateShiftRequest) Validate() error {
	return request.ShiftFields.validate()
}

// UpdateShiftRequest is the body for PUT /shifts/{id}. A blank status keeps the saved one.
type UpdateShiftRequest struct {
	ShiftFields
	Status string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateShiftRequest) Validate() error {
	if err := request.ShiftFields.validate(); err != nil {
		return err
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != ShiftStatusActive && request.Status != ShiftStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

func (fields *ShiftFields) validate() error {
	fields.Name = strings.TrimSpace(fields.Name)
	fields.ShiftType = strings.ToLower(strings.TrimSpace(fields.ShiftType))
	fields.StartTime = strings.TrimSpace(fields.StartTime)
	fields.EndTime = strings.TrimSpace(fields.EndTime)

	switch {
	case fields.UserID <= 0:
		return response.NewError(http.StatusBadRequest, "Please select customer")
	case !validText(fields.Name, shiftNameMaxLength):
		return response.NewError(http.StatusBadRequest, "Please enter valid shift name")
	case fields.ShiftType != ShiftTypeDay && fields.ShiftType != ShiftTypeNight:
		return response.NewError(http.StatusBadRequest, "Please select shift type")
	case !clockPattern.MatchString(fields.StartTime):
		return response.NewError(http.StatusBadRequest, "Please enter valid shift start time")
	case !clockPattern.MatchString(fields.EndTime):
		return response.NewError(http.StatusBadRequest, "Please enter valid shift end time")
	case fields.StartTime == fields.EndTime:
		return response.NewError(http.StatusBadRequest, "Shift end time must differ from the start time")
	case len(fields.BreakIDs) > shiftMaxBreaks:
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("A shift can have up to %d breaks", shiftMaxBreaks))
	}

	workingDays, ok := validWorkingDays(fields.WorkingDays)
	if !ok {
		return response.NewError(http.StatusBadRequest, "Please select valid working days")
	}
	fields.WorkingDays = workingDays

	breakIDs, ok := uniquePositiveIDs(fields.BreakIDs)
	if !ok {
		return response.NewError(http.StatusBadRequest, "Please select valid breaks")
	}
	fields.BreakIDs = breakIDs
	return nil
}

// MinutesBetween returns the minutes from start to end ("HH:MM"), past midnight when
// end is not after start: 22:00 → 06:00 is 480. Equal times give a whole day.
func MinutesBetween(start, end string) int {
	minutes := (clockMinutes(end) - clockMinutes(start) + minutesPerDay) % minutesPerDay
	if minutes == 0 {
		return minutesPerDay
	}
	return minutes
}

// clockMinutes returns the minutes after midnight of a valid "HH:MM".
func clockMinutes(clock string) int {
	hours, _ := strconv.Atoi(clock[:2])
	minutes, _ := strconv.Atoi(clock[3:])
	return hours*60 + minutes
}

// validWorkingDays sorts and de-duplicates days; ok is false when there is none or a
// day is not 1–7.
func validWorkingDays(days []int32) ([]int32, bool) {
	if len(days) == 0 {
		return nil, false
	}
	sorted := slices.Clone(days)
	slices.Sort(sorted)
	sorted = slices.Compact(sorted)
	if sorted[0] < 1 || sorted[len(sorted)-1] > 7 {
		return nil, false
	}
	return sorted, true
}
