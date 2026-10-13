package repository

import "errors"

var (
	// ErrNotFound is returned when no record matches the lookup.
	ErrNotFound = errors.New("record not found")
	// ErrDuplicateEmail is returned when the email is already used by another user.
	ErrDuplicateEmail = errors.New("duplicate email")
	// ErrDuplicateUsername is returned when the username is already used by another user.
	ErrDuplicateUsername = errors.New("duplicate username")
	// ErrDuplicateCode is returned when the code is already used by another branch.
	ErrDuplicateCode = errors.New("duplicate code")
	// ErrDuplicateName is returned when the name is already used by another record of the same table.
	ErrDuplicateName = errors.New("duplicate name")
	// ErrInvalidCategory is returned when category_id does not reference an existing category.
	ErrInvalidCategory = errors.New("invalid category")
	// ErrInvalidBrand is returned when brand_id does not reference an existing brand.
	ErrInvalidBrand = errors.New("invalid brand")
	// ErrInvalidRole is returned when role_id does not reference an existing role.
	ErrInvalidRole = errors.New("invalid role")
	// ErrInvalidCountry is returned when country_id does not reference an existing country.
	ErrInvalidCountry = errors.New("invalid country")
	// ErrInvalidState is returned when state_id does not reference an existing state.
	ErrInvalidState = errors.New("invalid state")
	// ErrInvalidParent is returned when an organization row's parent id does not reference an existing row.
	ErrInvalidParent = errors.New("invalid parent")
)
