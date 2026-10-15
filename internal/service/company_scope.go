package service

import (
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
)

// Company-owned master data (parts, clienteles) belongs to a customer (user_id).
// The superadmin works on every customer's rows; everyone else on their company's,
// and only managers change them (managersOnly routes).

// companyOwner returns the customer whose rows actor may read: the superadmin any
// (customerID narrows to one, nil = every customer), everyone else only their
// company's. Asking for another customer's rows, or having no company, is Forbidden.
func companyOwner(actor model.AuthUser, customerID *int64) (*int64, error) {
	if actor.Role == model.RoleSuperAdmin {
		return customerID, nil
	}
	company := actor.Company()
	if company == nil || (customerID != nil && *customerID != *company) {
		return nil, response.Forbidden
	}
	return company, nil
}

// ownCompanyCustomer sets *customerID to actor's company when actor is not the
// superadmin, whatever customer the request named. One without a company gets Forbidden.
func ownCompanyCustomer(actor model.AuthUser, customerID *int64) error {
	if actor.Role == model.RoleSuperAdmin {
		return nil
	}
	company := actor.Company()
	if company == nil {
		return response.Forbidden
	}
	*customerID = *company
	return nil
}
