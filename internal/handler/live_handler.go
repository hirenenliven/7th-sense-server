package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"sync"
	"time"

	"github.com/zishang520/socket.io/servers/engine/v3/transports"
	"github.com/zishang520/socket.io/servers/socket/v3"
	"github.com/zishang520/socket.io/v3/pkg/types"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/live"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// Socket.IO events of the live data besides live.EventValues / EventStatus / EventMachine
// (docs/live-socket-events.md).
const (
	liveEventSubscribe      = "subscribe"       // browser → server, acknowledged with dto.LiveSubscribeResponse
	liveEventUnsubscribe    = "unsubscribe"     // browser → server, acknowledged with {ok: true}
	liveEventAccessChanged  = "access:changed"  // server → browser: dto.LiveAccessChanged
	liveEventSessionExpired = "session:expired" // server → browser, then the socket is closed
)

const (
	liveQueryTimeout       = 5 * time.Second  // one database call of the live server
	liveMaxMessageBytes    = 64 * 1024        // a subscribe of 500 machines and 100 codes fits easily
	livePingInterval       = 25 * time.Second // Socket.IO heartbeat: a dead connection is dropped
	livePingTimeout        = 20 * time.Second //   after interval + timeout without a pong
	liveListenRetryMaximum = time.Minute      // longest wait before LISTEN is opened again
)

// LiveHandler serves the Socket.IO live data at /api/v1/socket.io/ (ADR 0009): it pushes the
// collector's machine values (live.Hub) to the screens. Every socket belongs to the user of the
// session cookie the handshake carried (the authenticate middleware ran on it) and joins a
// machine's rooms only when that user may read the machine (MachineService, the REST rule).
// The registry below tracks each socket's user, session, machines and subscriptions.
type LiveHandler struct {
	config         config.LiveConfig
	hub            *live.Hub
	authService    *service.AuthService
	machineService *service.MachineService
	userService    *service.UserService
	server         *socket.Server
	engine         http.Handler

	mu             sync.Mutex
	sockets        map[socket.SocketId]*liveSocket
	userSockets    map[int64][]*liveSocket // oldest first
	sessionSockets map[string][]*liveSocket
}

// liveSocket is one open connection (one browser tab) in the registry.
type liveSocket struct {
	socket      *socket.Socket
	user        model.AuthUser
	token       string // the session token, kept to check the session again; never logged
	sessionHash string
	allowed     map[int64]bool // machines the user may read, loaded at connect and on changes
	// subscriptions holds, per machine, the rooms this socket joined for it.
	subscriptions map[int64]liveSubscription
}

// liveSubscription is what one socket watches of one machine.
type liveSubscription struct {
	valuesRoom string // live.Hub room of its variable set; "" when it asked only for channels
	channels   []string
}

// NewLiveHandler creates the Socket.IO server (WebSocket transport only: no long-polling, so no
// sticky sessions and no session check per poll) and makes hub send through it.
func NewLiveHandler(liveConfig config.LiveConfig, hub *live.Hub, authService *service.AuthService,
	machineService *service.MachineService, userService *service.UserService) *LiveHandler {
	options := socket.DefaultServerOptions()
	options.SetServeClient(false)
	options.SetTransports(types.NewSet[transports.TransportCtor](&transports.WebSocketBuilder{}))
	options.SetPingInterval(livePingInterval)
	options.SetPingTimeout(livePingTimeout)
	options.SetMaxHttpBufferSize(liveMaxMessageBytes)

	liveHandler := &LiveHandler{
		config:         liveConfig,
		hub:            hub,
		authService:    authService,
		machineService: machineService,
		userService:    userService,
		server:         socket.NewServer(nil, options),
		sockets:        map[socket.SocketId]*liveSocket{},
		userSockets:    map[int64][]*liveSocket{},
		sessionSockets: map[string][]*liveSocket{},
	}
	liveHandler.engine = liveHandler.server.ServeHandler(options)
	_ = liveHandler.server.On("connection", liveHandler.onConnection)
	hub.SetEmitter(liveHandler.emit)
	return liveHandler
}

