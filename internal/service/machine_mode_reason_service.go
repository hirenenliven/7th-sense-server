package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// reasonActiveStatus is the status of a master reason that can be picked (both master lists).
const reasonActiveStatus = dto.DowntimeReasonStatusActive

// MachineModeReasonService holds the rules for recording why a machine was down or idle. Whoever
// may read a machine records reasons for it (operators do; MachineService decides who may). Only
// finished Downtime and Idle periods take reasons, each for a part of the period that overlaps no
// other reason of it, picked from the active master list of the machine's customer.
type MachineModeReasonService struct {
	dbPool                   *pgxpool.Pool
	reasonRepository         *repository.MachineModeReasonRepository
	machineService           *MachineService
	downtimeReasonRepository *repository.DowntimeReasonRepository
	idleReasonRepository     *repository.IdleReasonRepository
}

// NewMachineModeReasonService creates a MachineModeReasonService.
func NewMachineModeReasonService(dbPool *pgxpool.Pool, machineService *MachineService) *MachineModeReasonService {
	return &MachineModeReasonService{
		dbPool:                   dbPool,
		reasonRepository:         repository.NewMachineModeReasonRepository(dbPool),
		machineService:           machineService,
		downtimeReasonRepository: repository.NewDowntimeReasonRepository(dbPool),
		idleReasonRepository:     repository.NewIdleReasonRepository(dbPool),
	}
}

// GetByPeriod returns the whole period with its reasons. A period that takes no reason (another
// mode, or still running) is refused, so the popup does not open on it.
func (reasonService *MachineModeReasonService) GetByPeriod(ctx context.Context, actor model.AuthUser, periodID int64) (dto.MachineModeReasonsResponse, error) {
	period, machine, err := reasonService.readablePeriod(ctx, actor, periodID)
	if err != nil {
		return dto.MachineModeReasonsResponse{}, err
	}
	reasonType, err := reasonTypeOfPeriod(period)
	if err != nil {
		return dto.MachineModeReasonsResponse{}, err
	}
	reasons, err := reasonService.reasonRepository.FindAllByPeriod(ctx, periodID)
	if err != nil {
		return dto.MachineModeReasonsResponse{}, err
	}
	return dto.MachineModeReasonsResponse{Period: period, ReasonType: reasonType, CustomerID: machine.UserID, Reasons: reasons}, nil
}

// Create records a reason for a part of the period.
func (reasonService *MachineModeReasonService) Create(ctx context.Context, actor model.AuthUser, reasonRequest dto.CreateMachineModeReasonRequest) (model.MachineModeReason, error) {
	fields := reasonRequest.MachineModeReasonFields
	period, machine, err := reasonService.readablePeriod(ctx, actor, reasonRequest.PeriodID)
	if err != nil {
		return model.MachineModeReason{}, err
	}
	reasonType, err := reasonTypeOfPeriod(period)
	if err != nil {
		return model.MachineModeReason{}, err
	}
	if err := checkInsidePeriod(period, fields); err != nil {
		return model.MachineModeReason{}, err
	}
	if err := reasonService.checkReason(ctx, reasonType, machine.UserID, fields.ReasonID); err != nil {
		return model.MachineModeReason{}, err
	}

	var id int64
	err = reasonService.inPeriodLock(ctx, period.ID, func(txRepository *repository.MachineModeReasonRepository) error {
		if err := checkNoOverlap(ctx, txRepository, period.ID, fields, 0); err != nil {
			return err
		}
		var err error
		id, err = txRepository.Create(ctx, period.ID, reasonType, fields.ReasonID, fields.StartedAt, fields.EndedAt)
		return err
	})
	if err != nil {
		return model.MachineModeReason{}, err
	}
	return reasonService.find(ctx, id)
}

// Update replaces the reason and the time of a recorded row. The reason it already has stays
// allowed when it has since been made inactive or deleted from its master list.
func (reasonService *MachineModeReasonService) Update(ctx context.Context, actor model.AuthUser, id int64, reasonRequest dto.UpdateMachineModeReasonRequest) (model.MachineModeReason, error) {
	fields := reasonRequest.MachineModeReasonFields
	saved, err := reasonService.find(ctx, id)
	if err != nil {
		return model.MachineModeReason{}, err
	}
	period, machine, err := reasonService.readablePeriod(ctx, actor, saved.PeriodID)
	if err != nil {
		return model.MachineModeReason{}, err
	}
	if err := checkInsidePeriod(period, fields); err != nil {
		return model.MachineModeReason{}, err
	}
	if fields.ReasonID != saved.ReasonID {
		if err := reasonService.checkReason(ctx, saved.ReasonType, machine.UserID, fields.ReasonID); err != nil {
			return model.MachineModeReason{}, err
		}
	}

	err = reasonService.inPeriodLock(ctx, period.ID, func(txRepository *repository.MachineModeReasonRepository) error {
		if err := checkNoOverlap(ctx, txRepository, period.ID, fields, id); err != nil {
			return err
		}
		return toMachineModeReasonResponseError(txRepository.Update(ctx, id, saved.ReasonType, fields.ReasonID, fields.StartedAt, fields.EndedAt))
	})
	if err != nil {
		return model.MachineModeReason{}, err
	}
	return reasonService.find(ctx, id)
}

