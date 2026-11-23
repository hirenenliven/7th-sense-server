package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"math"
	"os"
	"runtime"
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/monitor"
	"7th-sense-server/internal/queue"
)

// The load test's machines: ids no real machine has, one block of 100,000 per run.
const telemetryLoadBaseMachine int64 = 2_100_100_000

// TestTelemetryLoad measures how much telemetry the RabbitMQ queue (ADR 0024) carries and
// stores, the way the product does it: CNC-style readings go through TelemetryQueue.Enqueue →
// the real publisher → RabbitMQ → the real writers → PostgreSQL, with the queue settings of
// config.Load (batch, prefetch, buffer, RABBITMQ_CONSUMERS) and the direct write when the queue
// refuses a reading. It runs only with TELEMETRY_LOAD=1 and uses RABBITMQ_TEST_URL, else
// RABBITMQ_URL. Every run uses its own exchange / queue and machine ids and removes them after.
//
// Phases (TELEMETRY_LOAD_PHASES, default all): broker (RabbitMQ alone, no database), worst
// (every variable changes on every reading), realistic (TELEMETRY_LOAD_CHANGE share of the
// variables changes, default 0.25, plus the 30 s live refresh, as the CNC worker stores them),
// stepup (worst case with more and more machines until the database falls behind).
// TELEMETRY_LOAD_MACHINES (1000), TELEMETRY_LOAD_SECONDS (60), TELEMETRY_LOAD_STEPS
// ("1000,1500,2000,3000"), TELEMETRY_LOAD_STEP_SECONDS (20), TELEMETRY_LOAD_BROKER_SECONDS (15),
// TELEMETRY_LOAD_RATE (readings per machine per second, default 1; 5 = every variable sends 5
// values a second, 1000 machines -> 200,000 values/s).
//
//	TELEMETRY_LOAD=1 go test ./tests -run TestTelemetryLoad -v -timeout 60m
func TestTelemetryLoad(t *testing.T) {
	if os.Getenv("TELEMETRY_LOAD") == "" {
		t.Skip("TELEMETRY_LOAD not set: the load test runs only on request")
	}
	brokerURL := envText("RABBITMQ_TEST_URL", os.Getenv("RABBITMQ_URL"))
	if brokerURL == "" {
		t.Skip("neither RABBITMQ_TEST_URL nor RABBITMQ_URL is set: no broker to test against")
	}
	phases := strings.Split(envText("TELEMETRY_LOAD_PHASES", "broker,worst,realistic,stepup"), ",")
	machines := envInt("TELEMETRY_LOAD_MACHINES", 1000)
	seconds := envInt("TELEMETRY_LOAD_SECONDS", 60)
	changeRatio := envRatio("TELEMETRY_LOAD_CHANGE", 0.25)
	rate := envInt("TELEMETRY_LOAD_RATE", 1)
	product := config.Load()
	t.Logf("settings: %d writers, %s queue, batch %d rows / %s, prefetch %d, publish buffer %d, DB pool %d, %d CPUs",
		product.RabbitMQ.Consumers, product.RabbitMQ.QueueType, product.RabbitMQ.BatchRows, product.RabbitMQ.BatchInterval,
		product.RabbitMQ.Prefetch, product.RabbitMQ.PublishBuffer, product.DB.MaxConns, runtime.NumCPU())

	var results []loadResult
	run := 0
	runPlan := func(plan loadPlan) loadResult {
		run++
		result := runTelemetryLoad(t, brokerURL, product, run, plan)
		results = append(results, result)
		t.Log(result.String())
		return result
	}
	for _, phase := range phases {
		switch strings.TrimSpace(phase) {
		case "broker":
			t.Run("broker bandwidth", func(t *testing.T) {
				measureBrokerBandwidth(t, brokerURL, product.RabbitMQ, envInt("TELEMETRY_LOAD_BROKER_SECONDS", 15))
			})
		case "worst":
			t.Run("worst case", func(t *testing.T) {
				runPlan(loadPlan{name: "worst", machines: machines, seconds: seconds, rate: rate, changeRatio: 1})
			})
		case "realistic":
			t.Run("realistic", func(t *testing.T) {
				runPlan(loadPlan{name: "realistic", machines: machines, seconds: seconds, rate: rate, changeRatio: changeRatio})
			})
		case "stepup":
			t.Run("step up", func(t *testing.T) {
				for _, step := range envInts("TELEMETRY_LOAD_STEPS", []int{1000, 1500, 2000, 3000}) {
					result := runPlan(loadPlan{name: "stepup", machines: step, seconds: envInt("TELEMETRY_LOAD_STEP_SECONDS", 20), rate: rate, changeRatio: 1})
					if !result.keptUp() {
						t.Logf("stopped stepping up: %d machines fell behind", step)
						break
					}
				}
			})
		}
	}
	t.Log("summary\n" + loadSummary(results))
}

