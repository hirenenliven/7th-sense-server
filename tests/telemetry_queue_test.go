package tests

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strconv"
	"sync"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/monitor"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/service"
)

// Machine ids no real machine has: the telemetry tables have no foreign keys, so the tests write
// rows for them and delete those rows afterwards.
const (
	telemetryTestMachineID     int64 = 2_100_000_001
	telemetryLockFirstMachine  int64 = 2_100_000_100
	telemetryLoadFirstMachine  int64 = 2_100_001_000
	telemetryLoadVariableCount       = 40
)

// deleteTelemetryRows removes the rows the tests wrote for machines from first to last.
func deleteTelemetryRows(t *testing.T, first, last int64) {
	t.Helper()
	ctx := context.Background()
	for _, table := range []string{"machine_telementry", "machine_variable_values"} {
		if _, err := testDBPool.Exec(ctx, `DELETE FROM `+table+` WHERE machine_id BETWEEN $1 AND $2`, first, last); err != nil {
			t.Errorf("clean up %s: %v", table, err)
		}
	}
}

func newTelemetryService() *service.MachineTelemetryService {
	return service.NewMachineTelemetryService(testDBPool, repository.NewMachineTelemetryRepository(testDBPool))
}

// TestTelemetryWritesAreRepeatable checks the two writes the telemetry queue may repeat or
// reorder (ADR 0024): the same telemetry rows written twice are stored once (migration 039), and
// a live value read earlier never replaces a newer one.
func TestTelemetryWritesAreRepeatable(t *testing.T) {
	ctx := context.Background()
	telemetryService := newTelemetryService()
	deleteTelemetryRows(t, telemetryTestMachineID, telemetryTestMachineID)
	t.Cleanup(func() { deleteTelemetryRows(t, telemetryTestMachineID, telemetryTestMachineID) })

	newer := time.Now().UTC().Truncate(time.Microsecond)
	older := newer.Add(-time.Second)
	newValue, oldValue := 2.0, 1.0

	t.Run("telemetry written twice is stored once", func(t *testing.T) {
		var indexExists bool
		if err := testDBPool.QueryRow(ctx, `SELECT EXISTS (SELECT 1 FROM pg_indexes
			WHERE tablename = 'machine_telementry' AND indexname = 'machine_telementry_machine_variable_reading_key')`).Scan(&indexExists); err != nil {
			t.Fatal(err)
		}
		if !indexExists {
			t.Skip("migration 039 not applied: machine_telementry has no unique reading index")
		}
		rows := []model.TelemetryValue{
			{MachineID: telemetryTestMachineID, VariableID: 1, Value: &newValue, RecordedAt: newer},
			{MachineID: telemetryTestMachineID, VariableID: 2, Value: &newValue, RecordedAt: newer},
		}
		for range 2 {
			if err := telemetryService.InsertTelemetry(ctx, rows); err != nil {
				t.Fatalf("insert telemetry: %v", err)
			}
		}
		var count int
		if err := testDBPool.QueryRow(ctx, `SELECT count(*) FROM machine_telementry WHERE machine_id = $1`,
			telemetryTestMachineID).Scan(&count); err != nil {
			t.Fatal(err)
		}
		if count != 2 {
			t.Fatalf("stored %d rows, want 2", count)
		}
	})

	t.Run("older live value never replaces a newer one", func(t *testing.T) {
		// One batch holding both readings, then the older one again on its own.
		err := telemetryService.UpsertLiveValues(ctx, []model.LiveValue{
			{MachineID: telemetryTestMachineID, VariableID: 1, Value: &newValue, UpdatedAt: newer},
			{MachineID: telemetryTestMachineID, VariableID: 1, Value: &oldValue, UpdatedAt: older},
		})
		if err != nil {
			t.Fatalf("upsert live values: %v", err)
		}
		if err := telemetryService.UpsertLiveValues(ctx, []model.LiveValue{
			{MachineID: telemetryTestMachineID, VariableID: 1, Value: &oldValue, UpdatedAt: older},
		}); err != nil {
			t.Fatalf("upsert live values: %v", err)
		}
		values, err := telemetryService.GetLiveValues(ctx, telemetryTestMachineID, []int{1})
		if err != nil {
			t.Fatal(err)
		}
		stored := values[1]
		if stored.Value == nil || *stored.Value != newValue || !stored.UpdatedAt.Equal(newer) {
			t.Fatalf("stored %+v, want value %v read at %s", stored, newValue, newer)
		}
	})
}

