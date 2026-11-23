package monitor

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/queue"
	"7th-sense-server/internal/response"
)

// fakeSession records what the writer writes and how it settles each message.
type fakeSession struct {
	mu         sync.Mutex
	writes     [][]model.TelemetryValue
	lives      [][]model.LiveValue
	acks       []string // "tag" or "tag+" (multiple)
	rejected   []uint64
	requeued   []uint64
	failWrites int                                         // this many writes fail with a database error first
	refuse     func(telemetry []model.TelemetryValue) bool // a write it returns true for is refused as invalid
}

func (session *fakeSession) actions() batchActions {
	return batchActions{
		write: func(ctx context.Context, telemetry []model.TelemetryValue, live []model.LiveValue) error {
			session.mu.Lock()
			defer session.mu.Unlock()
			if session.failWrites > 0 {
				session.failWrites--
				return errors.New("connection refused")
			}
			if session.refuse != nil && session.refuse(telemetry) {
				return response.InvalidTelemetry
			}
			session.writes = append(session.writes, telemetry)
			session.lives = append(session.lives, live)
			return nil
		},
		ack: func(tag uint64, multiple bool) error {
			session.mu.Lock()
			defer session.mu.Unlock()
			mark := ""
			if multiple {
				mark = "+"
			}
			session.acks = append(session.acks, string(rune('0'+tag))+mark)
			return nil
		},
		reject: func(tag uint64) error {
			session.mu.Lock()
			defer session.mu.Unlock()
			session.rejected = append(session.rejected, tag)
			return nil
		},
		requeue: func(tag uint64, multiple bool) error {
			session.mu.Lock()
			defer session.mu.Unlock()
			session.requeued = append(session.requeued, tag)
			return nil
		},
	}
}

func (session *fakeSession) snapshot() (writes int, acks []string, rejected []uint64) {
	session.mu.Lock()
	defer session.mu.Unlock()
	return len(session.writes), append([]string(nil), session.acks...), append([]uint64(nil), session.rejected...)
}

func newTestWriter(batchRows int, batchInterval time.Duration) *telemetryWriter {
	telemetryQueue := &TelemetryQueue{
		settings:    config.RabbitMQConfig{Queue: "telemetry.values", BatchRows: batchRows, BatchInterval: batchInterval, ReconnectMax: time.Second},
		stopTimeout: time.Second,
	}
	return &telemetryWriter{queue: telemetryQueue, name: "test"}
}

// testDelivery is a delivery of one reading of machineID with one live and one telemetry row
// of variableID.
func testDelivery(t *testing.T, tag uint64, machineID int64, variableID int, value float64, readAt time.Time) queue.Delivery {
	t.Helper()
	reading := model.NewTelemetryReading(machineID, readAt, model.TelemetrySourceCNC,
		[]model.LiveValue{{MachineID: machineID, VariableID: variableID, Value: &value}},
		[]model.TelemetryValue{{MachineID: machineID, VariableID: variableID, Value: &value}})
	body, err := json.Marshal(reading)
	if err != nil {
		t.Fatal(err)
	}
	return queue.Delivery{Body: body, Tag: tag}
}

// runWriter runs the writer on deliveries until stop is called; it returns stop.
func runWriter(writer *telemetryWriter, session *fakeSession, deliveries chan queue.Delivery) (stop func()) {
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		defer close(done)
		writer.consume(ctx, deliveries, session.actions())
	}()
	return func() {
		cancel()
		<-done
	}
}

func waitFor(t *testing.T, what string, condition func() bool) {
	t.Helper()
	deadline := time.Now().Add(3 * time.Second)
	for !condition() {
		if time.Now().After(deadline) {
			t.Fatalf("timed out waiting for %s", what)
		}
		time.Sleep(5 * time.Millisecond)
	}
}

func TestWriterFlushesBySize(t *testing.T) {
	writer := newTestWriter(4, time.Hour) // each reading is 2 rows: the second fills the batch
	session := &fakeSession{}
	deliveries := make(chan queue.Delivery, 10)
	stop := runWriter(writer, session, deliveries)
	defer stop()

	now := time.Now()
	deliveries <- testDelivery(t, 1, 1, 10, 1, now)
	deliveries <- testDelivery(t, 2, 1, 11, 2, now)
	waitFor(t, "a batch by size", func() bool { writes, _, _ := session.snapshot(); return writes == 1 })
	_, acks, _ := session.snapshot()
	if len(acks) != 1 || acks[0] != "2+" {
		t.Fatalf("acks = %v, want [2+] (every message up to the last)", acks)
	}
	if len(session.writes[0]) != 2 {
		t.Fatalf("telemetry rows written = %d, want 2", len(session.writes[0]))
	}
}

