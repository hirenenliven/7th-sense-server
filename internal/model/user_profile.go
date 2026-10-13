package model

import "time"

// UserProfile is a row of the user_profiles table.
type UserProfile struct {
	ID             int64      `db:"id" json:"id"`
	UserID         int64      `db:"user_id" json:"userId"`
	Designation    string     `db:"designation" json:"designation"`
	Grade          string     `db:"grade" json:"grade"`
	BirthDate      *time.Time `db:"birth_date" json:"birthDate"`
	JoiningDate    *time.Time `db:"joining_date" json:"joiningDate"`
	BloodGroup     string     `db:"blood_group" json:"bloodGroup"`
	ProfilePicture string     `db:"profile_picture" json:"profilePicture"`
	ContactPerson  string     `db:"contact_person" json:"contactPerson"`
	ContactPhone   string     `db:"contact_phone" json:"contactPhone"`
	Theme          string     `db:"theme" json:"theme"`
	Timezone       string     `db:"timezone" json:"timezone"`
	Language       string     `db:"language" json:"language"`
	StartDayTime   *time.Time `db:"start_day_time" json:"startDayTime"`
	Filters        []string   `db:"filters" json:"filters"`
	IsDeleted      bool       `db:"is_deleted" json:"isDeleted"`
	DeletedAt      *time.Time `db:"deleted_at" json:"deletedAt"`
	CreatedAt      time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt      time.Time  `db:"updated_at" json:"updatedAt"`
}
