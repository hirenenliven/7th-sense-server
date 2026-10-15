package tests

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"7th-sense-server/internal/model"
)

// clienteleBody is a valid Add Clientele form body; change fields per case.
func clienteleBody(customerID int64, code string) map[string]any {
	return map[string]any{
		"userId":        customerID,
		"clientCode":    code,
		"name":          "Client " + code,
		"email":         "Client." + code + "@Example.com",
		"contactNumber": "+91 98765-43210",
		"address":       "Plot 12, GIDC",
		"countryName":   " India ",
		"stateName":     "Gujarat",
		"city":          "Rajkot",
	}
}

func TestClienteleCRUD(t *testing.T) {
	uploadDir := t.TempDir()
	router := newPartRouter(asTestUser, uploadDir) // superadmin
	suffix := uniqueSuffix()
	fixture := newMachineFixture(t, suffix)
	otherFixture := newMachineFixture(t, suffix+"9")
	code := "CL-" + suffix
	clientelePath := func(id int64) string { return fmt.Sprintf("%s/clienteles/%d", apiV1, id) }

	var created model.Clientele
	t.Run("create", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/clienteles", clienteleBody(fixture.customerID, code))
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		switch {
		case created.ID == 0 || created.UserID != fixture.customerID || created.Status != "active" || created.CustomerName == "":
			t.Fatalf("created = %+v", created)
		case created.ContactNumber != "+919876543210" || created.CountryName != "India" || created.StateName != "Gujarat":
			t.Errorf("normalized = %q %q %q; want +919876543210 India Gujarat", created.ContactNumber, created.CountryName, created.StateName)
		case created.Email != strings.ToLower("client."+code+"@example.com"):
			t.Errorf("email = %q; want it lowercased", created.Email)
		}
	})
	if created.ID == 0 {
		t.FailNow()
	}

	t.Run("duplicate client ID", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/clienteles", clienteleBody(fixture.customerID, code))
		expectStatus(t, result, http.StatusBadRequest, "This customer already has a clientele with this client ID")
	})

	t.Run("same client ID for another customer", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/clienteles", clienteleBody(otherFixture.customerID, code))
		expectStatus(t, result, http.StatusCreated, "created")
	})

	t.Run("validation", func(t *testing.T) {
		cases := map[string]struct {
			field   string
			value   any
			message string
		}{
			"blank client ID":    {"clientCode", " ", "Please enter valid client ID"},
			"bad email":          {"email", "not-an-email", "Please enter valid client email"},
			"phone without code": {"contactNumber", "9876543210", "Please enter valid contact number"},
			"blank country":      {"countryName", " ", "Please select valid country"},
			"blank state":        {"stateName", "", "Please select valid state"},
			"blank city":         {"city", "", "Please enter valid city"},
		}
		for name, testCase := range cases {
			t.Run(name, func(t *testing.T) {
				body := with(clienteleBody(fixture.customerID, "V-"+suffix), testCase.field, testCase.value)
				result := sendRequest(t, router, http.MethodPost, apiV1+"/clienteles", body)
				expectStatus(t, result, http.StatusBadRequest, testCase.message)
			})
		}
	})

	t.Run("email is optional", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/clienteles", with(clienteleBody(fixture.customerID, "NE-"+suffix), "email", ""))
		expectStatus(t, result, http.StatusCreated, "created")
	})

	t.Run("update", func(t *testing.T) {
		body := with(clienteleBody(fixture.customerID, code), "name", "Renamed "+suffix)
		body["status"] = "inactive"
		result := sendRequest(t, router, http.MethodPut, clientelePath(created.ID), body)
		expectStatus(t, result, http.StatusOK, "success")
		var updated model.Clientele
		decodeData(t, result, &updated)
		if updated.Name != "Renamed "+suffix || updated.Status != "inactive" {
			t.Errorf("updated = %+v", updated)
		}
	})

	t.Run("company scope", func(t *testing.T) {
		ownCustomer := newPartRouter(actAs(fixture.customerID, model.RoleCustomer), uploadDir)
		otherCustomer := newPartRouter(actAs(otherFixture.customerID, model.RoleCustomer), uploadDir)

		result := sendRequest(t, ownCustomer, http.MethodGet, clientelePath(created.ID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		result = sendRequest(t, otherCustomer, http.MethodGet, clientelePath(created.ID), nil)
		expectStatus(t, result, http.StatusNotFound, "Clientele not found")
		result = sendRequest(t, otherCustomer, http.MethodPut, clientelePath(created.ID), clienteleBody(otherFixture.customerID, code))
		expectStatus(t, result, http.StatusNotFound, "Clientele not found")

		// A customer adds to its own company whatever userId says
		result = sendRequest(t, otherCustomer, http.MethodPost, apiV1+"/clienteles", clienteleBody(fixture.customerID, "OWN-"+suffix))
		expectStatus(t, result, http.StatusCreated, "created")
		var own model.Clientele
		decodeData(t, result, &own)
		if own.UserID != otherFixture.customerID {
			t.Errorf("userId = %d; want the acting customer %d", own.UserID, otherFixture.customerID)
		}
	})

	// Parts refer to a clientele of their own customer
	machineID := createTestMachine(t, fixture, "VMC-C-"+suffix)
	drawingResult := uploadDrawing(t, router, "drawing.pdf", testPDF)
	var drawing model.PartDrawing
	decodeData(t, drawingResult, &drawing)
	partBody := func(clienteleID int64) map[string]any {
		return map[string]any{
			"userId":      fixture.customerID,
			"clienteleId": clienteleID,
			"partNo":      "PN-C-" + suffix,
			"partName":    "Clientele part " + suffix,
			"partCount":   1,
			"processes":   []map[string]any{partProcessBody("P1", drawing, []int64{machineID})},
		}
	}

	var otherClientele model.Clientele
	result := sendRequest(t, router, http.MethodGet, apiV1+"/clienteles?userId="+fmt.Sprint(otherFixture.customerID), nil)
	var otherList []model.Clientele
	decodeData(t, result, &otherList)
	if len(otherList) > 0 {
		otherClientele = otherList[0]
	}

	t.Run("part rejects another customer's clientele", func(t *testing.T) {
		if otherClientele.ID == 0 {
			t.Skip("no clientele of the other customer")
		}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/parts", partBody(otherClientele.ID))
		expectStatus(t, result, http.StatusBadRequest, "Please select a clientele of the selected customer")
	})

	var part model.Part
	t.Run("part with clientele", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/parts", partBody(created.ID))
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &part)
		if part.ClienteleID == nil || *part.ClienteleID != created.ID || part.ClienteleName != "Renamed "+suffix {
			t.Errorf("part clientele = %v %q", part.ClienteleID, part.ClienteleName)
		}
	})

	t.Run("used clientele cannot be moved or deleted", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPut, clientelePath(created.ID), clienteleBody(otherFixture.customerID, "MV-"+suffix))
		expectStatus(t, result, http.StatusBadRequest, "This clientele is used by parts and cannot be moved to another customer")
		result = sendRequest(t, router, http.MethodDelete, clientelePath(created.ID), nil)
		expectStatus(t, result, http.StatusBadRequest, "This clientele is used by parts and cannot be deleted")
	})

	t.Run("delete once unused", func(t *testing.T) {
		if part.ID != 0 {
			result := sendRequest(t, router, http.MethodDelete, fmt.Sprintf("%s/parts/%d", apiV1, part.ID), nil)
			expectStatus(t, result, http.StatusOK, "Part deleted")
		}
		result := sendRequest(t, router, http.MethodDelete, clientelePath(created.ID), nil)
		expectStatus(t, result, http.StatusOK, "Clientele deleted")
		result = sendRequest(t, router, http.MethodGet, clientelePath(created.ID), nil)
		expectStatus(t, result, http.StatusNotFound, "Clientele not found")
	})
}
