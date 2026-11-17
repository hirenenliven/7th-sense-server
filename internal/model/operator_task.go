package model

import "time"

// MaintenanceTask is a row of the maintenance_tasks table: one check an operator can do on a
// machine ("Check for any Oil leakage"). The list is shared by every customer. UsageCount is the
// number of operator tasks (not deleted) that use it, read with a join.
type MaintenanceTask struct {
	ID          int64      `db:"id" json:"id"`
	Description string     `db:"description" json:"description"`
	Status      string     `db:"status" json:"status"`
	IsDeleted   bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt   *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt   time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt   time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (join)
	UsageCount int `db:"usage_count" json:"usageCount"`
}

// OperatorTask is a row of the operator_tasks table: a maintenance task that applies to a
// customer's machines of one model (the Machine Type; nil = every machine of the customer), plus
// the display names read with joins (never written back).
type OperatorTask struct {
	ID                int64      `db:"id" json:"id"`
	UserID            int64      `db:"user_id" json:"userId"`
	ModelID           *int64     `db:"model_id" json:"modelId"`
	MaintenanceTaskID int64      `db:"maintenance_task_id" json:"maintenanceTaskId"`
	IsDeleted         bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt         *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt         time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt         time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (joins)
	CustomerName   string `db:"customer_name" json:"customerName"`
	CompanyName    string `db:"company_name" json:"companyName"`
	ModelName      string `db:"model_name" json:"modelName"`
	ControllerName string `db:"controller_name" json:"controllerName"`
	BrandName      string `db:"brand_name" json:"brandName"`
	Task           string `db:"task" json:"task"`
}
