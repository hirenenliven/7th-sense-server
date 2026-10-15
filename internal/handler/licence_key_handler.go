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
