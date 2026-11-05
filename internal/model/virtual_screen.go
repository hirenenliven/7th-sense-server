package model

import "time"

// LiveVariable is a machine's latest value of one variable (machine_variable_values), by the
// variable's code (cnc_variables.code).
type LiveVariable struct {
	Code      string    `db:"code"`
	Value     *float64  `db:"value"`
	UpdatedAt time.Time `db:"updated_at"`
}

// MachineAlarm is the alarm a machine raised last (machine_events alarm_on). Active while no
// alarm_off followed it.
type MachineAlarm struct {
	Code    *int      `db:"code" json:"code"`
	Message string    `db:"message" json:"message"`
	At      time.Time `db:"created_at" json:"at"`
	Active  bool      `db:"active" json:"active"`
}

// JobProgress is a production schedule item (a job) with its part, process and clientele names,
// its planned cycle time and the cycles completed for it so far.
type JobProgress struct {
	ScheduleItemID   int64  `db:"schedule_item_id" json:"scheduleItemId"`
	JobID            string `db:"job_id" json:"jobId"`
	ScheduleDate     string `db:"schedule_date" json:"scheduleDate"` // YYYY-MM-DD
	StartTime        string `db:"start_time" json:"startTime"`       // HH:mm
	EndTime          string `db:"end_time" json:"endTime"`           // HH:mm, not after the start = past midnight
	Quantity         int    `db:"quantity" json:"quantity"`
	PartName         string `db:"part_name" json:"partName"`
	ProcessName      string `db:"process_name" json:"processName"`
	ClienteleName    string `db:"clientele_name" json:"clienteleName"`
	StandardCycleSec *int   `db:"standard_cycle_sec" json:"standardCycleSec"`
	CompletedCycles  int    `db:"completed_cycles" json:"completedCycles"`
}
