# 0019. Axis values: every stored value when it fits, and zoom by loading the window again

- **Date:** 2026-11-18
- **Status:** Accepted (supersedes the "always downsample" part of 0018)

## Context

ADR 0018 cut every range of the Machine Health X / Y / Z charts into 600 time buckets, so the 1-hour chart showed "1 point = 6 s" (the last value of each 6 s). The team wants to see the stored values themselves, and to zoom into the chart to read them properly. A busy machine stores about 12,000–15,000 values per axis per hour and about 35,000 per day (machine 195); 30 days is over a million per axis, more than a browser chart can draw.

## Decision

- **Every stored value while it fits**: the server first reads at most `dto.AxisValuesMaxRawValues + 1` (20,001) values per axis (`MachineHealthRepository.FindAxisRows`, one index range scan per axis with `LIMIT`). If no axis is over 20,000, it returns them as they are (`raw: true`, `bucketSec: 0`, each point min = max = value, count 1). Otherwise it falls back to the 600 buckets of 0018 (`raw: false`).
- **Zoom = an exact window**: `start` / `end` (RFC 3339) on the same endpoint (`dto.AxisValuesRequest.checkWindow`, ≤ 31 days, end cut to now). The client's drag-to-zoom asks for the dragged window, so zooming into a grouped day down to minutes shows every stored value there, not a stretched picture of buckets.
- **Client**: drag across any of the three charts selects a time window (`LineTrendChart onSelectRange`, shaded while dragging); all three charts follow; Zoom out (twice as wide, around the same centre) and Reset zoom. Range buttons 5 / 15 minutes added; a zoomed window does not refresh on the timer.

## Alternatives considered

- **Always raw.** 30 days would send millions of points per axis. Rejected.
- **Recharts Brush (range slider) only.** Zooms the picture of what was loaded; a grouped range stays grouped, and the overview chart doubles the drawing cost of thousands of points. Rejected in favour of loading the window.
- **Raise the bucket count** (e.g. 5,000). Still not the stored values. Rejected.

## Consequences

- Up to 12 hours on machine 195 is raw today (~80–115 ms database time for ~40,000 values over three axes); 1 day and longer are bucketed until zoomed in.
- A range just over the limit costs both queries (raw read of 20,001 per axis, then the buckets); still under 1 s locally for 1 day.
- Drawing 20,000 points per chart × 3 is the heaviest case for the browser; animation is off above 1,000 points.
