package monitor

import (
	"fmt"

	"7th-sense-server/internal/model"
)

// The DLL sends MachineMode and ProcStatus in the controller's own numbering, which differs per
// engine (vendor sample G0_V0.0.2_DLL main.go, modeNames / statusNames):
//
//	         MachineMode                          ProcStatus
//	Fanuc    0 MDI, 1 AUTO, 3 EDIT, 5 JOG, 9 REF  0 Reset, 1 Stop, 2 HOLD, 3 START, 4 MSTR
//	         (cnc_statinfo aut)                   (cnc_statinfo run)
//	Siemens  0 JOG, 1 MDA, 2 AUTO                 0 Reset, 1 Stop, 2 Active, 3 Waiting, 4 Interrupted
//	         (/Bag/State/opMode)                  (/Channel/State/acProg)
//
// These names are for the logs. The mode rules read the same numbering per engine in
// machine_mode.go (autoOpModes, programStatesByEngine).

// engineNone is machineState.engine of a machine whose controller has no cnc_engine.
const engineNone = -1

var opModeNamesByEngine = map[int]map[int]string{
	model.EngineFanuc:   {0: "MDI", 1: "AUTO", 3: "EDIT", 5: "JOG", 9: "REF"},
	model.EngineSiemens: {0: "JOG", 1: "MDA", 2: "AUTO"},
}

var procStatusNamesByEngine = map[int]map[int]string{
	model.EngineFanuc:   {0: "Reset", 1: "Stop", 2: "HOLD", 3: "START", 4: "MSTR"},
	model.EngineSiemens: {0: "Reset", 1: "Stop", 2: "Active", 3: "Waiting", 4: "Interrupted"},
}

// engineOf returns a machine's cnc_engine, engineNone when it has none.
func engineOf(engine *int) int {
	if engine == nil {
		return engineNone
	}
	return *engine
}

// opModeName names a MachineMode value the way engine numbers it, e.g. "AUTO"; a value the
// engine's list does not have is "unknown(<value>)".
func opModeName(engine, opMode int) string {
	return ncStatusName(opModeNamesByEngine, engine, opMode)
}

// procStatusName names a ProcStatus value the way engine numbers it, e.g. "START" on Fanuc
// and "Waiting" on Siemens for 3.
func procStatusName(engine, procStatus int) string {
	return ncStatusName(procStatusNamesByEngine, engine, procStatus)
}

func ncStatusName(namesByEngine map[int]map[int]string, engine, value int) string {
	if name, found := namesByEngine[engine][value]; found {
		return name
	}
	return fmt.Sprintf("unknown(%d)", value)
}