// loadPlan is one end-to-end run: machines each sending rate readings of 40 variables a second.
type loadPlan struct {
	name        string
	machines    int
	seconds     int
	rate        int     // readings per machine per second
	changeRatio float64 // share of variables 2..40 that changed per reading; 1 = all (worst case)
}

// loadResult is what one run measured.
type loadResult struct {
	plan            loadPlan
	readings        int64
	queued          int64 // readings handed to RabbitMQ
	fallbacks       int64 // readings the queue refused: written directly, as the worker does
	telemetryRows   int64 // machine_telementry rows offered
	liveRows        int64 // machine_variable_values rows offered
	bodyBytes       int64
	storedTelemetry int64
	storedLive      int64
	generate        time.Duration // how long sending took (planned: seconds)
	drain           time.Duration // from the last reading sent until every row was stored
	peakDepth       int           // most messages waiting in RabbitMQ
	ages            []time.Duration
}

// keptUp: nothing fell back to direct writes and the queue was empty within 5 s of the end.
func (result loadResult) keptUp() bool {
	return result.fallbacks == 0 && result.drain < 5*time.Second && result.lost() == 0
}

func (result loadResult) lost() int64 { return result.telemetryRows - result.storedTelemetry }

func (result loadResult) offeredRowsPerSecond() float64 {
	return float64(result.telemetryRows+result.liveRows) / float64(result.plan.seconds)
}

// sentValuesPerSecond is how many telemetry values the machines really sent a second (sending
// takes longer than planned when direct writes hold workers up).
func (result loadResult) sentValuesPerSecond() float64 {
	return float64(result.telemetryRows) / result.generate.Seconds()
}

func (result loadResult) storedRowsPerSecond() float64 {
	return float64(result.telemetryRows+result.liveRows) / (result.generate + result.drain).Seconds()
}

// agePercentile is the age of the newest data every machine has in the database, sampled
// once a second while sending: p50 / p95 / max.
func (result loadResult) agePercentile(percentile float64) time.Duration {
	if len(result.ages) == 0 {
		return 0
	}
	sorted := slices.Clone(result.ages)
	slices.Sort(sorted)
	return sorted[int(math.Ceil(percentile*float64(len(sorted))))-1]
}

