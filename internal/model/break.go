package model

import "time"

// Break is a row of the breaks table: a break of a customer (UserID), taken by the
// shift ShiftID (nil while no shift takes it), plus the display names of its customer
// and shift (joins, never written back). Times are "HH:MM" (24 h); a break whose end is
// not after its start runs past midnight.
type Break struct {
	ID        int64      `db:"id" json:"id"`
	UserID    int64      `db:"user_id" json:"userId"`
	ShiftID   *int64     `db:"shift_id" json:"shiftId"`
	Name      string     `db:"name" json:"name"`
	StartTime string     `db:"start_time" json:"startTime"`
	EndTime   string     `db:"end_time" json:"endTime"`
	Status    string     `db:"status" json:"status"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (joins)
	CustomerName string `db:"customer_name" json:"customerName"`
	CompanyName  string `db:"company_name" json:"companyName"`
	ShiftName    string `db:"shift_name" json:"shiftName"`
}
