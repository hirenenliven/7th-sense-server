package service

import (
	"context"
	"errors"
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// ClienteleService holds the business rules for clienteles: the clients a customer
// (user_id) makes parts for. Everyone reads their company's clienteles (the superadmin
// every customer's); managers change their company's (see companyOwner). A clientele
// that parts use cannot be deleted, nor moved to another customer.
type ClienteleService struct {
	clienteleRepository *repository.ClienteleRepository
	userRepository      *repository.UserRepository
}

// NewClienteleService creates a ClienteleService that stores clienteles with
// clienteleRepository and checks customers with userRepository.
func NewClienteleService(clienteleRepository *repository.ClienteleRepository, userRepository *repository.UserRepository) *ClienteleService {
	return &ClienteleService{clienteleRepository: clienteleRepository, userRepository: userRepository}
}

// GetAll returns the clienteles actor may read. A non-nil customerID keeps only that
// customer's clienteles.
func (clienteleService *ClienteleService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.Clientele, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	return clienteleService.clienteleRepository.FindAll(ctx, ownerID)
}

// GetByID returns the clientele when actor may read it, else "Clientele not found".
func (clienteleService *ClienteleService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.Clientele, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return model.Clientele{}, err
	}
	clientele, err := clienteleService.clienteleRepository.FindByID(ctx, ownerID, id)
	return clientele, toClienteleResponseError(err)
}

// CheckOwnedBy returns "Please select valid clientele" unless clientele id belongs to
// customer ownerID.
func (clienteleService *ClienteleService) CheckOwnedBy(ctx context.Context, ownerID, id int64) error {
	_, err := clienteleService.clienteleRepository.FindByID(ctx, &ownerID, id)
	if errors.Is(err, repository.ErrNotFound) {
		return response.NewError(http.StatusBadRequest, "Please select a clientele of the selected customer")
	}
	return err
}

// Create adds a clientele. A customer or admin (actor) adds it to their own company,
// whatever customer the request names; the superadmin to the customer named.
func (clienteleService *ClienteleService) Create(ctx context.Context, actor model.AuthUser, clienteleRequest dto.CreateClienteleRequest) (model.Clientele, error) {
	if err := ownCompanyCustomer(actor, &clienteleRequest.UserID); err != nil {
		return model.Clientele{}, err
	}
	if err := clienteleService.checkCustomer(ctx, clienteleRequest.UserID); err != nil {
		return model.Clientele{}, err
	}
	clientele := clienteleFromFields(clienteleRequest.ClienteleFields)
	clientele.Status = dto.ClienteleStatusActive

	id, err := clienteleService.clienteleRepository.Create(ctx, clientele)
	if err != nil {
		return model.Clientele{}, toClienteleResponseError(err)
	}
	return clienteleService.find(ctx, id)
}

// Update replaces the clientele. A blank status keeps the saved one. A customer or
// admin (actor) changes only their company's clienteles ("Clientele not found" for
// others) and cannot move one to another company.
func (clienteleService *ClienteleService) Update(ctx context.Context, actor model.AuthUser, id int64, clienteleRequest dto.UpdateClienteleRequest) (model.Clientele, error) {
	if err := ownCompanyCustomer(actor, &clienteleRequest.UserID); err != nil {
		return model.Clientele{}, err
	}
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return model.Clientele{}, err
	}
	if _, err := clienteleService.clienteleRepository.FindByID(ctx, ownerID, id); err != nil {
		return model.Clientele{}, toClienteleResponseError(err)
	}
	if err := clienteleService.checkCustomer(ctx, clienteleRequest.UserID); err != nil {
		return model.Clientele{}, err
	}
	usedElsewhere, err := clienteleService.clienteleRepository.IsUsedByPartsOfAnotherCustomer(ctx, id, clienteleRequest.UserID)
	if err != nil {
		return model.Clientele{}, err
	}
	if usedElsewhere {
		return model.Clientele{}, response.NewError(http.StatusBadRequest, "This clientele is used by parts and cannot be moved to another customer")
	}
	clientele := clienteleFromFields(clienteleRequest.ClienteleFields)
	clientele.ID = id
	clientele.Status = clienteleRequest.Status

	if err := clienteleService.clienteleRepository.Update(ctx, clientele); err != nil {
		return model.Clientele{}, toClienteleResponseError(err)
	}
	return clienteleService.find(ctx, id)
}

// Delete soft-deletes the clientele unless parts use it. A customer or admin (actor)
// deletes only their company's clienteles ("Clientele not found" for others).
func (clienteleService *ClienteleService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return err
	}
	if _, err := clienteleService.clienteleRepository.FindByID(ctx, ownerID, id); err != nil {
		return toClienteleResponseError(err)
	}
	used, err := clienteleService.clienteleRepository.IsUsedByParts(ctx, id)
	if err != nil {
		return err
	}
	if used {
		return response.NewError(http.StatusBadRequest, "This clientele is used by parts and cannot be deleted")
	}
	return toClienteleResponseError(clienteleService.clienteleRepository.Delete(ctx, ownerID, id))
}

// find returns clientele id whoever owns it (for the responses of Create and Update).
func (clienteleService *ClienteleService) find(ctx context.Context, id int64) (model.Clientele, error) {
	clientele, err := clienteleService.clienteleRepository.FindByID(ctx, nil, id)
	return clientele, toClienteleResponseError(err)
}

// checkCustomer returns "Please select valid customer" unless userID is a customer.
func (clienteleService *ClienteleService) checkCustomer(ctx context.Context, userID int64) error {
	isCustomer, err := clienteleService.userRepository.IsCustomer(ctx, userID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	return nil
}

func clienteleFromFields(fields dto.ClienteleFields) model.Clientele {
	return model.Clientele{
		UserID:        fields.UserID,
		ClientCode:    fields.ClientCode,
		Name:          fields.Name,
		Email:         fields.Email,
		ContactNumber: fields.ContactNumber,
		Address:       fields.Address,
		CountryName:   fields.CountryName,
		StateName:     fields.StateName,
		City:          fields.City,
	}
}

// toClienteleResponseError maps repository errors to clientele responses. Other
// errors pass through and are sent as ServerError.
func toClienteleResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Clientele not found")
	case errors.Is(err, repository.ErrDuplicateCode):
		return response.NewError(http.StatusBadRequest, "This customer already has a clientele with this client ID")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	default:
		return err
	}
}
