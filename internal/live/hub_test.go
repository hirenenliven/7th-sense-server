package live

import (
	"reflect"
	"testing"
	"time"

	"7th-sense-server/internal/model"
)

// sentMessage is one message the hub emitted.
type sentMessage struct {
	room    string
	event   string
	payload any
}

// newRecordingHub returns a hub with a 200 ms minimum interval and the messages it emits.
func newRecordingHub() (*Hub, *[]sentMessage) {
	hub := NewHub(200 * time.Millisecond)
	sent := &[]sentMessage{}
	hub.SetEmitter(func(room, event string, payload any) {
		*sent = append(*sent, sentMessage{room: room, event: event, payload: payload})
	})
	return hub, sent
}

func TestHubSendsOnlyTheRoomsCodes(t *testing.T) {
	hub, sent := newRecordingHub()
	start := time.Now()
	tile, _, _ := hub.JoinValues(225, []string{"SpindleLoad"}, time.Second)
	screen, _, _ := hub.JoinValues(225, []string{"axLoadZ", "SpindleLoad", "SpindleRpm"}, 200*time.Millisecond)
	health, _, _ := hub.JoinValues(225, []string{"SPBeTemp"}, time.Second)

	hub.PublishValues(225, start, map[string]any{"SpindleLoad": 44.5, "machPosX": 120.3, "axLoadZ": 38.0})
	hub.Flush(start.Add(time.Second))

	got := map[string]map[string]any{}
	for _, message := range *sent {
		got[message.room] = message.payload.(ValuesMessage).Values
	}
	want := map[string]map[string]any{
		tile:   {"SpindleLoad": 44.5},
		screen: {"SpindleLoad": 44.5, "axLoadZ": 38.0},
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("messages = %v, want %v (nothing for %s)", got, want, health)
	}
}

func TestHubKeepsNewestValueAndRespectsInterval(t *testing.T) {
	hub, sent := newRecordingHub()
	start := time.Now()
	room, _, _ := hub.JoinValues(225, []string{"SpindleLoad"}, time.Second)

	hub.PublishValues(225, start, map[string]any{"SpindleLoad": 43.0})
	hub.Flush(start.Add(time.Second))
	hub.PublishValues(225, start.Add(1100*time.Millisecond), map[string]any{"SpindleLoad": 44.0})
	hub.PublishValues(225, start.Add(1200*time.Millisecond), map[string]any{"SpindleLoad": 44.5})
	hub.Flush(start.Add(1500 * time.Millisecond)) // under 1 s since the last message: held back
	if len(*sent) != 1 {
		t.Fatalf("sent %d messages before the interval passed, want 1", len(*sent))
	}
	hub.Flush(start.Add(2 * time.Second))
	if len(*sent) != 2 {
		t.Fatalf("sent %d messages, want 2", len(*sent))
	}
	last := (*sent)[1]
	message := last.payload.(ValuesMessage)
	if last.room != room || message.Values["SpindleLoad"] != 44.5 || len(message.Values) != 1 {
		t.Fatalf("last message = %+v, want SpindleLoad 44.5 only", message)
	}
	if message.Seq != 3 {
		t.Fatalf("seq = %d, want 3 (one per changed reading)", message.Seq)
	}
}

func TestHubSkipsUnchangedValues(t *testing.T) {
	hub, sent := newRecordingHub()
	start := time.Now()
	hub.JoinValues(225, []string{"SpindleLoad", "progName"}, time.Second)
	hub.PublishValues(225, start, map[string]any{"SpindleLoad": 41.5, "progName": "O1234"})
	hub.Flush(start.Add(time.Second))
	hub.PublishValues(225, start.Add(time.Second), map[string]any{"SpindleLoad": 41.5, "progName": "O1234"})
	hub.Flush(start.Add(3 * time.Second))
	if len(*sent) != 1 {
		t.Fatalf("sent %d messages, want 1: an unchanged reading sends nothing", len(*sent))
	}
}

func TestHubSnapshotHasOnlyRequestedCodes(t *testing.T) {
	hub, _ := newRecordingHub()
	at := time.Now()
	hub.PublishValues(225, at, map[string]any{"SpindleLoad": 41.5, "machPosX": 12.0, "progName": "O1234"})
	hub.PublishStatus(225, Status{Online: true, ModeID: model.ModeInProduction, ModeSince: at, ProgramName: "O1234"})

	snapshot := hub.Snapshot(225, []string{"SpindleLoad", "progName", "toolNumber"})
	if !reflect.DeepEqual(snapshot.Values, map[string]any{"SpindleLoad": 41.5, "progName": "O1234"}) {
		t.Fatalf("snapshot values = %v", snapshot.Values)
	}
	if snapshot.Seq != 2 || snapshot.Status == nil || snapshot.ModeID != model.ModeInProduction {
		t.Fatalf("snapshot = %+v, want seq 2 with the status", snapshot)
	}
	if empty := hub.Snapshot(300, []string{"SpindleLoad"}); empty.Status != nil || len(empty.Values) != 0 || empty.At != nil {
		t.Fatalf("snapshot of an unread machine = %+v, want empty", empty)
	}
}

