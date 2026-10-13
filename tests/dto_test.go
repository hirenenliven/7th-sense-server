package tests

import (
	"testing"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/response"
)

func TestValidateCatalogFields(t *testing.T) {
	// Test catalog field validation through Category DTO
	tests := []struct {
		name    string
		request dto.CreateCategoryRequest
		wantErr error
	}{
		{
			name: "valid all fields",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "TC123",
				Description: "A valid description",
			},
			wantErr: nil,
		},
		{
			name: "valid with empty code",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "",
				Description: "A valid description",
			},
			wantErr: nil,
		},
		{
			name: "valid with empty description",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "TC123",
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "valid with only name",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "",
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "invalid empty name",
			request: dto.CreateCategoryRequest{
				Name:        "",
				Code:        "TC123",
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid name too long",
			request: dto.CreateCategoryRequest{
				Name:        string(make([]byte, 101)),
				Code:        "TC123",
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid code with special chars",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "TC@123",
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
		{
			name: "invalid code too short",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "T",
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
		{
			name: "invalid code too long",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        string(make([]byte, 51)),
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
		{
			name: "invalid description too long",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "TC123",
				Description: string(make([]byte, 1001)),
			},
			wantErr: response.ValidDescription,
		},
		{
			name: "valid code with underscores and hyphens",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "TC_123-456",
				Description: "A valid description",
			},
			wantErr: nil,
		},
		{
			name: "valid code with numbers only",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "12345",
				Description: "A valid description",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != tt.wantErr {
				t.Errorf("CreateCategoryRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestCreateBrandRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		request     dto.CreateBrandRequest
		wantErr     error
		checkFields func(*testing.T, dto.CreateBrandRequest)
	}{
		{
			name: "valid request with all fields",
			request: dto.CreateBrandRequest{
				Name:        "  Test Brand  ",
				Code:        "  tb123  ",
				Description: "  A valid description  ",
			},
			wantErr: nil,
			checkFields: func(t *testing.T, r dto.CreateBrandRequest) {
				if r.Name != "Test Brand" {
					t.Errorf("Name not trimmed, got %q", r.Name)
				}
				if r.Code != "TB123" {
					t.Errorf("Code not trimmed and uppercased, got %q", r.Code)
				}
				if r.Description != "A valid description" {
					t.Errorf("Description not trimmed, got %q", r.Description)
				}
			},
		},
		{
			name: "valid request with only name",
			request: dto.CreateBrandRequest{
				Name:        "Test Brand",
				Code:        "",
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "invalid empty name",
			request: dto.CreateBrandRequest{
				Name:        "  ",
				Code:        "TB123",
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid code format",
			request: dto.CreateBrandRequest{
				Name:        "Test Brand",
				Code:        "tb@123",
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != tt.wantErr {
				t.Errorf("CreateBrandRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFields != nil && err == nil {
				tt.checkFields(t, tt.request)
			}
		})
	}
}

func TestUpdateBrandRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		request     dto.UpdateBrandRequest
		wantErr     error
		checkFields func(*testing.T, dto.UpdateBrandRequest)
	}{
		{
			name: "valid request with all fields",
			request: dto.UpdateBrandRequest{
				Name:        "  Test Brand  ",
				Code:        "  tb123  ",
				Description: "  A valid description  ",
			},
			wantErr: nil,
			checkFields: func(t *testing.T, r dto.UpdateBrandRequest) {
				if r.Name != "Test Brand" {
					t.Errorf("Name not trimmed, got %q", r.Name)
				}
				if r.Code != "TB123" {
					t.Errorf("Code not trimmed and uppercased, got %q", r.Code)
				}
				if r.Description != "A valid description" {
					t.Errorf("Description not trimmed, got %q", r.Description)
				}
			},
		},
		{
			name: "valid request with only name",
			request: dto.UpdateBrandRequest{
				Name:        "Test Brand",
				Code:        "",
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "invalid empty name",
			request: dto.UpdateBrandRequest{
				Name:        "  ",
				Code:        "TB123",
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid code format",
			request: dto.UpdateBrandRequest{
				Name:        "Test Brand",
				Code:        "tb@123",
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != tt.wantErr {
				t.Errorf("UpdateBrandRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFields != nil && err == nil {
				tt.checkFields(t, tt.request)
			}
		})
	}
}

func TestCreateCategoryRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		request     dto.CreateCategoryRequest
		wantErr     error
		checkFields func(*testing.T, dto.CreateCategoryRequest)
	}{
		{
			name: "valid request with all fields",
			request: dto.CreateCategoryRequest{
				Name:        "  Test Category  ",
				Code:        "  tc123  ",
				Description: "  A valid description  ",
			},
			wantErr: nil,
			checkFields: func(t *testing.T, r dto.CreateCategoryRequest) {
				if r.Name != "Test Category" {
					t.Errorf("Name not trimmed, got %q", r.Name)
				}
				if r.Code != "TC123" {
					t.Errorf("Code not trimmed and uppercased, got %q", r.Code)
				}
				if r.Description != "A valid description" {
					t.Errorf("Description not trimmed, got %q", r.Description)
				}
			},
		},
		{
			name: "valid request with only name",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "",
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "invalid empty name",
			request: dto.CreateCategoryRequest{
				Name:        "  ",
				Code:        "TC123",
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid code format",
			request: dto.CreateCategoryRequest{
				Name:        "Test Category",
				Code:        "tc@123",
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != tt.wantErr {
				t.Errorf("CreateCategoryRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFields != nil && err == nil {
				tt.checkFields(t, tt.request)
			}
		})
	}
}

func TestUpdateCategoryRequestValidate(t *testing.T) {
	tests := []struct {
		name        string
		request     dto.UpdateCategoryRequest
		wantErr     error
		checkFields func(*testing.T, dto.UpdateCategoryRequest)
	}{
		{
			name: "valid request with all fields",
			request: dto.UpdateCategoryRequest{
				Name:        "  Test Category  ",
				Code:        "  tc123  ",
				Description: "  A valid description  ",
			},
			wantErr: nil,
			checkFields: func(t *testing.T, r dto.UpdateCategoryRequest) {
				if r.Name != "Test Category" {
					t.Errorf("Name not trimmed, got %q", r.Name)
				}
				if r.Code != "TC123" {
					t.Errorf("Code not trimmed and uppercased, got %q", r.Code)
				}
				if r.Description != "A valid description" {
					t.Errorf("Description not trimmed, got %q", r.Description)
				}
			},
		},
		{
			name: "valid request with only name",
			request: dto.UpdateCategoryRequest{
				Name:        "Test Category",
				Code:        "",
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "invalid empty name",
			request: dto.UpdateCategoryRequest{
				Name:        "  ",
				Code:        "TC123",
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid code format",
			request: dto.UpdateCategoryRequest{
				Name:        "Test Category",
				Code:        "tc@123",
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != tt.wantErr {
				t.Errorf("UpdateCategoryRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFields != nil && err == nil {
				tt.checkFields(t, tt.request)
			}
		})
	}
}

func TestCreateProductRequestValidate(t *testing.T) {
	validCategoryID := int64(1)
	invalidCategoryID := int64(0)

	tests := []struct {
		name        string
		request     dto.CreateProductRequest
		wantErr     error
		checkFields func(*testing.T, dto.CreateProductRequest)
	}{
		{
			name: "valid request with all fields",
			request: dto.CreateProductRequest{
				Name:        "  Test Product  ",
				Code:        "  tp123  ",
				CategoryID:  &validCategoryID,
				Image:       "  product.png  ",
				Description: "  A valid description  ",
			},
			wantErr: nil,
			checkFields: func(t *testing.T, r dto.CreateProductRequest) {
				if r.Name != "Test Product" {
					t.Errorf("Name not trimmed, got %q", r.Name)
				}
				if r.Code != "TP123" {
					t.Errorf("Code not trimmed and uppercased, got %q", r.Code)
				}
				if r.Image != "product.png" {
					t.Errorf("Image not trimmed, got %q", r.Image)
				}
				if r.Description != "A valid description" {
					t.Errorf("Description not trimmed, got %q", r.Description)
				}
			},
		},
		{
			name: "valid request with only name",
			request: dto.CreateProductRequest{
				Name:        "Test Product",
				Code:        "",
				CategoryID:  nil,
				Image:       "",
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "valid with empty image",
			request: dto.CreateProductRequest{
				Name:        "Test Product",
				Code:        "TP123",
				CategoryID:  &validCategoryID,
				Image:       "",
				Description: "A valid description",
			},
			wantErr: nil,
		},
		{
			name: "invalid empty name",
			request: dto.CreateProductRequest{
				Name:        "  ",
				Code:        "TP123",
				CategoryID:  &validCategoryID,
				Image:       "product.png",
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid category ID zero",
			request: dto.CreateProductRequest{
				Name:        "Test Product",
				Code:        "TP123",
				CategoryID:  &invalidCategoryID,
				Image:       "product.png",
				Description: "A valid description",
			},
			wantErr: response.ValidCategory,
		},
		{
			name: "invalid category ID negative",
			request: dto.CreateProductRequest{
				Name:        "Test Product",
				Code:        "TP123",
				CategoryID:  func() *int64 { i := int64(-1); return &i }(),
				Image:       "product.png",
				Description: "A valid description",
			},
			wantErr: response.ValidCategory,
		},
		{
			name: "invalid image too long",
			request: dto.CreateProductRequest{
				Name:        "Test Product",
				Code:        "TP123",
				CategoryID:  &validCategoryID,
				Image:       string(make([]byte, 256)),
				Description: "A valid description",
			},
			wantErr: response.ValidImage,
		},
		{
			name: "invalid code format",
			request: dto.CreateProductRequest{
				Name:        "Test Product",
				Code:        "tp@123",
				CategoryID:  &validCategoryID,
				Image:       "product.png",
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != tt.wantErr {
				t.Errorf("CreateProductRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFields != nil && err == nil {
				tt.checkFields(t, tt.request)
			}
		})
	}
}

func TestUpdateProductRequestValidate(t *testing.T) {
	validCategoryID := int64(1)
	invalidCategoryID := int64(0)

	tests := []struct {
		name        string
		request     dto.UpdateProductRequest
		wantErr     error
		checkFields func(*testing.T, dto.UpdateProductRequest)
	}{
		{
			name: "valid request with all fields",
			request: dto.UpdateProductRequest{
				Name:        "  Test Product  ",
				Code:        "  tp123  ",
				CategoryID:  &validCategoryID,
				Image:       "  product.png  ",
				Description: "  A valid description  ",
			},
			wantErr: nil,
			checkFields: func(t *testing.T, r dto.UpdateProductRequest) {
				if r.Name != "Test Product" {
					t.Errorf("Name not trimmed, got %q", r.Name)
				}
				if r.Code != "TP123" {
					t.Errorf("Code not trimmed and uppercased, got %q", r.Code)
				}
				if r.Image != "product.png" {
					t.Errorf("Image not trimmed, got %q", r.Image)
				}
				if r.Description != "A valid description" {
					t.Errorf("Description not trimmed, got %q", r.Description)
				}
			},
		},
		{
			name: "valid request with only name",
			request: dto.UpdateProductRequest{
				Name:        "Test Product",
				Code:        "",
				CategoryID:  nil,
				Image:       "",
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "valid with empty image",
			request: dto.UpdateProductRequest{
				Name:        "Test Product",
				Code:        "TP123",
				CategoryID:  &validCategoryID,
				Image:       "",
				Description: "A valid description",
			},
			wantErr: nil,
		},
		{
			name: "invalid empty name",
			request: dto.UpdateProductRequest{
				Name:        "  ",
				Code:        "TP123",
				CategoryID:  &validCategoryID,
				Image:       "product.png",
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid category ID zero",
			request: dto.UpdateProductRequest{
				Name:        "Test Product",
				Code:        "TP123",
				CategoryID:  &invalidCategoryID,
				Image:       "product.png",
				Description: "A valid description",
			},
			wantErr: response.ValidCategory,
		},
		{
			name: "invalid category ID negative",
			request: dto.UpdateProductRequest{
				Name:        "Test Product",
				Code:        "TP123",
				CategoryID:  func() *int64 { i := int64(-1); return &i }(),
				Image:       "product.png",
				Description: "A valid description",
			},
			wantErr: response.ValidCategory,
		},
		{
			name: "invalid image too long",
			request: dto.UpdateProductRequest{
				Name:        "Test Product",
				Code:        "TP123",
				CategoryID:  &validCategoryID,
				Image:       string(make([]byte, 256)),
				Description: "A valid description",
			},
			wantErr: response.ValidImage,
		},
		{
			name: "invalid code format",
			request: dto.UpdateProductRequest{
				Name:        "Test Product",
				Code:        "tp@123",
				CategoryID:  &validCategoryID,
				Image:       "product.png",
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != tt.wantErr {
				t.Errorf("UpdateProductRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFields != nil && err == nil {
				tt.checkFields(t, tt.request)
			}
		})
	}
}

func TestCreateControllerRequestValidate(t *testing.T) {
	validBrandID := int64(1)
	invalidBrandID := int64(0)

	tests := []struct {
		name        string
		request     dto.CreateControllerRequest
		wantErr     error
		checkFields func(*testing.T, dto.CreateControllerRequest)
	}{
		{
			name: "valid request with all fields",
			request: dto.CreateControllerRequest{
				Name:        "  Test Controller  ",
				Code:        "  tc123  ",
				Variant:     "  828D  ",
				BrandID:     &validBrandID,
				Description: "  A valid description  ",
			},
			wantErr: nil,
			checkFields: func(t *testing.T, r dto.CreateControllerRequest) {
				if r.Name != "Test Controller" {
					t.Errorf("Name not trimmed, got %q", r.Name)
				}
				if r.Code != "TC123" {
					t.Errorf("Code not trimmed and uppercased, got %q", r.Code)
				}
				if r.Variant != "828D" {
					t.Errorf("Variant not trimmed, got %q", r.Variant)
				}
				if r.Description != "A valid description" {
					t.Errorf("Description not trimmed, got %q", r.Description)
				}
			},
		},
		{
			name: "valid request with only name",
			request: dto.CreateControllerRequest{
				Name:        "Test Controller",
				Code:        "",
				Variant:     "",
				BrandID:     nil,
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "valid with empty variant",
			request: dto.CreateControllerRequest{
				Name:        "Test Controller",
				Code:        "TC123",
				Variant:     "",
				BrandID:     &validBrandID,
				Description: "A valid description",
			},
			wantErr: nil,
		},
		{
			name: "invalid empty name",
			request: dto.CreateControllerRequest{
				Name:        "  ",
				Code:        "TC123",
				Variant:     "828D",
				BrandID:     &validBrandID,
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid brand ID zero",
			request: dto.CreateControllerRequest{
				Name:        "Test Controller",
				Code:        "TC123",
				Variant:     "828D",
				BrandID:     &invalidBrandID,
				Description: "A valid description",
			},
			wantErr: response.ValidBrand,
		},
		{
			name: "invalid brand ID negative",
			request: dto.CreateControllerRequest{
				Name:        "Test Controller",
				Code:        "TC123",
				Variant:     "828D",
				BrandID:     func() *int64 { i := int64(-1); return &i }(),
				Description: "A valid description",
			},
			wantErr: response.ValidBrand,
		},
		{
			name: "invalid variant too long",
			request: dto.CreateControllerRequest{
				Name:        "Test Controller",
				Code:        "TC123",
				Variant:     string(make([]byte, 101)),
				BrandID:     &validBrandID,
				Description: "A valid description",
			},
			wantErr: response.ValidVariant,
		},
		{
			name: "invalid code format",
			request: dto.CreateControllerRequest{
				Name:        "Test Controller",
				Code:        "tc@123",
				Variant:     "828D",
				BrandID:     &validBrandID,
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != tt.wantErr {
				t.Errorf("CreateControllerRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFields != nil && err == nil {
				tt.checkFields(t, tt.request)
			}
		})
	}
}

func TestUpdateControllerRequestValidate(t *testing.T) {
	validBrandID := int64(1)
	invalidBrandID := int64(0)

	tests := []struct {
		name        string
		request     dto.UpdateControllerRequest
		wantErr     error
		checkFields func(*testing.T, dto.UpdateControllerRequest)
	}{
		{
			name: "valid request with all fields",
			request: dto.UpdateControllerRequest{
				Name:        "  Test Controller  ",
				Code:        "  tc123  ",
				Variant:     "  828D  ",
				BrandID:     &validBrandID,
				Description: "  A valid description  ",
			},
			wantErr: nil,
			checkFields: func(t *testing.T, r dto.UpdateControllerRequest) {
				if r.Name != "Test Controller" {
					t.Errorf("Name not trimmed, got %q", r.Name)
				}
				if r.Code != "TC123" {
					t.Errorf("Code not trimmed and uppercased, got %q", r.Code)
				}
				if r.Variant != "828D" {
					t.Errorf("Variant not trimmed, got %q", r.Variant)
				}
				if r.Description != "A valid description" {
					t.Errorf("Description not trimmed, got %q", r.Description)
				}
			},
		},
		{
			name: "valid request with only name",
			request: dto.UpdateControllerRequest{
				Name:        "Test Controller",
				Code:        "",
				Variant:     "",
				BrandID:     nil,
				Description: "",
			},
			wantErr: nil,
		},
		{
			name: "valid with empty variant",
			request: dto.UpdateControllerRequest{
				Name:        "Test Controller",
				Code:        "TC123",
				Variant:     "",
				BrandID:     &validBrandID,
				Description: "A valid description",
			},
			wantErr: nil,
		},
		{
			name: "invalid empty name",
			request: dto.UpdateControllerRequest{
				Name:        "  ",
				Code:        "TC123",
				Variant:     "828D",
				BrandID:     &validBrandID,
				Description: "A valid description",
			},
			wantErr: response.ValidName,
		},
		{
			name: "invalid brand ID zero",
			request: dto.UpdateControllerRequest{
				Name:        "Test Controller",
				Code:        "TC123",
				Variant:     "828D",
				BrandID:     &invalidBrandID,
				Description: "A valid description",
			},
			wantErr: response.ValidBrand,
		},
		{
			name: "invalid brand ID negative",
			request: dto.UpdateControllerRequest{
				Name:        "Test Controller",
				Code:        "TC123",
				Variant:     "828D",
				BrandID:     func() *int64 { i := int64(-1); return &i }(),
				Description: "A valid description",
			},
			wantErr: response.ValidBrand,
		},
		{
			name: "invalid variant too long",
			request: dto.UpdateControllerRequest{
				Name:        "Test Controller",
				Code:        "TC123",
				Variant:     string(make([]byte, 101)),
				BrandID:     &validBrandID,
				Description: "A valid description",
			},
			wantErr: response.ValidVariant,
		},
		{
			name: "invalid code format",
			request: dto.UpdateControllerRequest{
				Name:        "Test Controller",
				Code:        "tc@123",
				Variant:     "828D",
				BrandID:     &validBrandID,
				Description: "A valid description",
			},
			wantErr: response.ValidCode,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.request.Validate()
			if err != tt.wantErr {
				t.Errorf("UpdateControllerRequest.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.checkFields != nil && err == nil {
				tt.checkFields(t, tt.request)
			}
		})
	}
}
