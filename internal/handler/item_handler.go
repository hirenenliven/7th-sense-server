package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// ItemHandler handles HTTP requests for items.
type ItemHandler struct {
	service *service.ItemService
}

// NewItemHandler creates an ItemHandler that uses itemService.
func NewItemHandler(itemService *service.ItemService) *ItemHandler {
	return &ItemHandler{service: itemService}
}

// GetAll handles GET /api/v1/items.
func (itemHandler *ItemHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	items, err := itemHandler.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, items)
}

// GetByID handles GET /api/v1/items/{id}.
func (itemHandler *ItemHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	item, err := itemHandler.service.GetByID(request.Context(), request.PathValue("id"))
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, item)
}

// Create handles POST /api/v1/items.
func (itemHandler *ItemHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var itemRequest dto.ItemRequest
	if err := decodeAndValidate(writer, request, &itemRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	item, err := itemHandler.service.Create(request.Context(), itemRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, item)
}

// Update handles PUT /api/v1/items/{id}.
func (itemHandler *ItemHandler) Update(writer http.ResponseWriter, request *http.Request) {
	var itemRequest dto.ItemRequest
	if err := decodeAndValidate(writer, request, &itemRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	item, err := itemHandler.service.Update(request.Context(), request.PathValue("id"), itemRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, item)
}

// Delete handles DELETE /api/v1/items/{id}.
func (itemHandler *ItemHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	if err := itemHandler.service.Delete(request.Context(), request.PathValue("id")); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Item deleted")
}