func TestHubJoinValuesSharesRoomsAndReportsUnknownCodes(t *testing.T) {
	hub, _ := newRecordingHub()
	first, kept, unknown := hub.JoinValues(225, []string{"SpindleRpm", "SpindleLoad", "SpindleLoad"}, 0)
	if !reflect.DeepEqual(kept, []string{"SpindleLoad", "SpindleRpm"}) || unknown != nil {
		t.Fatalf("kept %v unknown %v: codes are accepted while the machine's codes are not known", kept, unknown)
	}
	second, _, _ := hub.JoinValues(225, []string{"SpindleLoad", "SpindleRpm"}, time.Second)
	if first != second {
		t.Fatalf("rooms %q and %q differ: same machine, codes and interval share one room", first, second)
	}
	if other, _, _ := hub.JoinValues(226, []string{"SpindleLoad", "SpindleRpm"}, time.Second); other == first {
		t.Fatalf("another machine got room %q", other)
	}

	hub.SetMachineCodes(225, []string{"SpindleLoad"})
	room, kept, unknown := hub.JoinValues(225, []string{"SpindleLoad", "SPBeTemp"}, 50*time.Millisecond)
	if !reflect.DeepEqual(kept, []string{"SpindleLoad"}) || !reflect.DeepEqual(unknown, []string{"SPBeTemp"}) {
		t.Fatalf("kept %v unknown %v", kept, unknown)
	}
	if room != valuesRoomName(225, []string{"SpindleLoad"}, 200*time.Millisecond) {
		t.Fatalf("room %q: the interval is raised to the minimum", room)
	}
	if room, _, _ := hub.JoinValues(225, []string{"SPBeTemp"}, time.Second); room != "" {
		t.Fatalf("room %q for unknown codes only, want none", room)
	}
}

func TestHubLeaveValuesStopsBuildingTheRoom(t *testing.T) {
	hub, sent := newRecordingHub()
	start := time.Now()
	room, _, _ := hub.JoinValues(225, []string{"SpindleLoad"}, time.Second)
	hub.JoinValues(225, []string{"SpindleLoad"}, time.Second)
	hub.LeaveValues(room)
	hub.PublishValues(225, start, map[string]any{"SpindleLoad": 1.0})
	hub.Flush(start.Add(time.Second))
	if len(*sent) != 1 {
		t.Fatalf("sent %d messages, want 1 while one member is left", len(*sent))
	}
	hub.LeaveValues(room)
	hub.PublishValues(225, start.Add(time.Second), map[string]any{"SpindleLoad": 2.0})
	hub.Flush(start.Add(3 * time.Second))
	if len(*sent) != 1 {
		t.Fatalf("sent %d messages, want none for a room without members", len(*sent))
	}
}

func TestHubStatusAndEvents(t *testing.T) {
	hub, sent := newRecordingHub()
	at := time.Now()
	status := Status{Online: true, ModeID: model.ModeDowntime, ModeSince: at}
	hub.PublishStatus(225, status)
	hub.PublishStatus(225, status) // unchanged: not sent again
	code := 25201
	hub.PublishEvent(225, model.EventAlarmOn, &code, "Axis Z drive fault", at)
	hub.PublishEvent(225, model.EventConnectOK, nil, "", at) // not a screen event
	hub.PublishEvent(225, EventCycleClosed, nil, "", at)
	hub.Flush(at)

	if len(*sent) != 3 {
		t.Fatalf("sent %d messages, want status + alarm_on + cycle_closed: %+v", len(*sent), *sent)
	}
	if (*sent)[0].room != StatusRoom(225) || (*sent)[0].event != EventStatus {
		t.Fatalf("first message = %+v, want the status", (*sent)[0])
	}
	alarm := (*sent)[1].payload.(EventMessage)
	if (*sent)[1].room != EventRoom(225) || alarm.Type != model.EventAlarmOn || *alarm.Code != code || alarm.Seq != 2 {
		t.Fatalf("alarm message = %+v", alarm)
	}
}

func TestNilHubIgnoresPublishing(t *testing.T) {
	var hub *Hub
	hub.SetMachineCodes(1, []string{"SpindleLoad"})
	hub.PublishValues(1, time.Now(), map[string]any{"SpindleLoad": 1.0})
	hub.PublishStatus(1, Status{})
	hub.PublishEvent(1, model.EventAlarmOn, nil, "", time.Now())
}
