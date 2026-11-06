package model

import "time"

// Kinds of reason a machine mode period takes (the master lists downtime_reasons, idle_reasons).
const (
	ReasonTypeDowntime = "downtime"
	ReasonTypeIdle     = "idle"
)

// ReasonPeriod is the whole machine_mode_log period reasons are recorded against. EndedAt is nil
// while the period is running.
type ReasonPeriod struct {
	ID        int64      `db:"id" json:"id"`
	MachineID int64      `db:"machine_id" json:"machineId"`
	ModeID    int        `db:"mode_id" json:"modeId"`
	ModeName  string     `db:"mode_name" json:"modeName"`
	StartedAt time.Time  `db:"started_at" json:"startedAt"`
	EndedAt   *time.Time `db:"ended_at" json:"endedAt"`
}

// MachineModeReason is a row of the machine_mode_reasons table: why a machine was down or idle
// for a part of a period, with the reason's master row. ReasonID, ReasonCode and Reason come from
// downtime_reasons or idle_reasons, whichever the row points at (a reason deleted from its master
// list keeps its text here).
type MachineModeReason struct {
	ID         int64      `db:"id" json:"id"`
	PeriodID   int64      `db:"period_id" json:"periodId"`
	ReasonType string     `db:"reason_type" json:"reasonType"`
	ReasonID   int64      `db:"reason_id" json:"reasonId"`
	ReasonCode string     `db:"reason_code" json:"reasonCode"`
	Reason     string     `db:"reason" json:"reason"`
	StartedAt  time.Time  `db:"started_at" json:"startedAt"`
	EndedAt    time.Time  `db:"ended_at" json:"endedAt"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
