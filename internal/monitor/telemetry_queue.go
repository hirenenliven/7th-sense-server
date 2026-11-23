package monitor

import (
	"context"
	"encoding/json"
	"log/slog"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/queue"
	"7th-sense-server/internal/service"
)

const (
	// queueStatsInterval is how often the queue's statistics are logged.
	queueStatsInterval = time.Minute
	// queueStopMargin is added to CNC_STOP_TIMEOUT_SEC for the writer to finish its last batch.
	queueStopMargin = 5 * time.Second
)

// TelemetryQueue hands the collectors' high-volume rows (machine_telementry,
// machine_variable_values) to RabbitMQ (ADR 0024) and runs the writers that store them in
// batches. Both the CNC and the PLC workers use it; a nil
// TelemetryQueue means RABBITMQ_ENABLED=false, and the workers write directly as before. Modes,
// events, cycles and statuses never go through it.
type TelemetryQueue struct {
	settings    config.RabbitMQConfig
	logs        config.CNCLogSettings
	stopTimeout time.Duration
	service     *service.MachineTelemetryService
	publisher   *queue.Publisher
	stats       writerStats

	// fallbackLogged: the first refusal after a quiet minute is logged at once, the rest are
	// counted in the statistics line.
	fallbackLogged atomic.Bool

	cancelPublisher context.CancelFunc
	cancelWriters   context.CancelFunc
	publisherDone   chan struct{}
	writersDone     sync.WaitGroup
	statsDone       chan struct{}
}

// NewTelemetryQueue creates the telemetry queue of rabbitConfig; Start runs it. cncConfig gives
// the log switches and the stop timeout.
func NewTelemetryQueue(rabbitConfig config.RabbitMQConfig, cncConfig config.CNCConfig, telemetryService *service.MachineTelemetryService) *TelemetryQueue {
	telemetryQueue := &TelemetryQueue{settings: rabbitConfig, logs: cncConfig.Logs, stopTimeout: cncConfig.StopTimeout,
		service: telemetryService}
	telemetryQueue.publisher = queue.NewPublisher(rabbitConfig, telemetryQueue.log)
	return telemetryQueue
}

// Start runs the publisher, the writers (RABBITMQ_CONSUMERS) and the statistics log in background
// goroutines and returns at once: a broker that is down does not stop the server, the publisher
// keeps reconnecting and the workers write directly meanwhile.
func (telemetryQueue *TelemetryQueue) Start() {
	publisherCtx, cancelPublisher := context.WithCancel(context.Background())
	writersCtx, cancelWriters := context.WithCancel(context.Background())
	telemetryQueue.cancelPublisher, telemetryQueue.cancelWriters = cancelPublisher, cancelWriters
	telemetryQueue.publisherDone, telemetryQueue.statsDone = make(chan struct{}), make(chan struct{})

	go func() {
		defer close(telemetryQueue.publisherDone)
		unsent := telemetryQueue.publisher.Run(publisherCtx, telemetryQueue.stopTimeout)
		telemetryQueue.storeUnsent(unsent)
	}()
	for index := range telemetryQueue.settings.Consumers {
		writer := &telemetryWriter{queue: telemetryQueue, name: strconv.Itoa(index + 1)}
		consumer := queue.NewConsumer(telemetryQueue.settings, writer.name, telemetryQueue.log)
		telemetryQueue.writersDone.Add(1)
		go func() {
			defer telemetryQueue.writersDone.Done()
			consumer.Run(writersCtx, writer.consumeSession)
		}()
	}
	go func() {
		defer close(telemetryQueue.statsDone)
		telemetryQueue.logStatsEveryMinute(writersCtx)
	}()
	telemetryQueue.log(slog.LevelInfo, "queue: telemetry queue started", telemetryQueue.SettingsLogAttributes()...)
}

