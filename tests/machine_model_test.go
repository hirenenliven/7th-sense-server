package tests

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"testing"

	"7th-sense-server/internal/model"
)

// machineModelRecord is the JSON shape of a model in API responses.
type machineModelRecord struct {
	ID             int64  `json:"id"`
	ControllerID   int64  `json:"controllerId"`
	Name           string `json:"name"`
	ControllerName string `json:"controllerName"`
	BrandID        *int64 `json:"brandId"`
	BrandName      string `json:"brandName"`
}

func TestMachineModelCRUD(t *testing.T) {
	router := newCatalogRouter()
	suffix := uniqueSuffix()
	brandID := createTestRecord(t, "/brands", "Brand")
	controllerID := createController(t, router, "Model Controller "+suffix, brandID)
	otherControllerID := createController(t, router, "Other Model Controller "+suffix, brandID)
	name := "Model " + suffix

	var created machineModelRecord
	t.Run("create", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/models", map[string]any{"controllerId": controllerID, "name": "  " + name + "  "})
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		switch {
		case created.ID == 0 || created.Name != name || created.ControllerID != controllerID:
			t.Fatalf("created = %+v; want trimmed name %q under controller %d", created, name, controllerID)
		case created.ControllerName != "Model Controller "+suffix || created.BrandID == nil || *created.BrandID != brandID || created.BrandName == "":
			t.Errorf("created = %+v; want controller and brand names", created)
		}
	})
	if created.ID == 0 {
		t.FailNow()
	}
	modelPath := fmt.Sprintf("%s/models/%d", apiV1, created.ID)

	t.Run("validation", func(t *testing.T) {
		for _, testCase := range []struct {
			body    map[string]any
			message string
		}{
			{map[string]any{"name": "No controller " + suffix}, "Please select controller"},
			{map[string]any{"controllerId": controllerID, "name": "   "}, "Please enter valid name"},
			{map[string]any{"controllerId": controllerID, "name": strings.Repeat("x", 101)}, "Please enter valid name"},
			{map[string]any{"controllerId": unknownID, "name": "Ghost " + suffix}, "Please select valid controller"},
		} {
			expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/models", testCase.body), http.StatusBadRequest, testCase.message)
		}
	})

	t.Run("name is unique per controller, case-insensitively", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/models", map[string]any{"controllerId": controllerID, "name": strings.ToUpper(name)}),
			http.StatusBadRequest, "This controller already has a model with this name")

		result := sendRequest(t, router, http.MethodPost, apiV1+"/models", map[string]any{"controllerId": otherControllerID, "name": name})
		expectStatus(t, result, http.StatusCreated, "created")
		var other machineModelRecord
		decodeData(t, result, &other)
		t.Cleanup(func() { sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/models/%d", apiV1, other.ID), nil) })
	})

	t.Run("list filters by controller", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/models?controllerId=%d", apiV1, controllerID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var list []machineModelRecord
		decodeData(t, result, &list)
		if len(list) != 1 || list[0].ID != created.ID {
			t.Errorf("list = %+v; want only model %d", list, created.ID)
		}
		expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/models?controllerId=abc", nil), http.StatusBadRequest, "Please enter valid controllerId")
	})

	t.Run("update", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPut, modelPath, map[string]any{"controllerId": controllerID, "name": name + " Updated"})
		expectStatus(t, result, http.StatusOK, "success")
		var updated machineModelRecord
		decodeData(t, result, &updated)
		if updated.Name != name+" Updated" {
			t.Errorf("updated = %+v; want the new name", updated)
		}
		unknownPath := fmt.Sprintf("%s/models/%d", apiV1, unknownID)
		expectStatus(t, sendRequest(t, router, http.MethodPut, unknownPath, map[string]any{"controllerId": controllerID, "name": "Ghost"}), http.StatusNotFound, "Model not found")
	})

	t.Run("controller with models cannot be deleted", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/controllers/%d", apiV1, controllerID), nil),
			http.StatusBadRequest, "This controller has models. Delete its models first")
	})

	t.Run("only the superadmin can change brands, controllers and models; everyone can read", func(t *testing.T) {
		controllerPath := fmt.Sprintf("%s/controllers/%d", apiV1, controllerID)
		brandPath := fmt.Sprintf("%s/brands/%d", apiV1, brandID)
		forbidden := "You do not have permission to do this"
		for _, role := range []string{model.RoleAdmin, model.RoleCustomer} {
			roleRouter := newCatalogRouterAs(actAs(testUserID, role))
			expectStatus(t, sendRequest(t, roleRouter, http.MethodGet, apiV1+"/models", nil), http.StatusOK, "success")
			expectStatus(t, sendRequest(t, roleRouter, http.MethodGet, apiV1+"/controllers", nil), http.StatusOK, "success")
			expectStatus(t, sendRequest(t, roleRouter, http.MethodPost, apiV1+"/models", map[string]any{"controllerId": controllerID, "name": "X"}), http.StatusForbidden, forbidden)
			expectStatus(t, sendRequest(t, roleRouter, http.MethodPut, modelPath, map[string]any{"controllerId": controllerID, "name": "X"}), http.StatusForbidden, forbidden)
			expectStatus(t, sendRequest(t, roleRouter, http.MethodDelete, modelPath, nil), http.StatusForbidden, forbidden)
			expectStatus(t, sendRequest(t, roleRouter, http.MethodPost, apiV1+"/controllers", map[string]any{"brandId": brandID, "name": "X"}), http.StatusForbidden, forbidden)
			expectStatus(t, sendRequest(t, roleRouter, http.MethodPut, controllerPath, map[string]any{"brandId": brandID, "name": "X"}), http.StatusForbidden, forbidden)
			expectStatus(t, sendRequest(t, roleRouter, http.MethodDelete, controllerPath, nil), http.StatusForbidden, forbidden)
			expectStatus(t, sendRequest(t, roleRouter, http.MethodGet, apiV1+"/brands", nil), http.StatusOK, "success")
			expectStatus(t, sendRequest(t, roleRouter, http.MethodPost, apiV1+"/brands", map[string]any{"name": "X"}), http.StatusForbidden, forbidden)
			expectStatus(t, sendRequest(t, roleRouter, http.MethodPut, brandPath, map[string]any{"name": "X"}), http.StatusForbidden, forbidden)
			expectStatus(t, sendRequest(t, roleRouter, http.MethodDelete, brandPath, nil), http.StatusForbidden, forbidden)
		}
	})

	t.Run("delete", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, modelPath, nil), http.StatusOK, "Model deleted")
		expectStatus(t, sendRequest(t, router, http.MethodGet, modelPath, nil), http.StatusNotFound, "Model not found")
		expectStatus(t, sendRequest(t, router, http.MethodDelete, modelPath, nil), http.StatusNotFound, "Model not found")
	})
}

