package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

const apiV1 = "/api/v1"

// apiResponse is the JSON body every endpoint returns.
type apiResponse struct {
	StatusCode int             `json:"statusCode"`
	Success    bool            `json:"success"`
	Message    string          `json:"message"`
	Data       json.RawMessage `json:"data"`
}

// asTestUser is the authenticate middleware of the test routers: every request
// acts as testUserID with the superadmin role, so suites test their entity without
// logging in. auth_test.go covers the real session middleware.
func asTestUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		user := model.AuthUser{ID: testUserID, Role: model.RoleSuperAdmin}
		next.ServeHTTP(writer, request.WithContext(middleware.WithAuthUser(request.Context(), user)))
	})
}

// newCatalogRouter builds the real router with the catalog handlers wired to testDBPool,
// acting as the superadmin test user.
func newCatalogRouter() http.Handler {
	return newCatalogRouterAs(asTestUser)
}

// newCatalogRouterAs is newCatalogRouter with authenticate deciding who the requests act as.
func newCatalogRouterAs(authenticate middleware.Middleware) http.Handler {
	brandService := service.NewBrandService(repository.NewBrandRepository(testDBPool))
	controllerService := service.NewControllerService(repository.NewControllerRepository(testDBPool), brandService)
	return routes.New(routes.Handlers{
		Category:     handler.NewCategoryHandler(service.NewCategoryService(repository.NewCategoryRepository(testDBPool))),
		Brand:        handler.NewBrandHandler(brandService),
		Product:      handler.NewProductHandler(service.NewProductService(repository.NewProductRepository(testDBPool))),
		Controller:   handler.NewControllerHandler(controllerService),
		MachineModel: handler.NewMachineModelHandler(service.NewMachineModelService(repository.NewMachineModelRepository(testDBPool), controllerService)),
	}, authenticate)
}

// sendRequest sends method path with body encoded as JSON (nil sends no body)
// and decodes the response envelope.
func sendRequest(t *testing.T, router http.Handler, method, path string, body any) apiResponse {
	t.Helper()

	var requestBody bytes.Buffer
	if body != nil {
		if err := json.NewEncoder(&requestBody).Encode(body); err != nil {
			t.Fatalf("encode request body: %v", err)
		}
	}
	request := httptest.NewRequest(method, path, &requestBody)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	var result apiResponse
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatalf("%s %s: decode response: %v", method, path, err)
	}
	if result.StatusCode != recorder.Code {
		t.Fatalf("%s %s: body statusCode %d != HTTP status %d", method, path, result.StatusCode, recorder.Code)
	}
	return result
}

// expectStatus fails the test when result does not have the wanted status and message.
func expectStatus(t *testing.T, result apiResponse, wantStatus int, wantMessage string) {
	t.Helper()
	if result.StatusCode != wantStatus || result.Message != wantMessage {
		t.Fatalf("got %d %q; want %d %q", result.StatusCode, result.Message, wantStatus, wantMessage)
	}
}

// decodeData decodes the data field of result into target.
func decodeData(t *testing.T, result apiResponse, target any) {
	t.Helper()
	if err := json.Unmarshal(result.Data, target); err != nil {
		t.Fatalf("decode data %s: %v", result.Data, err)
	}
}

// uniqueSuffix keeps names and codes unique across test runs, because the
// rows the tests create are kept (soft-deleted) for inspection.
// The counter keeps suffixes apart when the clock has not moved between two calls.
func uniqueSuffix() string {
	return fmt.Sprintf("%d%02d", time.Now().UnixNano()/1000%10_000_000, suffixCounter.Add(1)%100)
}

var suffixCounter atomic.Int64