// Stop runs after the collectors stopped: the publisher sends what is left (what it cannot send
// is written directly), then the writers store and acknowledge their current batch, then the
// connections close. Bounded by CNC_STOP_TIMEOUT_SEC plus a margin per step.
func (telemetryQueue *TelemetryQueue) Stop() {
	if telemetryQueue.cancelPublisher == nil {
		return
	}
	telemetryQueue.cancelPublisher()
	telemetryQueue.waitUntil(telemetryQueue.publisherDone, "publisher")

	telemetryQueue.cancelWriters()
	writersDone := make(chan struct{})
	go func() {
		telemetryQueue.writersDone.Wait()
		close(writersDone)
	}()
	telemetryQueue.waitUntil(writersDone, "writers")
	<-telemetryQueue.statsDone
	telemetryQueue.logStats()
	telemetryQueue.log(slog.LevelInfo, "queue: telemetry queue stopped")
}

// waitUntil waits for done at most CNC_STOP_TIMEOUT_SEC plus a margin.
func (telemetryQueue *TelemetryQueue) waitUntil(done <-chan struct{}, what string) {
	select {
	case <-done:
	case <-time.After(telemetryQueue.stopTimeout + queueStopMargin):
		telemetryQueue.log(slog.LevelError, "queue: "+what+" did not stop in time")
	}
}

// Enqueue hands one reading to the publisher without waiting. It returns false when the queue is
// off (nil), the reading cannot be encoded (a NaN or infinite number) or the publish buffer is
// full: the worker then writes the rows directly, so nothing is dropped.
func (telemetryQueue *TelemetryQueue) Enqueue(reading model.TelemetryReading) bool {
	if telemetryQueue == nil {
		return false
	}
	body, err := json.Marshal(reading)
	if err != nil {
		telemetryQueue.logFallback("queue: reading cannot be encoded, writing it directly", reading, err)
		return false
	}
	if !telemetryQueue.publisher.Publish(body) {
		telemetryQueue.logFallback("queue: publish buffer full or broker unavailable, writing directly", reading, nil)
		return false
	}
	return true
}

// logFallback logs the first fallback of a minute; the statistics line counts every one.
func (telemetryQueue *TelemetryQueue) logFallback(message string, reading model.TelemetryReading, err error) {
	if telemetryQueue.fallbackLogged.Swap(true) {
		return
	}
	attributes := []any{"machineId", reading.MachineID, "source", reading.Source, "buffer", telemetryQueue.settings.PublishBuffer}
	if err != nil {
		attributes = append(attributes, "error", err)
	}
	telemetryQueue.log(slog.LevelWarn, message, attributes...)
}

// storeUnsent writes directly the messages the publisher could not send before it stopped.
func (telemetryQueue *TelemetryQueue) storeUnsent(unsent [][]byte) {
	if len(unsent) == 0 {
		return
	}
	batch := telemetryBatch{}
	for _, body := range unsent {
		reading, err := decodeReading(body)
		if err != nil {
			continue // only bodies Enqueue encoded reach the publisher
		}
		batch.add(reading, 0)
	}
	ctx, cancel := context.WithTimeout(context.Background(), telemetryQueue.stopTimeout)
	defer cancel()
	if err := telemetryQueue.writeRows(ctx, batch.telemetryRows(), batch.liveRows()); err != nil {
		telemetryQueue.log(slog.LevelError, "queue: unsent readings could not be written at stop", "messages", len(unsent),
			"rows", batch.rows, "error", err)
		return
	}
	telemetryQueue.log(slog.LevelInfo, "queue: unsent readings written directly at stop", "messages", len(unsent), "rows", batch.rows)
}

// writeRows stores the rows of a batch: time series first, then the live values.
func (telemetryQueue *TelemetryQueue) writeRows(ctx context.Context, telemetry []model.TelemetryValue, live []model.LiveValue) error {
	if err := telemetryQueue.service.InsertTelemetry(ctx, telemetry); err != nil {
		return err
	}
	return telemetryQueue.service.UpsertLiveValues(ctx, live)
}

// logStatsEveryMinute logs the statistics once a minute until ctx ends.
func (telemetryQueue *TelemetryQueue) logStatsEveryMinute(ctx context.Context) {
	ticker := time.NewTicker(queueStatsInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			telemetryQueue.logStats()
		}
	}
}