// TestLiveValueUpsertsDoNotDeadlock runs what the telemetry queue writers do at the same time:
// several writers upserting the same machine_variable_values rows, each batch in another order.
// UpsertLiveValues sorts the rows, so they lock them in the same order and never deadlock
// (unsorted, the 5 values/s load test of ADR 0024 hit "deadlock detected" 51 times).
func TestLiveValueUpsertsDoNotDeadlock(t *testing.T) {
	const machines, writers, rounds = 50, 4, 20
	lastMachine := telemetryLockFirstMachine + machines - 1
	deleteTelemetryRows(t, telemetryLockFirstMachine, lastMachine)
	t.Cleanup(func() { deleteTelemetryRows(t, telemetryLockFirstMachine, lastMachine) })
	telemetryService := newTelemetryService()

	started := time.Now().UTC().Truncate(time.Microsecond)
	var done sync.WaitGroup
	errs := make(chan error, writers*rounds)
	for writer := range writers {
		done.Add(1)
		go func() {
			defer done.Done()
			for round := range rounds {
				readAt := started.Add(time.Duration(round*writers+writer) * time.Millisecond)
				value := float64(round)
				rows := make([]model.LiveValue, 0, machines*telemetryLoadVariableCount)
				for machine := range int64(machines) {
					for variable := 1; variable <= telemetryLoadVariableCount; variable++ {
						rows = append(rows, model.LiveValue{MachineID: telemetryLockFirstMachine + machine, VariableID: variable,
							Value: &value, UpdatedAt: readAt})
					}
				}
				if writer%2 == 1 {
					slices.Reverse(rows)
				}
				if err := telemetryService.UpsertLiveValues(context.Background(), rows); err != nil {
					errs <- err
				}
			}
		}()
	}
	done.Wait()
	close(errs)
	for err := range errs {
		t.Errorf("upsert live values: %v", err)
	}
}

