package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/live"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// liveWait bounds every wait for a socket message.
const liveWait = 5 * time.Second

// liveTestServer is a real HTTP server with the auth routes and the live data, on the test database.
type liveTestServer struct {
	server *httptest.Server
	router http.Handler
	hub    *live.Hub
}

// newLiveTestServer starts the router with the real session middleware, the live handler and its
// hub; both stop when the test ends.
func newLiveTestServer(t *testing.T) liveTestServer {
	t.Helper()
	ctx, cancel := context.WithCancel(context.Background())
	authService := service.NewAuthService(repository.NewAuthRepository(testDBPool))
	machineService := newMachineService()
	userCountryService := service.NewUserCountryService(repository.NewUserCountryRepository(testDBPool))
	userStateService := service.NewUserStateService(repository.NewUserStateRepository(testDBPool), userCountryService)
	userCityService := service.NewUserCityService(repository.NewUserCityRepository(testDBPool), userStateService)
	branchService := service.NewBranchService(repository.NewBranchRepository(testDBPool), userCityService)
	userService := service.NewUserService(testDBPool, branchService,
		service.NewRoleService(repository.NewRoleRepository(testDBPool)), machineService)

	liveConfig := config.LiveConfig{Enabled: true, MinInterval: 200 * time.Millisecond, Revalidate: time.Minute, MaxSocketsPerUser: 10}
	hub := live.NewHub(liveConfig.MinInterval)
	liveHandler := handler.NewLiveHandler(liveConfig, hub, authService, machineService, userService)
	authHandler := handler.NewAuthHandler(authService, false)
	authHandler.SetLiveHandler(liveHandler)
	router := routes.New(routes.Handlers{Auth: authHandler, Live: liveHandler}, middleware.Authenticate(authService))
	go hub.Run(ctx)
	go liveHandler.Run(ctx)

	server := httptest.NewServer(router)
	t.Cleanup(func() {
		liveHandler.Close()
		server.Close()
		cancel()
	})
	return liveTestServer{server: server, router: router, hub: hub}
}

// liveTestClient speaks the Socket.IO v4 protocol over a WebSocket, like socket.io-client with
// transports: ['websocket'].
type liveTestClient struct {
	t       *testing.T
	conn    *websocket.Conn
	writeMu sync.Mutex
	nextAck int
	events  chan liveTestEvent
	acks    chan json.RawMessage
	closed  chan struct{} // the server sent a disconnect packet or the connection ended
	once    sync.Once
}

// liveTestEvent is one event the server emitted.
type liveTestEvent struct {
	name    string
	payload json.RawMessage
}

// dialLive opens a socket with cookie (nil = none) and origin (empty = none) and returns the
// handshake status when it was refused.
func dialLive(t *testing.T, server liveTestServer, cookie *http.Cookie, origin string) (*liveTestClient, int) {
	t.Helper()
	header := http.Header{}
	if cookie != nil {
		header.Set("Cookie", cookie.String())
	}
	if origin != "" {
		header.Set("Origin", origin)
	}
	address := "ws" + strings.TrimPrefix(server.server.URL, "http") + apiV1 + "/socket.io/?EIO=4&transport=websocket"
	conn, response, err := websocket.DefaultDialer.Dial(address, header)
	if err != nil {
		if response == nil {
			t.Fatalf("dial live: %v", err)
		}
		return nil, response.StatusCode
	}
	client := &liveTestClient{t: t, conn: conn, events: make(chan liveTestEvent, 32),
		acks: make(chan json.RawMessage, 8), closed: make(chan struct{})}
	t.Cleanup(func() { _ = conn.Close() })
	go client.readPackets()
	client.write("40") // connect to the main namespace
	if event := client.expectEvent("connect"); event.name != "connect" {
		t.Fatalf("socket did not connect")
	}
	return client, http.StatusSwitchingProtocols
}

// readPackets decodes the server's packets until the connection closes: Engine.IO ping → pong,
// Socket.IO connect / event / ack → the channels.
func (client *liveTestClient) readPackets() {
	defer client.markClosed()
	for {
		_, data, err := client.conn.ReadMessage()
		if err != nil {
			return
		}
		packet := string(data)
		switch {
		case packet == "2":
			client.write("3")
		case strings.HasPrefix(packet, "40"):
			client.events <- liveTestEvent{name: "connect"}
		case strings.HasPrefix(packet, "41"):
			client.markClosed()
		case strings.HasPrefix(packet, "42"):
			var message []json.RawMessage
			if json.Unmarshal([]byte(packet[2:]), &message) == nil && len(message) > 0 {
				event := liveTestEvent{}
				_ = json.Unmarshal(message[0], &event.name)
				if len(message) > 1 {
					event.payload = message[1]
				}
				client.events <- event
			}
		case strings.HasPrefix(packet, "43"):
			body := strings.TrimLeft(packet[2:], "0123456789")
			var message []json.RawMessage
			if json.Unmarshal([]byte(body), &message) == nil && len(message) > 0 {
				client.acks <- message[0]
			}
		}
	}
}

