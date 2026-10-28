package monitor

import (
	"testing"
	"time"

	"7th-sense-server/internal/config"
)

func TestCheckConnectionLost(t *testing.T) {
	cncConfig := config.CNCConfig{Stale: 30 * time.Second, Quiet: 5 * time.Minute}
	lastReadingAt := time.Date(2026, 10, 28, 7, 4, 35, 0, time.UTC)
	tests := []struct {
		name       string
		connected  bool
		silence    time.Duration
		wantReason stopReason
		wantLost   bool
	}{
		// A held machine: the DLL said connected, then nothing changed so it sent nothing.
		{"held machine quiet past stale", true, 31 * time.Second, stopConnectionQuiet, false},
		{"held machine quiet past quiet limit", true, 5*time.Minute + time.Second, stopConnectionQuiet, true},
		{"controller lost, within stale", false, 10 * time.Second, stopConnectionStale, false},
		{"controller lost past stale", false, 31 * time.Second, stopConnectionStale, true},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			machineWorker := &worker{
				collector:              &Collector{config: cncConfig},
				lastConnectedReadingAt: lastReadingAt,
				lastReadingConnected:   test.connected,
			}
			reason, lost := machineWorker.checkConnectionLost(lastReadingAt.Add(test.silence))
			if reason != test.wantReason || lost != test.wantLost {
				t.Errorf("checkConnectionLost = (%v, %v), want (%v, %v)", reason, lost, test.wantReason, test.wantLost)
			}
		})
	}
}
