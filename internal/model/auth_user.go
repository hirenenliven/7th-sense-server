package model

import "slices"

// Role names of the roles table. Authorization checks use these names.
const (
	RoleSuperAdmin = "superadmin"
	RoleAdmin      = "admin"
	RoleSupervisor = "supervisor"
	RoleOperator   = "operator"
	RoleCustomer   = "customer"
	RoleQC         = "qc"
	RoleShopHead   = "shop_head"
)

// AdminRoles are the roles with full access (admin-only routes, acting for other users).
var AdminRoles = []string{RoleSuperAdmin, RoleAdmin}

// IsAdmin reports whether role is one of AdminRoles.
func IsAdmin(role string) bool {
	return slices.Contains(AdminRoles, role)
}

// AuthUser is a user as seen by login and sessions: the users row with its
// role name (roles.name). Password holds the bcrypt hash and is never sent to
// clients; it is empty when the user was read from a session.
type AuthUser struct {
	ID        int64  `db:"id" json:"id"`
	Username  string `db:"username" json:"username"`
	Email     string `db:"email" json:"email"`
	Password  string `db:"password" json:"-"`
	FirstName string `db:"first_name" json:"firstName"`
	LastName  string `db:"last_name" json:"lastName"`
	Status    string `db:"status" json:"status"`
	Role      string `db:"role" json:"role"`
}
