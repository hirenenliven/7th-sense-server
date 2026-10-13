package handler

import (
	"net/http"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
	"7th-sense-server/internal/service"
)

// ProductHandler handles HTTP requests for products.
type ProductHandler struct {
	service *service.ProductService
}

// NewProductHandler creates a ProductHandler that uses productService.
func NewProductHandler(productService *service.ProductService) *ProductHandler {
	return &ProductHandler{service: productService}
}

// GetAll handles GET /api/v1/products.
func (productHandler *ProductHandler) GetAll(writer http.ResponseWriter, request *http.Request) {
	products, err := productHandler.service.GetAll(request.Context())
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, products)
}

// GetByID handles GET /api/v1/products/{id}.
func (productHandler *ProductHandler) GetByID(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	product, err := productHandler.service.GetByID(request.Context(), id)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, product)
}

// Create handles POST /api/v1/products.
func (productHandler *ProductHandler) Create(writer http.ResponseWriter, request *http.Request) {
	var productRequest dto.CreateProductRequest
	if err := decodeAndValidate(writer, request, &productRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	product, err := productHandler.service.Create(request.Context(), productRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Created, product)
}

// Update handles PUT /api/v1/products/{id}.
func (productHandler *ProductHandler) Update(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	var productRequest dto.UpdateProductRequest
	if err := decodeAndValidate(writer, request, &productRequest); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	product, err := productHandler.service.Update(request.Context(), id, productRequest)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, product)
}

// Delete handles DELETE /api/v1/products/{id}.
func (productHandler *ProductHandler) Delete(writer http.ResponseWriter, request *http.Request) {
	id, err := pathID(request)
	if err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}

	if err := productHandler.service.Delete(request.Context(), id); err != nil {
		response.SendResponse(writer, request, err, nil)
		return
	}
	response.SendResponse(writer, request, response.Success, "Product deleted")
}
