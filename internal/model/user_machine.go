package model

import "time"

// UserMachine is a row of the user_machines table.
type UserMachine struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"userId"`
	MachineID  int64      `db:"machine_id" json:"machineId"`
	AssignedAt time.Time  `db:"assigned_at" json:"assignedAt"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
}