func (result loadResult) String() string {
	return fmt.Sprintf("%s: %d machines × 40 variables × %d/s × %d s | %d readings (%.0f msg/s, %.1f KB avg, %.2f MB/s) | rows offered %d telemetry + %d live = %.0f rows/s | sent %.0f values/s, %d readings to RabbitMQ | stored %d telemetry (lost %d), %d live rows | send %s, drain %s, stored %.0f rows/s | peak queue %d msgs | data age p50 %s p95 %s max %s | fallback %d | kept up %v",
		result.plan.name, result.plan.machines, result.plan.rate, result.plan.seconds, result.readings, float64(result.readings)/float64(result.plan.seconds),
		float64(result.bodyBytes)/float64(max(result.readings, 1))/1024, float64(result.bodyBytes)/float64(result.plan.seconds)/1e6,
		result.telemetryRows, result.liveRows, result.offeredRowsPerSecond(), result.sentValuesPerSecond(), result.queued, result.storedTelemetry, result.lost(), result.storedLive,
		result.generate.Round(time.Millisecond), result.drain.Round(time.Millisecond), result.storedRowsPerSecond(), result.peakDepth,
		result.agePercentile(0.5).Round(time.Millisecond), result.agePercentile(0.95).Round(time.Millisecond),
		result.agePercentile(1).Round(time.Millisecond), result.fallbacks, result.keptUp())
}

// loadSummary is a markdown table of the runs.
func loadSummary(results []loadResult) string {
	var text strings.Builder
	text.WriteString("| Run | Machines | Readings/s each | Msg/s | MB/s | Rows/s offered | Rows/s stored | Peak queue | Data age p50 / p95 / max | Drain | Fallback | Lost | Kept up |\n")
	text.WriteString("|---|---|---|---|---|---|---|---|---|---|---|---|---|\n")
	for _, result := range results {
		fmt.Fprintf(&text, "| %s | %d | %d | %.0f | %.2f | %.0f | %.0f | %d | %s / %s / %s | %s | %d | %d | %v |\n",
			result.plan.name, result.plan.machines, result.plan.rate, float64(result.readings)/float64(result.plan.seconds),
			float64(result.bodyBytes)/float64(result.plan.seconds)/1e6, result.offeredRowsPerSecond(), result.storedRowsPerSecond(),
			result.peakDepth, result.agePercentile(0.5).Round(100*time.Millisecond), result.agePercentile(0.95).Round(100*time.Millisecond),
			result.agePercentile(1).Round(100*time.Millisecond), result.drain.Round(100*time.Millisecond), result.fallbacks, result.lost(), result.keptUp())
	}
	return text.String()
}

// loadQueueSettings are the product's queue settings on an exchange / queue of the test's own.
func loadQueueSettings(brokerURL string, product config.RabbitMQConfig) config.RabbitMQConfig {
	settings := product
	settings.Enabled, settings.URL = true, brokerURL
	name := fmt.Sprintf("telemetry.load.%d", time.Now().UnixNano())
	settings.Exchange, settings.Queue, settings.RoutingKey = name, name, "values"
	return settings
}

