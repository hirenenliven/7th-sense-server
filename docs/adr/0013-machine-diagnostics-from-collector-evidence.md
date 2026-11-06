# 0013. Machine diagnostics: reasons from what the collector recorded, network tested from the server

- **Date:** 2026-11-06
- **Status:** Accepted

## Context

When a customer reports that a machine sends no data, the superadmin has to read server logs and tables to find out why. The cause can be anywhere: the machine's set-up, the server build, the network, the controller, or the browser's socket. The DLL also calls back only when something changes, so a healthy idle machine can be silent for minutes, and "no new data" alone is not a fault.

## Decision

- Superadmin-only `GET /machines/{id}/diagnostics` puts together what is already recorded: the machine, `machine_live_status`, the open `machine_mode_log` period, `machine_events`, the variables with their addresses, and the latest values in `live.Hub`. Reading it sends nothing to the machine, so the page reloads it every 5 s.
- The server returns each check with its result, the reason and what to do. The first failed check is the answer the page shows; later ones usually follow from it.
- Superadmin-only `GET /machines/{id}/diagnostics/network` runs one ping and one TCP connect, on opening and when asked, not on a timer. It runs from the server, the place the collector connects from. It only ever tests the machine's stored address, and ping receives a parsed or resolved IP, never the stored text. The TCP connect sends no data, so controller sessions are not disturbed. Any answer counts for the IP, since many plants block ping.
- A failed test while the collector is connected only warns: the test is what is in doubt.
- Connected with no recent reading is OK. A `machine_live_status` row not updated for 30 s does not count as connected: a stopped server leaves its last value.
- The page picks variables and subscribes through the normal socket and `MachineService` rule, with no new socket event. Only what the browser alone can know (the socket's state) is worked out in the client.

## Alternatives considered

- **Connect to the controller from this page through the DLL:** a real protocol test, but a second session on controllers with few sessions, and the collector already holds one and records the result.
- **Probe from the browser:** it says nothing about the server's route to the machine.
- **Expose the collector's workers in memory:** more exact, but it ties the API to the collector for what the status row and events already show.

## Consequences

- A reason is only as new as the last status (10 s) or event. Hub values are empty until a machine's first reading after a server start, which the page says.
- An open port does not prove the controller's protocol works: the Controller check answers that.
- A server without the collector's writes (a second instance beside the real one) reads the same database and reports the other instance's status.
- New failure causes need a new case in `explainNotConnected`.
