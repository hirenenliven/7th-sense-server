package service

import (
	"context"
	"errors"
	"math"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// shiftDateLayout is how a cycle's shift date is written (machine_cycles.shift_date).
const shiftDateLayout = "2006-01-02"

// MachineCycleService holds the rules of machine cycles (migration 030). The CNC collector
// (internal/monitor) opens and updates cycles; when one opens it is matched to the planned
// job, the part process of its program, the shift and the loading time before it. Users read
// the cycles of the machines they may access (MachineService decides which).
type MachineCycleService struct {
	cycleRepository *repository.MachineCycleRepository
	machineService  *MachineService
}

// NewMachineCycleService creates a MachineCycleService.
func NewMachineCycleService(cycleRepository *repository.MachineCycleRepository, machineService *MachineService) *MachineCycleService {
	return &MachineCycleService{cycleRepository: cycleRepository, machineService: machineService}
}

// GetAll returns the cycles of machine listRequest.MachineID that started in its date range,
// with their totals, when actor may access the machine (else "Machine not found").
func (cycleService *MachineCycleService) GetAll(ctx context.Context, actor model.AuthUser, listRequest dto.MachineCycleListRequest) (dto.MachineCyclesResponse, error) {
	machine, err := cycleService.machineService.GetByID(ctx, actor, nil, listRequest.MachineID)
	if err != nil {
		return dto.MachineCyclesResponse{}, err
	}
	from, to := listRequest.Range(time.Now())
	cycles, err := cycleService.cycleRepository.FindAll(ctx, machine.ID, from, to)
	if err != nil {
		return dto.MachineCyclesResponse{}, err
	}
	return dto.MachineCyclesResponse{
		MachineID: machine.ID, MachineName: machine.Name, From: from, To: to,
		Cycles: cycles, Summary: summarizeCycles(cycles),
	}, nil
}

// GetByID returns cycle id when actor may access its machine, else "Machine cycle not found".
func (cycleService *MachineCycleService) GetByID(ctx context.Context, actor model.AuthUser, id int64) (model.MachineCycle, error) {
	cycle, err := cycleService.cycleRepository.FindByID(ctx, id)
	if err != nil {
		return model.MachineCycle{}, toMachineCycleResponseError(err)
	}
	// Another company's cycle answers as missing, like its machine
	if _, err := cycleService.machineService.GetByID(ctx, actor, nil, cycle.MachineID); err != nil {
		if errors.Is(err, response.MachineNotFound) {
			return model.MachineCycle{}, response.MachineCycleNotFound
		}
		return model.MachineCycle{}, err
	}
	return cycle, nil
}

// summarizeCycles adds up cycles; the average cycle counts completed cycles only.
func summarizeCycles(cycles []model.MachineCycle) dto.MachineCycleSummary {
	summary := dto.MachineCycleSummary{CycleCount: len(cycles)}
	completedSec := 0
	for _, cycle := range cycles {
		switch cycle.Status {
		case model.CycleCompleted:
			summary.CompletedCount++
			completedSec += cycle.DurationSec
		case model.CycleUnfinished:
			summary.UnfinishedCount++
		default:
			summary.RunningCount++
		}
		summary.PartsMade += cycle.PartsMade
		summary.ProductionSec += cycle.ProductionSec
		summary.HoldSec += cycle.HoldSec
		summary.DowntimeSec += cycle.DowntimeSec
		summary.OfflineSec += cycle.OfflineSec
	}
	if summary.CompletedCount > 0 {
		summary.AverageCycleSec = completedSec / summary.CompletedCount
	}
	return summary
}

// GetOpenCycle returns machine machineID's running cycle, nil when it has none. Used by the
// collector when a worker starts, to go on with a cycle a previous run left open.
func (cycleService *MachineCycleService) GetOpenCycle(ctx context.Context, machineID int64) (*model.MachineCycle, error) {
	cycle, err := cycleService.cycleRepository.FindOpenByMachineID(ctx, machineID)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &cycle, nil
}

// Create opens cycle (MachineID, ProgramName, StartedAt, PartCountStart set by the collector)
// as a running cycle: it is matched to its job and shift, and the loading time since the
// machine's previous cycle is worked out. It returns the cycle with its id and job.
func (cycleService *MachineCycleService) Create(ctx context.Context, cycle model.MachineCycle) (model.MachineCycle, error) {
	job, err := cycleService.matchJob(ctx, cycle)
	if err != nil {
		return model.MachineCycle{}, err
	}
	cycle.ScheduleItemID, cycle.PartID, cycle.ProcessID = job.ScheduleItemID, job.PartID, job.ProcessID
	cycle.ShiftID, cycle.ShiftDate, cycle.StandardCycleSec = job.ShiftID, job.ShiftDate, job.StandardCycleSec

	cycle.LoadingSec, err = cycleService.loadingSeconds(ctx, cycle)
	if err != nil {
		return model.MachineCycle{}, err
	}
	cycle.ID, err = cycleService.cycleRepository.Create(ctx, cycle)
	if err != nil {
		return model.MachineCycle{}, toMachineCycleResponseError(err)
	}
	cycle.Status, cycle.PartCountEnd = model.CycleRunning, cycle.PartCountStart
	return cycle, nil
}

// Update saves what the collector counted in cycle (seconds per mode, part counter) and, once
// it has closed, its status, end reason and end time.
func (cycleService *MachineCycleService) Update(ctx context.Context, cycle model.MachineCycle) error {
	return toMachineCycleResponseError(cycleService.cycleRepository.Update(ctx, cycle))
}

// matchJob finds what cycle belongs to: the job planned on the machine when it started; without
// one, the part process whose program it runs when exactly one matches, and the shift running
// then. The shift date is the day the shift started, or the calendar day without a shift.
func (cycleService *MachineCycleService) matchJob(ctx context.Context, cycle model.MachineCycle) (model.CycleJob, error) {
	job, err := cycleService.cycleRepository.FindScheduledJob(ctx, cycle.MachineID, cycle.StartedAt, cycle.ProgramName)
	if err == nil {
		return job, nil
	}
	if !errors.Is(err, repository.ErrNotFound) {
		return model.CycleJob{}, err
	}

	job = model.CycleJob{}
	if cycle.ProgramName != "" {
		processes, err := cycleService.cycleRepository.FindProcessesByProgram(ctx, cycle.MachineID, cycle.ProgramName)
		if err != nil {
			return model.CycleJob{}, err
		}
		// Two parts with the same program on this machine: which one ran cannot be told
		if len(processes) == 1 {
			job = processes[0]
		}
	}

	shift, err := cycleService.cycleRepository.FindShift(ctx, cycle.MachineID, cycle.StartedAt)
	switch {
	case err == nil:
		job.ShiftID, job.ShiftDate = shift.ShiftID, shift.ShiftDate
	case errors.Is(err, repository.ErrNotFound):
		job.ShiftDate = cycle.StartedAt.In(time.Local).Format(shiftDateLayout)
	default:
		return model.CycleJob{}, err
	}
	return job, nil
}

// loadingSeconds returns the loading / unloading time before cycle: from the end of the
// machine's previous cycle to its start. nil when there is no previous cycle, when the gap
// crosses a shift (or a day without shifts), or when the machine was Offline in it: then the
// gap is not time spent loading.
func (cycleService *MachineCycleService) loadingSeconds(ctx context.Context, cycle model.MachineCycle) (*int, error) {
	previous, err := cycleService.cycleRepository.FindLatestClosed(ctx, cycle.MachineID, cycle.StartedAt)
	if errors.Is(err, repository.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if previous.EndedAt == nil || previous.ShiftDate != cycle.ShiftDate || !sameID(previous.ShiftID, cycle.ShiftID) {
		return nil, nil
	}
	offline, err := cycleService.cycleRepository.HasModeBetween(ctx, cycle.MachineID, model.ModeOffline, *previous.EndedAt, cycle.StartedAt)
	if err != nil || offline {
		return nil, err
	}
	seconds := int(math.Round(cycle.StartedAt.Sub(*previous.EndedAt).Seconds()))
	return &seconds, nil
}

// sameID reports whether two optional ids are both nil or equal.
func sameID(first, second *int64) bool {
	if first == nil || second == nil {
		return first == second
	}
	return *first == *second
}

// toMachineCycleResponseError maps repository errors to response errors.
func toMachineCycleResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.MachineCycleNotFound
	default:
		return err
	}
}
