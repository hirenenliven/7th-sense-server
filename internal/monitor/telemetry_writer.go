package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/queue"
	"7th-sense-server/internal/response"
)

// telemetryWriter is one consumer of the telemetry queue (RABBITMQ_CONSUMERS of them): it
// collects delivered readings until the batch holds BatchRows rows (5000) or BatchInterval
// (200 ms) passed since its first message, writes the whole batch in one
// machine_telementry insert and one machine_variable_values upsert (through
// MachineTelemetryService), and only then acknowledges its messages.
type telemetryWriter struct {
	queue *TelemetryQueue
	name  string
}

// batchActions are what a batch is written and settled with: the database write and the
// session's acknowledgements (plain functions, so the batching rules are tested without a broker).
type batchActions struct {
	write   func(ctx context.Context, telemetry []model.TelemetryValue, live []model.LiveValue) error
	ack     func(tag uint64, multiple bool) error
	reject  func(tag uint64) error
	requeue func(tag uint64, multiple bool) error
}

// consumeSession consumes one broker session until it is lost or the writer stops.
func (writer *telemetryWriter) consumeSession(ctx context.Context, session *queue.Session) {
	writer.consume(ctx, session.Deliveries(), batchActions{
		write: writer.queue.writeRows, ack: session.Ack, reject: session.Reject, requeue: session.Requeue,
	})
}

// consume batches deliveries until the channel closes (connection lost: the messages not yet
// acknowledged are delivered again) or ctx ends (the current batch is written first).
func (writer *telemetryWriter) consume(ctx context.Context, deliveries <-chan queue.Delivery, actions batchActions) {
	settings := writer.queue.settings
	batch := telemetryBatch{}
	var flushTimer *time.Timer
	var flushDue <-chan time.Time
	stopTimer := func() {
		if flushTimer != nil {
			flushTimer.Stop()
			flushTimer, flushDue = nil, nil
		}
	}
	defer stopTimer()

	for {
		select {
		case <-ctx.Done():
			writer.flushAtStop(&batch, actions)
			return
		case delivery, ok := <-deliveries:
			if !ok {
				return
			}
			reading, err := decodeReading(delivery.Body)
			if err != nil {
				writer.deadLetter(delivery.Tag, actions, "queue: message cannot be decoded, sent to the dead-letter queue", err)
				continue
			}
			batch.add(reading, delivery.Tag)
			if flushDue == nil {
				flushTimer = time.NewTimer(settings.BatchInterval)
				flushDue = flushTimer.C
			}
			if batch.rows >= settings.BatchRows {
				stopTimer()
				if !writer.flush(ctx, &batch, actions) {
					return
				}
			}
		case <-flushDue:
			stopTimer()
			if !writer.flush(ctx, &batch, actions) {
				return
			}
		}
	}
}

// flushAtStop writes the last batch with a deadline of its own: the writer's ctx has ended.
func (writer *telemetryWriter) flushAtStop(batch *telemetryBatch, actions batchActions) {
	stopCtx, cancel := context.WithTimeout(context.Background(), writer.queue.stopTimeout)
	defer cancel()
	writer.flush(stopCtx, batch, actions)
}

// flush writes the batch and acknowledges its messages, then empties it. A database error (not a
// refused value) keeps the messages unacknowledged and writes the same batch again after a wait
// (1 s doubling to ReconnectMax, 60 s), so they wait in RabbitMQ while the database is
// down without counting redeliveries. When ctx ends first the messages are handed back to the
// queue and flush returns false.
func (writer *telemetryWriter) flush(ctx context.Context, batch *telemetryBatch, actions batchActions) bool {
	if len(batch.entries) == 0 {
		return true
	}
	defer batch.reset()
	retry := backoff{max: writer.queue.settings.ReconnectMax}
	for {
		started := time.Now()
		err := actions.write(ctx, batch.telemetryRows(), batch.liveRows())
		if err == nil {
			writer.queue.stats.addBatch(len(batch.entries), batch.rows, time.Since(started))
			writer.settle(actions.ack(batch.lastTag(), true), "ack")
			return true
		}
		if errors.Is(err, response.InvalidTelemetry) {
			return writer.writeEachReading(ctx, batch, actions)
		}
		writer.queue.stats.addWriteError()
		if ctx.Err() != nil {
			writer.settle(actions.requeue(batch.lastTag(), true), "requeue")
			return false
		}
		wait := retry.nextWait()
		writer.queue.log(slog.LevelError, "queue: writing a batch failed, keeping it and trying again", "consumer", writer.name,
			"messages", len(batch.entries), "rows", batch.rows, "retryIn", wait.String(), "error", err)
		sleepUnlessCancelled(ctx, wait)
	}
}

