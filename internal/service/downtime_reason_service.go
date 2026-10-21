package service

import (
	"context"
	"errors"
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// DowntimeReasonService holds the business rules for downtime reasons: why a machine of a
// customer (user_id) was down. Everyone reads their company's reasons (the superadmin
// every customer's); managers change their company's (see companyOwner).
type DowntimeReasonService struct {
	downtimeReasonRepository *repository.DowntimeReasonRepository
	userRepository           *repository.UserRepository
}

// NewDowntimeReasonService creates a DowntimeReasonService that stores reasons with
// downtimeReasonRepository and checks customers with userRepository.
func NewDowntimeReasonService(downtimeReasonRepository *repository.DowntimeReasonRepository, userRepository *repository.UserRepository) *DowntimeReasonService {
	return &DowntimeReasonService{downtimeReasonRepository: downtimeReasonRepository, userRepository: userRepository}
}

// GetAll returns the downtime reasons actor may read. A non-nil customerID keeps only that
// customer's reasons.
func (downtimeReasonService *DowntimeReasonService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.DowntimeReason, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	return downtimeReasonService.downtimeReasonRepository.FindAll(ctx, ownerID)
}

// GetByID returns the downtime reason when actor may read it, else "Downtime reason not found".
func (downtimeReasonService *DowntimeReasonService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.DowntimeReason, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return model.DowntimeReason{}, err
	}
	reason, err := downtimeReasonService.downtimeReasonRepository.FindByID(ctx, ownerID, id)
	return reason, toDowntimeReasonResponseError(err)
}

// Create adds a downtime reason. A customer or admin (actor) adds it to their own company,
// whatever customer the request names; the superadmin to the customer named.
func (downtimeReasonService *DowntimeReasonService) Create(ctx context.Context, actor model.AuthUser, reasonRequest dto.CreateDowntimeReasonRequest) (model.DowntimeReason, error) {
	if err := ownCompanyCustomer(actor, &reasonRequest.UserID); err != nil {
		return model.DowntimeReason{}, err
	}
	if err := downtimeReasonService.checkCustomer(ctx, reasonRequest.UserID); err != nil {
		return model.DowntimeReason{}, err
	}
	reason := downtimeReasonFromFields(reasonRequest.DowntimeReasonFields)
	reason.Status = dto.DowntimeReasonStatusActive

	id, err := downtimeReasonService.downtimeReasonRepository.Create(ctx, reason)
	if err != nil {
		return model.DowntimeReason{}, toDowntimeReasonResponseError(err)
	}
	return downtimeReasonService.find(ctx, id)
}

// Update replaces the downtime reason. A blank status keeps the saved one. A customer or
// admin (actor) changes only their company's reasons ("Downtime reason not found" for
// others) and cannot move one to another company.
func (downtimeReasonService *DowntimeReasonService) Update(ctx context.Context, actor model.AuthUser, id int64, reasonRequest dto.UpdateDowntimeReasonRequest) (model.DowntimeReason, error) {
	if err := ownCompanyCustomer(actor, &reasonRequest.UserID); err != nil {
		return model.DowntimeReason{}, err
	}
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return model.DowntimeReason{}, err
	}
	if _, err := downtimeReasonService.downtimeReasonRepository.FindByID(ctx, ownerID, id); err != nil {
		return model.DowntimeReason{}, toDowntimeReasonResponseError(err)
	}
	if err := downtimeReasonService.checkCustomer(ctx, reasonRequest.UserID); err != nil {
		return model.DowntimeReason{}, err
	}
	reason := downtimeReasonFromFields(reasonRequest.DowntimeReasonFields)
	reason.ID = id
	reason.Status = reasonRequest.Status

	if err := downtimeReasonService.downtimeReasonRepository.Update(ctx, reason); err != nil {
		return model.DowntimeReason{}, toDowntimeReasonResponseError(err)
	}
	return downtimeReasonService.find(ctx, id)
}

// Delete soft-deletes the downtime reason. A customer or admin (actor) deletes only their
// company's reasons ("Downtime reason not found" for others).
func (downtimeReasonService *DowntimeReasonService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return err
	}
	return toDowntimeReasonResponseError(downtimeReasonService.downtimeReasonRepository.Delete(ctx, ownerID, id))
}

// find returns downtime reason id whoever owns it (for the responses of Create and Update).
func (downtimeReasonService *DowntimeReasonService) find(ctx context.Context, id int64) (model.DowntimeReason, error) {
	reason, err := downtimeReasonService.downtimeReasonRepository.FindByID(ctx, nil, id)
	return reason, toDowntimeReasonResponseError(err)
}

// checkCustomer returns "Please select valid customer" unless userID is a customer.
func (downtimeReasonService *DowntimeReasonService) checkCustomer(ctx context.Context, userID int64) error {
	isCustomer, err := downtimeReasonService.userRepository.IsCustomer(ctx, userID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	return nil
}

func downtimeReasonFromFields(fields dto.DowntimeReasonFields) model.DowntimeReason {
	return model.DowntimeReason{
		UserID:     fields.UserID,
		ReasonCode: fields.ReasonCode,
		Reason:     fields.Reason,
	}
}

// toDowntimeReasonResponseError maps repository errors to downtime reason responses. Other
// errors pass through and are sent as ServerError.
func toDowntimeReasonResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Downtime reason not found")
	case errors.Is(err, repository.ErrDuplicateCode):
		return response.NewError(http.StatusBadRequest, "This customer already has a downtime reason with this reason code")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "This customer already has this downtime reason")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	default:
		return err
	}
}
