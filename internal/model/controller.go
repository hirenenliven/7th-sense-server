package model

import "time"

// Controller is a row of the controllers table.
type Controller struct {
	ID          int64      `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Code        string     `db:"code" json:"code"`
	Variant     string     `db:"variant" json:"variant"`
	BrandID     *int64     `db:"brand_id" json:"brandId"`
	Description string     `db:"description" json:"description"`
	IsDeleted   bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}
