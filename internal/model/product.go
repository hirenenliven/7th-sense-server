package model

import "time"

// Product is a row of the products table.
type Product struct {
	ID          int64      `db:"id" json:"id"`
	Name        string     `db:"name" json:"name"`
	Code        string     `db:"code" json:"code"`
	CategoryID  *int64     `db:"category_id" json:"categoryId"`
	Image       string     `db:"image" json:"image"`
	Description string     `db:"description" json:"description"`
	IsDeleted   bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`
}
