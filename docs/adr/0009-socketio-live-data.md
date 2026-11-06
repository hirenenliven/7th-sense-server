# 0009. Socket.IO live data: the collector pushes each screen's variables, checked per user

- **Date:** 2026-11-05
- **Status:** Accepted

## Context

Screens such as the Virtual Screen polled a REST endpoint every 5 s: a value could be 5 s old and every open screen cost a set of queries every 5 s. The CNC collector runs inside the API process (ADR 0005), so every reading passes through this process before it is stored. Text values (program name, alarm text) were read but never stored, so no screen could show them live. Plan: `../../docs/socketio-live-data-plan.md` (workspace root).

## Decision

- **Transport:** Socket.IO v4 protocol with `github.com/zishang520/socket.io/servers/socket/v3` (v3.0.6), **WebSocket transport only** (no long-polling: no 10 s server timeouts, no session check per poll, no sticky sessions). The browser uses `socket.io-client@4`. A Phase 0 spike proved it builds under `GOARCH=386`, works behind our middleware and carries the request's `AuthUser`.
- **Mount:** `GET /api/v1/socket.io/` on the protected router, so the existing `authenticate` runs on the handshake (401 without a session); a "log in as" tab connects to `/api/as/<key>/v1/socket.io/` unchanged (`RouteLoginAs`). The handler checks `Origin` = `Host` (browsers send cookies on cross-site WebSocket requests and `CrossOriginProtection` does not cover GET) and clears the 10 s read / write deadlines before the upgrade. `middleware.statusRecorder` gained `Hijack()`, because the WebSocket library asserts `http.Hijacker` directly.
- **Hub** (`internal/live`): one in-memory store of the latest value of every variable per machine (numbers **and** text), fed by the workers **before** their database writes, under a mutex and never waiting for a socket. Its goroutine sends per room at most once per the room's interval (≥ `LIVE_MIN_INTERVAL_MS` 200 ms), only the codes that changed, newest value wins, nothing when nothing changed. Every change bumps a per-machine `seq`.
- **Rooms by exact variable set:** a screen subscribes with exact `cnc_variables` codes; the room is machine + sorted codes + interval (`m:<id>:v:<hash>@<ms>`), so identical screens share one message and no screen receives variables it does not show. `m:<id>:status` (online / mode / program) and `m:<id>:event` (alarm on/off, emergency on/off, cycle opened/closed) are channels. The browser never names rooms.
- **Access per socket, REST rule:** the socket's user comes from the session cookie, never from a message. Its allowed machines are `MachineService.GetAll(user)`, loaded at connect; a subscribe for any other machine is `denied` and joins nothing. Rooms are shared between users because only users who passed their own check get in.
- **Keeping access current:** `UserService.Update` / `Delete` send `NOTIFY user_access_changed '<userId>'` in their transaction (delivered at commit); the live server re-checks that user's sessions and machines at once, drops rooms no longer allowed and sends `access:changed`. `machines_changed` reloads every connected user's machines. Every `LIVE_REVALIDATE_SEC` (60 s) all sessions are re-authenticated (expired / deactivated → `session:expired` + close) as a safety net. `AuthHandler.Logout` closes that session's sockets only. At most `LIVE_MAX_SOCKETS_PER_USER` (10) sockets per user; the oldest are closed. Server-side closes send the Socket.IO disconnect packet, so the browser does not reconnect by itself.
- **Business rules stay in services:** derived values (cycle elapsed, job progress, OEE) are never pushed; an event makes the screen reload its REST endpoint. Events are published **after** their database write so that reload finds them.
- `LIVE_ENABLED=false` creates neither hub nor route; workers' hub calls are no-ops on a nil hub, screens fall back to REST polling.

## Alternatives considered

- **`googollee/go-socket.io`:** speaks only the old Socket.IO 1.x protocol; `socket.io-client@4` cannot connect.
- **Plain WebSocket (`coder/websocket`) with our own event names:** fewer dependencies, but reconnect, acknowledgements and rooms would have to be written by hand on both sides. Kept as the fallback had the spike failed.
- **Server-Sent Events:** one-way; subscriptions would need extra REST calls and the browser's 6-connection limit per host applies over HTTP/1.1.
- **Room per category or per variable:** a category sends unused variables (axis positions flood a spindle tile); one room per variable turns one reading into many messages per socket.
- **Access checked per message with a DB query:** correct but a query per subscribe; the cached allowed set plus notify / 60 s re-check gives the same rule without it.

## Consequences

- Values reach a screen within its interval (≥ 200 ms) of the reading, independent of database speed; open screens no longer poll every 5 s (the Virtual Screen keeps a 30 s REST refresh for derived data and the stale check).
- New dependency tree (Socket.IO server, `gorilla/websocket`, `quic-go` pulled in by the library's WebTransport support, unused here).
- OpenAPI cannot describe socket events: the contract lives in [docs/live-socket-events.md](../live-socket-events.md).
- The hub and registry are in-process memory: one API process only (true today, the DLL ties collection to one process). Several processes would need a Socket.IO adapter (e.g. Redis) at the hub boundary.
- A variable the collector already reads reaches a screen by adding its code to the client's `LIVE_CODES`; PLC / NC addresses still need a DLL read function first.
