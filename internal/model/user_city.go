package model

import "time"

// UserCity is a row of the user_city table: a city of one user, under one of the user's states.
type UserCity struct {
	ID          int64      `db:"id" json:"id"`
	UserID      int64      `db:"user_id" json:"userId"`
	UserStateID int64      `db:"user_state_id" json:"userStateId"`
	StateName   string     `db:"state_name" json:"stateName"` // read from user_state
	Name        string     `db:"name" json:"name"`
	IsDeleted   bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}
