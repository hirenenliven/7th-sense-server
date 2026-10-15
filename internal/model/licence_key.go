package model

import "time"

// LicenceKey is a row of the licence_keys table: a licence key the superadmin generated
// for a customer (UserID), plus the customer's display name (join, never written back).
type LicenceKey struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"userId"`
	LicenceKey string     `db:"licence_key" json:"licenceKey"`
	Duration   string     `db:"duration" json:"duration"`    // subscription plan (dto.LicenceKeyDurations)
	ExpiresOn  *string    `db:"expires_on" json:"expiresOn"` // last valid day "YYYY-MM-DD"; nil = never ends
	Status     string     `db:"status" json:"status"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (join / computed): IsExpired once ExpiresOn is before today
	CustomerName string `db:"customer_name" json:"customerName"`
	CompanyName  string `db:"company_name" json:"companyName"`
	IsExpired    bool   `db:"is_expired" json:"isExpired"`
}

// GeneratedLicenceKeys is the response of POST /licence-keys/generate: new keys, not saved yet.
type GeneratedLicenceKeys struct {
	LicenceKeys []string `json:"licenceKeys"`
}
