package tests

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// machineRecord holds the machine fields the tests check.
type machineRecord struct {
	ID                    int64   `json:"id"`
	UserID                *int64  `json:"userId"`
	ModelID               *int64  `json:"modelId"`
	BrandID               *int64  `json:"brandId"`
	ControllerID          *int64  `json:"controllerId"`
	LineID                *int64  `json:"lineId"`
	Name                  string  `json:"name"`
	SerialNumber          string  `json:"serialNumber"`
	Username              string  `json:"username"`
	HasPassword           bool    `json:"hasPassword"`
	InstallationDate      string  `json:"installationDate"`
	BillDate              *string `json:"billDate"`
	WarrantyDate          *string `json:"warrantyDate"`
	HourlyCost            string  `json:"hourlyCost"`
	Status                string  `json:"status"`
	ControllerVersionType *int    `json:"controllerVersionType"`
	HealthIDs             string  `json:"healthIds"`
	CustomerName          string  `json:"customerName"`
	ModelName             string  `json:"modelName"`
	BrandName             string  `json:"brandName"`
	ControllerName        string  `json:"controllerName"`
	LineName              string  `json:"lineName"`
}

// machineFixture holds the rows a machine refers to: a customer with one line, a brand
// with a controller and its model, plus a controller (and model) of another brand.
type machineFixture struct {
	customerID, lineID, brandID, controllerID, modelID, otherControllerID, otherModelID int64
}

// newMachineRouter builds the real router with the machine handler wired to testDBPool;
// authenticate decides who the requests act as.
func newMachineRouter(authenticate middleware.Middleware) http.Handler {
	userCountryService := service.NewUserCountryService(repository.NewUserCountryRepository(testDBPool))
	userStateService := service.NewUserStateService(repository.NewUserStateRepository(testDBPool), userCountryService)
	userCityService := service.NewUserCityService(repository.NewUserCityRepository(testDBPool), userStateService)
	branchService := service.NewBranchService(repository.NewBranchRepository(testDBPool), userCityService)
	userLocationService := service.NewUserLocationService(repository.NewUserLocationRepository(testDBPool), branchService)
	userShopService := service.NewUserShopService(repository.NewUserShopRepository(testDBPool), userLocationService)
	userLineService := service.NewUserLineService(repository.NewUserLineRepository(testDBPool), userShopService)
	brandService := service.NewBrandService(repository.NewBrandRepository(testDBPool))
	controllerService := service.NewControllerService(repository.NewControllerRepository(testDBPool), brandService)
	machineModelService := service.NewMachineModelService(repository.NewMachineModelRepository(testDBPool), controllerService)
	machineService := service.NewMachineService(
		repository.NewMachineRepository(testDBPool),
		repository.NewUserRepository(testDBPool),
		machineModelService,
		brandService,
		controllerService,
		userLineService,
	)
	return routes.New(routes.Handlers{Machine: handler.NewMachineHandler(machineService)}, authenticate)
}

// actAs is an authenticate middleware that makes every request come from userID with role.
func actAs(userID int64, role string) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			user := model.AuthUser{ID: userID, Role: role}
			next.ServeHTTP(writer, request.WithContext(middleware.WithAuthUser(request.Context(), user)))
		})
	}
}

