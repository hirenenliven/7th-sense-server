package tests

import (
	"fmt"
	"net/http"
	"testing"
)

// unknownID is an id no test row will reach.
const unknownID int64 = 2_000_000_000

func TestCategoryCRUD(t *testing.T) {
	runCatalogCRUD(t, catalogCRUDCase{
		path:            "/categories",
		entityName:      "Category",
		notFoundMessage: "Category not found",
	})
}

func TestBrandCRUD(t *testing.T) {
	runCatalogCRUD(t, catalogCRUDCase{
		path:            "/brands",
		entityName:      "Brand",
		notFoundMessage: "Brand not found",
	})
}

func TestProductCRUD(t *testing.T) {
	runCatalogCRUD(t, catalogCRUDCase{
		path:            "/products",
		entityName:      "Product",
		notFoundMessage: "Product not found",
		extraFields:     map[string]any{"categoryId": createTestRecord(t, "/categories", "Category"), "image": "lathe.png"},
	})
}

func TestProductCategoryAndImage(t *testing.T) {
	router := newCatalogRouter()

	t.Run("unknown category is rejected", func(t *testing.T) {
		body := map[string]any{"name": "Test Product " + uniqueSuffix(), "categoryId": unknownID}
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/products", body), http.StatusBadRequest, "Please enter valid category")
	})

	t.Run("non-positive category is rejected", func(t *testing.T) {
		body := map[string]any{"name": "Test Product " + uniqueSuffix(), "categoryId": 0}
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/products", body), http.StatusBadRequest, "Please enter valid category")
	})

	t.Run("empty image gets the default", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/products", map[string]any{"name": "Test Product " + uniqueSuffix()})
		expectStatus(t, result, http.StatusCreated, "created")
		var product struct {
			ID         int64  `json:"id"`
			Image      string `json:"image"`
			CategoryID *int64 `json:"categoryId"`
		}
		decodeData(t, result, &product)
		if product.Image != "no_image_available.png" || product.CategoryID != nil {
			t.Errorf("product = %+v; want default image and no category", product)
		}
		sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/products/%d", apiV1, product.ID), nil)
	})
}

func TestControllerCRUD(t *testing.T) {
	runCatalogCRUD(t, catalogCRUDCase{
		path:            "/controllers",
		entityName:      "Controller",
		notFoundMessage: "Controller not found",
		extraFields:     map[string]any{"brandId": createTestRecord(t, "/brands", "Brand"), "variant": "828D"},
	})
}

func TestControllerBrand(t *testing.T) {
	router := newCatalogRouter()

	t.Run("unknown brand is rejected", func(t *testing.T) {
		body := map[string]any{"name": "Test Controller " + uniqueSuffix(), "brandId": unknownID}
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/controllers", body), http.StatusBadRequest, "Please enter valid brand")
	})

	t.Run("brand and variant are saved", func(t *testing.T) {
		brandID := createTestRecord(t, "/brands", "Brand")
		body := map[string]any{"name": "Test Controller " + uniqueSuffix(), "brandId": brandID, "variant": " 840D "}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/controllers", body)
		expectStatus(t, result, http.StatusCreated, "created")
		var controller struct {
			ID      int64  `json:"id"`
			BrandID *int64 `json:"brandId"`
			Variant string `json:"variant"`
		}
		decodeData(t, result, &controller)
		if controller.BrandID == nil || *controller.BrandID != brandID || controller.Variant != "840D" {
			t.Errorf("controller = %+v; want brandId %d and trimmed variant", controller, brandID)
		}
		sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/controllers/%d", apiV1, controller.ID), nil)
	})
}

// createTestRecord creates a category or brand for a product or controller to
// reference, and soft-deletes it when the test finishes.
func createTestRecord(t *testing.T, path, entityName string) int64 {
	t.Helper()
	router := newCatalogRouter()

	result := sendRequest(t, router, http.MethodPost, apiV1+path, map[string]any{"name": "Test " + entityName + " Ref " + uniqueSuffix()})
	expectStatus(t, result, http.StatusCreated, "created")
	var record catalogRecord
	decodeData(t, result, &record)

	t.Cleanup(func() {
		sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s%s/%d", apiV1, path, record.ID), nil)
	})
	return record.ID
}
