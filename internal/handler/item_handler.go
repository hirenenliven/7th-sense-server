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
func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.service.GetAll(r.Context())
	if err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusOK, items)
}

// GetByID handles GET /api/v1/items/{id}.
func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	item, err := h.service.GetByID(r.Context(), r.PathValue("id"))
	if err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusOK, item)
}

// Create handles POST /api/v1/items.
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req dto.ItemRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		response.SendError(w, err)
		return
	}

	item, err := h.service.Create(r.Context(), req)
	if err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusCreated, item)
}

// Update handles PUT /api/v1/items/{id}.
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	var req dto.ItemRequest
	if err := decodeAndValidate(w, r, &req); err != nil {
		response.SendError(w, err)
		return
	}

	item, err := h.service.Update(r.Context(), r.PathValue("id"), req)
	if err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusOK, item)
}

// Delete handles DELETE /api/v1/items/{id}.
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if err := h.service.Delete(r.Context(), r.PathValue("id")); err != nil {
		response.SendError(w, err)
		return
	}
	response.Send(w, http.StatusOK, "Item deleted")
}
