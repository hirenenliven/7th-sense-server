package tests

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/db"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/service"
)

// TestTelemetryWriteIsolation times the two halves of a queue batch separately — the
// machine_telementry COPY and the machine_variable_values upsert — plus the combined
// WriteTelemetryBatch, at the batch shape the 40k values/s run produced (~10k rows).
// Runs only with TELEMETRY_WRITE_BENCH=1; uses its own machine ids and removes them after.
func TestTelemetryWriteIsolation(t *testing.T) {
	if os.Getenv("TELEMETRY_WRITE_BENCH") == "" {
		t.Skip("TELEMETRY_WRITE_BENCH not set")
	}
	pool := testDBPool

	const machines = 250 // 250 machines x 40 vars = 10,000 rows per table per batch
	first := telemetryLoadBaseMachine + 900_000
	last := first + machines - 1
	deleteTelemetryRows(t, first, last)
	defer deleteTelemetryRows(t, first, last)

	telemetryRepo := repository.NewMachineTelemetryRepository(pool)
	telemetryService := service.NewMachineTelemetryService(pool, telemetryRepo)

	batches := envInt("TELEMETRY_WRITE_BENCH_BATCHES", 12)
	rowsPerTable := machines * telemetryLoadVariableCount

	mkRows := func(batch int) ([]model.TelemetryValue, []model.LiveValue) {
		at := time.Now().UTC().Truncate(time.Microsecond).Add(time.Duration(batch) * time.Second)
		telemetry := make([]model.TelemetryValue, 0, rowsPerTable)
		live := make([]model.LiveValue, 0, rowsPerTable)
		for m := int64(0); m < machines; m++ {
			machineID := first + m
			for v := 1; v <= telemetryLoadVariableCount; v++ {
				value := float64(batch*1000 + v)
				telemetry = append(telemetry, model.TelemetryValue{MachineID: machineID, VariableID: v, Value: &value, RecordedAt: at})
				live = append(live, model.LiveValue{MachineID: machineID, VariableID: v, Value: &value, UpdatedAt: at})
			}
		}
		return telemetry, live
	}

	run := func(name string, write func(ctx context.Context, telemetry []model.TelemetryValue, live []model.LiveValue) error) {
		var total time.Duration
		var totalRows int64
		for b := 0; b < batches; b++ {
			telemetry, live := mkRows(b)
			start := time.Now()
			if err := write(context.Background(), telemetry, live); err != nil {
				t.Fatalf("%s batch %d: %v", name, b, err)
			}
			total += time.Since(start)
			totalRows += int64(len(telemetry) + len(live))
		}
		t.Logf("%-18s %d batches x %d rows | avg %.1f ms/batch | %.0f rows/s",
			name, batches, rowsPerTable*2, float64(total.Milliseconds())/float64(batches),
			float64(totalRows)/total.Seconds())
	}

	// Warm-up: first writes create chunks / grow the live table.
	telemetry, live := mkRows(-1)
	_ = telemetryService.WriteTelemetryBatch(context.Background(), telemetry, live)

	t.Run("telemetry COPY only", func(t *testing.T) {
		run("copy", func(ctx context.Context, telemetry []model.TelemetryValue, _ []model.LiveValue) error {
			return db.WithTransaction(ctx, pool, func(tx pgx.Tx) error {
				repo := repository.NewMachineTelemetryRepository(tx)
				if err := repo.UseAsynchronousCommit(ctx); err != nil {
					return err
				}
				return repo.CopyTelemetry(ctx, telemetry)
			})
		})
	})

	t.Run("live upsert only", func(t *testing.T) {
		run("upsert", func(ctx context.Context, _ []model.TelemetryValue, live []model.LiveValue) error {
			return db.WithTransaction(ctx, pool, func(tx pgx.Tx) error {
				repo := repository.NewMachineTelemetryRepository(tx)
				if err := repo.UseAsynchronousCommit(ctx); err != nil {
					return err
				}
				return repo.UpsertLiveValues(ctx, live)
			})
		})
	})

	t.Run("combined batch", func(t *testing.T) {
		run("combined", telemetryService.WriteTelemetryBatch)
	})

	fmt.Println("done")
}
