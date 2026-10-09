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

// NewItemHandler creates an ItemHandler that uses svc.
func NewItemHandler(svc *service.ItemService) *ItemHandler {
	return &ItemHandler{service: svc}
}

// GetAll handles GET /api/v1/items.
func (h *ItemHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	items, err := h.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, items)
}

// GetByID handles GET /api/v1/items/{id}.
func (h *ItemHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	item, err := h.service.GetByID(request.Context(), request.PathValue("id"))
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, item)
}

// Create handles POST /api/v1/items.
func (h *ItemHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var req dto.ItemRequest
	if err := decodeAndValidate(writer, request, &req); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	item, err := h.service.Create(request.Context(), req)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, item)
}

// Update handles PUT /api/v1/items/{id}.
func (h *ItemHandler) Update(writer http.ResponseWriter, request *http.Request) {
	var req dto.ItemRequest
	if err := decodeAndValidate(writer, request, &req); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	item, err := h.service.Update(request.Context(), request.PathValue("id"), req)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, item)
}

// Delete handles DELETE /api/v1/items/{id}.
func (h *ItemHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	if err := h.service.Delete(request.Context(), request.PathValue("id")); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Item deleted")
}