// markClosed records that the server closed the socket.
func (client *liveTestClient) markClosed() {
	client.once.Do(func() { close(client.closed) })
}

// write sends one raw packet.
func (client *liveTestClient) write(packet string) {
	client.writeMu.Lock()
	defer client.writeMu.Unlock()
	_ = client.conn.WriteMessage(websocket.TextMessage, []byte(packet))
}

// call emits event with payload and decodes its acknowledgement into answer.
func (client *liveTestClient) call(event string, payload, answer any) {
	client.t.Helper()
	body, _ := json.Marshal([]any{event, payload})
	client.nextAck++
	client.write("42" + strconv.Itoa(client.nextAck) + string(body))
	select {
	case raw := <-client.acks:
		if err := json.Unmarshal(raw, answer); err != nil {
			client.t.Fatalf("%s ack %s: %v", event, raw, err)
		}
	case <-time.After(liveWait):
		client.t.Fatalf("%s: no acknowledgement", event)
	}
}

// expectEvent waits for the next event called name, skipping others.
func (client *liveTestClient) expectEvent(name string) liveTestEvent {
	client.t.Helper()
	deadline := time.After(liveWait)
	for {
		select {
		case event := <-client.events:
			if event.name == name {
				return event
			}
		case <-deadline:
			client.t.Fatalf("no %s event", name)
			return liveTestEvent{}
		}
	}
}

// expectClosed waits until the server closed the socket.
func (client *liveTestClient) expectClosed() {
	client.t.Helper()
	select {
	case <-client.closed:
	case <-time.After(liveWait):
		client.t.Fatalf("socket still open")
	}
}

// liveValues is a machine:values message.
type liveValues struct {
	MachineID int64          `json:"machineId"`
	Seq       int64          `json:"seq"`
	Values    map[string]any `json:"values"`
}

// subscribeAnswer is the subscribe acknowledgement.
type subscribeAnswer struct {
	OK        bool    `json:"ok"`
	Message   string  `json:"message"`
	Denied    []int64 `json:"denied"`
	Unknown   []string
	Snapshots []struct {
		MachineID int64          `json:"machineId"`
		Seq       int64          `json:"seq"`
		Values    map[string]any `json:"values"`
	} `json:"snapshots"`
}