// Serve handles GET /api/v1/socket.io/ (and /api/as/<key>/v1/socket.io/ in a "log in as" tab):
// the Socket.IO WebSocket handshake. It runs behind authenticate, so a request without a
// session never gets here (401).
func (liveHandler *LiveHandler) Serve(writer http.ResponseWriter, request *http.Request) {
	// Browsers send cookies on cross-site WebSocket requests and CrossOriginProtection does not
	// cover GET, so without this another website could open a socket as the user.
	if !isSameOrigin(request) {
		response.SendResponse(writer, request, response.Forbidden, nil)
		return
	}
	// The socket outlives the server's 10 s read / write timeouts.
	controller := http.NewResponseController(writer)
	if err := errors.Join(controller.SetReadDeadline(time.Time{}), controller.SetWriteDeadline(time.Time{})); err != nil {
		slog.Warn("live: clear connection deadlines", "error", err)
	}
	liveHandler.engine.ServeHTTP(writer, request)
}

// Run follows session and access changes until ctx ends: LISTEN for user_access_changed /
// machines_changed, and every LIVE_REVALIDATE_SEC check every socket's session and machines.
func (liveHandler *LiveHandler) Run(ctx context.Context) {
	go liveHandler.followAccessNotifications(ctx)
	ticker := time.NewTicker(liveHandler.config.Revalidate)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			liveHandler.revalidateAll()
		}
	}
}

// Close disconnects every socket (server shutdown); the browsers reconnect when it is back.
func (liveHandler *LiveHandler) Close() {
	liveHandler.server.Close(nil)
}

// DisconnectSession closes every socket of token's session (logout); sockets of the user's other
// sessions (another device, a "log in as" tab) stay.
func (liveHandler *LiveHandler) DisconnectSession(token string) {
	if token == "" {
		return
	}
	liveHandler.closeSession(liveHandler.authService.SessionHash(token))
}

// onConnection registers a new socket for the session's user with the machines they may read.
func (liveHandler *LiveHandler) onConnection(clients ...any) {
	client, ok := clients[0].(*socket.Socket)
	if !ok {
		return
	}
	// The handshake request ended once upgraded: only its values are used, never its context
	// for database calls.
	request := client.Request().Request()
	user, ok := middleware.CurrentUser(request.Context())
	if !ok {
		closeLiveSocket(client)
		return
	}
	allowed, err := liveHandler.loadAllowedMachines(user)
	if err != nil {
		slog.Error("live: load machines of user", "userId", user.ID, "error", err)
		closeLiveSocket(client)
		return
	}
	token := middleware.SessionToken(request)
	entry := &liveSocket{
		socket: client, user: user, token: token, sessionHash: liveHandler.authService.SessionHash(token),
		allowed: allowed, subscriptions: map[int64]liveSubscription{},
	}
	liveHandler.register(entry)
	_ = client.On(liveEventSubscribe, func(args ...any) { liveHandler.subscribe(entry, args) })
	_ = client.On(liveEventUnsubscribe, func(args ...any) { liveHandler.unsubscribe(entry, args) })
	_ = client.On("disconnect", func(...any) { liveHandler.unregister(entry) })
	// A socket closed before its disconnect listener was added would stay registered.
	if !client.Connected() {
		liveHandler.unregister(entry)
	}
}

// register adds entry to the registry and its control rooms, closing the user's oldest sockets
// beyond LIVE_MAX_SOCKETS_PER_USER.
func (liveHandler *LiveHandler) register(entry *liveSocket) {
	liveHandler.mu.Lock()
	liveHandler.sockets[entry.socket.Id()] = entry
	liveHandler.userSockets[entry.user.ID] = append(liveHandler.userSockets[entry.user.ID], entry)
	liveHandler.sessionSockets[entry.sessionHash] = append(liveHandler.sessionSockets[entry.sessionHash], entry)
	entry.socket.Join(userRoom(entry.user.ID), sessionRoom(entry.sessionHash))
	userSockets := liveHandler.userSockets[entry.user.ID]
	var oldest []*liveSocket
	if extra := len(userSockets) - liveHandler.config.MaxSocketsPerUser; extra > 0 {
		oldest = slices.Clone(userSockets[:extra])
	}
	liveHandler.mu.Unlock()

	// Outside the lock: a disconnect runs unregister, which takes it.
	for _, closing := range oldest {
		closeLiveSocket(closing.socket)
	}
}