// runTelemetryLoad sends plan's readings through the product's telemetry queue and measures
// until every row is in the database.
func runTelemetryLoad(t *testing.T, brokerURL string, product config.Config, run int, plan loadPlan) loadResult {
	result := loadResult{plan: plan}
	first := telemetryLoadBaseMachine + int64(run)*100_000
	last := first + int64(plan.machines) - 1
	deleteTelemetryRows(t, first, last)
	defer deleteTelemetryRows(t, first, last)

	settings := loadQueueSettings(brokerURL, product.RabbitMQ)
	defer deleteQueueTopology(t, settings)
	telemetryService := newTelemetryService()
	telemetryQueue := monitor.NewTelemetryQueue(settings, product.CNC, telemetryService)
	telemetryQueue.Start()
	defer telemetryQueue.Stop()
	waitForLoadQueue(t, settings)

	// Watch the database and the queue once a second while sending.
	samplingCtx, stopSampling := context.WithCancel(context.Background())
	var samplingDone sync.WaitGroup
	var sampleMu sync.Mutex
	samplingDone.Add(1)
	go func() {
		defer samplingDone.Done()
		ticker := time.NewTicker(time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-samplingCtx.Done():
				return
			case <-ticker.C:
			}
			if machinesWithData, oldest, err := loadWatermark(first, last); err == nil && machinesWithData == plan.machines {
				sampleMu.Lock()
				result.ages = append(result.ages, time.Since(oldest))
				sampleMu.Unlock()
			}
			if waiting, _, err := queue.QueueDepth(settings); err == nil {
				sampleMu.Lock()
				result.peakDepth = max(result.peakDepth, waiting)
				sampleMu.Unlock()
			}
		}
	}()

	// One goroutine per machine, as the product has one worker per machine: it reads rate times
	// a second at its own moment (machines spread over ten offsets of the period, as workers
	// started at different moments are) and enqueues the reading; when the queue refuses it, that
	// machine alone writes it directly and is late for its next readings, like a real worker.
	const offsets = 10
	var readings, queued, fallbacks, telemetryRows, liveRows, bodyBytes atomic.Int64
	period := time.Second / time.Duration(plan.rate)
	steps := plan.seconds * plan.rate
	started := time.Now()
	var machinesDone sync.WaitGroup
	for machine := range plan.machines {
		machinesDone.Add(1)
		go func() {
			defer machinesDone.Done()
			machineID := first + int64(machine)
			offset := time.Duration(machine%offsets) * period / offsets
			for step := range steps {
				at := started.Add(offset + time.Duration(step)*period)
				time.Sleep(time.Until(at))
				readAt := at.UTC().Truncate(time.Microsecond)
				live, telemetry := loadReading(machineID, machine, readAt, step, plan.rate, plan.changeRatio)
				reading := model.NewTelemetryReading(machineID, readAt, model.TelemetrySourceCNC, live, telemetry)
				if body, err := json.Marshal(reading); err == nil {
					bodyBytes.Add(int64(len(body)))
				}
				readings.Add(1)
				telemetryRows.Add(int64(len(telemetry)))
				liveRows.Add(int64(len(live)))
				if telemetryQueue.Enqueue(reading) {
					queued.Add(1)
					continue
				}
				// Refused (buffer full): the worker writes the reading itself.
				fallbacks.Add(1)
				ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
				if err := telemetryService.InsertTelemetry(ctx, telemetry); err != nil {
					t.Errorf("direct telemetry write: %v", err)
				} else if err := telemetryService.UpsertLiveValues(ctx, live); err != nil {
					t.Errorf("direct live write: %v", err)
				}
				cancel()
			}
		}()
	}
	machinesDone.Wait()
	result.generate = time.Since(started)
	sentAt := time.Now()

	// Done when every machine's last reading is stored (its slot-0 time or later) and the
	// telemetry row count matches. The wait allows for the database writing 20,000 rows/s.
	lastStep := started.Add(time.Duration(steps-1) * period).UTC().Truncate(time.Microsecond)
	result.telemetryRows, result.liveRows = telemetryRows.Load(), liveRows.Load()
	deadline := time.Now().Add(5*time.Minute + time.Duration((result.telemetryRows+result.liveRows)/20_000)*time.Second)
	for time.Now().Before(deadline) {
		var done int
		if err := testDBPool.QueryRow(context.Background(), `SELECT count(*) FROM machine_variable_values
			WHERE machine_id BETWEEN $1 AND $2 AND variable_id = 1 AND updated_at >= $3`, first, last, lastStep).Scan(&done); err != nil {
			t.Fatal(err)
		}
		if done == plan.machines {
			if err := testDBPool.QueryRow(context.Background(), `SELECT count(*) FROM machine_telementry WHERE machine_id BETWEEN $1 AND $2`,
				first, last).Scan(&result.storedTelemetry); err != nil {
				t.Fatal(err)
			}
			if result.storedTelemetry >= result.telemetryRows {
				break
			}
		}
		time.Sleep(250 * time.Millisecond)
	}
	result.drain = time.Since(sentAt)
	stopSampling()
	samplingDone.Wait()
	if err := testDBPool.QueryRow(context.Background(), `SELECT count(*) FROM machine_variable_values WHERE machine_id BETWEEN $1 AND $2`,
		first, last).Scan(&result.storedLive); err != nil {
		t.Fatal(err)
	}
	result.readings, result.queued, result.fallbacks, result.bodyBytes = readings.Load(), queued.Load(), fallbacks.Load(), bodyBytes.Load()
	if result.lost() != 0 {
		t.Errorf("%s: %d telemetry rows offered, %d stored", plan.name, result.telemetryRows, result.storedTelemetry)
	}
	if result.storedLive != int64(plan.machines*telemetryLoadVariableCount) {
		t.Errorf("%s: %d live rows, want %d", plan.name, result.storedLive, plan.machines*telemetryLoadVariableCount)
	}
	return result
}

