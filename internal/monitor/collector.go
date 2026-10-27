// Package monitor collects CNC machine data (docs/cnc-telemetry-collector-plan.md). A
// Collector runs inside the API server (internal/collector) for one shard of the monitored machines: one
// worker per machine reads the DLL and stores every reading right away, and a listener
// follows machine and variable changes in the database. It never runs SQL itself: it calls
// service.MachineTelemetryService.
package monitor

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// Collector collects the machines of one shard (machines.id % shards = shard).
type Collector struct {
	config         config.CNCConfig
	service        *service.MachineTelemetryService
	driver         *cnc.Driver
	shard          int
	shards         int
	connectLimiter chan struct{} // semaphore: at most CNC_MAX_CONNECTING CNCConnect calls at once

	mu        sync.Mutex        // guards workers and workerCtx
	workers   map[int64]*worker // running workers by machine id
	workerCtx context.Context   // parent of every worker's context; cancelled on stop
}

// NewCollector creates a Collector for shard of shards.
func NewCollector(cncConfig config.CNCConfig, telemetryService *service.MachineTelemetryService, driver *cnc.Driver, shard, shards int) *Collector {
	return &Collector{
		config:         cncConfig,
		service:        telemetryService,
		driver:         driver,
		shard:          shard,
		shards:         shards,
		connectLimiter: make(chan struct{}, cncConfig.MaxConnecting),
		workers:        map[int64]*worker{},
	}
}

// RunCollection collects the shard's machines and blocks until ctx is cancelled, then stops
// cleanly: every machine is disconnected and every machine of the shard is left offline.
// It returns an error only when the collector cannot start (the supervisor then retries).
func (collector *Collector) RunCollection(ctx context.Context) error {
	logger := slog.With("shard", collector.shard, "shards", collector.shards)

	// DATA_COLLECTION_STEP_5: Prepare the shard: check the hypertable and start every machine
	// offline, closing mode periods a previous run left open.
	hypertable, err := collector.service.IsHypertable(ctx)
	if err != nil {
		return fmt.Errorf("check machine_telementry: %w", err)
	}
	if !hypertable {
		logger.Error("machine_telementry is not a TimescaleDB hypertable: fine for development, too slow for production")
	}
	if err := collector.service.InitializeShardForCollection(ctx, collector.shard, collector.shards, time.Now()); err != nil {
		return fmt.Errorf("initialize shard: %w", err)
	}

	workerCtx, stopWorkers := context.WithCancel(ctx)
	collector.mu.Lock()
	collector.workerCtx = workerCtx
	collector.mu.Unlock()

	// DATA_COLLECTION_STEP_6: Record collector_start and start the loop that follows the
	// database (machine list and change notifications).
	collector.recordCollectorEvent(model.EventCollectorStart,
		fmt.Sprintf("shard %d of %d, %s", collector.shard, collector.shards, collector.driver.Name()))
	logger.Info("cnc collector started", "driver", collector.driver.Name())
	logDebugMachines(collector.config.DebugMachineIDs)

	syncLoopDone := make(chan struct{})
	go func() {
		defer close(syncLoopDone)
		collector.runDatabaseSynchronizationLoop(workerCtx)
	}()

	<-ctx.Done()
	// DATA_COLLECTION_STEP_20: Stop the sync loop and every worker (each disconnects and writes
	// its Offline period), then mark the shard offline and record collector_stop.
	logger.Info("cnc collector stopping")
	stopWorkers()
	<-syncLoopDone
	collector.waitForWorkersToStop()

	stopCtx, cancel := context.WithTimeout(context.Background(), collector.config.StopTimeout)
	defer cancel()
	if err := collector.service.FinalizeShardCollection(stopCtx, collector.shard, collector.shards, time.Now()); err != nil {
		logger.Error("cnc collector: mark machines offline", "error", err)
	}
	collector.recordCollectorEvent(model.EventCollectorStop, fmt.Sprintf("shard %d of %d", collector.shard, collector.shards))
	logger.Info("cnc collector stopped")
	return nil
}

