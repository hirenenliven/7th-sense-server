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
	shard  = 0
	shards = 1
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

// ErrWrongBuild is returned by Start when CNC_ENABLED=true but the server is not a 32-bit build:
// the server then does not start.
var ErrWrongBuild = errors.New("cnc data collection needs a 32-bit build: start the server with " +
	`.\run.ps1 or go tool air, or build it with .\build.ps1 (GOARCH=386); CNC_ENABLED=false serves the API only`)

// Start loads the DLL and starts collecting, then returns at once. It does nothing when
// CNC_ENABLED=false and returns ErrWrongBuild for a build that cannot load the 32-bit DLL;
// when the DLL itself cannot be loaded the error is logged and the API keeps running.
func (supervisor *Supervisor) Start() error {
	if !supervisor.config.Enabled {
		slog.Info("cnc data collection disabled (CNC_ENABLED=false)")
		return nil
	}
	driver, err := cnc.Open(supervisor.config)
	if err != nil {
		if errors.Is(err, cnc.ErrNotSupported) {
			return ErrWrongBuild
		}
		slog.Error("cnc data collection not started: cnc dll unavailable", "error", err)
		return nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	supervisor.cancel = cancel
	telemetryService := service.NewMachineTelemetryService(supervisor.dbPool,
		repository.NewMachineTelemetryRepository(supervisor.dbPool))
	supervisor.wg.Add(1)
	go func() {
		defer supervisor.wg.Done()
		supervisor.keepRunning(ctx, telemetryService, driver)
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
	supervisor.cancel()
	done := make(chan struct{})
	go func() {
		supervisor.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		slog.Info("cnc data collection stopped")
	case <-time.After(supervisor.config.StopTimeout + 5*time.Second):
		slog.Error("cnc data collection did not stop in time")
	}
}

// keepRunning runs the collector until ctx ends; a collector that could not start is retried
// with backoff.
func (supervisor *Supervisor) keepRunning(ctx context.Context, telemetryService *service.MachineTelemetryService, driver *cnc.Driver) {
	wait := time.Second
	for ctx.Err() == nil {
		collector := monitor.NewCollector(supervisor.config, telemetryService, driver, shard, shards)
		err := collector.Run(ctx)
		if ctx.Err() != nil {
			return
		}
		slog.Error("cnc collector exited, restarting", "error", err, "retryIn", wait.String())
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
		wait = min(wait*2, time.Minute)
	}
}