// TestTelemetryQueueEndToEnd publishes simulated readings through the real publisher → RabbitMQ
// → writer → PostgreSQL and checks every row arrives once. It runs only with RABBITMQ_TEST_URL;
// RABBITMQ_LOAD_MACHINES (default 50) machines × 40 variables send one reading a second for
// RABBITMQ_LOAD_SECONDS (default 5), through RABBITMQ_LOAD_CONSUMERS writers (default 2) on a
// RABBITMQ_LOAD_QUEUE_TYPE queue (default quorum). For the load test of ADR 0024 use e.g. 2000
// machines for 120 s and watch the queue in the management UI.
func TestTelemetryQueueEndToEnd(t *testing.T) {
	brokerURL := os.Getenv("RABBITMQ_TEST_URL")
	if brokerURL == "" {
		t.Skip("RABBITMQ_TEST_URL not set: no broker to test against")
	}
	machines := envInt("RABBITMQ_LOAD_MACHINES", 50)
	seconds := envInt("RABBITMQ_LOAD_SECONDS", 5)
	lastMachine := telemetryLoadFirstMachine + int64(machines) - 1
	deleteTelemetryRows(t, telemetryLoadFirstMachine, lastMachine)
	t.Cleanup(func() { deleteTelemetryRows(t, telemetryLoadFirstMachine, lastMachine) })

	name := fmt.Sprintf("telemetry.test.%d", time.Now().UnixNano())
	rabbitConfig := config.RabbitMQConfig{Enabled: true, URL: brokerURL, Exchange: name, RoutingKey: "values", Queue: name,
		QueueType: envText("RABBITMQ_LOAD_QUEUE_TYPE", config.QueueTypeQuorum), Consumers: envInt("RABBITMQ_LOAD_CONSUMERS", 2), Prefetch: 1000, BatchRows: 5000, BatchInterval: 200 * time.Millisecond,
		PublishBuffer: 20000, ReconnectMax: 5 * time.Second}
	t.Cleanup(func() { deleteQueueTopology(t, rabbitConfig) })
	cncConfig := config.CNCConfig{StopTimeout: 10 * time.Second}

	telemetryQueue := monitor.NewTelemetryQueue(rabbitConfig, cncConfig, newTelemetryService())
	telemetryQueue.Start()
	defer telemetryQueue.Stop()

	started := time.Now()
	fallbacks := 0
	readAt := time.Now().UTC().Truncate(time.Second)
	for second := range seconds {
		for machine := range machines {
			if !telemetryQueue.Enqueue(simulatedReading(telemetryLoadFirstMachine+int64(machine), readAt, second)) {
				fallbacks++
			}
		}
		readAt = readAt.Add(time.Second)
		time.Sleep(time.Until(started.Add(time.Duration(second+1) * time.Second)))
	}
	if fallbacks > 0 {
		t.Fatalf("%d readings were refused by the queue (publish buffer full)", fallbacks)
	}

	want := machines * seconds * telemetryLoadVariableCount
	deadline := time.Now().Add(time.Minute)
	var stored int
	for {
		if err := testDBPool.QueryRow(context.Background(), `SELECT count(*) FROM machine_telementry WHERE machine_id BETWEEN $1 AND $2`,
			telemetryLoadFirstMachine, lastMachine).Scan(&stored); err != nil {
			t.Fatal(err)
		}
		if stored >= want || time.Now().After(deadline) {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	elapsed := time.Since(started)
	t.Logf("%d machines × %d variables × %d s: %d telemetry rows stored in %s (%.0f rows/s incl. live values)",
		machines, telemetryLoadVariableCount, seconds, stored, elapsed.Round(time.Millisecond), float64(2*stored)/elapsed.Seconds())
	if stored != want {
		t.Fatalf("stored %d telemetry rows, want %d", stored, want)
	}
	var liveRows int
	if err := testDBPool.QueryRow(context.Background(), `SELECT count(*) FROM machine_variable_values WHERE machine_id BETWEEN $1 AND $2`,
		telemetryLoadFirstMachine, lastMachine).Scan(&liveRows); err != nil {
		t.Fatal(err)
	}
	if liveRows != machines*telemetryLoadVariableCount {
		t.Fatalf("live rows %d, want %d", liveRows, machines*telemetryLoadVariableCount)
	}
}

// simulatedReading is one reading of machineID: every variable changed (the worst case).
func simulatedReading(machineID int64, readAt time.Time, step int) model.TelemetryReading {
	live := make([]model.LiveValue, telemetryLoadVariableCount)
	telemetry := make([]model.TelemetryValue, telemetryLoadVariableCount)
	for index := range telemetryLoadVariableCount {
		value := float64(step*1000 + index)
		live[index] = model.LiveValue{MachineID: machineID, VariableID: index + 1, Value: &value, UpdatedAt: readAt}
		telemetry[index] = model.TelemetryValue{MachineID: machineID, VariableID: index + 1, Value: &value, RecordedAt: readAt}
	}
	return model.NewTelemetryReading(machineID, readAt, model.TelemetrySourceCNC, live, telemetry)
}

// deleteQueueTopology removes the exchanges and queues the test declared.
func deleteQueueTopology(t *testing.T, settings config.RabbitMQConfig) {
	connection, err := amqp.Dial(settings.URL)
	if err != nil {
		t.Logf("clean up broker: %v", err)
		return
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		t.Logf("clean up broker: %v", err)
		return
	}
	channel.QueueDelete(settings.Queue, false, false, false)
	channel.QueueDelete(settings.Queue+".dead", false, false, false)
	channel.ExchangeDelete(settings.Exchange, false, false)
	channel.ExchangeDelete(settings.Exchange+".dlx", false, false)
}

// envText reads a value from the environment, else fallback.
func envText(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

// envInt reads a positive number from the environment, else fallback.
func envInt(key string, fallback int) int {
	value, err := strconv.Atoi(os.Getenv(key))
	if err != nil || value <= 0 {
		return fallback
	}
	return value
}
