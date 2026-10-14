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

// A customer is a company. Its own users (members) have users.company_id = the
// customer's id. Only the superadmin works across companies; a customer and its
// admins manage their company's users, and see all of its machines.

// MemberRoles are the roles of a customer's own users.
var MemberRoles = []string{RoleAdmin, RoleSupervisor, RoleOperator, RoleQC, RoleShopHead}

// UserManagerRoles may use the user routes (/users). Which users and roles each
// of them may manage is decided by AssignableRoles and the company (service.UserService).
var UserManagerRoles = []string{RoleSuperAdmin, RoleCustomer, RoleAdmin}

// AssignableRoles returns the roles a user with role may give: the superadmin
// customers and every member role, a customer every member role, an admin every
// member role but admin, nobody else any. Nobody gives the superadmin role.
func AssignableRoles(role string) []string {
	switch role {
	case RoleSuperAdmin:
		return append([]string{RoleCustomer}, MemberRoles...)
	case RoleCustomer:
		return MemberRoles
	case RoleAdmin:
		return slices.DeleteFunc(slices.Clone(MemberRoles), func(memberRole string) bool { return memberRole == RoleAdmin })
	default:
		return nil
	}
}

// CanAssignRole reports whether a user with actorRole may give role.
func CanAssignRole(actorRole, role string) bool {
	return slices.Contains(AssignableRoles(actorRole), role)
}

// SeesAllCompanyMachines reports whether a user with role sees every machine of
// their company: the customer and its admins. Other members see only the machines
// assigned to them (user_machine_access).
func SeesAllCompanyMachines(role string) bool {
	return role == RoleCustomer || role == RoleAdmin
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
	// CompanyID is the customer a member belongs to (users.company_id); nil for
	// customers and the superadmin.
	CompanyID *int64 `db:"company_id" json:"companyId"`
}

// Company returns the customer whose data the user works on: the customer
// itself, or a member's company. It is nil for the superadmin and for users that
// belong to no company.
func (user AuthUser) Company() *int64 {
	if user.Role == RoleCustomer {
		return &user.ID
	}
	return user.CompanyID
}
