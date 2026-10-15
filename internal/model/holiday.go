package model

import "time"

// Holiday is a row of the holidays table: a customer's (UserID) holiday or week off,
// plus the customer's display name (join, never written back).
//
// A holiday (HolidayType "holiday") has a Name and StartDate / EndDate ("YYYY-MM-DD");
// a week off ("week_off") has a WeekDay (ISO, 1 = Monday … 7 = Sunday) and an optional
// Name. The fields the kind does not use are empty / nil.
type Holiday struct {
	ID          int64      `db:"id" json:"id"`
	UserID      int64      `db:"user_id" json:"userId"`
	HolidayType string     `db:"holiday_type" json:"holidayType"`
	Name        string     `db:"name" json:"name"`
	StartDate   *string    `db:"start_date" json:"startDate"`
	EndDate     *string    `db:"end_date" json:"endDate"`
	WeekDay     *int16     `db:"week_day" json:"weekDay"`
	Status      string     `db:"status" json:"status"`
	IsDeleted   bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (join)
	CustomerName string `db:"customer_name" json:"customerName"`
	CompanyName  string `db:"company_name" json:"companyName"`
}
