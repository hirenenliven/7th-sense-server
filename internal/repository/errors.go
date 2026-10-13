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
)
