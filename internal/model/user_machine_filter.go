package model

import "time"

// UserMachineFilter is a row of the user_machine_filters table: the organization levels
// (Country → … → Line, nil = any) that narrow the machines of user UserID's company. With
// AllMachines the user may access every company machine on a line under every level named,
// machines added later included; without it the row only keeps the form's selection.
type UserMachineFilter struct {
	ID             int64     `db:"id" json:"id"`
	UserID         int64     `db:"user_id" json:"userId"`
	UserCountryID  *int64    `db:"user_country_id" json:"userCountryId"`
	UserStateID    *int64    `db:"user_state_id" json:"userStateId"`
	UserCityID     *int64    `db:"user_city_id" json:"userCityId"`
	BranchID       *int64    `db:"branch_id" json:"branchId"`
	UserLocationID *int64    `db:"user_location_id" json:"userLocationId"`
	UserShopID     *int64    `db:"user_shop_id" json:"userShopId"`
	UserLineID     *int64    `db:"user_line_id" json:"userLineId"`
	AllMachines    bool      `db:"all_machines" json:"allMachines"`
	CreatedAt      time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time `db:"updated_at" json:"updatedAt"`
}

// HasLevels reports whether the filter names at least one organization level.
func (filter UserMachineFilter) HasLevels() bool {
	for _, level := range []*int64{filter.UserCountryID, filter.UserStateID, filter.UserCityID, filter.BranchID,
		filter.UserLocationID, filter.UserShopID, filter.UserLineID} {
		if level != nil {
			return true
		}
	}
	return false
}
