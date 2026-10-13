package model

import "time"

// UserLine is a row of the user_line table: a line of one user, under one of the user's shops.
type UserLine struct {
	ID         int64      `db:"id" json:"id"`
	UserID     int64      `db:"user_id" json:"userId"`
	UserShopID int64      `db:"user_shop_id" json:"userShopId"`
	ShopName   string     `db:"shop_name" json:"shopName"` // read from user_shop
	Name       string     `db:"name" json:"name"`
	IsDeleted  bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt  *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt  time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt  time.Time  `db:"updated_at" json:"updatedAt"`
}
