package monitor

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"sync"
	"time"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/queue"
	"7th-sense-server/internal/response"
)

// telemetryWriter is one consumer of the telemetry queue (RABBITMQ_CONSUMERS of them): it
// collects delivered readings until the batch holds BatchRows rows (RABBITMQ_BATCH_ROWS) or
// BatchInterval (RABBITMQ_BATCH_MS) passed since its first message, writes the whole batch in one
// machine_telementry COPY and one machine_variable_values upsert (through
// MachineTelemetryService), and only then acknowledges its messages. With LiveEveryN > 1 the
// upsert runs once every N batches: the newest live value of each machine-variable is kept in
// pendingLive meanwhile, so the live table ends up the same, refreshed less often.
type telemetryWriter struct {
	queue *TelemetryQueue
	name  string

	pendingLive map[liveKey]model.LiveValue // newest live row per machine-variable not yet upserted
	liveWait    int                         // batches since the last live upsert
	liveFlush   bool                        // this flush carries the pending live rows
	flushLive   bool                        // set at stop: the pending live rows go out with the last batch
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

	// With LiveEveryN > 1 the pending live rows also go out on this timer: batches may have
	// stopped arriving, and live values must not wait for the next one (or for stop).
	var liveTicker *time.Ticker
	var liveDue <-chan time.Time
	if settings.LiveEveryN > 1 {
		liveTicker = time.NewTicker(settings.BatchInterval * time.Duration(settings.LiveEveryN))
		liveDue = liveTicker.C
		defer liveTicker.Stop()
	}

	for {
		select {
		case <-ctx.Done():
			writer.flushAtStop(&batch, actions)
			return
		case <-liveDue:
			writer.flushPendingLive(ctx, actions)
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

// flushAtStop writes the last batch with a deadline of its own: the writer's ctx has ended. The
// pending live rows go with it, so a stop never leaves live values unwritten.
func (writer *telemetryWriter) flushAtStop(batch *telemetryBatch, actions batchActions) {
	writer.flushLive = true
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
	// The live rows are decided once, before the retries: a failed write keeps them pending.
	live := writer.liveRowsForWrite(batch.liveRows())
	defer batch.reset()
	if len(batch.entries) == 0 {
		// Only reachable at stop with live rows still pending (flushLive): write them alone.
		if len(live) == 0 {
			return true
		}
		if err := actions.write(ctx, nil, live); err != nil {
			writer.queue.stats.addWriteError()
			writer.queue.log(slog.LevelError, "queue: pending live values could not be written at stop",
				"consumer", writer.name, "rows", len(live), "error", err)
			return true
		}
		writer.liveWritten()
		return true
	}
	telemetry := batch.telemetryRows()
	retry := backoff{max: writer.queue.settings.ReconnectMax}
	for {
		started := time.Now()
		err := actions.write(ctx, telemetry, live)
		if err == nil {
			writer.liveWritten()
			writer.queue.stats.addBatch(len(batch.entries), len(telemetry)+len(live), time.Since(started))
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

// gzip writers and readers are reused: a new compressor allocates about a megabyte, and the queue
// encodes one message per reading (1,000 a second at 40,000 values/s).
var (
	gzipWriters = sync.Pool{New: func() any {
		writer, _ := gzip.NewWriterLevel(nil, gzip.BestSpeed) // BestSpeed is a valid level: no error
		return writer
	}}
	gzipReaders sync.Pool // *gzip.Reader; gzip.NewReader needs a stream, so the pool starts empty
)

// EncodeReading encodes one queue message (ADR 0030): the reading as JSON, gzip-compressed. The
// JSON repeats the same keys for every value, so gzip makes it about 6 times smaller, and RabbitMQ
// writes every message to disk twice (quorum log and segment file).
func EncodeReading(reading model.TelemetryReading) ([]byte, error) {
	var body bytes.Buffer
	writer := gzipWriters.Get().(*gzip.Writer)
	defer gzipWriters.Put(writer)
	writer.Reset(&body)
	if err := json.NewEncoder(writer).Encode(reading); err != nil {
		return nil, fmt.Errorf("encode reading: %w", err)
	}
	if err := writer.Close(); err != nil {
		return nil, fmt.Errorf("encode reading: %w", err)
	}
	return body.Bytes(), nil
}

// decodeReading reads one queue message: gzip-compressed JSON (EncodeReading), or plain JSON, as
// messages queued before ADR 0030 are.
func decodeReading(body []byte) (model.TelemetryReading, error) {
	var reading model.TelemetryReading
	if isGzip(body) {
		var err error
		if body, err = gunzip(body); err != nil {
			return reading, fmt.Errorf("decode reading: %w", err)
		}
	}
	if err := json.Unmarshal(body, &reading); err != nil {
		return reading, fmt.Errorf("decode reading: %w", err)
	}
	if reading.MachineID <= 0 || reading.ReadAt.IsZero() {
		return reading, errors.New("decode reading: no machineId or readAt")
	}
	return reading, nil
}

// isGzip reports whether body starts with the gzip magic number (JSON never does).
func isGzip(body []byte) bool {
	return len(body) >= 2 && body[0] == 0x1f && body[1] == 0x8b
}

// gunzip returns the decompressed body.
func gunzip(body []byte) ([]byte, error) {
	reader, _ := gzipReaders.Get().(*gzip.Reader)
	var err error
	if reader == nil {
		reader, err = gzip.NewReader(bytes.NewReader(body))
	} else {
		err = reader.Reset(bytes.NewReader(body))
	}
	if err != nil {
		return nil, err
	}
	defer gzipReaders.Put(reader)
	return io.ReadAll(reader)
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

// liveKey names one machine-variable of the live-values table.
type liveKey struct {
	machineID  int64
	variableID int
}

// liveRowsForWrite returns the live rows this flush upserts. With LiveEveryN <= 1 it is the
// batch's own live rows, as before. With LiveEveryN > 1 the batch's rows are merged into
// pendingLive (newest per machine-variable wins) and the pending set is returned only every
// LiveEveryN-th batch or at stop (flushLive); liveWritten drops it after the write succeeds, so
// a failed or requeued batch never loses live rows.
func (writer *telemetryWriter) liveRowsForWrite(rows []model.LiveValue) []model.LiveValue {
	if writer.queue.settings.LiveEveryN <= 1 {
		return rows
	}
	if writer.pendingLive == nil {
		writer.pendingLive = map[liveKey]model.LiveValue{}
	}
	for _, row := range rows {
		key := liveKey{row.MachineID, row.VariableID}
		if kept, found := writer.pendingLive[key]; !found || !row.UpdatedAt.Before(kept.UpdatedAt) {
			writer.pendingLive[key] = row
		}
	}
	writer.liveWait++
	if !writer.flushLive && writer.liveWait < writer.queue.settings.LiveEveryN {
		writer.liveFlush = false
		return nil
	}
	writer.liveFlush = true
	pending := make([]model.LiveValue, 0, len(writer.pendingLive))
	for _, row := range writer.pendingLive {
		pending = append(pending, row)
	}
	return pending
}

// liveWritten drops the pending live rows after their upsert succeeded and starts the count over.
func (writer *telemetryWriter) liveWritten() {
	if !writer.liveFlush {
		return
	}
	writer.pendingLive = map[liveKey]model.LiveValue{}
	writer.liveWait, writer.liveFlush, writer.flushLive = 0, false, false
}

// flushPendingLive writes the accumulated live rows on their own, without a batch: it keeps live
// values fresh while deliveries are quiet. A failed write keeps them pending for the next tick.
func (writer *telemetryWriter) flushPendingLive(ctx context.Context, actions batchActions) {
	if len(writer.pendingLive) == 0 {
		return
	}
	writer.liveFlush = true
	pending := make([]model.LiveValue, 0, len(writer.pendingLive))
	for _, row := range writer.pendingLive {
		pending = append(pending, row)
	}
	started := time.Now()
	if err := actions.write(ctx, nil, pending); err != nil {
		writer.liveFlush = false
		writer.queue.stats.addWriteError()
		writer.queue.log(slog.LevelError, "queue: live values could not be written, keeping them pending",
			"consumer", writer.name, "rows", len(pending), "error", err)
		return
	}
	writer.liveWritten()
	writer.queue.stats.addBatch(0, len(pending), time.Since(started))
}

// liveRows returns the newest machine_variable_values row of each machine-variable of the batch.
func (batch *telemetryBatch) liveRows() []model.LiveValue {
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
