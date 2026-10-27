package service

import (
	"context"
	"time"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
)

// HealthService serves the verification endpoints that show what the CNC collectors recorded.
type HealthService struct {
	machineService   *MachineService
	telemetryService *MachineTelemetryService
}

// NewHealthService creates a HealthService. machineService finds the machine.
func NewHealthService(machineService *MachineService, telemetryService *MachineTelemetryService) *HealthService {
	return &HealthService{machineService: machineService, telemetryService: telemetryService}
}

// GetMachineModeHistory returns machine machineID's mode periods of the last
// dto.MachineModeHistoryDays days and the time per mode. The route is public, so every
// machine that is not deleted can be read (an empty scope); others answer 404.
func (healthService *HealthService) GetMachineModeHistory(ctx context.Context, machineID int64) (dto.MachineModeHistoryResponse, error) {
	machine, err := healthService.machineService.find(ctx, model.MachineScope{}, machineID)
	if err != nil {
		return dto.MachineModeHistoryResponse{}, err
	}

	to := time.Now()
	from := to.AddDate(0, 0, -dto.MachineModeHistoryDays)
	periods, err := healthService.telemetryService.GetModePeriods(ctx, machineID, from)
	if err != nil {
		return dto.MachineModeHistoryResponse{}, err
	}

	return dto.MachineModeHistoryResponse{
		MachineID:     machine.ID,
		MachineName:   machine.Name,
		CurrentModeID: machine.ModeID,
		OnlineStatus:  machine.OnlineStatus,
		From:          from,
		To:            to,
		Periods:       periods,
		Totals:        modeTotals(periods),
	}, nil
}

// modeTotals sums the periods per mode, in the order each mode first appears. Periods are
// counted whole, so one that started before the window adds its full duration.
func modeTotals(periods []model.MachineModePeriod) []dto.MachineModeTotal {
	totals := []dto.MachineModeTotal{}
	indexByMode := map[int]int{}
	for _, period := range periods {
		index, found := indexByMode[period.ModeID]
		if !found {
			index = len(totals)
			indexByMode[period.ModeID] = index
			totals = append(totals, dto.MachineModeTotal{ModeID: period.ModeID, ModeName: period.ModeName})
		}
		totals[index].PeriodCount++
		totals[index].DurationSec += period.DurationSec
	}
	return totals
}
