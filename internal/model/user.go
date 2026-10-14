package model

import "time"

// User is a row of the users table. The db tags map columns to fields;
// Password holds the bcrypt hash and is never sent to clients.
type User struct {
	ID        int64     `db:"id" json:"id"`
	Username  string    `db:"username" json:"username"`
	Email     string    `db:"email" json:"email"`
	Password  string    `db:"password" json:"-"`
	FirstName string    `db:"first_name" json:"firstName"`
	LastName  string    `db:"last_name" json:"lastName"`
	Mobile    string    `db:"mobile" json:"mobile"`
	RoleID    int64     `db:"role_id" json:"roleId"`
	Role      string    `db:"role" json:"role"` // roles.name of RoleID (read only)
	BranchID  *int64    `db:"branch_id" json:"branchId"`
	CompanyID *int64    `db:"company_id" json:"companyId"` // customer a member belongs to
	Status    string    `db:"status" json:"status"`
	CreatedAt time.Time `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time `db:"updated_at" json:"updatedAt"`
}
