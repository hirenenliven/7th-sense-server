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
//
// @Summary      Get clienteles
// @Description  Returns the clienteles the user may read. A non-nil customerID keeps only that customer's clienteles.
// @Description
// @Description  Used by pages: Clientele, Add / edit part, Operations Center, Plant Intelligence
// @Tags         Clienteles
// @ID           getClienteles
// @Produce      json
// @Security     cookieAuth
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=[]model.Clientele}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/clienteles [get]
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
//
// @Summary      Get clientele by id
// @Description  Returns the clientele when the user may read it, else "Clientele not found".
// @Description
// @Description  Used by pages: Clientele
// @Tags         Clienteles
// @ID           getClienteleById
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        userId query int false "User id"
// @Success      200 {object} response.Envelope{data=model.Clientele}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/clienteles/{id} [get]
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
//
// @Summary      Create clientele
// @Description  Adds a clientele. A customer or admin (the logged-in user) adds it to their own company, whatever customer the request names; the superadmin to the customer named.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Clientele
// @Tags         Clienteles
// @ID           createClientele
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateClienteleRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.Clientele}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/clienteles [post]
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
//
// @Summary      Update clientele
// @Description  Replaces the clientele. A blank status keeps the saved one.
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description  Validation: Please enter valid status.
// @Description
// @Description  Used by pages: Clientele
// @Tags         Clienteles
// @ID           updateClientele
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateClienteleRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.Clientele}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/clienteles/{id} [put]
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
//
// @Summary      Delete clientele
// @Description  Soft-deletes the clientele unless parts use it. A customer or admin (the logged-in user) deletes only their company's clienteles ("Clientele not found" for others).
// @Description  Roles: superadmin, customer, admin (others get 403).
// @Description
// @Description  Used by pages: Clientele
// @Tags         Clienteles
// @ID           deleteClientele
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/clienteles/{id} [delete]
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
