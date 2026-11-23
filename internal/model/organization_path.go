package model

// Organization levels of a company, top-down. The names are also the filters a customer picks on
// the customer form (user_profiles.filters) and the levels the Add / Edit Machine form shows.
const (
	OrganizationLevelCountry    = "country"
	OrganizationLevelState      = "state"
	OrganizationLevelCity       = "city"
	OrganizationLevelBranch     = "branch"
	OrganizationLevelLocation   = "location"
	OrganizationLevelShop       = "shop"
	OrganizationLevelLine       = "line"
	OrganizationLevelCell       = "cell"
	OrganizationLevelSupervisor = "supervisor"
)

// OrganizationLevels lists the levels in hierarchy order (each one's parent comes before it).
var OrganizationLevels = []string{
	OrganizationLevelCountry, OrganizationLevelState, OrganizationLevelCity, OrganizationLevelBranch,
	OrganizationLevelLocation, OrganizationLevelShop, OrganizationLevelLine, OrganizationLevelCell,
	OrganizationLevelSupervisor,
}

// OrganizationPath names one row per organization level (nil = none): a machine's place in its
// customer's organization (machines, migration 040). The JSON names match the organization APIs.
type OrganizationPath struct {
	UserCountryID    *int64 `json:"userCountryId"`
	UserStateID      *int64 `json:"userStateId"`
	UserCityID       *int64 `json:"userCityId"`
	BranchID         *int64 `json:"branchId"`
	UserLocationID   *int64 `json:"userLocationId"`
	UserShopID       *int64 `json:"userShopId"`
	UserLineID       *int64 `json:"userLineId"`
	UserCellID       *int64 `json:"userCellId"`
	UserSupervisorID *int64 `json:"userSupervisorId"`
}

// levelIDs returns the path's ids in the order of OrganizationLevels.
func (path OrganizationPath) levelIDs() []*int64 {
	return []*int64{path.UserCountryID, path.UserStateID, path.UserCityID, path.BranchID,
		path.UserLocationID, path.UserShopID, path.UserLineID, path.UserCellID, path.UserSupervisorID}
}

// Level returns the id at level (one of OrganizationLevels), nil when none or unknown.
func (path OrganizationPath) Level(level string) *int64 {
	for index, name := range OrganizationLevels {
		if name == level {
			return path.levelIDs()[index]
		}
	}
	return nil
}

// IsEmpty reports whether the path names no level.
func (path OrganizationPath) IsEmpty() bool {
	return path.Deepest() == OrganizationPath{}
}

// Deepest keeps only the lowest level the path names: a machine is saved there, the levels above
// come from the organization itself.
func (path OrganizationPath) Deepest() OrganizationPath {
	switch {
	case path.UserSupervisorID != nil:
		return OrganizationPath{UserSupervisorID: path.UserSupervisorID}
	case path.UserCellID != nil:
		return OrganizationPath{UserCellID: path.UserCellID}
	case path.UserLineID != nil:
		return OrganizationPath{UserLineID: path.UserLineID}
	case path.UserShopID != nil:
		return OrganizationPath{UserShopID: path.UserShopID}
	case path.UserLocationID != nil:
		return OrganizationPath{UserLocationID: path.UserLocationID}
	case path.BranchID != nil:
		return OrganizationPath{BranchID: path.BranchID}
	case path.UserCityID != nil:
		return OrganizationPath{UserCityID: path.UserCityID}
	case path.UserStateID != nil:
		return OrganizationPath{UserStateID: path.UserStateID}
	case path.UserCountryID != nil:
		return OrganizationPath{UserCountryID: path.UserCountryID}
	}
	return OrganizationPath{}
}

// SortOrganizationLevels returns the known level names of names in hierarchy order, dropping
// unknown ones and repeats.
func SortOrganizationLevels(names []string) []string {
	wanted := make(map[string]bool, len(names))
	for _, name := range names {
		wanted[name] = true
	}
	sorted := make([]string, 0, len(names))
	for _, level := range OrganizationLevels {
		if wanted[level] {
			sorted = append(sorted, level)
		}
	}
	return sorted
}