// loadReading is reading number step of a machine reading rate times a second, as the CNC worker
// builds it (worker.go): variable 1 changes on every reading (so the test can see how fresh each
// machine's data is), the others with changeRatio. Changed values go to machine_telementry; live
// values hold the changed ones, and every value on the first reading and every 30 s
// (LiveRefresh; machines started at different moments refresh at different seconds).
func loadReading(machineID int64, machine int, readAt time.Time, step, rate int, changeRatio float64) ([]model.LiveValue, []model.TelemetryValue) {
	second := step / rate
	refresh := step == 0 || (step%rate == 0 && (second+machine)%30 == 0)
	live := make([]model.LiveValue, 0, telemetryLoadVariableCount)
	telemetry := make([]model.TelemetryValue, 0, telemetryLoadVariableCount)
	for variable := 1; variable <= telemetryLoadVariableCount; variable++ {
		changed := variable == 1 || step == 0 || changeRatio >= 1 || loadChance(machine, step, variable) < changeRatio
		value := float64(step*1000 + variable)
		if changed || refresh {
			live = append(live, model.LiveValue{MachineID: machineID, VariableID: variable, Value: &value, UpdatedAt: readAt})
		}
		if changed {
			telemetry = append(telemetry, model.TelemetryValue{MachineID: machineID, VariableID: variable, Value: &value, RecordedAt: readAt})
		}
	}
	return live, telemetry
}

// loadChance is a repeatable number in [0, 1) for a machine, reading and variable.
func loadChance(machine, step, variable int) float64 {
	x := uint64(machine)*0x9E3779B97F4A7C15 ^ uint64(step)*0xBF58476D1CE4E5B9 ^ uint64(variable)*0x94D049BB133111EB
	x ^= x >> 31
	x *= 0xD6E8FEB86659FD93
	x ^= x >> 32
	return float64(x>>11) / float64(1<<53)
}

// loadWatermark gives how many machines have data and the read time of the oldest "newest
// reading" among them: every machine has data in the database at least that fresh.
func loadWatermark(first, last int64) (int, time.Time, error) {
	var machines int
	var oldest *time.Time
	err := testDBPool.QueryRow(context.Background(), `SELECT count(*), min(updated_at) FROM machine_variable_values
		WHERE machine_id BETWEEN $1 AND $2 AND variable_id = 1`, first, last).Scan(&machines, &oldest)
	if err != nil || oldest == nil {
		return machines, time.Time{}, err
	}
	return machines, *oldest, nil
}

