// Package collector runs the data collection inside the API server: starting the server
// connects the machines and starts storing their data; there is no other process or service.
// No DLL code and no SQL here.
//
// The CNC DLL is 32-bit, so the server must be built with GOARCH=386 to collect snapshot data (Air
// does this, see .air.toml, run.ps1 and build.ps1). The PLC reader (OPC UA, gopcua) is pure Go and
// runs in any build (PLC_ENABLED).
package collector

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/cert"
	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/db"
	"7th-sense-server/internal/live"
	"7th-sense-server/internal/monitor"
	"7th-sense-server/internal/opcua"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/service"
)

// One process collects every machine: shard 0 of 1.
const (
	collectorShard      = 0
	collectorShardCount = 1
)

// Supervisor runs the CNC collector and the PLC collector in the server process and restarts the
// CNC collector when it cannot start (for example the database is not ready yet). With
// RABBITMQ_ENABLED it also runs the telemetry queue both collectors store their values through.
type Supervisor struct {
	config         config.CNCConfig
	rabbitConfig   config.RabbitMQConfig
	dbConfig       config.DBConfig
	dbPool         *pgxpool.Pool
	writerPool     *pgxpool.Pool // the telemetry queue writers' own connections; nil = none opened
	cycleService   *service.MachineCycleService
	hub            *live.Hub
	telemetryQueue *monitor.TelemetryQueue // nil = direct writes (RABBITMQ_ENABLED=false)
	ctx            context.Context         // cancelled by Stop; parent of both collectors
	cancel         context.CancelFunc
	wg             sync.WaitGroup
}

// NewSupervisor creates a Supervisor for cncConfig that stores through dbPool, and through the
// telemetry queue of rabbitConfig when it is enabled (its writers get a pool of their own on
// dbConfig); the collectors write machine cycles through cycleService and publish their readings
// to hub (nil = live data off).
func NewSupervisor(cncConfig config.CNCConfig, rabbitConfig config.RabbitMQConfig, dbConfig config.DBConfig,
	dbPool *pgxpool.Pool, cycleService *service.MachineCycleService, hub *live.Hub) *Supervisor {
	ctx, cancel := context.WithCancel(context.Background())
	return &Supervisor{config: cncConfig, rabbitConfig: rabbitConfig, dbConfig: dbConfig, dbPool: dbPool,
		cycleService: cycleService, hub: hub, ctx: ctx, cancel: cancel}
}

// StartTelemetryQueue starts the telemetry queue (ADR 0024) when RABBITMQ_ENABLED: the publisher
// and the RABBITMQ_CONSUMERS writers. Call it before the collectors start, so they get the
// queue. A broker that is down does not stop the server: the publisher keeps reconnecting and the
// workers write directly meanwhile.
func (supervisor *Supervisor) StartTelemetryQueue() {
	if !supervisor.rabbitConfig.Enabled {
		slog.Info("telemetry queue off (RABBITMQ_ENABLED=false): collectors write directly")
		return
	}
	writerPool := supervisor.openWriterPool()
	telemetryService := service.NewMachineTelemetryService(writerPool, repository.NewMachineTelemetryRepository(writerPool))
	supervisor.telemetryQueue = monitor.NewTelemetryQueue(supervisor.rabbitConfig, supervisor.config, telemetryService)
	supervisor.telemetryQueue.Start()
}

// openWriterPool opens RABBITMQ_CONSUMERS connections for the queue writers alone (ADR 0027): on
// the shared pool, the direct writes of every worker under overload took the connections and the
// writers slowed down exactly when the queue grew. When it cannot be opened the writers share
// the main pool, as before.
func (supervisor *Supervisor) openWriterPool() *pgxpool.Pool {
	writerConfig := supervisor.dbConfig
	writerConfig.MaxConns = supervisor.rabbitConfig.Consumers
	writerPool, err := db.Connect(supervisor.ctx, writerConfig)
	if err != nil {
		slog.Error("telemetry queue writers share the main database pool: own pool not opened", "error", err)
		return supervisor.dbPool
	}
	supervisor.writerPool = writerPool
	return writerPool
}

