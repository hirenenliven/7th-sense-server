package dto

import (
	"net/http"
	"strings"

	"7th-sense-server/internal/response"
)

// Holiday kinds (holidays.holiday_type).
const (
	HolidayTypeHoliday = "holiday"
	HolidayTypeWeekOff = "week_off"
)

// Holiday statuses accepted by UpdateHolidayRequest.
const (
	HolidayStatusActive   = "active"
	HolidayStatusInactive = "inactive"
)

const holidayNameMaxLength = 100

// HolidayFields are the Add Holiday / Week Off form fields shared by create and update.
// A holiday needs Name, StartDate and EndDate; a week off needs WeekDay (Name optional).
// The fields the kind does not use are cleared.
type HolidayFields struct {
	UserID      int64  `json:"userId"`      // customer the holiday belongs to
	HolidayType string `json:"holidayType"` // HolidayTypeHoliday or HolidayTypeWeekOff
	Name        string `json:"name"`
	StartDate   string `json:"startDate"` // "YYYY-MM-DD", holiday only
	EndDate     string `json:"endDate"`   // "YYYY-MM-DD", holiday only
	WeekDay     int16  `json:"weekDay"`   // 1 (Mon) … 7 (Sun), week off only
}

// CreateHolidayRequest is the body for POST /holidays.
type CreateHolidayRequest struct {
	HolidayFields
}

// Validate trims the fields and checks they are usable.
func (request *CreateHolidayRequest) Validate() error {
	return request.HolidayFields.validate()
}

// UpdateHolidayRequest is the body for PUT /holidays/{id}. A blank status keeps the saved one.
type UpdateHolidayRequest struct {
	HolidayFields
	Status string `json:"status"`
}

// Validate trims the fields and checks they are usable.
func (request *UpdateHolidayRequest) Validate() error {
	if err := request.HolidayFields.validate(); err != nil {
		return err
	}
	request.Status = strings.ToLower(strings.TrimSpace(request.Status))
	if request.Status != "" && request.Status != HolidayStatusActive && request.Status != HolidayStatusInactive {
		return response.NewError(http.StatusBadRequest, "Please enter valid status")
	}
	return nil
}

func (fields *HolidayFields) validate() error {
	fields.HolidayType = strings.ToLower(strings.TrimSpace(fields.HolidayType))
	fields.Name = strings.TrimSpace(fields.Name)
	fields.StartDate = strings.TrimSpace(fields.StartDate)
	fields.EndDate = strings.TrimSpace(fields.EndDate)

	if fields.UserID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please select customer")
	}

	switch fields.HolidayType {
	case HolidayTypeHoliday:
		fields.WeekDay = 0
		switch {
		case !validText(fields.Name, holidayNameMaxLength):
			return response.NewError(http.StatusBadRequest, "Please enter valid holiday name")
		case !validDate(fields.StartDate):
			return response.NewError(http.StatusBadRequest, "Please enter valid start date")
		case !validDate(fields.EndDate):
			return response.NewError(http.StatusBadRequest, "Please enter valid end date")
		case fields.EndDate < fields.StartDate: // YYYY-MM-DD sorts as text
			return response.NewError(http.StatusBadRequest, "End date cannot be before the start date")
		}
	case HolidayTypeWeekOff:
		fields.StartDate, fields.EndDate = "", ""
		switch {
		case fields.WeekDay < 1 || fields.WeekDay > 7:
			return response.NewError(http.StatusBadRequest, "Please select week off day")
		case fields.Name != "" && !validText(fields.Name, holidayNameMaxLength):
			return response.NewError(http.StatusBadRequest, "Please enter valid name")
		}
	default:
		return response.NewError(http.StatusBadRequest, "Please select holiday or week off")
	}
	return nil
}
