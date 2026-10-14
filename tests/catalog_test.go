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
		duplicateName:   "This brand already has a controller with this name",
		extraFields:     map[string]any{"brandId": createTestRecord(t, "/brands", "Brand")},
		noCode:          true,
	})
}

func TestControllerBrand(t *testing.T) {
	router := newCatalogRouter()

	t.Run("brand is required and must exist", func(t *testing.T) {
		name := "Test Controller " + uniqueSuffix()
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/controllers", map[string]any{"name": name}), http.StatusBadRequest, "Please select brand")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/controllers", map[string]any{"name": name, "brandId": unknownID}), http.StatusBadRequest, "Please enter valid brand")
	})

	t.Run("name is unique per brand", func(t *testing.T) {
		name := "Test Controller " + uniqueSuffix()
		firstID := createController(t, router, name, createTestRecord(t, "/brands", "Brand"))
		otherBrandID := createTestRecord(t, "/brands", "Brand")
		secondID := createController(t, router, name, otherBrandID)
		if firstID == secondID {
			t.Fatalf("both brands got controller id %d", firstID)
		}

		result := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/controllers?brandId=%d", apiV1, otherBrandID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var list []catalogRecord
		decodeData(t, result, &list)
		if len(list) != 1 || list[0].ID != secondID {
			t.Errorf("brand filter = %+v; want only controller %d", list, secondID)
		}
	})

	t.Run("brand is saved and code / variant are gone", func(t *testing.T) {
		brandID := createTestRecord(t, "/brands", "Brand")
		name := "Test Controller " + uniqueSuffix() + " 840D"
		body := map[string]any{"name": " " + name + " ", "brandId": brandID, "code": "IGNORED", "variant": "IGNORED"}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/controllers", body)
		expectStatus(t, result, http.StatusCreated, "created")
		var controller map[string]any
		decodeData(t, result, &controller)
		if controller["name"] != name || controller["brandId"] != float64(brandID) {
			t.Errorf("controller = %+v; want trimmed name %q and brandId %d", controller, name, brandID)
		}
		for _, field := range []string{"code", "variant"} {
			if _, found := controller[field]; found {
				t.Errorf("controller has %q; controllers no longer have it", field)
			}
		}
		sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/controllers/%v", apiV1, controller["id"]), nil)
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
