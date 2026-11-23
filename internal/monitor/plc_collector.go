package monitor

import (
	"context"
	"errors"
	"log/slog"
	"strconv"
	"strings"
	"sync"
	"time"

	"7th-sense-server/internal/cert"
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/live"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/opcua"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// The PLC collector reads every machine (no sharding): shard 0 of 1.
const (
	plcShard      = 0
	plcShardCount = 1
)

// PlcCollector reads the PLC variables of every active Siemens machine over OPC UA (ADR 0021):
// one plcWorker per machine subscribes to the machine's PLC node ids and stores what changes,
// like the CNC collector stores the DLL's snapshot variables. It runs beside the CNC collector,
// in any build (gopcua is pure Go), and follows machine and variable edits through the same
// database notifications.
type PlcCollector struct {
	config      config.CNCConfig
	service     *service.MachineTelemetryService
	hub         *live.Hub          // nil = live data off
	queue       *TelemetryQueue    // takes the changes' rows to RabbitMQ; nil = direct writes
	board       *opcua.StatusBoard // what each worker is doing, for the diagnostics
	certificate cert.Paths

	mu        sync.Mutex
	workers   map[int64]*plcWorker
	workerCtx context.Context
	running   sync.WaitGroup // every worker goroutine
}

// NewPlcCollector creates a PlcCollector. Its workers store their values through telemetryQueue
// (nil = direct writes); certificate is presented to controllers whose machine uses Basic256Sha256.
func NewPlcCollector(cncConfig config.CNCConfig, telemetryService *service.MachineTelemetryService, hub *live.Hub,
	telemetryQueue *TelemetryQueue, board *opcua.StatusBoard, certificate cert.Paths) *PlcCollector {
	return &PlcCollector{config: cncConfig, service: telemetryService, hub: hub, queue: telemetryQueue, board: board,
		certificate: certificate, workers: map[int64]*plcWorker{}}
}

// RunCollection reads the machines until ctx is cancelled, then closes every connection.
func (collector *PlcCollector) RunCollection(ctx context.Context) {
	workerCtx, stopWorkers := context.WithCancel(ctx)
	collector.mu.Lock()
	collector.workerCtx = workerCtx
	collector.mu.Unlock()

	plc := collector.config.PLC
	collector.log(slog.LevelInfo, "plc collector started", "publishIntervalMs", plc.PublishInterval.Milliseconds(),
		"samplingIntervalMs", plc.SamplingInterval.Milliseconds(), "healthSec", plc.HealthInterval.Seconds(),
		"certificate", collector.certificate.CertificateFile)

	syncLoopDone := make(chan struct{})
	go func() {
		defer close(syncLoopDone)
		collector.runDatabaseSynchronizationLoop(workerCtx)
	}()

	<-ctx.Done()
	stopWorkers()
	<-syncLoopDone
	collector.waitForWorkersToStop()
	collector.log(slog.LevelInfo, "plc collector stopped")
}

// runDatabaseSynchronizationLoop keeps the workers in step with the database, as the CNC
// collector does: LISTEN, a full reload each time listening starts, reconnect with backoff.
func (collector *PlcCollector) runDatabaseSynchronizationLoop(ctx context.Context) {
	retry := backoff{max: time.Minute}
	for ctx.Err() == nil {
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
		collector.log(slog.LevelError, "plc collector: database listener lost, reconnecting", "retryIn", wait.String(), "error", err)
		sleepUnlessCancelled(ctx, wait)
	}
}

// isRead reports whether the PLC collector reads machine: in CNC_MACHINE_IDS and a controller
// with an OPC UA server.
func (collector *PlcCollector) isRead(machine model.MonitoredMachine) bool {
	return collector.config.Machines.Includes(machine.ID) && model.ReadsPLCOverOPCUA(machine.CNCEngine)
}

// syncWorkersWithActiveMachines starts or updates a worker per active machine read over OPC UA
// and stops the others.
func (collector *PlcCollector) syncWorkersWithActiveMachines(ctx context.Context) {
	machines, err := collector.service.GetActiveMachines(ctx, plcShard, plcShardCount)
	if err != nil {
		collector.log(slog.LevelError, "plc collector: list active machines", "error", err)
		return
	}
	readMachineIDs := map[int64]bool{}
	machineIDs := []int64{}
	for _, machine := range machines {
		if !collector.isRead(machine) {
			continue
		}
		readMachineIDs[machine.ID] = true
		machineIDs = append(machineIDs, machine.ID)
		collector.startOrUpdateWorker(machine)
	}
	collector.mu.Lock()
	var removed []int64
	for machineID := range collector.workers {
		if !readMachineIDs[machineID] {
			removed = append(removed, machineID)
		}
	}
	collector.mu.Unlock()
	for _, machineID := range removed {
		collector.stopWorker(machineID, true)
	}
	collector.log(slog.LevelInfo, "plc collector: machines loaded", "machines", len(machineIDs), "machineIds", machineIDs)
}

// handleChangeNotification follows machines_changed (start, update or stop that machine's worker)
// and cnc_config_changed (reload the variables of one worker, "machine:<id>", or of every worker).
func (collector *PlcCollector) handleChangeNotification(ctx context.Context, channel, payload string) {
	switch channel {
	case service.MachinesChangedChannel:
		machineID, err := strconv.ParseInt(payload, 10, 64)
		if err != nil {
			return
		}
		machine, err := collector.service.GetActiveMachine(ctx, machineID)
		if errors.Is(err, response.MachineNotFound) {
			collector.stopWorker(machineID, true)
			return
		}
		if err != nil {
			collector.log(slog.LevelError, "plc collector: reload machine", "machineId", machineID, "error", err)
			return
		}
		if !collector.isRead(machine) {
			collector.stopWorker(machineID, true)
			return
		}
		collector.startOrUpdateWorker(machine)

	case service.CNCConfigChangedChannel:
		collector.mu.Lock()
		defer collector.mu.Unlock()
		if machineText, ok := strings.CutPrefix(payload, "machine:"); ok {
			if machineID, err := strconv.ParseInt(machineText, 10, 64); err == nil && collector.workers[machineID] != nil {
				collector.workers[machineID].requestReload()
			}
			return
		}
		for _, plcWorker := range collector.workers {
			plcWorker.requestReload()
		}
	}
}

// startOrUpdateWorker starts a worker for machine; a running one reconnects when the way to reach
// the controller changed, else takes the new details (name, controller version) and resubscribes.
func (collector *PlcCollector) startOrUpdateWorker(machine model.MonitoredMachine) {
	collector.mu.Lock()
	existing := collector.workers[machine.ID]
	collector.mu.Unlock()
	if existing != nil {
		applied := existing.appliedMachine
		if samePLCConnection(applied, machine) {
			if applied.Name != machine.Name || applied.ControllerVersionType != machine.ControllerVersionType {
				existing.appliedMachine = machine
				existing.queueMachineUpdate(machine)
			}
			return
		}
		collector.stopWorker(machine.ID, false)
	}

	collector.mu.Lock()
	defer collector.mu.Unlock()
	if collector.workerCtx == nil || collector.workerCtx.Err() != nil {
		return
	}
	workerCtx, cancel := context.WithCancel(collector.workerCtx)
	plcWorker := newPlcWorker(collector, machine, cancel)
	collector.workers[machine.ID] = plcWorker
	collector.running.Add(1)
	go func() {
		defer collector.running.Done()
		plcWorker.run(workerCtx)
	}()
}

// samePLCConnection reports whether a and b reach the controller the same way.
func samePLCConnection(a, b model.MonitoredMachine) bool {
	return a.LiveIP == b.LiveIP && a.Port == b.Port && a.Username == b.Username && a.Password == b.Password &&
		a.SecurityMode == b.SecurityMode
}

// stopWorker cancels the machine's worker without waiting (it closes its connection on its own
// goroutine). forget drops its diagnostics status: the machine is no longer read (not a restart).
func (collector *PlcCollector) stopWorker(machineID int64, forget bool) {
	collector.mu.Lock()
	plcWorker := collector.workers[machineID]
	delete(collector.workers, machineID)
	collector.mu.Unlock()
	if plcWorker != nil {
		plcWorker.cancel()
	}
	if forget {
		collector.board.Remove(machineID)
		collector.hub.SetMachineCodes(machineID, live.ReaderPLC, nil)
	}
}

// waitForWorkersToStop waits for every worker to close its connection, at most CNC_STOP_TIMEOUT_SEC.
func (collector *PlcCollector) waitForWorkersToStop() {
	done := make(chan struct{})
	go func() {
		collector.running.Wait()
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(collector.config.StopTimeout):
		collector.log(slog.LevelError, "plc collector: workers did not stop in time")
	}
}

// log writes a PLC collector line about no single machine (CNC_LOG_PLC).
func (collector *PlcCollector) log(level slog.Level, message string, attributes ...any) {
	writeLog(collector.config.Logs, logTypePLC, level, nil, message, attributes...)
}
