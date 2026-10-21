package model

import "time"

// DowntimeReason is a row of the downtime_reasons table: why a machine of a customer
// (UserID) was down, plus the customer's display name (join, never written back).
// ReasonCode is optional (empty when none).
type DowntimeReason struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"userId"`
	ReasonCode string     `db:"reason_code" json:"reasonCode"`
	Reason     string     `db:"reason" json:"reason"`
	Status     string     `db:"status" json:"status"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (join)
	CustomerName string `db:"customer_name" json:"customerName"`
	CompanyName  string `db:"company_name" json:"companyName"`
}
