package model

import "time"

// UserState is a row of the user_state table: a state of one user, under one of the user's countries.
type UserState struct {
	ID            int64      `db:"id" json:"id"`
	UserID        int64      `db:"user_id" json:"userId"`
	UserCountryID int64      `db:"user_country_id" json:"userCountryId"`
	CountryName   string     `db:"country_name" json:"countryName"` // read from user_country
	Name          string     `db:"name" json:"name"`
	IsDeleted     bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updatedAt"`
}