// unregister removes a closed socket from the registry and the hub rooms it held.
func (liveHandler *LiveHandler) unregister(entry *liveSocket) {
	liveHandler.mu.Lock()
	defer liveHandler.mu.Unlock()
	if liveHandler.sockets[entry.socket.Id()] != entry {
		return
	}
	delete(liveHandler.sockets, entry.socket.Id())
	liveHandler.userSockets[entry.user.ID] = withoutSocket(liveHandler.userSockets[entry.user.ID], entry)
	if len(liveHandler.userSockets[entry.user.ID]) == 0 {
		delete(liveHandler.userSockets, entry.user.ID)
	}
	liveHandler.sessionSockets[entry.sessionHash] = withoutSocket(liveHandler.sessionSockets[entry.sessionHash], entry)
	if len(liveHandler.sessionSockets[entry.sessionHash]) == 0 {
		delete(liveHandler.sessionSockets, entry.sessionHash)
	}
	for _, subscription := range entry.subscriptions {
		if subscription.valuesRoom != "" {
			liveHandler.hub.LeaveValues(subscription.valuesRoom)
		}
	}
	entry.subscriptions = nil
}

// subscribe handles the "subscribe" message: for every machine the user may read, the socket
// joins the room of exactly the codes asked (replacing its earlier subscription of that
// machine) and the channels asked; the answer carries their current values.
func (liveHandler *LiveHandler) subscribe(entry *liveSocket, args []any) {
	acknowledge := liveAck(args)
	var subscribeRequest dto.LiveSubscribeRequest
	if err := decodeLiveMessage(args, &subscribeRequest); err != nil {
		acknowledge(dto.LiveSubscribeResponse{Message: err.Error()})
		return
	}
	acknowledge(liveHandler.applySubscription(entry, subscribeRequest))
}

// applySubscription joins entry to the rooms of subscribeRequest; see subscribe.
func (liveHandler *LiveHandler) applySubscription(entry *liveSocket, subscribeRequest dto.LiveSubscribeRequest) dto.LiveSubscribeResponse {
	answer := dto.LiveSubscribeResponse{OK: true, Denied: []int64{}, Unknown: []string{}, Snapshots: []live.Snapshot{}}
	liveHandler.mu.Lock()
	defer liveHandler.mu.Unlock()
	if entry.subscriptions == nil { // closed meanwhile
		return answer
	}
	for _, machineID := range uniqueIDs(subscribeRequest.MachineIDs) {
		if !entry.allowed[machineID] {
			answer.Denied = append(answer.Denied, machineID)
			continue
		}
		liveHandler.leaveMachine(entry, machineID)
		subscription := liveSubscription{channels: slices.Compact(slices.Sorted(slices.Values(subscribeRequest.Channels)))}
		room, kept, unknown := liveHandler.hub.JoinValues(machineID, subscribeRequest.Codes, subscribeRequest.Interval())
		if room != "" {
			subscription.valuesRoom = room
			entry.socket.Join(socket.Room(room))
		}
		for _, channel := range subscription.channels {
			entry.socket.Join(channelRoom(machineID, channel))
		}
		entry.subscriptions[machineID] = subscription
		for _, code := range unknown {
			if !slices.Contains(answer.Unknown, code) {
				answer.Unknown = append(answer.Unknown, code)
			}
		}
		answer.Snapshots = append(answer.Snapshots, liveHandler.hub.Snapshot(machineID, kept))
	}
	return answer
}

// unsubscribe handles the "unsubscribe" message: the socket leaves every room of those machines.
func (liveHandler *LiveHandler) unsubscribe(entry *liveSocket, args []any) {
	acknowledge := liveAck(args)
	var unsubscribeRequest dto.LiveUnsubscribeRequest
	if err := decodeLiveMessage(args, &unsubscribeRequest); err != nil {
		acknowledge(map[string]any{"ok": false, "message": err.Error()})
		return
	}
	liveHandler.mu.Lock()
	if entry.subscriptions != nil {
		for _, machineID := range unsubscribeRequest.MachineIDs {
			liveHandler.leaveMachine(entry, machineID)
		}
	}
	liveHandler.mu.Unlock()
	acknowledge(map[string]any{"ok": true})
}

// leaveMachine drops entry's subscription of machineID. The caller holds liveHandler.mu.
func (liveHandler *LiveHandler) leaveMachine(entry *liveSocket, machineID int64) {
	subscription, found := entry.subscriptions[machineID]
	if !found {
		return
	}
	if subscription.valuesRoom != "" {
		entry.socket.Leave(socket.Room(subscription.valuesRoom))
		liveHandler.hub.LeaveValues(subscription.valuesRoom)
	}
	for _, channel := range subscription.channels {
		entry.socket.Leave(channelRoom(machineID, channel))
	}
	delete(entry.subscriptions, machineID)
}

