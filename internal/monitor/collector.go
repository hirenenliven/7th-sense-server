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
	config     config.CNCConfig
	service    *service.MachineTelemetryService
	driver     *cnc.Driver
	shard      int
	shards     int
	connecting chan struct{} // CNC_MAX_CONNECTING

	mu        sync.Mutex
	workers   map[int64]*worker
	workerCtx context.Context
}

// NewCollector creates a Collector for shard of shards.
func NewCollector(cncConfig config.CNCConfig, telemetryService *service.MachineTelemetryService, driver *cnc.Driver, shard, shards int) *Collector {
	return &Collector{
		config:     cncConfig,
		service:    telemetryService,
		driver:     driver,
		shard:      shard,
		shards:     shards,
		connecting: make(chan struct{}, cncConfig.MaxConnecting),
		workers:    map[int64]*worker{},
	}
}

// Run collects until ctx is cancelled, then stops cleanly: every machine is disconnected and
// every machine of the shard is left offline. It returns an error
// only when the collector cannot start.
func (collector *Collector) Run(ctx context.Context) error {
	logger := slog.With("shard", collector.shard, "shards", collector.shards)

	hypertable, err := collector.service.IsHypertable(ctx)
	if err != nil {
		return fmt.Errorf("check machine_telementry: %w", err)
	}
	if !hypertable {
		logger.Error("machine_telementry is not a TimescaleDB hypertable: fine for development, too slow for production")
	}
	if err := collector.service.StartShard(ctx, collector.shard, collector.shards, time.Now()); err != nil {
		return fmt.Errorf("start shard: %w", err)
	}

	workerCtx, stopWorkers := context.WithCancel(ctx)
	collector.mu.Lock()
	collector.workerCtx = workerCtx
	collector.mu.Unlock()

	collector.addCollectorEvent(model.EventCollectorStart,
		fmt.Sprintf("shard %d of %d, %s", collector.shard, collector.shards, collector.driver.Name()))
	logger.Info("cnc collector started", "driver", collector.driver.Name())

	listenerDone := make(chan struct{})
	go func() {
		defer close(listenerDone)
		collector.listen(workerCtx)
	}()

	<-ctx.Done()
	logger.Info("cnc collector stopping")
	stopWorkers()
	<-listenerDone
	collector.waitWorkers()

	stopCtx, cancel := context.WithTimeout(context.Background(), collector.config.StopTimeout)
	defer cancel()
	if err := collector.service.StopShard(stopCtx, collector.shard, collector.shards, time.Now()); err != nil {
		logger.Error("cnc collector: mark machines offline", "error", err)
	}
	collector.addCollectorEvent(model.EventCollectorStop, fmt.Sprintf("shard %d of %d", collector.shard, collector.shards))
	logger.Info("cnc collector stopped")
	return nil
}

// listen follows machines_changed / cnc_config_changed. Each time listening starts (also after
// a lost connection) every machine is reloaded, so no change is missed.
func (collector *Collector) listen(ctx context.Context) {
	retry := backoff{max: time.Minute}
	for ctx.Err() == nil {
		err := collector.service.Listen(ctx, func() {
			retry.reset()
			collector.resync(ctx)
		}, func(channel, payload string) {
			collector.onNotification(ctx, channel, payload)
		})
		if ctx.Err() != nil {
			return
		}
		wait := retry.next()
		slog.Error("cnc collector: database listener lost, reconnecting", "shard", collector.shard, "retryIn", wait.String(), "error", err)
		sleep(ctx, wait)
	}
}

// resync makes the running workers match the active machines of the shard.
func (collector *Collector) resync(ctx context.Context) {
	machines, err := collector.service.GetActiveMachines(ctx, collector.shard, collector.shards)
	if err != nil {
		slog.Error("cnc collector: list active machines", "shard", collector.shard, "error", err)
		return
	}
	active := make(map[int64]bool, len(machines))
	for _, machine := range machines {
		active[machine.ID] = true
		collector.applyMachine(machine)
	}
	collector.mu.Lock()
	var gone []int64
	for machineID := range collector.workers {
		if !active[machineID] {
			gone = append(gone, machineID)
		}
	}
	collector.mu.Unlock()
	for _, machineID := range gone {
		collector.stopWorker(machineID)
	}
	slog.Info("cnc collector: machines loaded", "shard", collector.shard, "machines", len(machines))
}

func (collector *Collector) onNotification(ctx context.Context, channel, payload string) {
	switch channel {
	case service.MachinesChangedChannel:
		machineID, err := strconv.ParseInt(payload, 10, 64)
		if err != nil || machineID%int64(collector.shards) != int64(collector.shard) {
			return
		}
		machine, err := collector.service.GetActiveMachine(ctx, machineID)
		if errors.Is(err, response.MachineNotFound) {
			collector.stopWorker(machineID)
			return
		}
		if err != nil {
			slog.Error("cnc collector: reload machine", "machineId", machineID, "error", err)
			return
		}
		collector.applyMachine(machine)

	case service.CNCConfigChangedChannel:
		collector.mu.Lock()
		defer collector.mu.Unlock()
		if machineText, ok := strings.CutPrefix(payload, "machine:"); ok {
			if machineID, err := strconv.ParseInt(machineText, 10, 64); err == nil {
				if machineWorker := collector.workers[machineID]; machineWorker != nil {
					machineWorker.reloadConfig()
				}
			}
			return
		}
		for _, machineWorker := range collector.workers {
			machineWorker.reloadConfig()
		}
	}
}

// applyMachine starts a worker for machine, restarts it when the way to reach the machine
// changed, or reloads its variables (its controller version may have changed). It runs only on
// the listener goroutine (resync and notifications), which owns worker.applied.
func (collector *Collector) applyMachine(machine model.MonitoredMachine) {
	collector.mu.Lock()
	existing := collector.workers[machine.ID]
	collector.mu.Unlock()

	if existing != nil {
		if sameMachine(existing.applied, machine) {
			return
		}
		if existing.applied.SameConnection(machine) {
			existing.updateMachine(machine)
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
	if collector.workerCtx.Err() != nil || collector.workers[machine.ID] != nil {
		return
	}
	workerCtx, cancel := context.WithCancel(collector.workerCtx)
	machineWorker := newWorker(collector, machine, cancel)
	collector.workers[machine.ID] = machineWorker
	go machineWorker.run(workerCtx)
}

// sameMachine reports whether nothing the collector uses changed.
func sameMachine(a, b model.MonitoredMachine) bool {
	return a.SameConnection(b) && a.Name == b.Name && a.ControllerVersionType == b.ControllerVersionType
}

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

// waitWorkers waits until every worker has disconnected, at most CNC_STOP_TIMEOUT_SEC.
func (collector *Collector) waitWorkers() {
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

// addCollectorEvent stores a collector_start / collector_stop event directly.
func (collector *Collector) addCollectorEvent(eventType, message string) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	event := model.MachineEvent{Shard: collector.shard, Type: eventType, Message: message, At: time.Now()}
	if err := collector.service.AddEvent(ctx, event); err != nil {
		slog.Error("cnc collector: store event", "event", eventType, "error", err)
	}
}
