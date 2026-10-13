package model

import "time"

// UserBranch is a row of the user_branch table: a branch of one user, under one of the user's cities.
type UserBranch struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"userId"`
	UserCityID int64      `db:"user_city_id" json:"userCityId"`
	CityName   string     `db:"city_name" json:"cityName"` // read from user_city
	Name       string     `db:"name" json:"name"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
