# 0028. Telemetry queue tuning settings come from the environment

- **Date:** 2026-11-25
- **Status:** Proposed

## Context

[ADR 0024](0024-rabbitmq-telemetry-queue.md) kept only `RABBITMQ_ENABLED`, `RABBITMQ_URL` and `RABBITMQ_CONSUMERS` as settings and fixed the queue name, queue type, prefetch, batch size, publish buffer and reconnect wait in `config.Load`. Since then the batch values were changed in code during the overload work ([ADR 0026](0026-telemetry-postgres-write-tuning.md), [ADR 0027](0027-telemetry-overload-and-capacity.md)), and every try needed a rebuild. Different machines (disk, broker host) want different values.

## Decision

`config.Load` reads these from the environment; the defaults are today's values, so an unchanged `.env` behaves as before:

| Variable | Default | Field |
|---|---|---|
| `RABBITMQ_QUEUE` | `telemetry.values` | `Queue` (dead letters `<queue>.dead`) |
| `RABBITMQ_QUEUE_TYPE` | `quorum` | `QueueType`: `quorum` or `classic`; anything else is logged and quorum |
| `RABBITMQ_PREFETCH` | `1000` | `Prefetch` |
| `RABBITMQ_BATCH_ROWS` | `10000` | `BatchRows` |
| `RABBITMQ_BATCH_MS` | `1000` | `BatchInterval` |
| `RABBITMQ_PUBLISH_BUFFER` | `20000` | `PublishBuffer` |
| `RABBITMQ_RECONNECT_MAX_SEC` | `60` | `ReconnectMax` |

A missing, non-numeric or non-positive number falls back to its default (`getEnvInt`). The exchange `telemetry` and routing key `values` stay fixed. A restart applies changes.

## Alternatives considered

- **Keep them fixed (ADR 0024):** every tuning try is a code change and a rebuild.
- **Validate ranges (e.g. prefetch ≥ batch rows):** not done; the defaults are sized together and the comments in `RabbitMQConfig` say how they relate.

## Consequences

- Tuning on a site needs only `.env` and a restart.
- `RABBITMQ_QUEUE_TYPE` on an existing queue of the other type makes the declare fail (PRECONDITION_FAILED) and the publisher keeps reconnecting while readings are written directly: delete the queue or use another `RABBITMQ_QUEUE`.
- Seven more variables to document (`.env.example`, `docs/cnc-data-collection-flow.md`, `.claude/CLAUDE.md`).
