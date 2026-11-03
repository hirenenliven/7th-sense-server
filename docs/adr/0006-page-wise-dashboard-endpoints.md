# 0006. Page-wise dashboard endpoints returning only chart fields

- **Date:** 2026-11-03
- **Status:** Accepted

## Context

The Productivity and Machine Health pages need aggregated chart data; returning whole model rows wasted bandwidth.

## Decision

Serve charts at `/productivity/*` and `/machine-health/*` (`DashboardService`). Each response row is a DTO filled by one constructor (`dto.NewTimelinePeriod`, `NewFinishedCycle`, `NewFinishCyclesSummary`), so adding or dropping a field is a one-file change.

## Alternatives considered

- Generic endpoints returning full models
- GraphQL

## Consequences

- Small responses shaped for the charts.
- A chart change may need a server change; documented per page.