// TestCatalogInUseByMachine checks that a model or controller a machine uses keeps
// what the machine was saved with.
func TestCatalogInUseByMachine(t *testing.T) {
	router := newCatalogRouter()
	suffix := uniqueSuffix()
	brandID := createTestRecord(t, "/brands", "Brand")
	otherBrandID := createTestRecord(t, "/brands", "Brand")
	controllerID := createController(t, router, "Used Controller "+suffix, brandID)
	otherControllerID := createController(t, router, "Spare Controller "+suffix, brandID)
	modelID := insertTestRow(t, "models", `INSERT INTO models (controller_id, name) VALUES ($1, $2) RETURNING id`, controllerID, "Used Model "+suffix)
	insertTestRow(t, "machines", `INSERT INTO machines (name, serial_number, live_ip, port, username, password, local_ip,
			brand_id, controller_id, model_id, controller_version, installation_date, hourly_cost)
		VALUES ($1, $1, '10.0.0.1', 4840, 'user', 'secret', '10.0.0.2', $2, $3, $4, 'V1', '2026-01-01', 0) RETURNING id`,
		"In use "+suffix, brandID, controllerID, modelID)

	modelPath := fmt.Sprintf("%s/models/%d", apiV1, modelID)
	controllerPath := fmt.Sprintf("%s/controllers/%d", apiV1, controllerID)

	expectStatus(t, sendRequest(t, router, http.MethodPut, modelPath, map[string]any{"controllerId": otherControllerID, "name": "Used Model " + suffix}),
		http.StatusBadRequest, "Machines use this model, so its controller cannot be changed")
	expectStatus(t, sendRequest(t, router, http.MethodPut, modelPath, map[string]any{"controllerId": controllerID, "name": "Renamed Model " + suffix}),
		http.StatusOK, "success")
	expectStatus(t, sendRequest(t, router, http.MethodDelete, modelPath, nil),
		http.StatusBadRequest, "Machines use this model. Change or delete those machines first")

	expectStatus(t, sendRequest(t, router, http.MethodPut, controllerPath, map[string]any{"brandId": otherBrandID, "name": "Used Controller " + suffix}),
		http.StatusBadRequest, "Machines use this controller, so its brand cannot be changed")
	expectStatus(t, sendRequest(t, router, http.MethodPut, controllerPath, map[string]any{"brandId": brandID, "name": "Renamed Controller " + suffix}),
		http.StatusOK, "success")
}

// createController creates a controller of brandID through the API and deletes it
// (and soft-deletes its models first) when the test finishes.
func createController(t *testing.T, router http.Handler, name string, brandID int64) int64 {
	t.Helper()
	result := sendRequest(t, router, http.MethodPost, apiV1+"/controllers", map[string]any{"name": name, "brandId": brandID})
	expectStatus(t, result, http.StatusCreated, "created")
	var controller catalogRecord
	decodeData(t, result, &controller)

	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `UPDATE models SET is_deleted = true, deleted_at = NOW() WHERE controller_id = $1 AND is_deleted = false`, controller.ID)
		_, _ = testDBPool.Exec(context.Background(), `UPDATE controllers SET is_deleted = true, deleted_at = NOW() WHERE id = $1 AND is_deleted = false`, controller.ID)
	})
	return controller.ID
}
