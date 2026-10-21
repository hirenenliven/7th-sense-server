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

// IdleReasonService holds the business rules for idle reasons: why a machine of a
// customer (user_id) stood idle. Everyone reads their company's reasons (the superadmin
// every customer's); managers change their company's (see companyOwner).
type IdleReasonService struct {
	idleReasonRepository *repository.IdleReasonRepository
	userRepository           *repository.UserRepository
}

// NewIdleReasonService creates a IdleReasonService that stores reasons with
// idleReasonRepository and checks customers with userRepository.
func NewIdleReasonService(idleReasonRepository *repository.IdleReasonRepository, userRepository *repository.UserRepository) *IdleReasonService {
	return &IdleReasonService{idleReasonRepository: idleReasonRepository, userRepository: userRepository}
}

// GetAll returns the idle reasons actor may read. A non-nil customerID keeps only that
// customer's reasons.
func (idleReasonService *IdleReasonService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.IdleReason, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	return idleReasonService.idleReasonRepository.FindAll(ctx, ownerID)
}

// GetByID returns the idle reason when actor may read it, else "Idle reason not found".
func (idleReasonService *IdleReasonService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.IdleReason, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return model.IdleReason{}, err
	}
	reason, err := idleReasonService.idleReasonRepository.FindByID(ctx, ownerID, id)
	return reason, toIdleReasonResponseError(err)
}

// Create adds a idle reason. A customer or admin (actor) adds it to their own company,
// whatever customer the request names; the superadmin to the customer named.
func (idleReasonService *IdleReasonService) Create(ctx context.Context, actor model.AuthUser, reasonRequest dto.CreateIdleReasonRequest) (model.IdleReason, error) {
	if err := ownCompanyCustomer(actor, &reasonRequest.UserID); err != nil {
		return model.IdleReason{}, err
	}
	if err := idleReasonService.checkCustomer(ctx, reasonRequest.UserID); err != nil {
		return model.IdleReason{}, err
	}
	reason := idleReasonFromFields(reasonRequest.IdleReasonFields)
	reason.Status = dto.IdleReasonStatusActive

	id, err := idleReasonService.idleReasonRepository.Create(ctx, reason)
	if err != nil {
		return model.IdleReason{}, toIdleReasonResponseError(err)
	}
	return idleReasonService.find(ctx, id)
}

// Update replaces the idle reason. A blank status keeps the saved one. A customer or
// admin (actor) changes only their company's reasons ("Idle reason not found" for
// others) and cannot move one to another company.
func (idleReasonService *IdleReasonService) Update(ctx context.Context, actor model.AuthUser, id int64, reasonRequest dto.UpdateIdleReasonRequest) (model.IdleReason, error) {
	if err := ownCompanyCustomer(actor, &reasonRequest.UserID); err != nil {
		return model.IdleReason{}, err
	}
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return model.IdleReason{}, err
	}
	if _, err := idleReasonService.idleReasonRepository.FindByID(ctx, ownerID, id); err != nil {
		return model.IdleReason{}, toIdleReasonResponseError(err)
	}
	if err := idleReasonService.checkCustomer(ctx, reasonRequest.UserID); err != nil {
		return model.IdleReason{}, err
	}
	reason := idleReasonFromFields(reasonRequest.IdleReasonFields)
	reason.ID = id
	reason.Status = reasonRequest.Status

	if err := idleReasonService.idleReasonRepository.Update(ctx, reason); err != nil {
		return model.IdleReason{}, toIdleReasonResponseError(err)
	}
	return idleReasonService.find(ctx, id)
}

// Delete soft-deletes the idle reason. A customer or admin (actor) deletes only their
// company's reasons ("Idle reason not found" for others).
func (idleReasonService *IdleReasonService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return err
	}
	return toIdleReasonResponseError(idleReasonService.idleReasonRepository.Delete(ctx, ownerID, id))
}

// find returns idle reason id whoever owns it (for the responses of Create and Update).
func (idleReasonService *IdleReasonService) find(ctx context.Context, id int64) (model.IdleReason, error) {
	reason, err := idleReasonService.idleReasonRepository.FindByID(ctx, nil, id)
	return reason, toIdleReasonResponseError(err)
}

// checkCustomer returns "Please select valid customer" unless userID is a customer.
func (idleReasonService *IdleReasonService) checkCustomer(ctx context.Context, userID int64) error {
	isCustomer, err := idleReasonService.userRepository.IsCustomer(ctx, userID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	return nil
}

func idleReasonFromFields(fields dto.IdleReasonFields) model.IdleReason {
	return model.IdleReason{
		UserID:     fields.UserID,
		ReasonCode: fields.ReasonCode,
		Reason:     fields.Reason,
	}
}

// toIdleReasonResponseError maps repository errors to idle reason responses. Other
// errors pass through and are sent as ServerError.
func toIdleReasonResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Idle reason not found")
	case errors.Is(err, repository.ErrDuplicateCode):
		return response.NewError(http.StatusBadRequest, "This customer already has a idle reason with this reason code")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "This customer already has this idle reason")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	default:
		return err
	}
}
