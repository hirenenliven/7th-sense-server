package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// BreakService holds the business rules for breaks: a customer's (user_id) breaks, each
// taken by at most one of the customer's shifts (shift_id, optional). Everyone reads their
// company's breaks (the superadmin every customer's); managers change their company's
// (see companyOwner). A break on a shift must lie inside it and not overlap the shift's
// other breaks (see checkBreaksFit).
type BreakService struct {
	breakRepository *repository.BreakRepository
	shiftRepository *repository.ShiftRepository
	userRepository  *repository.UserRepository
}

// NewBreakService creates a BreakService that stores breaks in dbPool and checks
// customers with userRepository.
func NewBreakService(dbPool *pgxpool.Pool, userRepository *repository.UserRepository) *BreakService {
	return &BreakService{
		breakRepository: repository.NewBreakRepository(dbPool),
		shiftRepository: repository.NewShiftRepository(dbPool),
		userRepository:  userRepository,
	}
}

// GetAll returns the breaks actor may read. A non-nil customerID keeps only that
// customer's breaks.
func (breakService *BreakService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.Break, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	return breakService.breakRepository.FindAll(ctx, ownerID)
}

// GetByID returns the break when actor may read it, else "Break not found".
func (breakService *BreakService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.Break, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return model.Break{}, err
	}
	found, err := breakService.breakRepository.FindByID(ctx, ownerID, id)
	return found, toBreakResponseError(err)
}

// Create adds a break, on the shift named or on none. A customer or admin (actor) adds
// it to their own company, whatever customer the request names; the superadmin to the
// customer named.
func (breakService *BreakService) Create(ctx context.Context, actor model.AuthUser, breakRequest dto.CreateBreakRequest) (model.Break, error) {
	if err := ownCompanyCustomer(actor, &breakRequest.UserID); err != nil {
		return model.Break{}, err
	}
	shiftBreak := breakFromFields(breakRequest.BreakFields)
	shiftBreak.Status = dto.BreakStatusActive
	if err := breakService.checkReferences(ctx, shiftBreak); err != nil {
		return model.Break{}, err
	}

	id, err := breakService.breakRepository.Create(ctx, shiftBreak)
	if err != nil {
		return model.Break{}, toBreakResponseError(err)
	}
	return breakService.find(ctx, id)
}

// Update replaces the break, its shift included (nil takes it off its shift). A blank
// status keeps the saved one. A customer or admin (actor) changes only their company's
// breaks ("Break not found" for others) and cannot move one to another company.
func (breakService *BreakService) Update(ctx context.Context, actor model.AuthUser, id int64, breakRequest dto.UpdateBreakRequest) (model.Break, error) {
	if err := ownCompanyCustomer(actor, &breakRequest.UserID); err != nil {
		return model.Break{}, err
	}
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return model.Break{}, err
	}
	if _, err := breakService.breakRepository.FindByID(ctx, ownerID, id); err != nil {
		return model.Break{}, toBreakResponseError(err)
	}
	shiftBreak := breakFromFields(breakRequest.BreakFields)
	shiftBreak.ID = id
	shiftBreak.Status = breakRequest.Status
	if err := breakService.checkReferences(ctx, shiftBreak); err != nil {
		return model.Break{}, err
	}

	if err := breakService.breakRepository.Update(ctx, shiftBreak); err != nil {
		return model.Break{}, toBreakResponseError(err)
	}
	return breakService.find(ctx, id)
}

// Delete soft-deletes the break (its shift no longer takes it). A customer or admin
// (actor) deletes only their company's breaks ("Break not found" for others).
func (breakService *BreakService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return err
	}
	return toBreakResponseError(breakService.breakRepository.Delete(ctx, ownerID, id))
}

// checkReferences returns a 400 unless the customer exists and, when the break is on a
// shift, the shift is one of the customer's and the break fits it beside the shift's
// other breaks.
func (breakService *BreakService) checkReferences(ctx context.Context, shiftBreak model.Break) error {
	isCustomer, err := breakService.userRepository.IsCustomer(ctx, shiftBreak.UserID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	if shiftBreak.ShiftID == nil {
		return nil
	}

	shift, err := breakService.shiftRepository.FindByID(ctx, &shiftBreak.UserID, *shiftBreak.ShiftID)
	if errors.Is(err, repository.ErrNotFound) {
		return response.NewError(http.StatusBadRequest, "Please select a shift of the selected customer")
	}
	if err != nil {
		return err
	}
	shiftBreaks, err := breakService.breakRepository.FindByShifts(ctx, []int64{shift.ID})
	if err != nil {
		return err
	}
	breaks := []model.Break{shiftBreak}
	for _, other := range shiftBreaks {
		if other.ID != shiftBreak.ID {
			breaks = append(breaks, model.Break{ID: other.ID, Name: other.Name, StartTime: other.StartTime, EndTime: other.EndTime})
		}
	}
	if err := checkBreaksFit(shift.StartTime, shift.EndTime, breaks); err != nil {
		return response.NewError(http.StatusBadRequest, fmt.Sprintf("Shift %s: %s", shift.Name, err.Error()))
	}
	return nil
}

// find returns break id whoever owns it (for the responses of Create and Update).
func (breakService *BreakService) find(ctx context.Context, id int64) (model.Break, error) {
	found, err := breakService.breakRepository.FindByID(ctx, nil, id)
	return found, toBreakResponseError(err)
}

func breakFromFields(fields dto.BreakFields) model.Break {
	return model.Break{
		UserID:    fields.UserID,
		ShiftID:   fields.ShiftID,
		Name:      fields.Name,
		StartTime: fields.StartTime,
		EndTime:   fields.EndTime,
	}
}

// toBreakResponseError maps repository errors to break responses. Other errors pass
// through and are sent as ServerError.
func toBreakResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Break not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "This customer already has a break with this name")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	case errors.Is(err, repository.ErrInvalidParent):
		return response.NewError(http.StatusBadRequest, "Please select valid shift")
	default:
		return err
	}
}
