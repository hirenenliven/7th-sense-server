package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// ClienteleHandler handles HTTP requests for clienteles. The superadmin sees and
// changes every customer's clienteles; everyone else reads their company's and
// managers change them (service.ClienteleService decides which).
type ClienteleHandler struct {
	service *service.ClienteleService
}

// NewClienteleHandler creates a ClienteleHandler that uses clienteleService.
func NewClienteleHandler(clienteleService *service.ClienteleService) *ClienteleHandler {
	return &ClienteleHandler{service: clienteleService}
}

// GetAll handles GET /api/v1/clienteles[?userId=<customer id>].
func (clienteleHandler *ClienteleHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
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

	clienteles, err := clienteleHandler.service.GetAll(request.Context(), actor, customerID)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, clienteles)
}

// GetByID handles GET /api/v1/clienteles/{id}[?userId=<customer id>].
func (clienteleHandler *ClienteleHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
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

	clientele, err := clienteleHandler.service.GetByID(request.Context(), actor, customerID, id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, clientele)
}

// Create handles POST /api/v1/clienteles.
func (clienteleHandler *ClienteleHandler) Create(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var clienteleRequest dto.CreateClienteleRequest
	if err := decodeAndValidate(writer, request, &clienteleRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	clientele, err := clienteleHandler.service.Create(request.Context(), actor, clienteleRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, clientele)
}

// Update handles PUT /api/v1/clienteles/{id}.
func (clienteleHandler *ClienteleHandler) Update(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var clienteleRequest dto.UpdateClienteleRequest
	if err := decodeAndValidate(writer, request, &clienteleRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	clientele, err := clienteleHandler.service.Update(request.Context(), actor, id, clienteleRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, clientele)
}

// Delete handles DELETE /api/v1/clienteles/{id}.
func (clienteleHandler *ClienteleHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	actor, err := currentUser(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := clienteleHandler.service.Delete(request.Context(), actor, id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Clientele deleted")
}
