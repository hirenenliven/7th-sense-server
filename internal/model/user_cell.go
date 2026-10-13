package model

import "time"

// UserCell is a row of the user_cell table: a cell of one user, under one of the user's lines.
type UserCell struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"userId"`
	UserLineID int64      `db:"user_line_id" json:"userLineId"`
	LineName   string     `db:"line_name" json:"lineName"` // read from user_line
	Name       string     `db:"name" json:"name"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
