package model

import "time"

// Cycle statuses (machine_cycles.status).
const (
	CycleRunning    = "running"
	CycleCompleted  = "completed"
	CycleUnfinished = "unfinished"
)

// Why a cycle closed (machine_cycles.end_reason, migration 030). The first three complete the
// cycle, the others leave it unfinished (CycleStatusFor).
const (
	CycleEndPartCount      = "part_count"      // ProgramPartCount stepped up
	CycleEndProgramLoop    = "program_loop"    // CurProgLine jumped back to the start (M99)
	CycleEndProgramEnd     = "program_end"     // running -> reset without an alarm (M30 / M02)
	CycleEndAlarmReset     = "alarm_reset"     // reset while an alarm or E-stop was active
	CycleEndOperatorReset  = "operator_reset"  // held -> reset: the operator aborted the part
	CycleEndProgramChanged = "program_changed" // another program name
	CycleEndOfflineReset   = "offline_reset"   // reset found after the machine was Offline
	CycleEndCollectorStop  = "collector_stop"  // reset found after the server restarted
	CycleEndTimeout        = "timeout"         // open longer than CNC_CYCLE_MAX_HOURS
)

// CycleStatusFor returns the status of a cycle closed for endReason.
func CycleStatusFor(endReason string) string {
	switch endReason {
	case CycleEndPartCount, CycleEndProgramLoop, CycleEndProgramEnd:
		return CycleCompleted
	}
	return CycleUnfinished
}

// MachineCycle is one row of machine_cycles: one run of the part program on a machine, written
// by the CNC collector. The job fields are nil when no job matched at cycle start; EndedAt and
// EndReason are empty while it runs. MachineName, JobID, PartName, ProcessName, ShiftName and
// DurationSec are read for the API (DurationSec runs up to now while the cycle is open).
type MachineCycle struct {
	ID               int64      `db:"id" json:"id"`
	MachineID        int64      `db:"machine_id" json:"machineId"`
	MachineName      string     `db:"machine_name" json:"machineName"`
	ProgramName      string     `db:"program_name" json:"programName"`
	ScheduleItemID   *int64     `db:"schedule_item_id" json:"scheduleItemId"`
	JobID            string     `db:"job_id" json:"jobId"`
	PartID           *int64     `db:"part_id" json:"partId"`
	PartName         string     `db:"part_name" json:"partName"`
	ProcessID        *int64     `db:"process_id" json:"processId"`
	ProcessName      string     `db:"process_name" json:"processName"`
	ShiftID          *int64     `db:"shift_id" json:"shiftId"`
	ShiftName        string     `db:"shift_name" json:"shiftName"`
	ShiftDate        string     `db:"shift_date" json:"shiftDate"` // YYYY-MM-DD
	StandardCycleSec *int       `db:"standard_cycle_sec" json:"standardCycleSec"`
	Status           string     `db:"status" json:"status"`
	EndReason        string     `db:"end_reason" json:"endReason"`
	StartedAt        time.Time  `db:"started_at" json:"startedAt"`
	EndedAt          *time.Time `db:"ended_at" json:"endedAt"`
	DurationSec      int        `db:"duration_sec" json:"durationSec"`
	ProductionSec    int        `db:"production_sec" json:"productionSec"`
	HoldSec          int        `db:"hold_sec" json:"holdSec"`
	DowntimeSec      int        `db:"downtime_sec" json:"downtimeSec"`
	OfflineSec       int        `db:"offline_sec" json:"offlineSec"`
	LoadingSec       *int       `db:"loading_sec" json:"loadingSec"`
	PartCountStart   *int       `db:"part_count_start" json:"partCountStart"`
	PartCountEnd     *int       `db:"part_count_end" json:"partCountEnd"`
	PartsMade        int        `db:"parts_made" json:"partsMade"`
	CreatedAt        time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt        time.Time  `db:"updated_at" json:"updatedAt"`
}

// CycleJob is what a cycle is matched to when it opens: a planned job (production schedule
// item) or, without one, the part process whose program the machine runs, and the shift.
type CycleJob struct {
	ScheduleItemID   *int64 `db:"schedule_item_id"`
	PartID           *int64 `db:"part_id"`
	ProcessID        *int64 `db:"process_id"`
	ShiftID          *int64 `db:"shift_id"`
	ShiftDate        string `db:"shift_date"` // YYYY-MM-DD, the day the shift started
	StandardCycleSec *int   `db:"standard_cycle_sec"`
}