// StartCNCDataCollection loads the CNC DLL and starts the collector in a background goroutine,
// then returns at once. If the build is not 32-bit or the DLL is missing, the error is logged and
// the API keeps running without collection.
func (supervisor *Supervisor) StartCNCDataCollection() error {
	// DATA_COLLECTION_STEP_3: Load the CNC DLL and register its callback slots (one per machine).
	driver, err := cnc.LoadDriver(supervisor.config)
	if err != nil {
		if errors.Is(err, cnc.ErrNotSupported) {
			slog.Warn("cnc data collection skipped: not a 32-bit build (build with GOARCH=386 for CNC collection)")
			return nil
		}
		slog.Error("cnc data collection not started: cnc dll unavailable", "error", err)
		return nil
	}

	// DATA_COLLECTION_STEP_3.2: Initializing/creating the service and repository objects - nothing more than that.
	telemetryService := service.NewMachineTelemetryService(supervisor.dbPool,
		repository.NewMachineTelemetryRepository(supervisor.dbPool))

	// DATA_COLLECTION_STEP_4: Run the collector in the background so the HTTP server starts without waiting for the database or the machines.
	// Background goroutine pattern in Go: wg.Add(1) before starting goroutine, defer wg.Done() inside goroutine,
	// and wg.Wait() to wait for completion during shutdown.
	supervisor.wg.Add(1) // "I'm starting 1 background task"
	go func() {
		defer supervisor.wg.Done() // "I'll say I'm done when finished"
		supervisor.runCollectorWithRetry(supervisor.ctx, telemetryService, driver)
	}()
	slog.Info("cnc data collection started")
	return nil
}

// StartPLCDataCollection starts the PLC collector in a background goroutine when PLC_ENABLED:
// every active Siemens machine is connected over OPC UA and its PLC variables are read. board
// receives each machine's connection status (the diagnostics read it); certificate is presented
// to controllers that need Basic256Sha256.
func (supervisor *Supervisor) StartPLCDataCollection(board *opcua.StatusBoard, certificate cert.Paths) {
	if !supervisor.config.PLC.Enabled {
		slog.Info("plc data collection off (PLC_ENABLED=false)")
		return
	}
	telemetryService := service.NewMachineTelemetryService(supervisor.dbPool,
		repository.NewMachineTelemetryRepository(supervisor.dbPool))
	collector := monitor.NewPlcCollector(supervisor.config, telemetryService, supervisor.hub, supervisor.telemetryQueue, board, certificate)
	supervisor.wg.Add(1)
	go func() {
		defer supervisor.wg.Done()
		collector.RunCollection(supervisor.ctx)
	}()
	slog.Info("plc data collection started")
}

// Stop stops collecting: every machine is disconnected, what was collected is written and the
// machines are left offline. Then the telemetry queue sends what it still holds and its writers
// store their last batch, before the caller closes the database pool. Each step waits at most
// CNC_STOP_TIMEOUT_SEC (plus a margin).
func (supervisor *Supervisor) Stop() {
	defer supervisor.stopTelemetryQueue()
	// Cancelling makes Collector.RunCollection run its shutdown; the wait is bounded so a
	// worker stuck inside the DLL cannot hang the server's shutdown.
	supervisor.cancel()
	done := make(chan struct{})
	go func() {
		supervisor.wg.Wait() // "Wait until all tasks say they're done"
		close(done)
	}()
	select {
	case <-done:
		slog.Info("data collection stopped")
	case <-time.After(supervisor.config.StopTimeout + 5*time.Second):
		slog.Error("data collection did not stop in time")
	}
}

// stopTelemetryQueue stops the telemetry queue once the collectors have stopped publishing, then
// closes the writers' pool.
func (supervisor *Supervisor) stopTelemetryQueue() {
	if supervisor.telemetryQueue != nil {
		supervisor.telemetryQueue.Stop()
	}
	if supervisor.writerPool != nil {
		supervisor.writerPool.Close()
	}
}

// runCollectorWithRetry runs the collector until ctx is cancelled. A collector that could not
// start (for example the database is not ready yet) is created again after a wait that
// doubles from 1 s up to 1 min.
func (supervisor *Supervisor) runCollectorWithRetry(ctx context.Context, telemetryService *service.MachineTelemetryService, driver *cnc.Driver) {
	retryWait := time.Second
	for ctx.Err() == nil {
		collector := monitor.NewCollector(supervisor.config, telemetryService, supervisor.cycleService, supervisor.hub,
			supervisor.telemetryQueue, driver, collectorShard, collectorShardCount)
		err := collector.RunCollection(ctx)
		if ctx.Err() != nil {
			return
		}
		slog.Error("cnc collector exited, restarting", "error", err, "retryIn", retryWait.String())
		select {
		case <-ctx.Done():
			return
		case <-time.After(retryWait):
		}
		retryWait = min(retryWait*2, time.Minute)
	}
}