// logStats logs the publisher's counters, the writers' work since the last line and the queue
// depth. A growing depth means the writers are slower than the collectors.
func (telemetryQueue *TelemetryQueue) logStats() {
	publisher := telemetryQueue.publisher.Stats()
	writers := telemetryQueue.stats.takeWindow()
	attributes := []any{
		"connected", publisher.Connected, "buffered", publisher.Buffered, "published", publisher.Published,
		"confirmed", publisher.Confirmed, "publishFailed", publisher.Failed, "fallback", publisher.Fallback,
		"reconnects", publisher.Reconnects,
		"messagesWritten", writers.messages, "rowsWritten", writers.rows, "batches", writers.batches,
		"averageBatchRows", writers.averageBatchRows(), "averageWriteMs", writers.averageWriteMs(),
		"writeErrors", writers.writeErrors, "deadLetters", writers.deadLetters,
	}
	if publisher.LastError != "" {
		attributes = append(attributes, "lastError", publisher.LastError)
	}
	if waiting, dead, err := queue.QueueDepth(telemetryQueue.settings); err == nil {
		attributes = append(attributes, "queueDepth", waiting, "deadLetterDepth", dead)
	} else {
		attributes = append(attributes, "queueDepthError", err.Error())
	}
	telemetryQueue.fallbackLogged.Store(false)
	telemetryQueue.log(slog.LevelInfo, "queue: statistics", attributes...)
}

// SettingsLogAttributes describes the queue's settings for the start logs (no password).
func (telemetryQueue *TelemetryQueue) SettingsLogAttributes() []any {
	if telemetryQueue == nil {
		return []any{"rabbitmq", false}
	}
	settings := telemetryQueue.settings
	return []any{"rabbitmq", true, "rabbitmqUrl", queue.SafeURL(settings.URL), "exchange", settings.Exchange,
		"routingKey", settings.RoutingKey, "queue", settings.Queue, "queueType", settings.QueueType,
		"consumers", settings.Consumers,
		"prefetch", settings.Prefetch, "batchRows", settings.BatchRows, "batchMs", settings.BatchInterval.Milliseconds(),
		"publishBuffer", settings.PublishBuffer}
}

// log writes a queue line (logType queue); it is the queue package's Logger.
func (telemetryQueue *TelemetryQueue) log(level slog.Level, message string, attributes ...any) {
	writeLog(telemetryQueue.logs, logTypeQueue, level, nil, message, attributes...)
}

// writerStats counts the writers' work; the statistics line reports and resets it.
type writerStats struct {
	mu     sync.Mutex
	counts writerCounts
}

// writerCounts is the writers' work over one statistics window.
type writerCounts struct {
	messages    int64
	rows        int64
	batches     int64
	writeTime   time.Duration
	writeErrors int64
	deadLetters int64
}

// addBatch counts one stored batch.
func (stats *writerStats) addBatch(messages, rows int, took time.Duration) {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.counts.messages += int64(messages)
	stats.counts.rows += int64(rows)
	stats.counts.batches++
	stats.counts.writeTime += took
}

// addWriteError counts a failed write (the batch is tried again).
func (stats *writerStats) addWriteError() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.counts.writeErrors++
}

// addDeadLetter counts a message sent to the dead-letter queue.
func (stats *writerStats) addDeadLetter() {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	stats.counts.deadLetters++
}

// takeWindow returns the counts since the last call and starts counting from zero.
func (stats *writerStats) takeWindow() writerCounts {
	stats.mu.Lock()
	defer stats.mu.Unlock()
	window := stats.counts
	stats.counts = writerCounts{}
	return window
}

func (counts writerCounts) averageBatchRows() float64 {
	if counts.batches == 0 {
		return 0
	}
	return roundToTenth(float64(counts.rows) / float64(counts.batches))
}

func (counts writerCounts) averageWriteMs() float64 {
	if counts.batches == 0 {
		return 0
	}
	return roundToTenth(float64(counts.writeTime.Microseconds()) / 1000 / float64(counts.batches))
}
