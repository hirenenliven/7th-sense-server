package dto

import (
	"math"
	"net/http"
	"time"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
)

// MachineModeReasonFields are the fields of the Update Downtime / Idle Reason form shared by
// create and update: the reason (of the master list the period's mode uses) and the part of the
// period it explains.
type MachineModeReasonFields struct {
	ReasonID  int64     `json:"reasonId"`
	StartedAt time.Time `json:"startedAt"`
	EndedAt   time.Time `json:"endedAt"`
}

// CreateMachineModeReasonRequest is the body for POST /machine-mode-reasons.
type CreateMachineModeReasonRequest struct {
	PeriodID int64 `json:"periodId"`
	MachineModeReasonFields
}

// Validate checks the period and the fields are usable.
func (request *CreateMachineModeReasonRequest) Validate() error {
	if request.PeriodID <= 0 {
		return response.NewError(http.StatusBadRequest, "Please select machine mode period")
	}
	return request.MachineModeReasonFields.validate()
}

// UpdateMachineModeReasonRequest is the body for PUT /machine-mode-reasons/{id}.
type UpdateMachineModeReasonRequest struct {
	MachineModeReasonFields
}

// Validate checks the fields are usable.
func (request *UpdateMachineModeReasonRequest) Validate() error {
	return request.MachineModeReasonFields.validate()
}

func (fields *MachineModeReasonFields) validate() error {
	switch {
	case fields.ReasonID <= 0 || fields.ReasonID > math.MaxInt32: // the reason columns are integer
		return response.NewError(http.StatusBadRequest, "Please select reason")
	case fields.StartedAt.IsZero():
		return response.NewError(http.StatusBadRequest, "Please select start time")
	case fields.EndedAt.IsZero():
		return response.NewError(http.StatusBadRequest, "Please select end time")
	case !fields.EndedAt.After(fields.StartedAt):
		return response.NewError(http.StatusBadRequest, "End time must be after start time")
	}
	return nil
}

// MachineModeReasonsResponse is the response of GET /machine-mode-reasons: the whole period, what
// the popup needs to offer its reasons (the master list ReasonType, of the machine's customer
// CustomerID) and the reasons recorded so far, in time order.
type MachineModeReasonsResponse struct {
	Period     model.ReasonPeriod        `json:"period"`
	ReasonType string                    `json:"reasonType"`
	CustomerID *int64                    `json:"customerId"`
	Reasons    []model.MachineModeReason `json:"reasons"`
}
