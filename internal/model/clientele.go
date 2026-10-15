package model

import "time"

// Clientele is a row of the clienteles table: a client a customer (UserID) makes parts
// for, plus the customer's display name (join, never written back). CountryName and
// StateName are names ("India", "Gujarat"); ContactNumber is +<code><number>.
type Clientele struct {
	ID            int64      `db:"id" json:"id"`
	UserID        int64      `db:"user_id" json:"userId"`
	ClientCode    string     `db:"client_code" json:"clientCode"`
	Name          string     `db:"name" json:"name"`
	Email         string     `db:"email" json:"email"`
	ContactNumber string     `db:"contact_number" json:"contactNumber"`
	Address       string     `db:"address" json:"address"`
	CountryName   string     `db:"country_name" json:"countryName"`
	StateName     string     `db:"state_name" json:"stateName"`
	City          string     `db:"city" json:"city"`
	Status        string     `db:"status" json:"status"`
	IsDeleted     bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (join)
	CustomerName string `db:"customer_name" json:"customerName"`
}
