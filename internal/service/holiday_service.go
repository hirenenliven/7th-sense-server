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

// HolidayService holds the business rules for a customer's (user_id) holidays and week
// offs. Everyone reads their company's (the superadmin every customer's); managers change
// their company's (see companyOwner). A weekday is a customer's week off once.
type HolidayService struct {
	holidayRepository *repository.HolidayRepository
	userRepository    *repository.UserRepository
}

// NewHolidayService creates a HolidayService that stores holidays with
// holidayRepository and checks customers with userRepository.
func NewHolidayService(holidayRepository *repository.HolidayRepository, userRepository *repository.UserRepository) *HolidayService {
	return &HolidayService{holidayRepository: holidayRepository, userRepository: userRepository}
}

// GetAll returns the holidays and week offs actor may read. A non-nil customerID keeps
// only that customer's.
func (holidayService *HolidayService) GetAll(ctx context.Context, actor model.AuthUser, customerID *int64) ([]model.Holiday, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return nil, err
	}
	return holidayService.holidayRepository.FindAll(ctx, ownerID)
}

// GetByID returns the holiday when actor may read it, else "Holiday not found".
func (holidayService *HolidayService) GetByID(ctx context.Context, actor model.AuthUser, customerID *int64, id int64) (model.Holiday, error) {
	ownerID, err := companyOwner(actor, customerID)
	if err != nil {
		return model.Holiday{}, err
	}
	holiday, err := holidayService.holidayRepository.FindByID(ctx, ownerID, id)
	return holiday, toHolidayResponseError(err)
}

// Create adds a holiday or week off. A customer or admin (actor) adds it to their own
// company, whatever customer the request names; the superadmin to the customer named.
func (holidayService *HolidayService) Create(ctx context.Context, actor model.AuthUser, holidayRequest dto.CreateHolidayRequest) (model.Holiday, error) {
	if err := ownCompanyCustomer(actor, &holidayRequest.UserID); err != nil {
		return model.Holiday{}, err
	}
	if err := holidayService.checkCustomer(ctx, holidayRequest.UserID); err != nil {
		return model.Holiday{}, err
	}
	holiday := holidayFromFields(holidayRequest.HolidayFields)
	holiday.Status = dto.HolidayStatusActive

	id, err := holidayService.holidayRepository.Create(ctx, holiday)
	if err != nil {
		return model.Holiday{}, toHolidayResponseError(err)
	}
	return holidayService.find(ctx, id)
}

// Update replaces the holiday (it may change between holiday and week off). A blank
// status keeps the saved one. A customer or admin (actor) changes only their company's
// ("Holiday not found" for others) and cannot move one to another company.
func (holidayService *HolidayService) Update(ctx context.Context, actor model.AuthUser, id int64, holidayRequest dto.UpdateHolidayRequest) (model.Holiday, error) {
	if err := ownCompanyCustomer(actor, &holidayRequest.UserID); err != nil {
		return model.Holiday{}, err
	}
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return model.Holiday{}, err
	}
	if _, err := holidayService.holidayRepository.FindByID(ctx, ownerID, id); err != nil {
		return model.Holiday{}, toHolidayResponseError(err)
	}
	if err := holidayService.checkCustomer(ctx, holidayRequest.UserID); err != nil {
		return model.Holiday{}, err
	}
	holiday := holidayFromFields(holidayRequest.HolidayFields)
	holiday.ID = id
	holiday.Status = holidayRequest.Status

	if err := holidayService.holidayRepository.Update(ctx, holiday); err != nil {
		return model.Holiday{}, toHolidayResponseError(err)
	}
	return holidayService.find(ctx, id)
}

// Delete soft-deletes the holiday. A customer or admin (actor) deletes only their
// company's ("Holiday not found" for others).
func (holidayService *HolidayService) Delete(ctx context.Context, actor model.AuthUser, id int64) error {
	ownerID, err := companyOwner(actor, nil)
	if err != nil {
		return err
	}
	return toHolidayResponseError(holidayService.holidayRepository.Delete(ctx, ownerID, id))
}

// find returns holiday id whoever owns it (for the responses of Create and Update).
func (holidayService *HolidayService) find(ctx context.Context, id int64) (model.Holiday, error) {
	holiday, err := holidayService.holidayRepository.FindByID(ctx, nil, id)
	return holiday, toHolidayResponseError(err)
}

// checkCustomer returns "Please select valid customer" unless userID is a customer.
func (holidayService *HolidayService) checkCustomer(ctx context.Context, userID int64) error {
	isCustomer, err := holidayService.userRepository.IsCustomer(ctx, userID)
	if err != nil {
		return err
	}
	if !isCustomer {
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	}
	return nil
}

// holidayFromFields copies the form fields; the fields the kind does not use become NULL.
func holidayFromFields(fields dto.HolidayFields) model.Holiday {
	holiday := model.Holiday{
		UserID:      fields.UserID,
		HolidayType: fields.HolidayType,
		Name:        fields.Name,
		StartDate:   optionalText(fields.StartDate),
		EndDate:     optionalText(fields.EndDate),
	}
	if fields.HolidayType == dto.HolidayTypeWeekOff {
		weekDay := fields.WeekDay
		holiday.WeekDay = &weekDay
	}
	return holiday
}

// toHolidayResponseError maps repository errors to holiday responses. Other errors pass
// through and are sent as ServerError.
func toHolidayResponseError(err error) error {
	switch {
	case errors.Is(err, repository.ErrNotFound):
		return response.NewError(http.StatusNotFound, "Holiday not found")
	case errors.Is(err, repository.ErrDuplicateWeekOff):
		return response.NewError(http.StatusBadRequest, "This customer already has a week off on this day")
	case errors.Is(err, repository.ErrDuplicateName):
		return response.NewError(http.StatusBadRequest, "This customer already has this holiday on this date")
	case errors.Is(err, repository.ErrInvalidUser):
		return response.NewError(http.StatusBadRequest, "Please select valid customer")
	default:
		return err
	}
}