// followAccessNotifications LISTENs for user_access_changed and machines_changed until ctx ends,
// reopening a lost connection with backoff (1 s doubling to 1 min).
func (liveHandler *LiveHandler) followAccessNotifications(ctx context.Context) {
	wait := time.Second
	for ctx.Err() == nil {
		err := liveHandler.userService.ListenForAccessChanges(ctx, func() { wait = time.Second },
			func(channel, payload string) { liveHandler.handleAccessNotification(channel, payload) })
		if ctx.Err() != nil {
			return
		}
		slog.Error("live: access listener lost, reconnecting", "retryIn", wait.String(), "error", err)
		select {
		case <-ctx.Done():
			return
		case <-time.After(wait):
		}
		wait = min(wait*2, liveListenRetryMaximum)
	}
}

// handleAccessNotification follows one notification: a user's role, status or machines changed
// (payload: user id) → check that user's sessions and machines now; a machine was added,
// deleted or changed (payload: machine id) → reload every connected user's machines.
func (liveHandler *LiveHandler) handleAccessNotification(channel, payload string) {
	switch channel {
	case service.UserAccessChangedChannel:
		userID, err := strconv.ParseInt(payload, 10, 64)
		if err != nil {
			return
		}
		liveHandler.revalidateUser(userID)
	case service.MachinesChangedChannel:
		for _, userID := range liveHandler.connectedUserIDs() {
			liveHandler.reloadAccess(userID)
		}
	}
}

// revalidateAll checks every connected user's sessions and machines (LIVE_REVALIDATE_SEC), in
// case a notification was missed or a session expired.
func (liveHandler *LiveHandler) revalidateAll() {
	for _, userID := range liveHandler.connectedUserIDs() {
		liveHandler.revalidateUser(userID)
	}
}

// revalidateUser closes the user's sessions that ended (logged out elsewhere, expired, user
// deactivated or deleted) and reloads the machines the user may read.
func (liveHandler *LiveHandler) revalidateUser(userID int64) {
	liveHandler.mu.Lock()
	sessions := map[string]string{} // session hash → token
	for _, entry := range liveHandler.userSockets[userID] {
		sessions[entry.sessionHash] = entry.token
	}
	liveHandler.mu.Unlock()

	for sessionHash, token := range sessions {
		ctx, cancel := context.WithTimeout(context.Background(), liveQueryTimeout)
		user, err := liveHandler.authService.Authenticate(ctx, token)
		cancel()
		switch {
		case errors.Is(err, response.LoginRequired):
			liveHandler.closeSession(sessionHash)
		case err != nil:
			slog.Error("live: check session", "userId", userID, "error", err)
		default:
			liveHandler.updateSessionUser(sessionHash, user) // the role or company may have changed
		}
	}
	liveHandler.reloadAccess(userID)
}

// updateSessionUser stores the user the session resolves to now on its sockets.
func (liveHandler *LiveHandler) updateSessionUser(sessionHash string, user model.AuthUser) {
	liveHandler.mu.Lock()
	defer liveHandler.mu.Unlock()
	for _, entry := range liveHandler.sessionSockets[sessionHash] {
		entry.user = user
	}
}

// reloadAccess reloads the machines userID may read on all of the user's sockets: the sockets
// leave the rooms of machines no longer allowed, and the user's sockets get access:changed.
func (liveHandler *LiveHandler) reloadAccess(userID int64) {
	liveHandler.mu.Lock()
	entries := slices.Clone(liveHandler.userSockets[userID])
	var user model.AuthUser
	if len(entries) > 0 {
		user = entries[0].user
	}
	liveHandler.mu.Unlock()
	if len(entries) == 0 {
		return
	}
	allowed, err := liveHandler.loadAllowedMachines(user)
	if err != nil {
		slog.Error("live: reload machines of user", "userId", userID, "error", err)
		return
	}

	change := dto.LiveAccessChanged{Removed: []int64{}, Added: []int64{}}
	liveHandler.mu.Lock()
	for machineID := range entries[0].allowed {
		if !allowed[machineID] {
			change.Removed = append(change.Removed, machineID)
		}
	}
	for machineID := range allowed {
		if !entries[0].allowed[machineID] {
			change.Added = append(change.Added, machineID)
		}
	}
	for _, entry := range entries {
		for _, machineID := range change.Removed {
			if entry.subscriptions != nil {
				liveHandler.leaveMachine(entry, machineID)
			}
		}
		entry.allowed = allowed
	}
	liveHandler.mu.Unlock()

	if len(change.Removed) > 0 || len(change.Added) > 0 {
		slices.Sort(change.Removed)
		slices.Sort(change.Added)
		liveHandler.emit(string(userRoom(userID)), liveEventAccessChanged, change)
	}
}

