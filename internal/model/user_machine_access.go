package model

import "time"

// UserMachineAccess is a row of the user_machine_access table: user UserID may
// access machine MachineID.
type UserMachineAccess struct {
	ID        int64     `db:"id" json:"id"`
	UserID    int64     `db:"user_id" json:"userId"`
	MachineID int64     `db:"machine_id" json:"machineId"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
