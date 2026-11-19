package model

import "time"

// MachineAccuracy is a row of the machine_accuracies table: one accuracy level a machine or a
// part process can have ("Good"). The list is shared by every customer. UsageCount is the number
// of machines and part processes (not deleted) that use it, read with a join.
type MachineAccuracy struct {
	ID        int64      `db:"id" json:"id"`
	Name      string     `db:"name" json:"name"`
	IsDeleted bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`

	// Read-only (join)
	UsageCount int `db:"usage_count" json:"usageCount"`
}
