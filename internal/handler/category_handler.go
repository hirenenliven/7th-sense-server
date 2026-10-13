package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// CategoryHandler handles HTTP requests for categories.
type CategoryHandler struct {
	service *service.CategoryService
}

// NewCategoryHandler creates a CategoryHandler that uses categoryService.
func NewCategoryHandler(categoryService *service.CategoryService) *CategoryHandler {
	return &CategoryHandler{service: categoryService}
}

// GetAll handles GET /api/v1/categories.
func (categoryHandler *CategoryHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	categories, err := categoryHandler.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, categories)
}

// GetByID handles GET /api/v1/categories/{id}.
func (categoryHandler *CategoryHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	category, err := categoryHandler.service.GetByID(request.Context(), id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, category)
}

// Create handles POST /api/v1/categories.
func (categoryHandler *CategoryHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var categoryRequest dto.CreateCategoryRequest
	if err := decodeAndValidate(writer, request, &categoryRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	category, err := categoryHandler.service.Create(request.Context(), categoryRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, category)
}

// Update handles PUT /api/v1/categories/{id}.
func (categoryHandler *CategoryHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var categoryRequest dto.UpdateCategoryRequest
	if err := decodeAndValidate(writer, request, &categoryRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	category, err := categoryHandler.service.Update(request.Context(), id, categoryRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, category)
}

// Delete handles DELETE /api/v1/categories/{id}.
func (categoryHandler *CategoryHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := categoryHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Category deleted")
}
