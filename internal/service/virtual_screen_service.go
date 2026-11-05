package service

import (
	"context"
	"errors"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
)

// GetVirtualScreen returns what machine machineID is doing now (the Virtual Screen): its mode and
// program (open machine_mode_log period), latest CNC values, last alarm, running cycle and the
// job it works on (the running cycle's, else the one scheduled now). "Machine not found" when
// actor may not read it.
func (dashboardService *DashboardService) GetVirtualScreen(ctx context.Context, actor model.AuthUser, machineID int64) (dto.VirtualScreenResponse, error) {
	now := time.Now()
	machines, err := dashboardService.machinesFor(ctx, actor, []int64{machineID})
	if err != nil {
		return dto.VirtualScreenResponse{}, err
	}
	machine := machines[0]
	telemetryRepository := dashboardService.telemetryService.telemetryRepository
	cycleRepository := dashboardService.cycleService.cycleRepository

	screen := dto.VirtualScreenResponse{MachineID: machine.ID, ServerTime: now, Online: machine.OnlineStatus}
	periods, err := telemetryRepository.FindModePeriods(ctx, machine.ID, now)
	if err != nil {
		return dto.VirtualScreenResponse{}, err
	}
	if len(periods) > 0 {
		current := periods[len(periods)-1]
		screen.ModeID, screen.ModeName, screen.ProgramName = &current.ModeID, current.ModeName, current.ProgramName
		screen.ModeSince = &current.StartedAt
	} else if machine.ModeID != nil {
		modeID := int(*machine.ModeID)
		screen.ModeID = &modeID
	}
	if screen.OpMode, err = dashboardService.telemetryService.GetOpModeName(ctx, machine.ID); err != nil {
		return dto.VirtualScreenResponse{}, err
	}

	variables, err := telemetryRepository.FindLiveVariables(ctx, machine.ID, dto.VirtualScreenVariableCodes)
	if err != nil {
		return dto.VirtualScreenResponse{}, err
	}
	screen.Live = dto.NewVirtualScreenLive(variables)
	if screen.Alarm, err = telemetryRepository.FindLatestAlarm(ctx, machine.ID); err != nil {
		return dto.VirtualScreenResponse{}, err
	}

	var scheduleItemID *int64
	cycle, err := cycleRepository.FindOpenByMachineID(ctx, machine.ID)
	switch {
	case err == nil:
		screen.Cycle = &dto.VirtualScreenCycle{ProgramName: cycle.ProgramName, StartedAt: cycle.StartedAt,
			ElapsedSec: cycle.DurationSec, StandardCycleSec: cycle.StandardCycleSec}
		scheduleItemID = cycle.ScheduleItemID
		if screen.ProgramName == "" {
			screen.ProgramName = cycle.ProgramName
		}
	case !errors.Is(err, repository.ErrNotFound):
		return dto.VirtualScreenResponse{}, err
	}
	if scheduleItemID == nil {
		scheduled, err := cycleRepository.FindScheduledJob(ctx, machine.ID, now, screen.ProgramName)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return dto.VirtualScreenResponse{}, err
		}
		scheduleItemID = scheduled.ScheduleItemID
	}
	if scheduleItemID != nil {
		job, err := cycleRepository.FindJobProgress(ctx, *scheduleItemID)
		if err != nil && !errors.Is(err, repository.ErrNotFound) {
			return dto.VirtualScreenResponse{}, err
		}
		if err == nil {
			screen.Job = &job
		}
	}
	return screen, nil
}