// runDatabaseSynchronizationLoop keeps the workers in step with the database: it LISTENs for
// machines_changed / cnc_config_changed and, each time listening starts, reloads every active
// machine. A lost LISTEN connection is reopened with backoff (1 s doubling to 1 min, jittered).
// Runs on its own goroutine, the only one that starts or updates workers.
func (collector *Collector) runDatabaseSynchronizationLoop(ctx context.Context) {
	retry := backoff{max: time.Minute}
	for ctx.Err() == nil {
		// DATA_COLLECTION_STEP_7: LISTEN for machine and variable changes. The full reload runs
		// after LISTEN is active, so no change made in between is missed (also after a reconnect).
		err := collector.service.ListenForConfigurationChanges(ctx, func() {
			retry.reset()
			collector.syncWorkersWithActiveMachines(ctx)
		}, func(channel, payload string) {
			collector.handleChangeNotification(ctx, channel, payload)
		})
		if ctx.Err() != nil {
			return
		}
		wait := retry.nextWait()
		slog.Error("cnc collector: database listener lost, reconnecting", "shard", collector.shard, "retryIn", wait.String(), "error", err)
		sleepUnlessCancelled(ctx, wait)
	}
}

// syncWorkersWithActiveMachines loads the shard's active machines and makes the running workers
// match them: it starts or updates a worker per active machine and stops the others.
func (collector *Collector) syncWorkersWithActiveMachines(ctx context.Context) {
	// DATA_COLLECTION_STEP_8: Load all active machines of the shard from the database.
	machines, err := collector.service.GetActiveMachines(ctx, collector.shard, collector.shards)
	if err != nil {
		slog.Error("cnc collector: list active machines", "shard", collector.shard, "error", err)
		return
	}
	activeMachineIDs := make(map[int64]bool, len(machines))
	machineIDs := make([]int64, 0, len(machines))
	for _, machine := range machines {
		activeMachineIDs[machine.ID] = true
		machineIDs = append(machineIDs, machine.ID)
		// DATA_COLLECTION_STEP_9: Start (or update) one worker per active machine.
		collector.startOrUpdateWorker(machine)
	}
	collector.mu.Lock()
	var removedMachineIDs []int64
	for machineID := range collector.workers {
		if !activeMachineIDs[machineID] {
			removedMachineIDs = append(removedMachineIDs, machineID)
		}
	}
	collector.mu.Unlock()
	for _, machineID := range removedMachineIDs {
		collector.stopWorker(machineID)
	}
	slog.Info("cnc collector: machines loaded", "shard", collector.shard, "machines", len(machines), "machineIds", machineIDs)
}

// handleChangeNotification follows one database notification: machines_changed (payload: the
// machine id) starts, updates or stops that machine's worker; cnc_config_changed asks one worker
// ("machine:<id>") or every worker to reload its variables.
func (collector *Collector) handleChangeNotification(ctx context.Context, channel, payload string) {
	// DATA_COLLECTION_STEP_18: Follow machine and variable edits without polling.
	switch channel {
	case service.MachinesChangedChannel:
		machineID, err := strconv.ParseInt(payload, 10, 64)
		if err != nil || machineID%int64(collector.shards) != int64(collector.shard) {
			return
		}
		machine, err := collector.service.GetActiveMachine(ctx, machineID)
		slog.Info("cnc collector: processing machine change notification", "machineId", machineID, "shard", collector.shard, "shards", collector.shards)
		if errors.Is(err, response.MachineNotFound) {
			collector.stopWorker(machineID)
			return
		}
		if err != nil {
			slog.Error("cnc collector: reload machine", "machineId", machineID, "error", err)
			return
		}
		collector.startOrUpdateWorker(machine)

	case service.CNCConfigChangedChannel:
		slog.Info("cnc collector: variable config changed, reloading variables", "payload", payload)
		collector.mu.Lock()
		defer collector.mu.Unlock()
		if machineText, ok := strings.CutPrefix(payload, "machine:"); ok {
			if machineID, err := strconv.ParseInt(machineText, 10, 64); err == nil {
				if machineWorker := collector.workers[machineID]; machineWorker != nil {
					machineWorker.requestVariableReload()
				}
			}
			return
		}
		for _, machineWorker := range collector.workers {
			machineWorker.requestVariableReload()
		}
	}
}

