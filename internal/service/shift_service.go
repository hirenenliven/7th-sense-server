package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// ShiftService holds the business rules for shifts. A shift belongs to a customer
// (user_id) and takes some of the customer's breaks (Masters › Breaks: breaks.shift_id),
// assigned with the shift in one transaction; they must lie inside the shift, not
// overlap and not be taken by another shift. Everyone
// reads their company's shifts (the superadmin every customer's); managers change their
// company's (see companyOwner). The minutes of a shift and of its breaks are worked out here.
type ShiftService struct {
	dbPool          *pgxpool.Pool
	shiftRepository *repository.ShiftRepository
	breakRepository *repository.BreakRepository
	userRepository  *repository.UserRepository
}

// NewShiftService creates a ShiftService that stores shifts in dbPool and checks
// customers with userRepository.
func NewShiftService(dbPool *pgxpool.Pool, userRepository *repository.UserRepository) *ShiftService {
	return &ShiftService{
		dbPool:          dbPool,
		shiftRepository: repository.NewShiftRepository(dbPool),
		breakRepository: repository.NewBreakRepository(dbPool),
		userRepository:  userRepository,
	}
}

// GetAll returns the shifts actor may read, each with its breaks. A non-nil
// customerID keeps only that customer's shifts.
func (shiftService *ShiftService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.Shift, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	shifts, err := shiftService.shiftRepository.FindAll(ctx, ownerID)
	if err != nil {
		return nil, err
	}
	return shiftsWithBreaks(ctx, shiftService.breakRepository, shifts)
}

// GetByID returns the shift with its breaks when actor may read it, else "Shift not found".
func (shiftService *ShiftService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.Shift, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return model.Shift{}, err
	}
	return shiftDetails(ctx, shiftService.shiftRepository, shiftService.breakRepository, ownerID, id)
}

// Create adds a shift with its breaks. A customer or admin (actor) adds it to their
// own company, whatever customer the request names; the superadmin to the customer named.
func (shiftService *ShiftService) Create(ctx context.Context, actor model.AuthUser, shiftRequest dto.CreateShiftRequest) (model.Shift, error) {
	if err := ownCompanyCustomer(actor, &shiftRequest.UserID); err != nil {
		return model.Shift{}, err
	}
	if err := shiftService.checkReferences(ctx, shiftRequest.ShiftFields, 0); err != nil {
		return model.Shift{}, err
	}
	shift := shiftFromFields(shiftRequest.ShiftFields)
	shift.Status = dto.ShiftStatusActive

	var created model.Shift
	err := db.WithTransaction(ctx, shiftService.dbPool, func(tx pgx.Tx) error {
		shiftRepository := repository.NewShiftRepository(tx)
		breakRepository := repository.NewBreakRepository(tx)

		id, err := shiftRepository.Create(ctx, shift)
		if err != nil {
			return toShiftResponseError(err)
		}
		if err := breakRepository.AssignShift(ctx, id, shiftRequest.BreakIDs); err != nil {
			return err
		}
		created, err = shiftDetails(ctx, shiftRepository, breakRepository, nil, id)
		return err
	})
	return created, err
}

// Update replaces the shift and its breaks: breaks left out are released (no shift).
// A blank status keeps the saved one. A
// customer or admin (actor) changes only their company's shifts ("Shift not found" for
// others) and cannot move one to another company.
func (shiftService *ShiftService) Update(ctx context.Context, actor model.AuthUser, id int64, shiftRequest dto.UpdateShiftRequest) (model.Shift, error) {
	if err := ownCompanyCustomer(actor, &shiftRequest.UserID); err != nil {
		return model.Shift{}, err
	}
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return model.Shift{}, err
	}
	if _, err := shiftService.shiftRepository.FindByID(ctx, ownerID, id); err != nil {
		return model.Shift{}, toShiftResponseError(err)
	}
	if err := shiftService.checkReferences(ctx, shiftRequest.ShiftFields, id); err != nil {
		return model.Shift{}, err
	}
	shift := shiftFromFields(shiftRequest.ShiftFields)
	shift.ID = id
	shift.Status = shiftRequest.Status

	var updated model.Shift
	err = db.WithTransaction(ctx, shiftService.dbPool, func(tx pgx.Tx) error {
		shiftRepository := repository.NewShiftRepository(tx)
		breakRepository := repository.NewBreakRepository(tx)

		if err := shiftRepository.Update(ctx, shift); err != nil {
			return toShiftResponseError(err)
		}
		if err := breakRepository.AssignShift(ctx, id, shiftRequest.BreakIDs); err != nil {
			return err
		}
		updated, err = shiftDetails(ctx, shiftRepository, breakRepository, nil, id)
		return err
	})
	return updated, err
}

// Delete soft-deletes the shift and releases its breaks (they stay in the customer's
// list without a shift). A customer or admin (actor) deletes only their company's shifts
// ("Shift not found" for others).
func (shiftService *ShiftService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return err
	}
	return db.WithTransaction(ctx, shiftService.dbPool, func(tx pgx.Tx) error {
		if err := repository.NewShiftRepository(tx).Delete(ctx, ownerID, id); err != nil {
			return toShiftResponseError(err)
		}
		return repository.NewBreakRepository(tx).ReleaseShift(ctx, id)
	})
}

