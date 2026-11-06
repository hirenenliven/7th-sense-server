package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// LicenceKeyHandler handles HTTP requests for licence keys. The superadmin generates,
// saves, edits and deletes them; a customer and its admins read their company's
// (service.LicenceKeyService decides which).
type LicenceKeyHandler struct {
	service *service.LicenceKeyService
}

// NewLicenceKeyHandler creates a LicenceKeyHandler that uses licenceKeyService.
func NewLicenceKeyHandler(licenceKeyService *service.LicenceKeyService) *LicenceKeyHandler {
	return &LicenceKeyHandler{service: licenceKeyService}
}

// GetAll handles GET /api/v1/licence-keys[?userId=<customer id>].
//
// @Summary      Get licence keys
// @Description  Returns the licence keys the user may read. A non-nil customerID keeps only that customer's.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Licence keys, Operations Center, Plant Intelligence
// @Tags         Licence keys
// @ID           getLicenceKeys
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=[]model.LicenceKey}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/licence-keys [get]
func (licenceKeyHandler *LicenceKeyHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	customerID, err := queryID(request, "userId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	licenceKeys, err := licenceKeyHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, licenceKeys)
}

// GetByID handles GET /api/v1/licence-keys/{id}[?userId=<customer id>].
func (licenceKeyHandler *LicenceKeyHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	customerID, err := queryID(request, "userId")
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	licenceKey, err := licenceKeyHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, licenceKey)
}

// Generate handles POST /api/v1/licence-keys/generate: new keys, not saved.
//
// @Summary      Generate licence keys
// @Description  Returns count new random keys. They are not saved: the superadmin reviews them and saves them with Create.
// @Description  Roles: superadmin only (others get 403).
// @Description  Validation: Please enter valid number of keys (1 to 100).
// @Description
// @Description  Used by pages: Licence keys
// @Tags         Licence keys
// @ID           generateLicenceKeys
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.GenerateLicenceKeysRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.GeneratedLicenceKeys}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/licence-keys/generate [post]
func (licenceKeyHandler *LicenceKeyHandler) Generate(writer http.ResponseWriter, request *http.Request) {
	var generateRequest dto.GenerateLicenceKeysRequest
	if err := decodeAndValidate(writer, request, &generateRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	generated, err := licenceKeyHandler.service.Generate(generateRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, generated)
}

// Create handles POST /api/v1/licence-keys: saves generated keys for a customer.
//
// @Summary      Create licence keys
// @Description  Saves the keys for the customer named, all or none, and returns them. Each licence runs for its duration from today.
// @Description  Roles: superadmin only (others get 403).
// @Description  Validation: Please select customer; Please generate licence keys; Please save at most 100 licence keys at once; Please enter valid licence key; Licence key is repeated.
// @Description
// @Description  Used by pages: Licence keys
// @Tags         Licence keys
// @ID           createLicenceKeys
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateLicenceKeysRequest true "Request body"
// @Success      201 {object} response.Envelope{data=[]model.LicenceKey}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/licence-keys [post]
func (licenceKeyHandler *LicenceKeyHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var createRequest dto.CreateLicenceKeysRequest
	if err := decodeAndValidate(writer, request, &createRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	licenceKeys, err := licenceKeyHandler.service.Create(request.Context(), createRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, licenceKeys)
}

// Update handles PUT /api/v1/licence-keys/{id}.
//
// @Summary      Update licence key
// @Description  Changes the customer, duration and status of the key (a blank status keeps the saved one). The end date is counted again from the day the key was saved.
// @Description  Roles: superadmin only (others get 403).
// @Description  Validation: Please select customer; Please enter valid status.
// @Description
// @Description  Used by pages: Licence keys
// @Tags         Licence keys
// @ID           updateLicenceKey
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateLicenceKeyRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.LicenceKey}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/licence-keys/{id} [put]
func (licenceKeyHandler *LicenceKeyHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var updateRequest dto.UpdateLicenceKeyRequest
	if err := decodeAndValidate(writer, request, &updateRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	licenceKey, err := licenceKeyHandler.service.Update(request.Context(), id, updateRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, licenceKey)
}

// Delete handles DELETE /api/v1/licence-keys/{id}.
//
// @Summary      Delete licence key
// @Description  Soft-deletes the licence key.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Licence keys
// @Tags         Licence keys
// @ID           deleteLicenceKey
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/licence-keys/{id} [delete]
func (licenceKeyHandler *LicenceKeyHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := licenceKeyHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Licence key deleted")
}
