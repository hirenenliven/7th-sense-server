package model

import "time"

// UserLocation is a row of the user_location table: a location of one user, under a branch (branches table).
type UserLocation struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"userId"`
	BranchID   int64      `db:"branch_id" json:"branchId"`
	BranchName string     `db:"branch_name" json:"branchName"` // read from branches
	Name       string     `db:"name" json:"name"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