// startOrUpdateWorker starts a worker for machine, restarts it when the way to reach the
// machine changed (IP, port, engine, login), or hands it the new name / controller version and
// reloads its variables. It runs only on the sync loop goroutine, which owns worker.appliedMachine.
func (collector *Collector) startOrUpdateWorker(machine model.MonitoredMachine) {
	collector.mu.Lock()
	existing := collector.workers[machine.ID]
	collector.mu.Unlock()

	if existing != nil {
		if isMachineUnchanged(existing.appliedMachine, machine) {
			return
		}
		if existing.appliedMachine.SameConnection(machine) {
			existing.queueMachineUpdate(machine)
			return
		}
		slog.Info("cnc collector: machine connection changed, reconnecting", "machineId", machine.ID)
		collector.stopWorker(machine.ID)
		select { // the old worker writes its Offline period before the new one starts
		case <-existing.done:
		case <-time.After(collector.config.StopTimeout):
			slog.Warn("cnc collector: old worker did not stop in time", "machineId", machine.ID)
		}
	}

	collector.mu.Lock()
	defer collector.mu.Unlock()
	// No new worker once the collector is stopping.
	if collector.workerCtx.Err() != nil || collector.workers[machine.ID] != nil {
		return
	}
	workerCtx, cancel := context.WithCancel(collector.workerCtx)
	machineWorker := newMachineWorker(collector, machine, cancel)
	collector.workers[machine.ID] = machineWorker
	slog.Info("cnc collector: worker started", "machineId", machine.ID, "machine", machine.Name,
		"ip", machine.LiveIP, "port", machine.Port)
	go machineWorker.runConnectionLoop(workerCtx)
}

// isMachineUnchanged reports whether nothing the collector uses changed.
func isMachineUnchanged(a, b model.MonitoredMachine) bool {
	return a.SameConnection(b) && a.Name == b.Name && a.ControllerVersionType == b.ControllerVersionType
}

// stopWorker removes the machine's worker and cancels it without waiting: the worker
// disconnects and records the machine offline on its own goroutine.
func (collector *Collector) stopWorker(machineID int64) {
	collector.mu.Lock()
	machineWorker := collector.workers[machineID]
	delete(collector.workers, machineID)
	collector.mu.Unlock()
	if machineWorker != nil {
		slog.Info("cnc collector: machine no longer collected", "machineId", machineID)
		machineWorker.cancel()
	}
}

// waitForWorkersToStop waits until every worker has disconnected, at most CNC_STOP_TIMEOUT_SEC
// for all of them together.
func (collector *Collector) waitForWorkersToStop() {
	collector.mu.Lock()
	workers := make([]*worker, 0, len(collector.workers))
	for _, machineWorker := range collector.workers {
		workers = append(workers, machineWorker)
	}
	collector.mu.Unlock()

	deadline := time.After(collector.config.StopTimeout)
	for _, machineWorker := range workers {
		select {
		case <-machineWorker.done:
		case <-deadline:
			slog.Warn("cnc collector: workers did not stop in time", "shard", collector.shard)
			return
		}
	}
}

// recordCollectorEvent stores a collector_start / collector_stop event in machine_events. Its
// own timeout lets it run after ctx is cancelled (at stop).
func (collector *Collector) recordCollectorEvent(eventType, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	event := model.MachineEvent{Shard: collector.shard, Type: eventType, Message: message, At: time.Now()}
	if err := collector.service.AddEvent(ctx, event); err != nil {
		slog.Error("cnc collector: store event", "event", eventType, "error", err)
	}
}