// catalogRecord holds the fields shared by categories, brands, products and controllers.
type catalogRecord struct {
	ID          int64  `json:"id"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Description string `json:"description"`
	IsDeleted   bool   `json:"isDeleted"`
}

// catalogCRUDCase describes one catalog entity for runCatalogCRUD.
type catalogCRUDCase struct {
	path            string         // collection path, e.g. "/categories"
	entityName      string         // "Category", used in the delete message
	notFoundMessage string         // e.g. "Category not found"
	duplicateName   string         // duplicate-name message; "" = "Name already exists"
	extraFields     map[string]any // entity-specific fields sent on create and update
	noCode          bool           // the entity has no code (controllers): code checks are skipped
}

// runCatalogCRUD runs the full create, read, update and delete flow shared by
// every catalog entity, including validation, duplicate and soft-delete rules.
func runCatalogCRUD(t *testing.T, testCase catalogCRUDCase) {
	router := newCatalogRouter()
	suffix := uniqueSuffix()
	name := "Test " + testCase.entityName + " " + suffix
	code := "TC" + suffix
	if testCase.noCode {
		code = ""
	}

	requestBody := func(name, code, description string) map[string]any {
		body := map[string]any{"name": name, "code": code, "description": description}
		for key, value := range testCase.extraFields {
			body[key] = value
		}
		return body
	}

	var created catalogRecord
	t.Run("create", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+testCase.path, requestBody("  "+name+"  ", code, "first description"))
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		if created.ID == 0 || created.Name != name || created.Code != code || created.Description != "first description" || created.IsDeleted {
			t.Fatalf("created = %+v; want trimmed name %q, code %q", created, name, code)
		}
	})
	if created.ID == 0 {
		t.FailNow()
	}
	itemPath := fmt.Sprintf("%s%s/%d", apiV1, testCase.path, created.ID)

	t.Run("create rejects invalid input", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+testCase.path, requestBody("   ", "", "")), http.StatusBadRequest, "Please enter valid name")
		if !testCase.noCode {
			expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+testCase.path, requestBody("Other "+suffix, "bad code!", "")), http.StatusBadRequest, "Please enter valid code")
		}
	})

	t.Run("create rejects duplicate name and code", func(t *testing.T) {
		duplicateName := testCase.duplicateName
		if duplicateName == "" {
			duplicateName = "Name already exists"
		}
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+testCase.path, requestBody(name, "", "")), http.StatusBadRequest, duplicateName)
		if !testCase.noCode {
			expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+testCase.path, requestBody("Other "+suffix, code, "")), http.StatusBadRequest, "Code already exists")
		}
	})

	t.Run("get by id", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, itemPath, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var found catalogRecord
		decodeData(t, result, &found)
		if found.ID != created.ID || found.Name != name {
			t.Errorf("found = %+v; want %+v", found, created)
		}
	})

	t.Run("get all contains created", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, apiV1+testCase.path, nil)
		expectStatus(t, result, http.StatusOK, "success")
		var all []catalogRecord
		decodeData(t, result, &all)
		if !containsID(all, created.ID) {
			t.Errorf("list does not contain id %d", created.ID)
		}
	})

	t.Run("update", func(t *testing.T) {
		updatedName := name + " Updated"
		result := sendRequest(t, router, http.MethodPut, itemPath, requestBody(updatedName, code, "")) // empty description is cleared
		expectStatus(t, result, http.StatusOK, "success")
		var updated catalogRecord
		decodeData(t, result, &updated)
		if updated.Name != updatedName || updated.Code != code || updated.Description != "" {
			t.Errorf("updated = %+v; want name %q, code %q and empty description", updated, updatedName, code)
		}
	})

	t.Run("update and get unknown id return not found", func(t *testing.T) {
		unknownPath := fmt.Sprintf("%s%s/%d", apiV1, testCase.path, unknownID)
		expectStatus(t, sendRequest(t, router, http.MethodPut, unknownPath, requestBody("Ghost "+suffix, "", "")), http.StatusNotFound, testCase.notFoundMessage)
		expectStatus(t, sendRequest(t, router, http.MethodGet, unknownPath, nil), http.StatusNotFound, testCase.notFoundMessage)
		expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+testCase.path+"/abc", nil), http.StatusBadRequest, "Please enter valid id")
	})

	t.Run("delete", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, itemPath, nil), http.StatusOK, testCase.entityName+" deleted")
		expectStatus(t, sendRequest(t, router, http.MethodGet, itemPath, nil), http.StatusNotFound, testCase.notFoundMessage)
		expectStatus(t, sendRequest(t, router, http.MethodDelete, itemPath, nil), http.StatusNotFound, testCase.notFoundMessage)
	})

	t.Run("deleted name and code can be reused", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+testCase.path, requestBody(name+" Updated", code, ""))
		expectStatus(t, result, http.StatusCreated, "created")
		var recreated catalogRecord
		decodeData(t, result, &recreated)
		expectStatus(t, sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s%s/%d", apiV1, testCase.path, recreated.ID), nil), http.StatusOK, testCase.entityName+" deleted")
	})
}

func containsID(records []catalogRecord, id int64) bool {
	for _, record := range records {
		if record.ID == id {
			return true
		}
	}
	return false
}