// checkReferences returns a 400 unless the customer exists, every break is one of the
// customer's, free or already taken by shiftID (0 for a new shift), and the breaks fit
// the shift (see checkBreaksFit).
func (shiftService *ShiftService) checkReferences(ctx context.Context, fields dto.ShiftFields, shiftID int64) error {
	isCustomer, err := shiftService.userRepository.IsCustomer(ctx, fields.UserID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	if len(fields.BreakIDs) == 0 {
		return nil
	}

	breaks, err := shiftService.breakRepository.FindByIDs(ctx, fields.UserID, fields.BreakIDs)
	if err != nil {
		return err
	}
	if len(breaks) != len(fields.BreakIDs) {
		return response.NewError(http.StatusBadRequest, "Please select breaks of the selected customer")
	}
	for _, shiftBreak := range breaks {
		if shiftBreak.ShiftID != nil && *shiftBreak.ShiftID != shiftID {
			return response.NewError(http.StatusBadRequest,
				fmt.Sprintf("Break %s is already taken by shift %s", shiftBreak.Name, shiftBreak.ShiftName))
		}
	}
	return checkBreaksFit(fields.StartTime, fields.EndTime, breaks)
}

// checkBreaksFit returns a 400 naming the break unless every break lies inside the
// shift startTime–endTime and no two overlap. Times are compared as minutes from the
// shift start, so a shift or break past midnight works too.
func checkBreaksFit(startTime, endTime string, breaks []model.Break) error {
	type span struct {
		name     string
		from, to int
	}
	shiftLength := dto.MinutesBetween(startTime, endTime)
	spans := make([]span, 0, len(breaks))
	for _, shiftBreak := range breaks {
		from := minutesFromShiftStart(startTime, shiftBreak.StartTime)
		to := from + dto.MinutesBetween(shiftBreak.StartTime, shiftBreak.EndTime)
		if to > shiftLength {
			return response.NewError(http.StatusBadRequest,
				fmt.Sprintf("Break %s (%s to %s) is not within the shift time", shiftBreak.Name, shiftBreak.StartTime, shiftBreak.EndTime))
		}
		spans = append(spans, span{shiftBreak.Name, from, to})
	}

	slices.SortFunc(spans, func(a, b span) int { return a.from - b.from })
	for index := 1; index < len(spans); index++ {
		if spans[index].from < spans[index-1].to {
			return response.NewError(http.StatusBadRequest,
				fmt.Sprintf("Breaks %s and %s overlap", spans[index-1].name, spans[index].name))
		}
	}
	return nil
}

// minutesFromShiftStart returns the minutes from the shift start to clock (0 at the start).
func minutesFromShiftStart(shiftStart, clock string) int {
	if clock == shiftStart {
		return 0
	}
	return dto.MinutesBetween(shiftStart, clock)
}

// shiftDetails returns shift id of ownerID with its breaks and minutes.
func shiftDetails(ctx context.Context, shiftRepository *repository.ShiftRepository, breakRepository *repository.BreakRepository,
	ownerID *int64, id int64) (model.Shift, error) {
	shift, err := shiftRepository.FindByID(ctx, ownerID, id)
	if err != nil {
		return model.Shift{}, toShiftResponseError(err)
	}
	shifts, err := shiftsWithBreaks(ctx, breakRepository, []model.Shift{shift})
	if err != nil {
		return model.Shift{}, err
	}
	return shifts[0], nil
}

// shiftsWithBreaks fills the breaks and minutes of shifts (one query for all their breaks).
func shiftsWithBreaks(ctx context.Context, breakRepository *repository.BreakRepository, shifts []model.Shift) ([]model.Shift, error) {
	ids := make([]int64, len(shifts))
	for index, shift := range shifts {
		ids[index] = shift.ID
	}
	breaks, err := breakRepository.FindByShifts(ctx, ids)
	if err != nil {
		return nil, err
	}
	byShift := make(map[int64][]model.ShiftBreak, len(shifts))
	for _, shiftBreak := range breaks {
		byShift[shiftBreak.ShiftID] = append(byShift[shiftBreak.ShiftID], shiftBreak)
	}

	for index := range shifts {
		shift := &shifts[index]
		shift.Breaks = byShift[shift.ID]
		if shift.Breaks == nil {
			shift.Breaks = []model.ShiftBreak{}
		}
		slices.SortStableFunc(shift.Breaks, func(a, b model.ShiftBreak) int {
			return minutesFromShiftStart(shift.StartTime, a.StartTime) - minutesFromShiftStart(shift.StartTime, b.StartTime)
		})
		shift.DurationMinutes = dto.MinutesBetween(shift.StartTime, shift.EndTime)
		shift.BreakMinutes = 0
		for _, shiftBreak := range shift.Breaks {
			shift.BreakMinutes += dto.MinutesBetween(shiftBreak.StartTime, shiftBreak.EndTime)
		}
	}
	return shifts, nil
}

// shiftFromFields copies the form fields into a shift (the breaks are linked by id).
func shiftFromFields(fields dto.ShiftFields) model.Shift {
	return model.Shift{
		UserID:      fields.UserID,
		Name:        fields.Name,
		ShiftType:   fields.ShiftType,
		StartTime:   fields.StartTime,
		EndTime:     fields.EndTime,
		WorkingDays: fields.WorkingDays,
	}
}

// toShiftResponseError maps repository errors to shift responses. Other errors pass
// through and are sent as ServerError.
func toShiftResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Shift not found")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "This customer already has a shift with this name")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	default:
		return err
	}
}
