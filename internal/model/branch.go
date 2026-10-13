package model

import "time"

// Branch is a row of the branches table. A branch can belong to a city of the
// organization hierarchy (user_city); its locations are in user_location.
type Branch struct {
	ID         int64      `db:"id" json:"id"`
	Code       string     `db:"code" json:"code"`
	Name       string     `db:"name" json:"name"`
	Email      string     `db:"email" json:"email"`
	Address    string     `db:"address" json:"address"`
	City       string     `db:"city" json:"city"`
	StateID    *int64     `db:"state_id" json:"stateId"`
	UserCityID *int64     `db:"user_city_id" json:"userCityId"`
	CityName   string     `db:"city_name" json:"cityName"` // read from user_city
	CustomerID *int64     `db:"customer_id" json:"customerId"`
	GSTNumber  string     `db:"gst_number" json:"gstNumber"`
	PANNumber  string     `db:"pan_number" json:"panNumber"`
	ColorCode  string     `db:"color_code" json:"colorCode"`
	Status     string     `db:"status" json:"status"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
