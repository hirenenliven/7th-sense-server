package service

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// Notification channels the CNC collectors listen to (migration 028 triggers).
const (
	MachinesChangedChannel  = "machines_changed"   // payload: machine id
	CNCConfigChangedChannel = "cnc_config_changed" // payload: "variables" | "machine:<id>"
)

// MachineTelemetryService holds the business rules of CNC data collection: which machines the
// collectors (internal/monitor) connect to (active, not deleted), which variables they read,
// and how what they collect is stored (one transaction per batch, mode periods in order).
type MachineTelemetryService struct {
	dbPool              *pgxpool.Pool
	telemetryRepository *repository.MachineTelemetryRepository
}

// NewMachineTelemetryService creates a MachineTelemetryService. dbPool runs the write
// transactions and the LISTEN connection.
func NewMachineTelemetryService(dbPool *pgxpool.Pool, telemetryRepository *repository.MachineTelemetryRepository) *MachineTelemetryService {
	return &MachineTelemetryService{dbPool: dbPool, telemetryRepository: telemetryRepository}
}

// GetActiveMachines returns the active machines of shard (machines.id % shards = shard).
func (telemetryService *MachineTelemetryService) GetActiveMachines(ctx context.Context, shard, shards int) ([]model.MonitoredMachine, error) {
	return telemetryService.telemetryRepository.FindActiveMachines(ctx, shard, shards)
}

// GetActiveMachine returns the machine when it is active, else response.MachineNotFound
// (deleted or inactive: the collector disconnects it).
func (telemetryService *MachineTelemetryService) GetActiveMachine(ctx context.Context, id int64) (model.MonitoredMachine, error) {
	machine, err := telemetryService.telemetryRepository.FindActiveMachineByID(ctx, id)
	return machine, toMachineTelemetryResponseError(err)
}

// GetVariableConfig returns the variables machine reads from the given sources.
func (telemetryService *MachineTelemetryService) GetVariableConfig(ctx context.Context, machine model.MonitoredMachine, sources []string) ([]model.VariableConfig, error) {
	if machine.CNCEngine == nil {
		return nil, nil
	}
	return telemetryService.telemetryRepository.FindVariableConfig(ctx, machine, sources)
}

// IsHypertable reports whether machine_telementry is a TimescaleDB hypertable.
func (telemetryService *MachineTelemetryService) IsHypertable(ctx context.Context) (bool, error) {
	return telemetryService.telemetryRepository.IsHypertable(ctx)
}

// UpsertLiveValues upserts machine_variable_values rows: one per machine-variable,
// with the latest value (value or value_text depending on data type).
func (telemetryService *MachineTelemetryService) UpsertLiveValues(ctx context.Context, values []model.LiveValue) error {
	return telemetryService.telemetryRepository.UpsertLiveValues(ctx, values)
}

// GetLiveValues returns the current live values for a machine's variables.
func (telemetryService *MachineTelemetryService) GetLiveValues(ctx context.Context, machineID int64, variableIDs []int) (map[int]model.LiveValue, error) {
	return telemetryService.telemetryRepository.GetLiveValues(ctx, machineID, variableIDs)
}

// StartShard prepares shard for a collector run: periods left open by a previous run end at
// the machine's last_seen_at (the time the collector was down counts as offline) and every
// machine of the shard starts offline until it connects.
func (telemetryService *MachineTelemetryService) StartShard(ctx context.Context, shard, shards int, at time.Time) error {
	return telemetryService.markShardOffline(ctx, shard, shards, at, true)
}

// StopShard ends a collector run: every machine of shard goes offline at at.
func (telemetryService *MachineTelemetryService) StopShard(ctx context.Context, shard, shards int, at time.Time) error {
	return telemetryService.markShardOffline(ctx, shard, shards, at, false)
}

func (telemetryService *MachineTelemetryService) markShardOffline(ctx context.Context, shard, shards int, at time.Time, starting bool) error {
	return db.WithTransaction(ctx, telemetryService.dbPool, func(tx pgx.Tx) error {
		return repository.NewMachineTelemetryRepository(tx).MarkShardOffline(ctx, shard, shards, at, starting)
	})
}

// InsertTelemetry stores the rows of one reading in machine_telementry.
func (telemetryService *MachineTelemetryService) InsertTelemetry(ctx context.Context, values []model.TelemetryValue) error {
	return telemetryService.telemetryRepository.InsertTelemetry(ctx, values)
}

// ChangeMode ends the machine's open mode period and opens one in change.ModeID, in one
// transaction (machine_mode_log and machines.mode_id stay in step).
func (telemetryService *MachineTelemetryService) ChangeMode(ctx context.Context, change model.ModeChange) error {
	return db.WithTransaction(ctx, telemetryService.dbPool, func(tx pgx.Tx) error {
		return repository.NewMachineTelemetryRepository(tx).ChangeMode(ctx, change)
	})
}

// SaveStatus writes the machine's machine_live_status row.
func (telemetryService *MachineTelemetryService) SaveStatus(ctx context.Context, status model.MachineLiveStatus) error {
	return telemetryService.telemetryRepository.UpsertLiveStatuses(ctx, []model.MachineLiveStatus{status})
}

// toMachineTelemetryResponseError maps repository errors to response errors.
func toMachineTelemetryResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.MachineNotFound
	default:
		return err
	}
}

// AddEvent stores one event.
func (telemetryService *MachineTelemetryService) AddEvent(ctx context.Context, event model.MachineEvent) error {
	return telemetryService.telemetryRepository.InsertEvents(ctx, []model.MachineEvent{event})
}

// Listen calls notify for every machines_changed / cnc_config_changed notification until ctx
// ends or the connection fails (then it returns the error; the caller reconnects and reloads).
// ready is called once listening has started.
func (telemetryService *MachineTelemetryService) Listen(ctx context.Context, ready func(), notify func(channel, payload string)) error {
	return db.Listen(ctx, telemetryService.dbPool,
		[]string{MachinesChangedChannel, CNCConfigChangedChannel}, ready, notify)
}
