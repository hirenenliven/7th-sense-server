package model

import "time"

// UserShop is a row of the user_shop table: a shop of one user, under one of the user's locations.
type UserShop struct {
	ID             int64      `db:"id" json:"id"`
	UserID         int64      `db:"user_id" json:"userId"`
	UserLocationID int64      `db:"user_location_id" json:"userLocationId"`
	LocationName   string     `db:"location_name" json:"locationName"` // read from user_location
	Name           string     `db:"name" json:"name"`
	IsDeleted      bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
}