func TestWriterFlushesByTime(t *testing.T) {
	writer := newTestWriter(5000, 20*time.Millisecond)
	session := &fakeSession{}
	deliveries := make(chan queue.Delivery, 10)
	stop := runWriter(writer, session, deliveries)
	defer stop()

	deliveries <- testDelivery(t, 1, 1, 10, 1, time.Now())
	waitFor(t, "a batch by time", func() bool { writes, _, _ := session.snapshot(); return writes == 1 })
}

func TestWriterKeepsNewestLiveValue(t *testing.T) {
	newer := time.Now()
	older := newer.Add(-time.Second)
	batch := telemetryBatch{}
	for tag, delivery := range []queue.Delivery{testDelivery(t, 1, 7, 10, 2, newer), testDelivery(t, 2, 7, 10, 1, older)} {
		reading, err := decodeReading(delivery.Body)
		if err != nil {
			t.Fatal(err)
		}
		batch.add(reading, uint64(tag+1))
	}
	live := batch.liveRows()
	if len(live) != 1 || *live[0].Value != 2 || !live[0].UpdatedAt.Equal(newer) {
		t.Fatalf("live rows = %+v, want one row with the newer value 2", live)
	}
	if rows := batch.telemetryRows(); len(rows) != 2 {
		t.Fatalf("telemetry rows = %d, want both readings' rows", len(rows))
	}
}

func TestWriterAcksOnlyAfterSuccess(t *testing.T) {
	writer := newTestWriter(1, time.Hour)
	session := &fakeSession{failWrites: 1}
	deliveries := make(chan queue.Delivery, 10)
	stop := runWriter(writer, session, deliveries)
	defer stop()

	deliveries <- testDelivery(t, 1, 1, 10, 1, time.Now())
	time.Sleep(100 * time.Millisecond) // the first write failed; the retry waits about 1 s
	if _, acks, _ := session.snapshot(); len(acks) != 0 {
		t.Fatalf("acked %v before the batch was written", acks)
	}
	waitFor(t, "the retried write", func() bool { writes, acks, _ := session.snapshot(); return writes == 1 && len(acks) == 1 })
	if session.requeued != nil {
		t.Fatalf("requeued %v: a database error must keep the batch, not requeue it", session.requeued)
	}
}

func TestWriterDeadLettersRefusedAndUndecodableReadings(t *testing.T) {
	writer := newTestWriter(6, time.Hour)
	session := &fakeSession{refuse: func(telemetry []model.TelemetryValue) bool {
		for _, row := range telemetry {
			if row.MachineID == 99 {
				return true
			}
		}
		return false
	}}
	deliveries := make(chan queue.Delivery, 10)
	stop := runWriter(writer, session, deliveries)
	defer stop()

	now := time.Now()
	deliveries <- queue.Delivery{Body: []byte("not json"), Tag: 1}
	deliveries <- testDelivery(t, 2, 1, 10, 1, now)
	deliveries <- testDelivery(t, 3, 99, 10, 1, now) // refused by the database
	deliveries <- testDelivery(t, 4, 1, 11, 1, now)  // fills the batch (6 rows)
	waitFor(t, "the batch written reading by reading", func() bool { writes, _, _ := session.snapshot(); return writes == 2 })
	_, acks, rejected := session.snapshot()
	if len(rejected) != 2 || rejected[0] != 1 || rejected[1] != 3 {
		t.Fatalf("rejected = %v, want [1 3]", rejected)
	}
	if len(acks) != 2 || acks[0] != "2" || acks[1] != "4" {
		t.Fatalf("acks = %v, want [2 4] one by one", acks)
	}
}

func TestWriterWritesLastBatchAtStop(t *testing.T) {
	writer := newTestWriter(5000, time.Hour)
	session := &fakeSession{}
	deliveries := make(chan queue.Delivery, 10)
	stop := runWriter(writer, session, deliveries)

	deliveries <- testDelivery(t, 1, 1, 10, 1, time.Now())
	time.Sleep(50 * time.Millisecond)
	stop()
	if writes, acks, _ := session.snapshot(); writes != 1 || len(acks) != 1 {
		t.Fatalf("writes %d, acks %v: the last batch must be written and acked at stop", writes, acks)
	}
}
