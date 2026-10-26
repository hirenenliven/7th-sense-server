//go:build !(windows && 386)

package cnc

import (
	"7th-sense-server/internal/config"
	"7th-sense-server/internal/model"
)

// dllDriver is unavailable outside windows/386: the DLL is 32-bit Windows code.
type dllDriver struct{}

func loadDLL(config.CNCConfig) (*dllDriver, error) {
	return nil, ErrNotSupported
}

func (*dllDriver) connectMachine(model.MonitoredMachine) (*Connection, int, error) {
	return nil, 0, ErrNotSupported
}
