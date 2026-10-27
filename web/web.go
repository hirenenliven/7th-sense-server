// Package web holds the static pages the server serves without a session.
package web

import "embed"

// ModeChart is the machine mode chart page (mode-chart/index.html and mode-chart.js), served
// at /mode-chart/. It reads GET /api/v1/health/machines/{id}/modes.
//
//go:embed mode-chart
var ModeChart embed.FS
