package model

import "time"

// MachineModel is a row of the models table: a machine model of a controller
// (Brand → Controller → Model). ControllerName, BrandID and BrandName are read
// with JOINs for display and never written back.
type MachineModel struct {
	ID           int64      `db:"id" json:"id"`
	ControllerID int64      `db:"controller_id" json:"controllerId"`
	Name         string     `db:"name" json:"name"`
	IsDeleted    bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt    *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt    time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt    time.Time  `db:"updated_at" json:"updatedAt"`

	// Display fields (JOIN, read-only)
	ControllerName string `db:"controller_name" json:"controllerName"`
	BrandID        *int64 `db:"brand_id" json:"brandId"`
	BrandName      string `db:"brand_name" json:"brandName"`
}