// TestLiveData covers the Socket.IO live data end to end: the handshake needs a session and
// this site's Origin, a socket gets only the machines its user may read and only the codes it
// asked for, losing access drops the rooms, and logout closes that session's sockets only.
func TestLiveData(t *testing.T) {
	server := newLiveTestServer(t)
	suffix := uniqueSuffix()
	company := newTestCompany(t, "live"+suffix, 2)
	otherCompany := newTestCompany(t, "liveo"+suffix, 1)
	assigned, unassigned, foreign := company.machineIDs[0], company.machineIDs[1], otherCompany.machineIDs[0]

	operatorName := "live_operator_" + suffix
	operatorID := createAuthTestUser(t, operatorName, model.RoleOperator)
	ctx := context.Background()
	if _, err := testDBPool.Exec(ctx, `UPDATE users SET company_id = $2 WHERE id = $1`, operatorID, company.customerID); err != nil {
		t.Fatalf("set operator company: %v", err)
	}
	if _, err := testDBPool.Exec(ctx, `INSERT INTO user_machine_access (user_id, machine_id) VALUES ($1, $2)`, operatorID, assigned); err != nil {
		t.Fatalf("assign machine: %v", err)
	}
	t.Cleanup(func() { _, _ = testDBPool.Exec(ctx, `DELETE FROM user_machine_access WHERE user_id = $1`, operatorID) })
	server.hub.SetMachineCodes(assigned, []string{"SpindleLoad", "SpindleRpm", "progName"})
	server.hub.PublishValues(assigned, time.Now(), map[string]any{"SpindleLoad": 41.5, "SpindleRpm": 2400.0, "progName": "O1234"})

	session := login(t, server.router, operatorName, false)

	t.Run("handshake needs a session and this site", func(t *testing.T) {
		if _, status := dialLive(t, server, nil, ""); status != http.StatusUnauthorized {
			t.Fatalf("without a session: status %d, want 401", status)
		}
		if _, status := dialLive(t, server, session, "https://evil.example.com"); status != http.StatusForbidden {
			t.Fatalf("from another site: status %d, want 403", status)
		}
	})

	client, _ := dialLive(t, server, session, server.server.URL)

	t.Run("subscribe gets only readable machines and known codes", func(t *testing.T) {
		var answer subscribeAnswer
		client.call("subscribe", map[string]any{
			"machineIds": []int64{assigned, unassigned, foreign}, "codes": []string{"SpindleLoad", "progName", "SPBeTemp"},
			"channels": []string{"status", "event"}, "intervalMs": 200, "userId": 1, "role": "superadmin",
		}, &answer)
		if !answer.OK || fmt.Sprint(answer.Denied) != fmt.Sprint([]int64{unassigned, foreign}) || fmt.Sprint(answer.Unknown) != "[SPBeTemp]" {
			t.Fatalf("answer = %+v, want %d and %d denied, SPBeTemp unknown", answer, unassigned, foreign)
		}
		if len(answer.Snapshots) != 1 || answer.Snapshots[0].MachineID != assigned ||
			fmt.Sprint(answer.Snapshots[0].Values) != "map[SpindleLoad:41.5 progName:O1234]" {
			t.Fatalf("snapshots = %+v, want only machine %d's SpindleLoad and progName", answer.Snapshots, assigned)
		}

		var invalid subscribeAnswer
		client.call("subscribe", map[string]any{"machineIds": []int64{}}, &invalid)
		if invalid.OK || invalid.Message == "" {
			t.Fatalf("invalid subscribe answer = %+v, want an error", invalid)
		}
	})

	t.Run("values arrive filtered to the subscribed codes", func(t *testing.T) {
		server.hub.PublishValues(foreign, time.Now(), map[string]any{"SpindleLoad": 99.0})
		server.hub.PublishValues(assigned, time.Now(), map[string]any{"SpindleLoad": 44.5, "SpindleRpm": 2500.0, "machPosX": 1.0})
		var values liveValues
		_ = json.Unmarshal(client.expectEvent(live.EventValues).payload, &values)
		if values.MachineID != assigned || fmt.Sprint(values.Values) != "map[SpindleLoad:44.5]" {
			t.Fatalf("values = %+v, want machine %d SpindleLoad 44.5 only", values, assigned)
		}

		code := 25201
		server.hub.PublishEvent(assigned, model.EventAlarmOn, &code, "Axis Z drive fault", time.Now())
		var event struct {
			Type string `json:"type"`
			Code int    `json:"code"`
		}
		_ = json.Unmarshal(client.expectEvent(live.EventMachine).payload, &event)
		if event.Type != model.EventAlarmOn || event.Code != code {
			t.Fatalf("event = %+v, want alarm_on %d", event, code)
		}
	})

	t.Run("losing access leaves the machine's rooms", func(t *testing.T) {
		if _, err := testDBPool.Exec(ctx, `DELETE FROM user_machine_access WHERE user_id = $1`, operatorID); err != nil {
			t.Fatalf("remove access: %v", err)
		}
		if _, err := testDBPool.Exec(ctx, `SELECT pg_notify($1, $2)`, service.UserAccessChangedChannel, strconv.FormatInt(operatorID, 10)); err != nil {
			t.Fatalf("notify: %v", err)
		}
		var change struct {
			Removed []int64 `json:"removed"`
		}
		_ = json.Unmarshal(client.expectEvent("access:changed").payload, &change)
		if fmt.Sprint(change.Removed) != fmt.Sprint([]int64{assigned}) {
			t.Fatalf("access change = %+v, want %d removed", change, assigned)
		}
		var answer subscribeAnswer
		client.call("subscribe", map[string]any{"machineIds": []int64{assigned}, "codes": []string{"SpindleLoad"}}, &answer)
		if fmt.Sprint(answer.Denied) != fmt.Sprint([]int64{assigned}) {
			t.Fatalf("answer = %+v, want %d denied now", answer, assigned)
		}
		if _, err := testDBPool.Exec(ctx, `INSERT INTO user_machine_access (user_id, machine_id) VALUES ($1, $2)`, operatorID, assigned); err != nil {
			t.Fatalf("assign machine again: %v", err)
		}
		_, _ = testDBPool.Exec(ctx, `SELECT pg_notify($1, $2)`, service.UserAccessChangedChannel, strconv.FormatInt(operatorID, 10))
		client.expectEvent("access:changed")
	})

	t.Run("logout closes the session's sockets only", func(t *testing.T) {
		secondTab, _ := dialLive(t, server, session, "")
		otherSession := login(t, server.router, operatorName, false)
		otherDevice, _ := dialLive(t, server, otherSession, "")
		var answer subscribeAnswer
		otherDevice.call("subscribe", map[string]any{"machineIds": []int64{assigned}, "codes": []string{"SpindleRpm"}}, &answer)
		if !answer.OK || len(answer.Denied) != 0 {
			t.Fatalf("other device answer = %+v", answer)
		}

		result, _ := authRequest(t, server.router, http.MethodPost, apiV1+"/auth/logout", nil, session, nil)
		expectStatus(t, result, http.StatusOK, "Logged out")
		client.expectClosed()
		secondTab.expectClosed()

		server.hub.PublishValues(assigned, time.Now(), map[string]any{"SpindleRpm": 2600.0})
		var values liveValues
		_ = json.Unmarshal(otherDevice.expectEvent(live.EventValues).payload, &values)
		if fmt.Sprint(values.Values) != "map[SpindleRpm:2600]" {
			t.Fatalf("other device values = %+v, want SpindleRpm 2600", values)
		}
	})
}
