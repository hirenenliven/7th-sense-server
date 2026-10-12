package service

import (
	"context"
	"errors"

	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/response"
)

// AddressService holds the business rules for addresses.
type AddressService struct {
	addressRepository *repository.AddressRepository
}

// NewAddressService creates an AddressService backed by addressRepository.
func NewAddressService(addressRepository *repository.AddressRepository) *AddressService {
	return &AddressService{addressRepository: addressRepository}
}

func (addressService *AddressService) GetByUserID(ctx context.Context, userID int64) ([]model.Address, error) {
	addresses, err := addressService.addressRepository.FindByUserID(ctx, userID)
	return addresses, toAddressResponseError(err)
}

func (addressService *AddressService) Create(ctx context.Context, address model.Address) (model.Address, error) {
	createdAddress, err := addressService.addressRepository.Create(ctx, address)
	return createdAddress, toAddressResponseError(err)
}

func (addressService *AddressService) Update(ctx context.Context, address model.Address) (model.Address, error) {
	updatedAddress, err := addressService.addressRepository.Update(ctx, address)
	return updatedAddress, toAddressResponseError(err)
}

func (addressService *AddressService) Delete(ctx context.Context, id int64) error {
	return toAddressResponseError(addressService.addressRepository.Delete(ctx, id))
}

// toAddressResponseError maps repository errors to address responses. Other errors
// pass through and are sent as ServerError.
func toAddressResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NotFound
	default:
		return err
	}
}