// Delete soft-deletes a recorded row of a machine actor may read.
func (reasonService *MachineModeReasonService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	saved, err := reasonService.find(ctx, id)
	if err != nil {
		return err
	}
	if _, _, err := reasonService.readablePeriod(ctx, actor, saved.PeriodID); err != nil {
		return err
	}
	return toMachineModeReasonResponseError(reasonService.reasonRepository.Delete(ctx, id))
}

// readablePeriod returns the period and its machine when actor may read the machine, else
// "Machine not found" like every other read of a machine.
func (reasonService *MachineModeReasonService) readablePeriod(ctx context.Context, actor model.AuthUser, periodID int64) (model.ReasonPeriod, model.Machine, error) {
	period, err := reasonService.reasonRepository.FindPeriodByID(ctx, periodID)
	if err != nil {
		return model.ReasonPeriod{}, model.Machine{}, toMachineModeReasonResponseError(err)
	}
	machine, err := reasonService.machineService.GetByID(ctx, actor, nil, period.MachineID)
	if err != nil {
		return model.ReasonPeriod{}, model.Machine{}, err
	}
	return period, machine, nil
}

// inPeriodLock runs save in a transaction that holds the period's row.
func (reasonService *MachineModeReasonService) inPeriodLock(ctx context.Context, periodID int64, save func(*repository.MachineModeReasonRepository) error) error {
	return db.WithTransaction(ctx, reasonService.dbPool, func(tx pgx.Tx) error {
		txRepository := repository.NewMachineModeReasonRepository(tx)
		if err := txRepository.LockPeriod(ctx, periodID); err != nil {
			return toMachineModeReasonResponseError(err)
		}
		return save(txRepository)
	})
}

// checkReason returns "Please select valid … reason" unless reasonID is an active reason of the
// customer in the master list of reasonType.
func (reasonService *MachineModeReasonService) checkReason(ctx context.Context, reasonType string, customerID *int64, reasonID int64) error {
	invalid := response.NewError(http.StatusBadRequest, fmt.Sprintf("Please select valid %s reason", reasonType))
	// A nil owner would match any customer's reason in the repositories
	if customerID == nil {
		return invalid
	}
	var status string
	var err error
	if reasonType == model.ReasonTypeDowntime {
		var reason model.DowntimeReason
		reason, err = reasonService.downtimeReasonRepository.FindByID(ctx, customerID, reasonID)
		status = reason.Status
	} else {
		var reason model.IdleReason
		reason, err = reasonService.idleReasonRepository.FindByID(ctx, customerID, reasonID)
		status = reason.Status
	}
	if errors.Is(err, repository.ErrNotFound) {
		return invalid
	}
	if err != nil {
		return err
	}
	if status != reasonActiveStatus {
		return invalid
	}
	return nil
}

// find returns reason row id.
func (reasonService *MachineModeReasonService) find(ctx context.Context, id int64) (model.MachineModeReason, error) {
	reason, err := reasonService.reasonRepository.FindByID(ctx, id)
	return reason, toMachineModeReasonResponseError(err)
}

// reasonTypeOfPeriod returns which master list the period's mode uses, or the error the user sees
// for a period that takes no reason. The collector records only Idle 5 for a machine standing
// ready and Downtime 6 for an alarm, so these two modes are all there is to explain.
func reasonTypeOfPeriod(period model.ReasonPeriod) (string, error) {
	var reasonType string
	switch period.ModeID {
	case model.ModeDowntime:
		reasonType = model.ReasonTypeDowntime
	case model.ModeIdle:
		reasonType = model.ReasonTypeIdle
	default:
		return "", response.NewError(http.StatusBadRequest, fmt.Sprintf("You can not add a reason in %s mode.", period.ModeName))
	}
	// A running period has no end yet, so its time cannot be explained
	if period.EndedAt == nil {
		return "", response.NewError(http.StatusBadRequest, fmt.Sprintf("You can not add %s reason in running mode.", reasonType))
	}
	return reasonType, nil
}

// checkInsidePeriod returns an error unless the fields' time lies inside the period.
func checkInsidePeriod(period model.ReasonPeriod, fields dto.MachineModeReasonFields) error {
	if fields.StartedAt.Before(period.StartedAt) || period.EndedAt == nil || fields.EndedAt.After(*period.EndedAt) {
		return response.NewError(http.StatusBadRequest, "Reason time must be inside the period")
	}
	return nil
}

// checkNoOverlap returns an error when another reason of the period covers the fields' time.
func checkNoOverlap(ctx context.Context, reasonRepository *repository.MachineModeReasonRepository, periodID int64, fields dto.MachineModeReasonFields, exceptID int64) error {
	overlaps, err := reasonRepository.HasOverlap(ctx, periodID, fields.StartedAt, fields.EndedAt, exceptID)
	if err != nil {
		return err
	}
	if overlaps {
		return response.NewError(http.StatusBadRequest, "This time overlaps another reason of this period")
	}
	return nil
}

// toMachineModeReasonResponseError maps repository errors to responses. Other errors pass through
// and are sent as ServerError.
func toMachineModeReasonResponseError(err error) error {
	if errors.Is(err, repository.ErrNotFound) {
		return response.NewError(http.StatusNotFound, "Machine mode reason not found")
	}
	return err
}