// writeEachReading writes a batch PostgreSQL refused one reading at a time: the readings it
// accepts are acknowledged, the refused ones go to the dead-letter queue. On any other error the
// readings left are handed back to the queue.
func (writer *telemetryWriter) writeEachReading(ctx context.Context, batch *telemetryBatch, actions batchActions) bool {
	for index, entry := range batch.entries {
		single := telemetryBatch{}
		single.add(entry.reading, entry.tag)
		err := actions.write(ctx, single.telemetryRows(), single.liveRows())
		switch {
		case err == nil:
			writer.queue.stats.addBatch(1, single.rows, 0)
			writer.settle(actions.ack(entry.tag, false), "ack")
		case errors.Is(err, response.InvalidTelemetry):
			writer.deadLetter(entry.tag, actions, "queue: reading refused by the database, sent to the dead-letter queue",
				fmt.Errorf("machine %d at %s: %w", entry.reading.MachineID, entry.reading.ReadAt.Format(time.RFC3339Nano), err))
		default:
			writer.queue.stats.addWriteError()
			for _, left := range batch.entries[index:] {
				writer.settle(actions.requeue(left.tag, false), "requeue")
			}
			return ctx.Err() == nil
		}
	}
	return true
}

// deadLetter sends message tag to the dead-letter queue and logs why.
func (writer *telemetryWriter) deadLetter(tag uint64, actions batchActions, message string, reason error) {
	writer.queue.stats.addDeadLetter()
	writer.settle(actions.reject(tag), "reject")
	writer.queue.log(slog.LevelError, message, "consumer", writer.name, "deadLetterQueue",
		queue.DeadLetterQueue(writer.queue.settings), "error", reason)
}

// settle logs a failed acknowledgement. The connection is gone then, and the broker delivers the
// messages again; the writes are idempotent, so they are stored once.
func (writer *telemetryWriter) settle(err error, what string) {
	if err != nil {
		writer.queue.log(slog.LevelWarn, "queue: "+what+" failed, the messages will be delivered again", "consumer", writer.name,
			"error", err)
	}
}

// decodeReading reads one queue message.
func decodeReading(body []byte) (model.TelemetryReading, error) {
	var reading model.TelemetryReading
	if err := json.Unmarshal(body, &reading); err != nil {
		return reading, fmt.Errorf("decode reading: %w", err)
	}
	if reading.MachineID <= 0 || reading.ReadAt.IsZero() {
		return reading, errors.New("decode reading: no machineId or readAt")
	}
	return reading, nil
}

// batchEntry is one delivered reading of a batch.
type batchEntry struct {
	reading model.TelemetryReading
	tag     uint64
}

// telemetryBatch collects readings until they are written together.
type telemetryBatch struct {
	entries []batchEntry
	rows    int
}

// add appends reading, delivered with tag.
func (batch *telemetryBatch) add(reading model.TelemetryReading, tag uint64) {
	batch.entries = append(batch.entries, batchEntry{reading: reading, tag: tag})
	batch.rows += reading.RowCount()
}

// lastTag is the tag of the newest delivery: acknowledging up to it settles the whole batch.
func (batch *telemetryBatch) lastTag() uint64 {
	return batch.entries[len(batch.entries)-1].tag
}

// telemetryRows returns every machine_telementry row of the batch.
func (batch *telemetryBatch) telemetryRows() []model.TelemetryValue {
	var rows []model.TelemetryValue
	for _, entry := range batch.entries {
		rows = append(rows, entry.reading.TelemetryValues()...)
	}
	return rows
}

// liveRows returns the newest machine_variable_values row of each machine-variable of the batch.
func (batch *telemetryBatch) liveRows() []model.LiveValue {
	type liveKey struct {
		machineID  int64
		variableID int
	}
	newest := map[liveKey]int{}
	var rows []model.LiveValue
	for _, entry := range batch.entries {
		for _, row := range entry.reading.LiveValues() {
			key := liveKey{row.MachineID, row.VariableID}
			index, found := newest[key]
			switch {
			case !found:
				newest[key] = len(rows)
				rows = append(rows, row)
			case !row.UpdatedAt.Before(rows[index].UpdatedAt):
				rows[index] = row
			}
		}
	}
	return rows
}

// reset empties the batch.
func (batch *telemetryBatch) reset() {
	batch.entries, batch.rows = nil, 0
}