// waitForLoadQueue waits until the publisher declared the queue and the writers can consume.
func waitForLoadQueue(t *testing.T, settings config.RabbitMQConfig) {
	t.Helper()
	deadline := time.Now().Add(15 * time.Second)
	for {
		if _, _, err := queue.QueueDepth(settings); err == nil {
			time.Sleep(time.Second) // the writers connect right after
			return
		} else if time.Now().After(deadline) {
			t.Fatalf("queue %s not declared: %v", settings.Queue, err)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

// measureBrokerBandwidth pushes worst-case readings (40 changed variables) through the
// product's publisher as fast as it takes them, with a plain consumer that only acknowledges:
// RabbitMQ's own limit on this machine, without the database.
func measureBrokerBandwidth(t *testing.T, brokerURL string, product config.RabbitMQConfig, seconds int) {
	settings := loadQueueSettings(brokerURL, product)
	defer deleteQueueTopology(t, settings)
	publisher := queue.NewPublisher(settings, func(level slog.Level, message string, attributes ...any) {
		if level >= slog.LevelWarn {
			t.Log(append([]any{message}, attributes...)...)
		}
	})
	publisherCtx, stopPublisher := context.WithCancel(context.Background())
	publisherDone := make(chan struct{})
	go func() {
		defer close(publisherDone)
		publisher.Run(publisherCtx, 10*time.Second)
	}()
	defer func() {
		stopPublisher()
		<-publisherDone
	}()
	waitForLoadQueue(t, settings)

	connection, err := amqp.Dial(brokerURL)
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		t.Fatal(err)
	}
	if err := channel.Qos(settings.Prefetch, 0, false); err != nil {
		t.Fatal(err)
	}
	deliveries, err := channel.Consume(settings.Queue, "load-test", false, false, false, false, nil)
	if err != nil {
		t.Fatal(err)
	}
	var consumed, consumedBytes atomic.Int64
	consumerDone := make(chan struct{})
	go func() {
		defer close(consumerDone)
		count := 0
		for delivery := range deliveries {
			consumed.Add(1)
			consumedBytes.Add(int64(len(delivery.Body)))
			if count++; count%200 == 0 {
				_ = delivery.Ack(true)
			}
		}
	}()

	live, telemetry := loadReading(telemetryLoadBaseMachine, 0, time.Now().UTC(), 0, 1, 1)
	body, err := json.Marshal(model.NewTelemetryReading(telemetryLoadBaseMachine, time.Now().UTC(), model.TelemetrySourceCNC, live, telemetry))
	if err != nil {
		t.Fatal(err)
	}
	started := time.Now()
	deadline := started.Add(time.Duration(seconds) * time.Second)
	var accepted, refused int64
	for time.Now().Before(deadline) {
		if publisher.Publish(body) {
			accepted++
		} else {
			refused++
			runtime.Gosched()
		}
	}
	sendTime := time.Since(started)
	confirmedInTime := publisher.Stats().Confirmed
	consumedInTime := consumed.Load()
	drainDeadline := time.Now().Add(time.Minute)
	for consumed.Load() < accepted && time.Now().Before(drainDeadline) {
		time.Sleep(100 * time.Millisecond)
	}
	channel.Cancel("load-test", false)
	<-consumerDone

	perSecond := float64(confirmedInTime) / sendTime.Seconds()
	t.Logf("broker only: message %d bytes (40 telemetry + 40 live values) | confirmed %.0f msg/s = %.1f MB/s = %.0f rows/s | consumed %.0f msg/s during the burst, %d of %d after drain | %d publishes refused while the buffer was full | ≈ %.0f machines at 1 reading/s (worst case)",
		len(body), perSecond, perSecond*float64(len(body))/1e6, perSecond*2*telemetryLoadVariableCount,
		float64(consumedInTime)/sendTime.Seconds(), consumed.Load(), accepted, refused, perSecond)
}

// envRatio reads a share between 0 and 1 from the environment, else fallback.
func envRatio(key string, fallback float64) float64 {
	value, err := strconv.ParseFloat(os.Getenv(key), 64)
	if err != nil || value < 0 || value > 1 {
		return fallback
	}
	return value
}

// envInts reads a comma-separated list of positive numbers, else fallback.
func envInts(key string, fallback []int) []int {
	var values []int
	for _, part := range strings.Split(os.Getenv(key), ",") {
		if value, err := strconv.Atoi(strings.TrimSpace(part)); err == nil && value > 0 {
			values = append(values, value)
		}
	}
	if len(values) == 0 {
		return fallback
	}
	return values
}
