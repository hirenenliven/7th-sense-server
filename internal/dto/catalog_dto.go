package dto

import (
	"regexp"

	"7th-sense-server/internal/response"
)

// catalogCodePattern matches the optional code of categories, brands and products.
var catalogCodePattern = regexp.MustCompile(`^[A-Z0-9_-]{2,50}$`)

const (
	maxCatalogNameLength = 100
	maxDescriptionLength = 1000
)

// validateCatalogFields checks the name, code and description shared by
// categories, brands, products and controllers (controllers have no code: pass "").
// Name is required; code and description are optional.
func validateCatalogFields(name, code, description string) error {
	if name == "" || len(name) > maxCatalogNameLength {
		return response.ValidName
	}
	if code != "" && !catalogCodePattern.MatchString(code) {
		return response.ValidCode
	}
	if len(description) > maxDescriptionLength {
		return response.ValidDescription
	}
	return nil
}
