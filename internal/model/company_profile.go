package model

import "time"

// CompanyProfile is a row of the company_profiles table.
type CompanyProfile struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"userId"`
	GSTNumber  string     `db:"gst_number" json:"gstNumber"`
	PANNumber  string     `db:"pan_number" json:"panNumber"`
	Website    string     `db:"website" json:"website"`
	FinixID    *int64     `db:"finix_id" json:"finixId"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