// closeSession tells the session's sockets it ended and closes them.
func (liveHandler *LiveHandler) closeSession(sessionHash string) {
	liveHandler.mu.Lock()
	entries := slices.Clone(liveHandler.sessionSockets[sessionHash])
	liveHandler.mu.Unlock()
	if len(entries) == 0 {
		return
	}
	liveHandler.emit(string(sessionRoom(sessionHash)), liveEventSessionExpired, map[string]any{})
	for _, entry := range entries {
		closeLiveSocket(entry.socket)
	}
}

// connectedUserIDs returns the users with at least one open socket.
func (liveHandler *LiveHandler) connectedUserIDs() []int64 {
	liveHandler.mu.Lock()
	defer liveHandler.mu.Unlock()
	userIDs := make([]int64, 0, len(liveHandler.userSockets))
	for userID := range liveHandler.userSockets {
		userIDs = append(userIDs, userID)
	}
	return userIDs
}

// loadAllowedMachines returns the machines user may read, with the REST rule (MachineService).
func (liveHandler *LiveHandler) loadAllowedMachines(user model.AuthUser) (map[int64]bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), liveQueryTimeout)
	defer cancel()
	machines, err := liveHandler.machineService.GetAll(ctx, user, nil)
	if err != nil {
		return nil, err
	}
	allowed := make(map[int64]bool, len(machines))
	for _, machine := range machines {
		allowed[machine.ID] = true
	}
	return allowed, nil
}

// emit sends event with payload to room (the hub's sender too).
func (liveHandler *LiveHandler) emit(room, event string, payload any) {
	if err := liveHandler.server.To(socket.Room(room)).Emit(event, payload); err != nil {
		slog.Error("live: emit", "event", event, "error", err)
	}
}

// closeLiveSocket ends client with a Socket.IO disconnect packet: the browser then does not
// reconnect by itself ("io server disconnect"), so a socket closed for a tab limit or an ended
// session does not come straight back. Its connection closes once the browser drops it.
func closeLiveSocket(client *socket.Socket) {
	client.Disconnect(false)
}

// decodeLiveMessage decodes the first argument of a Socket.IO message into target and validates it.
func decodeLiveMessage(args []any, target interface{ Validate() error }) error {
	if len(args) == 0 {
		return response.InvalidJSON
	}
	if _, isAck := args[0].(socket.Ack); isAck {
		return response.InvalidJSON
	}
	encoded, err := json.Marshal(args[0])
	if err != nil {
		return response.InvalidJSON
	}
	if err := json.Unmarshal(encoded, target); err != nil {
		return response.InvalidJSON
	}
	return target.Validate()
}

// liveAck returns the acknowledgement callback of a Socket.IO message, or one that does nothing
// when the browser asked for none.
func liveAck(args []any) func(payload any) {
	if len(args) > 0 {
		if acknowledge, ok := args[len(args)-1].(socket.Ack); ok {
			return func(payload any) { acknowledge([]any{payload}, nil) }
		}
	}
	return func(any) {}
}

// isSameOrigin reports whether a browser request comes from this site: its Origin, when sent,
// names the request's host. Non-browser clients send no Origin.
func isSameOrigin(request *http.Request) bool {
	origin := request.Header.Get("Origin")
	if origin == "" {
		return true
	}
	parsed, err := url.Parse(origin)
	return err == nil && parsed.Host == request.Host
}

// channelRoom is machineID's room of channel (live.ChannelStatus or live.ChannelEvent).
func channelRoom(machineID int64, channel string) socket.Room {
	if channel == live.ChannelStatus {
		return socket.Room(live.StatusRoom(machineID))
	}
	return socket.Room(live.EventRoom(machineID))
}

// userRoom holds all sockets of a user; it carries no machine data.
func userRoom(userID int64) socket.Room {
	return socket.Room("u:" + strconv.FormatInt(userID, 10))
}

// sessionRoom holds the sockets of one session; it carries no machine data.
func sessionRoom(sessionHash string) socket.Room {
	return socket.Room("s:" + sessionHash)
}

// withoutSocket returns entries without entry.
func withoutSocket(entries []*liveSocket, entry *liveSocket) []*liveSocket {
	return slices.DeleteFunc(entries, func(other *liveSocket) bool { return other == entry })
}

// uniqueIDs returns ids without duplicates, in their first order.
func uniqueIDs(ids []int64) []int64 {
	unique := make([]int64, 0, len(ids))
	for _, id := range ids {
		if !slices.Contains(unique, id) {
			unique = append(unique, id)
		}
	}
	return unique
}