// insertTestRow runs an INSERT … RETURNING id and soft-deletes the row from table when the test ends.
func insertTestRow(t *testing.T, table, query string, args ...any) int64 {
	t.Helper()
	var id int64
	if err := testDBPool.QueryRow(context.Background(), query, args...).Scan(&id); err != nil {
		t.Fatalf("insert %s: %v", table, err)
	}
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `UPDATE `+table+` SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, id)
	})
	return id
}

func newMachineFixture(t *testing.T, suffix string) machineFixture {
	t.Helper()
	var fixture machineFixture
	fixture.customerID = createAuthTestUser(t, "machine_customer_"+suffix, model.RoleCustomer)
	fixture.brandID = insertTestRow(t, "brands", `INSERT INTO brands (name) VALUES ($1) RETURNING id`, "Brand "+suffix)
	otherBrandID := insertTestRow(t, "brands", `INSERT INTO brands (name) VALUES ($1) RETURNING id`, "Other brand "+suffix)
	fixture.controllerID = insertTestRow(t, "controllers", `INSERT INTO controllers (name, brand_id) VALUES ($1, $2) RETURNING id`, "Controller "+suffix, fixture.brandID)
	fixture.otherControllerID = insertTestRow(t, "controllers", `INSERT INTO controllers (name, brand_id) VALUES ($1, $2) RETURNING id`, "Other controller "+suffix, otherBrandID)
	fixture.modelID = insertTestRow(t, "models", `INSERT INTO models (controller_id, name) VALUES ($1, $2) RETURNING id`, fixture.controllerID, "Model "+suffix)
	fixture.otherModelID = insertTestRow(t, "models", `INSERT INTO models (controller_id, name) VALUES ($1, $2) RETURNING id`, fixture.otherControllerID, "Other model "+suffix)

	// The customer's organization down to one line
	customer := fixture.customerID
	countryID := insertTestRow(t, "user_country", `INSERT INTO user_country (user_id, name) VALUES ($1, $2) RETURNING id`, customer, "Country "+suffix)
	stateID := insertTestRow(t, "user_state", `INSERT INTO user_state (user_id, user_country_id, name) VALUES ($1, $2, $3) RETURNING id`, customer, countryID, "State "+suffix)
	cityID := insertTestRow(t, "user_city", `INSERT INTO user_city (user_id, user_state_id, name) VALUES ($1, $2, $3) RETURNING id`, customer, stateID, "City "+suffix)
	branchID := insertTestRow(t, "branches", `INSERT INTO branches (code, name, user_city_id) VALUES ($1, $2, $3) RETURNING id`, "M"+suffix, "Branch "+suffix, cityID)
	locationID := insertTestRow(t, "user_location", `INSERT INTO user_location (user_id, branch_id, name) VALUES ($1, $2, $3) RETURNING id`, customer, branchID, "Location "+suffix)
	shopID := insertTestRow(t, "user_shop", `INSERT INTO user_shop (user_id, user_location_id, name) VALUES ($1, $2, $3) RETURNING id`, customer, locationID, "Shop "+suffix)
	fixture.lineID = insertTestRow(t, "user_line", `INSERT INTO user_line (user_id, user_shop_id, name) VALUES ($1, $2, $3) RETURNING id`, customer, shopID, "Line "+suffix)
	return fixture
}

// machineBody is a valid Add Machine form body for fixture; change fields per case.
func machineBody(fixture machineFixture, name string) map[string]any {
	return map[string]any{
		"userId":                fixture.customerID,
		"modelId":               fixture.modelID,
		"name":                  name,
		"serialNumber":          "SN-" + name,
		"brandId":               fixture.brandID,
		"controllerId":          fixture.controllerID,
		"controllerVersion":     "V04.07",
		"accuracyId":            1,
		"lineId":                fixture.lineID,
		"isInsured":             true,
		"maintenanceRequired":   false,
		"installationDate":      "2026-01-15",
		"billDate":              "",
		"warrantyDate":          "2027-01-15",
		"hourlyCost":            "800.50",
		"amc":                   "",
		"machineSystemType":     "NewPlc",
		"securityMode":          "Basic256Sha256",
		"localIp":               "192.168.1.10",
		"liveIp":                "103.81.119.106",
		"port":                  4840,
		"username":              "operator",
		"password":              "secret123",
		"controllerVersionType": 1,
		"longitude":             "70.671297",
		"latitude":              "22.244054",
		"healthIds":             "1,5,6",
	}
}

func with(body map[string]any, key string, value any) map[string]any {
	copied := make(map[string]any, len(body))
	for field, fieldValue := range body {
		copied[field] = fieldValue
	}
	copied[key] = value
	return copied
}

func TestMachineCRUD(t *testing.T) {
	router := newMachineRouter(asTestUser) // admin
	suffix := uniqueSuffix()
	fixture := newMachineFixture(t, suffix)
	name := "VMC-" + suffix
	machinePath := func(id int64) string { return fmt.Sprintf("%s/machines/%d", apiV1, id) }

	var created machineRecord
	t.Run("create", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/machines", machineBody(fixture, name))
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		switch {
		case created.ID == 0 || created.Name != name || *created.UserID != fixture.customerID:
			t.Fatalf("created = %+v; want name %q for customer %d", created, name, fixture.customerID)
		case !created.HasPassword || created.Status != "active" || created.HourlyCost != "800.50":
			t.Errorf("created = %+v; want saved password, active, cost 800.50", created)
		case created.BillDate != nil || created.WarrantyDate == nil || *created.WarrantyDate != "2027-01-15":
			t.Errorf("dates = %v / %v; want no bill date and warranty 2027-01-15", created.BillDate, created.WarrantyDate)
		case created.ModelName != "Model "+suffix || created.BrandName != "Brand "+suffix ||
			created.ControllerName != "Controller "+suffix || created.LineName != "Line "+suffix || created.CustomerName == "":
			t.Errorf("display names = %+v", created)
		}
	})
	if created.ID == 0 {
		t.FailNow()
	}

	t.Run("password is never returned", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, machinePath(created.ID), nil)
		var raw map[string]json.RawMessage
		decodeData(t, result, &raw)
		if _, found := raw["password"]; found {
			t.Fatal("GET /machines/{id} returned the password")
		}
	})

	t.Run("validation", func(t *testing.T) {
		body := machineBody(fixture, "Other-"+suffix)
		for _, testCase := range []struct {
			field   string
			value   any
			message string
		}{
			{"userId", 0, "Please select customer"},
			{"modelId", 0, "Please select model"},
			{"name", "  ", "Please enter valid machine name"},
			{"installationDate", "2026-02-30", "Please enter valid installation date"},
			{"warrantyDate", "2025-01-01", "Warranty end date cannot be before the installation date"},
			{"hourlyCost", "12.345", "Please enter valid hourly cost"},
			{"liveIp", "300.1.1.1", "Please enter valid live IP"},
			{"port", 70000, "Please enter valid port"},
			{"controllerVersionType", 3, "Please select controller version"},
			{"latitude", "95", "Please enter valid latitude"},
			{"healthIds", "1,x", "Please select valid machine health"},
			{"securityMode", "Open", "Please select valid security type"},
			{"password", "", "Please enter valid password"},
		} {
			expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", with(body, testCase.field, testCase.value)), http.StatusBadRequest, testCase.message)
		}
	})

	t.Run("references are checked", func(t *testing.T) {
		body := machineBody(fixture, "Ref-"+suffix)
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", with(body, "userId", testUserID)), http.StatusBadRequest, "Please select valid customer")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", with(body, "brandId", unknownID)), http.StatusBadRequest, "Please select valid brand")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", with(body, "controllerId", fixture.otherControllerID)), http.StatusBadRequest, "Please select a controller of the selected brand")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", with(body, "modelId", unknownID)), http.StatusBadRequest, "Please select valid model")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", with(body, "modelId", fixture.otherModelID)), http.StatusBadRequest, "Please select a model of the selected controller")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", with(body, "lineId", unknownID)), http.StatusBadRequest, "Please select valid line")
	})

	t.Run("name is unique per customer, case-insensitively", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", machineBody(fixture, "vmc-"+suffix)), http.StatusBadRequest, "This customer already has a machine with this name")
	})

	t.Run("list filters by customer", func(t *testing.T) {
		var all, mine, none []machineRecord
		decodeData(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines", nil), &all)
		decodeData(t, sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/machines?userId=%d", apiV1, fixture.customerID), nil), &mine)
		decodeData(t, sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/machines?userId=%d", apiV1, unknownID), nil), &none)
		if !containsMachineID(all, created.ID) || len(mine) != 1 || mine[0].ID != created.ID || len(none) != 0 {
			t.Fatalf("all has machine = %v, customer's = %+v, unknown's = %+v", containsMachineID(all, created.ID), mine, none)
		}
	})

	t.Run("update keeps the password when blank", func(t *testing.T) {
		body := with(with(machineBody(fixture, name+"-B"), "password", ""), "billDate", "2026-01-10")
		result := sendRequest(t, router, http.MethodPut, machinePath(created.ID), body)
		expectStatus(t, result, http.StatusOK, "success")
		var updated machineRecord
		decodeData(t, result, &updated)
		if updated.Name != name+"-B" || updated.BillDate == nil || *updated.BillDate != "2026-01-10" {
			t.Fatalf("updated = %+v", updated)
		}
		var password string
		if err := testDBPool.QueryRow(context.Background(), `SELECT password FROM machines WHERE id = $1`, created.ID).Scan(&password); err != nil || password != "secret123" {
			t.Fatalf("saved password = %q, %v; want the original", password, err)
		}
		expectStatus(t, sendRequest(t, router, http.MethodPut, machinePath(unknownID), body), http.StatusNotFound, "Machine not found")
	})

	t.Run("customers read only their own machines and cannot change them", func(t *testing.T) {
		customerRouter := newMachineRouter(actAs(fixture.customerID, model.RoleCustomer))
		var own []machineRecord
		decodeData(t, sendRequest(t, customerRouter, http.MethodGet, apiV1+"/machines", nil), &own)
		if len(own) != 1 || own[0].ID != created.ID {
			t.Fatalf("customer's list = %+v; want only %d", own, created.ID)
		}
		expectStatus(t, sendRequest(t, customerRouter, http.MethodGet, fmt.Sprintf("%s/machines?userId=%d", apiV1, testUserID), nil), http.StatusForbidden, "You do not have permission to do this")
		expectStatus(t, sendRequest(t, customerRouter, http.MethodDelete, machinePath(created.ID), nil), http.StatusForbidden, "You do not have permission to do this")

		strangerRouter := newMachineRouter(actAs(testUserID, model.RoleCustomer))
		expectStatus(t, sendRequest(t, strangerRouter, http.MethodGet, machinePath(created.ID), nil), http.StatusNotFound, "Machine not found")
	})

	t.Run("delete", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, machinePath(created.ID), nil), http.StatusOK, "Machine deleted")
		expectStatus(t, sendRequest(t, router, http.MethodGet, machinePath(created.ID), nil), http.StatusNotFound, "Machine not found")
		expectStatus(t, sendRequest(t, router, http.MethodDelete, machinePath(created.ID), nil), http.StatusNotFound, "Machine not found")
	})
}

func containsMachineID(records []machineRecord, id int64) bool {
	for _, record := range records {
		if record.ID == id {
			return true
		}
	}
	return false
}
