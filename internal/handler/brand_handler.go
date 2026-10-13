package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// BrandHandler handles HTTP requests for brands.
type BrandHandler struct {
	service *service.BrandService
}

// NewBrandHandler creates a BrandHandler that uses brandService.
func NewBrandHandler(brandService *service.BrandService) *BrandHandler {
	return &BrandHandler{service: brandService}
}

// GetAll handles GET /api/v1/brands.
func (brandHandler *BrandHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	brands, err := brandHandler.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, brands)
}

// GetByID handles GET /api/v1/brands/{id}.
func (brandHandler *BrandHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	brand, err := brandHandler.service.GetByID(request.Context(), id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, brand)
}

// Create handles POST /api/v1/brands.
func (brandHandler *BrandHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var brandRequest dto.CreateBrandRequest
	if err := decodeAndValidate(writer, request, &brandRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	brand, err := brandHandler.service.Create(request.Context(), brandRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, brand)
}

// Update handles PUT /api/v1/brands/{id}.
func (brandHandler *BrandHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var brandRequest dto.UpdateBrandRequest
	if err := decodeAndValidate(writer, request, &brandRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	brand, err := brandHandler.service.Update(request.Context(), id, brandRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, brand)
}

// Delete handles DELETE /api/v1/brands/{id}.
func (brandHandler *BrandHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := brandHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Brand deleted")
}
