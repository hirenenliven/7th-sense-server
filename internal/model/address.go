package model

import "time"

// Address is a row of the addresses table.
type Address struct {
	ID            int64      `db:"id" json:"id"`
	UserID        int64      `db:"user_id" json:"userId"`
	AddressType   string     `db:"address_type" json:"addressType"`
	AddressLine1  string     `db:"address_line1" json:"addressLine1"`
	AddressLine2  string     `db:"address_line2" json:"addressLine2"`
	City          string     `db:"city" json:"city"`
	StateName     string     `db:"state_name" json:"stateName"`
	CountryName   string     `db:"country_name" json:"countryName"`
	PinCode       string     `db:"pin_code" json:"pinCode"`
	ContactPerson string     `db:"contact_person" json:"contactPerson"`
	ContactPhone  string     `db:"contact_phone" json:"contactPhone"`
	IsDefault     bool       `db:"is_default" json:"isDefault"`
	IsDeleted     bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt     *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt     time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt     time.Time  `db:"updated_at" json:"updatedAt"`
}
