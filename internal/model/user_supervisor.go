package model

import "time"

// UserSupervisor is a row of the user_supervisor table: a supervisor of one user, under one of the user's cells.
type UserSupervisor struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"userId"`
	UserCellID int64      `db:"user_cell_id" json:"userCellId"`
	CellName   string     `db:"cell_name" json:"cellName"` // read from user_cell
	Name       string     `db:"name" json:"name"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
