package model

import "time"

// IdleReason is a row of the idle_reasons table: why a machine of a customer
// (UserID) stood idle, plus the customer's display name (join, never written back).
// ReasonCode is optional (empty when none).
type IdleReason struct {
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
