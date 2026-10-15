package model

import "time"

// Shift is a row of the shifts table, plus its customer's display name (join, never
// written back). Times are "HH:MM" (24 h); a shift whose end is not after its start
// runs past midnight. WorkingDays are ISO weekdays (1 = Monday … 7 = Sunday).
//
// Breaks (the customer's breaks the shift takes, from the shift start on),
// DurationMinutes and BreakMinutes are filled by the service: the minutes are worked
// out from the times and never stored.
type Shift struct {
	ID          int64      `db:"id" json:"id"`
	UserID      int64      `db:"user_id" json:"userId"`
	Name        string     `db:"name" json:"name"`
	ShiftType   string     `db:"shift_type" json:"shiftType"`
	StartTime   string     `db:"start_time" json:"startTime"`
	EndTime     string     `db:"end_time" json:"endTime"`
	WorkingDays []int32    `db:"working_days" json:"workingDays"`
	Status      string     `db:"status" json:"status"`
	IsDeleted   bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (join)
	CustomerName string `db:"customer_name" json:"customerName"`
	CompanyName  string `db:"company_name" json:"companyName"`

	Breaks          []ShiftBreak `db:"-" json:"breaks"`
	DurationMinutes int          `db:"-" json:"durationMinutes"` // start to end
	BreakMinutes    int          `db:"-" json:"breakMinutes"`    // sum of the breaks
}

// ShiftBreak is a break a shift takes (a breaks row whose shift_id is the shift).
// ID is the break's id. Times are "HH:MM" (24 h).
type ShiftBreak struct {
	ID        int64  `db:"id" json:"id"`
	ShiftID   int64  `db:"shift_id" json:"shiftId"`
	Name      string `db:"name" json:"name"`
	StartTime string `db:"start_time" json:"startTime"`
	EndTime   string `db:"end_time" json:"endTime"`
}
