// Package collector runs the CNC data collection inside the API server: starting the server
// connects the machines and starts storing their data; there is no other process or service.
// No DLL code and no SQL here.
//
// The CNC DLL is 32-bit, so the server must be built with GOARCH=386 to collect (Air does this,
// see .air.toml, run.ps1 and build.ps1). Any other build refuses to start while CNC_ENABLED=true.
package collector

import (
	"context"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"7th-sense-server/internal/cnc"
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/monitor"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/service"
)

// One process collects every machine: shard 0 of 1.
const (
	collectorShard      = 0
	collectorShardCount = 1
)

// Supervisor runs the collector in the server process and restarts it when it cannot start
// (for example the database is not ready yet).
type Supervisor struct {
	config config.CNCConfig
	dbPool *pgxpool.Pool
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

// NewSupervisor creates a Supervisor for cncConfig that stores through dbPool.
func NewSupervisor(cncConfig config.CNCConfig, dbPool *pgxpool.Pool) *Supervisor {
	return &Supervisor{config: cncConfig, dbPool: dbPool}
}

// ErrWrongBuild is returned by StartCNCDataCollection when CNC_ENABLED=true but the server is not a 32-bit build:
// the server then does not start.
var ErrWrongBuild = errors.New("cnc data collection needs a 32-bit build: start the server with " +
	`.\run.ps1 or go tool air, or build it with .\build.ps1 (GOARCH=386); CNC_ENABLED=false serves the API only`)

// StartCNCDataCollection loads the CNC DLL and starts the collector in a background goroutine,
// then returns at once. It does nothing when CNC_ENABLED=false and returns ErrWrongBuild for a
// build that cannot load the 32-bit DLL; when the DLL itself cannot be loaded the error is
// logged and the API keeps running without collection.
func (supervisor *Supervisor) StartCNCDataCollection() error {
	// DATA_COLLECTION_STEP_2: Check that collection is enabled (CNC_ENABLED).
	// if !supervisor.config.Enabled {
	// 	slog.Info("cnc data collection disabled (CNC_ENABLED=false)")
	// 	return nil
	// }
	// DATA_COLLECTION_STEP_3: Load the CNC DLL and register its callback slots (one per machine).
	driver, err := cnc.LoadDriver(supervisor.config)
	if err != nil {
		if errors.Is(err, cnc.ErrNotSupported) {
			return ErrWrongBuild
		}
		slog.Error("cnc data collection not started: cnc dll unavailable", "error", err)
		return nil
	}

	// DATA_COLLECTION_STEP_3.1: Create a cancellation context to gracefully stop the collector goroutine on shutdown.
	ctx, cancel := context.WithCancel(context.Background())
	supervisor.cancel = cancel

	// DATA_COLLECTION_STEP_3.2: Initializing/creating the service and repository objects - nothing more than that.
	telemetryService := service.NewMachineTelemetryService(supervisor.dbPool,
		repository.NewMachineTelemetryRepository(supervisor.dbPool))

	// DATA_COLLECTION_STEP_4: Run the collector in the background so the HTTP server starts without waiting for the database or the machines.
	// Background goroutine pattern in Go: wg.Add(1) before starting goroutine, defer wg.Done() inside goroutine,
	// and wg.Wait() to wait for completion during shutdown.
	supervisor.wg.Add(1) // "I'm starting 1 background task"
	go func() {
		defer supervisor.wg.Done() // "I'll say I'm done when finished"
		supervisor.runCollectorWithRetry(ctx, telemetryService, driver)
	}()
	slog.Info("cnc data collection started")
	return nil
}

// Stop stops collecting: every machine is disconnected, what was collected is written and the
// machines are left offline. It waits at most CNC_STOP_TIMEOUT_SEC (plus a margin).
func (supervisor *Supervisor) Stop() {
	if supervisor.cancel == nil {
		return
	}
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
		slog.Info("cnc data collection stopped")
	case <-time.After(supervisor.config.StopTimeout + 5*time.Second):
		slog.Error("cnc data collection did not stop in time")
	}
}

// runCollectorWithRetry runs the collector until ctx is cancelled. A collector that could not
// start (for example the database is not ready yet) is created again after a wait that
// doubles from 1 s up to 1 min.
func (supervisor *Supervisor) runCollectorWithRetry(ctx context.Context, telemetryService *service.MachineTelemetryService, driver *cnc.Driver) {
	retryWait := time.Second
	for ctx.Err() == nil {
		collector := monitor.NewCollector(supervisor.config, telemetryService, driver, collectorShard, collectorShardCount)
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
