package opcua

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"
	"testing"
	"time"

	"github.com/gopcua/opcua/server"
	"github.com/gopcua/opcua/ua"
)

// startTestServer runs a gopcua server on a free localhost port (security None, anonymous; the
// gopcua server implements no other security) with a map namespace, until the test ends.
func startTestServer(t *testing.T) (*server.Server, *server.MapNamespace, int) {
	t.Helper()
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	port := listener.Addr().(*net.TCPAddr).Port
	listener.Close()

	opcServer := server.New(
		server.EnableSecurity("None", ua.MessageSecurityModeNone),
		server.EnableAuthMode(ua.UserTokenTypeAnonymous),
		server.EndPoint("localhost", port),
		server.SetLogger(slog.New(slog.DiscardHandler)),
	)
	namespace := server.NewMapNamespace(opcServer, "PlcTest")
	// SetValue notifies subscribers, which needs a started server: seed the map directly
	namespace.Data["HydPressure"] = 21.5
	namespace.Data["DriveReady"] = true
	if err := opcServer.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = opcServer.Close() })
	return opcServer, namespace, port
}

// changeReader hands out the changes of one subscription by handle; a publish carries several
// handles, so the ones not asked for yet are kept.
type changeReader struct {
	changes <-chan []Change
	waiting []Change
}

// next waits for the next change of handle.
func (reader *changeReader) next(t *testing.T, handle uint32) Change {
	t.Helper()
	timeout := time.After(10 * time.Second)
	for {
		for index, change := range reader.waiting {
			if change.Handle == handle {
				reader.waiting = append(reader.waiting[:index], reader.waiting[index+1:]...)
				return change
			}
		}
		select {
		case batch := <-reader.changes:
			reader.waiting = append(reader.waiting, batch...)
		case <-timeout:
			t.Fatalf("no change of handle %d", handle)
		}
	}
}

// TestConnectionSubscribe connects, subscribes to good, unknown and invalid node ids, receives the
// first values and a later change, checks the connection, and notices the server going away.
func TestConnectionSubscribe(t *testing.T) {
	opcServer, namespace, port := startTestServer(t)
	ctx := t.Context()

	connection, err := Connect(ctx, Settings{Endpoint: Endpoint("localhost", port), SecurityMode: SecurityModeNone, ConnectTimeout: 5 * time.Second})
	if err != nil {
		t.Fatal(err)
	}
	defer connection.Close(context.Background())

	namespaceIndex := namespace.ID()
	items := []Item{
		{Handle: 1, NodeID: fmt.Sprintf("ns=%d;s=HydPressure", namespaceIndex)},
		{Handle: 2, NodeID: fmt.Sprintf("ns=%d;s=DriveReady", namespaceIndex)},
		{Handle: 3, NodeID: "ns=2;s=/Random/Plc/DB9999.DBX0.0"},
		{Handle: 4, NodeID: "not a node id"},
	}
	results, changes, _, err := connection.Subscribe(ctx, 100*time.Millisecond, 50*time.Millisecond, items)
	if err != nil {
		t.Fatal(err)
	}
	statusByHandle := map[uint32]string{}
	for _, result := range results {
		statusByHandle[result.Handle] = result.Status
	}
	if statusByHandle[1] != "" || statusByHandle[2] != "" {
		t.Errorf("good nodes refused: %v", statusByHandle)
	}
	if !strings.HasPrefix(statusByHandle[4], "invalid node id") {
		t.Errorf("invalid node id accepted: %v", statusByHandle)
	}

	reader := &changeReader{changes: changes}
	if change := reader.next(t, 1); !change.Good || change.Value.Number != 21.5 {
		t.Errorf("first HydPressure %+v", change)
	}
	if change := reader.next(t, 2); !change.Good || change.Value.Number != 1 {
		t.Errorf("first DriveReady (bool → 1) %+v", change)
	}
	// An unknown node is refused at subscribe (Siemens: BadNodeIdUnknown) or, as here, its first
	// value is bad
	if statusByHandle[3] == "" {
		if change := reader.next(t, 3); change.Good || change.Status == "" {
			t.Errorf("unknown node gave a good value: %+v", change)
		}
	}

	namespace.SetValue("HydPressure", 22.25)
	if change := reader.next(t, 1); change.Value.Number != 22.25 {
		t.Errorf("changed HydPressure %+v", change)
	}

	if err := connection.CheckAlive(ctx); err != nil {
		t.Errorf("check alive: %v", err)
	}

	// A server that goes away may not close the channel (a machine switched off sends no reset):
	// the health read must then fail, which is what makes the PLC worker reconnect.
	_ = opcServer.Close()
	select {
	case <-connection.Lost():
	case <-time.After(2 * time.Second):
		checkCtx, cancel := context.WithTimeout(ctx, 3*time.Second)
		defer cancel()
		if err := connection.CheckAlive(checkCtx); err == nil {
			t.Error("server closed but the connection is neither lost nor failing")
		}
	}
}

// TestValueOf checks the conversion of the variant types the collector stores.
func TestValueOf(t *testing.T) {
	cases := []struct {
		value  any
		number float64
		text   string
		ok     bool
	}{
		{true, 1, "", true},
		{false, 0, "", true},
		{int16(-12), -12, "", true},
		{uint32(4000000000), 4000000000, "", true},
		{float32(1.5), 1.5, "", true},
		{"O1234", 0, "O1234", true},
		{[]byte{1, 2}, 0, "", false},
		{[]float32{21.5}, 21.5, "", true}, // e.g. DB9004.DBD8:REAL[1]
		{[]int16{-3}, -3, "", true},
		{[]string{"O1"}, 0, "O1", true},
		{[]float32{1, 2}, 0, "", false},
		{[]float32{}, 0, "", false},
		{[][]byte{{1}}, 0, "", false},
	}
	for _, testCase := range cases {
		variant, err := ua.NewVariant(testCase.value)
		if err != nil {
			t.Fatalf("%T: %v", testCase.value, err)
		}
		value, ok := ValueOf(variant)
		if ok != testCase.ok || value.Number != testCase.number || value.Text != testCase.text {
			t.Errorf("%T %v: got %+v ok=%v", testCase.value, testCase.value, value, ok)
		}
	}
	if _, ok := ValueOf(nil); ok {
		t.Error("nil variant converted")
	}
}
