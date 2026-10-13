package model

import "time"

// UserCountry is a row of the user_country table: a country of one user.
type UserCountry struct {
	ID        int64      `db:"id" json:"id"`
	UserID    int64      `db:"user_id" json:"userId"`
	Name      string     `db:"name" json:"name"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
}
