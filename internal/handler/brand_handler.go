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
//
// @Summary      Get brands
// @Description  GetAll handles GET /api/v1/brands.
// @Description
// @Description  Used by pages: Add / edit machine, Controllers, Brands
// @Tags         Machine catalog
// @ID           getBrands
// @Produce      json
// @Security     cookieAuth
// @Success      200 {object} response.Envelope{data=[]model.Brand}
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/brands [get]
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
//
// @Summary      Create brand
// @Description  Adds a machine brand. Name required (unique), code optional (unique), description optional.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Brands
// @Tags         Machine catalog
// @ID           createBrand
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        request body dto.CreateBrandRequest true "Request body"
// @Success      201 {object} response.Envelope{data=model.Brand}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/brands [post]
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
//
// @Summary      Update brand
// @Description  Replaces the brand's name, code and description.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Brands
// @Tags         Machine catalog
// @ID           updateBrand
// @Accept       json
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Param        request body dto.UpdateBrandRequest true "Request body"
// @Success      200 {object} response.Envelope{data=model.Brand}
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/brands/{id} [put]
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
//
// @Summary      Delete brand
// @Description  Soft-deletes the brand. A brand that still has controllers or is used by machines is kept: those must be deleted or changed first.
// @Description  Roles: superadmin only (others get 403).
// @Description
// @Description  Used by pages: Brands
// @Tags         Machine catalog
// @ID           deleteBrand
// @Produce      json
// @Security     cookieAuth
// @Param        id path int true "Id"
// @Success      200 {object} response.Envelope
// @Failure      400 {object} response.Envelope "Validation error"
// @Failure      401 {object} response.Envelope "Not logged in"
// @Failure      403 {object} response.Envelope "Role not allowed"
// @Failure      404 {object} response.Envelope "Not found"
// @Failure      500 {object} response.Envelope "Server error"
// @Router       /v1/brands/{id} [delete]
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
