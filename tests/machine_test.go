package tests

import (
	"fmt"
	"net/http"
	"testing"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// machineRecord represents a machine as returned by the API.
type machineRecord struct {
	ID                    int64   `json:"id"`
	UserID                *int64  `json:"userId"`
	MachineTypeID         *int64  `json:"machineTypeId"`
	ProductID             *int64  `json:"productId"`
	Name                  string  `json:"name"`
	SerialNumber          string  `json:"serialNumber"`
	LiveIP                string  `json:"liveIp"`
	Port                  int     `json:"port"`
	Username              string  `json:"username"`
	Password              string  `json:"password"`
	LocalIP               string  `json:"localIp"`
	BrandID               *int64  `json:"brandId"`
	ControllerID          *int64  `json:"controllerId"`
	ControllerVersion     string  `json:"controllerVersion"`
	AccuracyID            *int64  `json:"accuracyId"`
	LineID                *int64  `json:"lineId"`
	InstallationDate      string  `json:"installationDate"`
	BillDate              string  `json:"billDate"`
	WarrantyDate          string  `json:"warrantyDate"`
	AMC                   string  `json:"amc"`
	IsInsured             bool    `json:"isInsured"`
	OnlineStatus          bool    `json:"onlineStatus"`
	MaintenanceRequired   bool    `json:"maintenanceRequired"`
	HourlyCost            string  `json:"hourlyCost"`
	Status                string  `json:"status"`
	IsMonitor             bool    `json:"isMonitor"`
	OnlineDate            *string `json:"onlineDate"`
	OfflineDate           string  `json:"offlineDate"`
	ModeID                *int64  `json:"modeId"`
	ControllerVersionType *int    `json:"controllerVersionType"`
	HealthIDs             string  `json:"healthIds"`
	Longitude             string  `json:"longitude"`
	Latitude              string  `json:"latitude"`
	CycleAnalysis         bool    `json:"cycleAnalysis"`
	MachineSystemType     string  `json:"machineSystemType"`
	SecurityMode          string  `json:"securityMode"`
	IsDeleted             bool    `json:"isDeleted"`
	DeletedAt             string  `json:"deletedAt"`
	CreatedAt             string  `json:"createdAt"`
	UpdatedAt             string  `json:"updatedAt"`
}

// newMachineRouter builds the real router with the machine handler wired to testDBPool.
func newMachineRouter() http.Handler {
	return routes.New(routes.Handlers{
		Machine: handler.NewMachineHandler(service.NewMachineService(repository.NewMachineRepository(testDBPool))),
	})
}

func TestMachineCRUD(t *testing.T) {
	router := newMachineRouter()
	suffix := uniqueSuffix()
	name := "Machine-" + suffix
	serialNumber := "SN-" + suffix

	// Minimal valid machine request
	requestBody := map[string]any{
		"name":              name,
		"serialNumber":      serialNumber,
		"liveIp":            "192.168.1.100",
		"port":              8080,
		"username":          "admin",
		"password":          "password123",
		"localIp":           "192.168.1.101",
		"controllerVersion": "1.0.0",
		"installationDate":  "2024-01-01",
		"billDate":          "2024-01-01",
		"warrantyDate":      "2025-01-01",
		"hourlyCost":        "100.50",
		"status":            "active",
		"machineSystemType": "NewPlc",
		"securityMode":      "Basic256Sha256",
	}

	var created machineRecord
	t.Run("create", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/machines", requestBody)
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		if created.ID == 0 || created.Name != name || created.SerialNumber != serialNumber {
			t.Fatalf("created = %+v; want name %q, serialNumber %q", created, name, serialNumber)
		}
		if created.LiveIP != "192.168.1.100" || created.Port != 8080 || created.Username != "admin" {
			t.Errorf("created fields not saved correctly")
		}
		if created.IsDeleted {
			t.Error("created isDeleted = true; want false")
		}
	})
	if created.ID == 0 {
		t.FailNow()
	}

	t.Run("create rejects invalid input", func(t *testing.T) {
		// Empty name
		invalidBody := map[string]any{
			"name":              "",
			"serialNumber":      serialNumber + "2",
			"liveIp":            "192.168.1.100",
			"port":              8080,
			"username":          "admin",
			"password":          "password123",
			"localIp":           "192.168.1.101",
			"controllerVersion": "1.0.0",
			"installationDate":  "2024-01-01",
			"billDate":          "2024-01-01",
			"warrantyDate":      "2025-01-01",
			"hourlyCost":        "100.50",
			"status":            "active",
			"machineSystemType": "NewPlc",
			"securityMode":      "Basic256Sha256",
		}
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", invalidBody), http.StatusBadRequest, "Please enter valid name")

		// Invalid IP
		invalidBody["name"] = "Test-Machine-2"
		invalidBody["serialNumber"] = serialNumber + "2"
		invalidBody["liveIp"] = "invalid-ip"
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", invalidBody), http.StatusBadRequest, "Please enter valid IP address")

		// Invalid port
		invalidBody["liveIp"] = "192.168.1.100"
		invalidBody["port"] = 99999
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", invalidBody), http.StatusBadRequest, "Please enter valid port")

		// Invalid date
		invalidBody["port"] = 8080
		invalidBody["installationDate"] = "invalid-date"
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", invalidBody), http.StatusBadRequest, "Please enter valid date")

		// Invalid machine system type
		invalidBody["installationDate"] = "2024-01-01"
		invalidBody["machineSystemType"] = "InvalidType"
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", invalidBody), http.StatusBadRequest, "Please enter valid machine system type")

		// Invalid security mode
		invalidBody["machineSystemType"] = "NewPlc"
		invalidBody["securityMode"] = "InvalidMode"
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", invalidBody), http.StatusBadRequest, "Please enter valid security mode")
	})

	t.Run("create rejects duplicate name", func(t *testing.T) {
		duplicateBody := map[string]any{
			"name":              name,
			"serialNumber":      serialNumber + "3",
			"liveIp":            "192.168.1.100",
			"port":              8080,
			"username":          "admin",
			"password":          "password123",
			"localIp":           "192.168.1.101",
			"controllerVersion": "1.0.0",
			"installationDate":  "2024-01-01",
			"billDate":          "2024-01-01",
			"warrantyDate":      "2025-01-01",
			"hourlyCost":        "100.50",
			"status":            "active",
			"machineSystemType": "NewPlc",
			"securityMode":      "Basic256Sha256",
		}
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/machines", duplicateBody), http.StatusBadRequest, "Name already exists")
	})

	t.Run("get by id", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, apiV1+"/machines/"+jsonNumber(created.ID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var found machineRecord
		decodeData(t, result, &found)
		if found.ID != created.ID || found.Name != name {
			t.Errorf("found = %+v; want id %d, name %q", found, created.ID, name)
		}
	})

	t.Run("get all contains created", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, apiV1+"/machines", nil)
		expectStatus(t, result, http.StatusOK, "success")
		var all []machineRecord
		decodeData(t, result, &all)
		if !containsMachineID(all, created.ID) {
			t.Errorf("list does not contain id %d", created.ID)
		}
	})

	t.Run("update", func(t *testing.T) {
		updatedName := name + "-Updated"
		updateBody := map[string]any{
			"name":              updatedName,
			"serialNumber":      serialNumber,
			"liveIp":            "192.168.1.102",
			"port":              8081,
			"username":          "admin2",
			"password":          "newpassword",
			"localIp":           "192.168.1.103",
			"controllerVersion": "2.0.0",
			"installationDate":  "2024-02-01",
			"billDate":          "2024-02-01",
			"warrantyDate":      "2025-02-01",
			"hourlyCost":        "150.75",
			"status":            "inactive",
			"machineSystemType": "OldPlc",
			"securityMode":      "None",
		}
		result := sendRequest(t, router, http.MethodPut, apiV1+"/machines/"+jsonNumber(created.ID), updateBody)
		expectStatus(t, result, http.StatusOK, "success")
		var updated machineRecord
		decodeData(t, result, &updated)
		if updated.Name != updatedName || updated.LiveIP != "192.168.1.102" || updated.Port != 8081 {
			t.Errorf("updated = %+v; fields not updated", updated)
		}
	})

	t.Run("update and get unknown id return not found", func(t *testing.T) {
		unknownPath := apiV1 + "/machines/" + jsonNumber(unknownID)
		updateBody := map[string]any{
			"name":              "Ghost-Machine",
			"serialNumber":      "SN-GHOST",
			"liveIp":            "192.168.1.100",
			"port":              8080,
			"username":          "admin",
			"password":          "password123",
			"localIp":           "192.168.1.101",
			"controllerVersion": "1.0.0",
			"installationDate":  "2024-01-01",
			"billDate":          "2024-01-01",
			"warrantyDate":      "2025-01-01",
			"hourlyCost":        "100.50",
			"status":            "active",
			"machineSystemType": "NewPlc",
			"securityMode":      "Basic256Sha256",
		}
		expectStatus(t, sendRequest(t, router, http.MethodPut, unknownPath, updateBody), http.StatusNotFound, "Machine not found")
		expectStatus(t, sendRequest(t, router, http.MethodGet, unknownPath, nil), http.StatusNotFound, "Machine not found")
		expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines/abc", nil), http.StatusBadRequest, "Please enter valid id")
	})

	t.Run("delete", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, apiV1+"/machines/"+jsonNumber(created.ID), nil), http.StatusOK, "Machine deleted")
		expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines/"+jsonNumber(created.ID), nil), http.StatusNotFound, "Machine not found")
		expectStatus(t, sendRequest(t, router, http.MethodDelete, apiV1+"/machines/"+jsonNumber(created.ID), nil), http.StatusNotFound, "Machine not found")
	})

	t.Run("deleted name can be reused", func(t *testing.T) {
		recreateBody := map[string]any{
			"name":              name,
			"serialNumber":      serialNumber + "-REUSE",
			"liveIp":            "192.168.1.100",
			"port":              8080,
			"username":          "admin",
			"password":          "password123",
			"localIp":           "192.168.1.101",
			"controllerVersion": "1.0.0",
			"installationDate":  "2024-01-01",
			"billDate":          "2024-01-01",
			"warrantyDate":      "2025-01-01",
			"hourlyCost":        "100.50",
			"status":            "active",
			"machineSystemType": "NewPlc",
			"securityMode":      "Basic256Sha256",
		}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/machines", recreateBody)
		expectStatus(t, result, http.StatusCreated, "created")
		var recreated machineRecord
		decodeData(t, result, &recreated)
		expectStatus(t, sendRequest(t, router, http.MethodDelete, apiV1+"/machines/"+jsonNumber(recreated.ID), nil), http.StatusOK, "Machine deleted")
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

func jsonNumber(id int64) string {
	return fmt.Sprintf("%d", id)
}
