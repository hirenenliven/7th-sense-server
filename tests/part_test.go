package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"testing"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// testPDF is the smallest content the drawing upload accepts as a PDF.
var testPDF = []byte("%PDF-1.4\n%test drawing\n%%EOF\n")

// newPartRouter builds the real router with the part and clientele handlers wired to
// testDBPool, keeping drawings in uploadDir; authenticate decides who the requests act as.
func newPartRouter(authenticate middleware.Middleware, uploadDir string) http.Handler {
	clienteleService := service.NewClienteleService(repository.NewClienteleRepository(testDBPool), repository.NewUserRepository(testDBPool))
	partService := service.NewPartService(testDBPool, repository.NewUserRepository(testDBPool),
		newMachineService(), clienteleService, repository.NewPartDrawingStorage(uploadDir))
	return routes.New(routes.Handlers{
		Part:      handler.NewPartHandler(partService),
		Clientele: handler.NewClienteleHandler(clienteleService),
	}, authenticate)
}

// uploadDrawing sends content as the "file" field of POST /parts/drawings.
func uploadDrawing(t *testing.T, router http.Handler, fileName string, content []byte) apiResponse {
	t.Helper()
	var body bytes.Buffer
	form := multipart.NewWriter(&body)
	fileField, err := form.CreateFormFile("file", fileName)
	if err != nil {
		t.Fatalf("create form file: %v", err)
	}
	if _, err := fileField.Write(content); err != nil {
		t.Fatalf("write form file: %v", err)
	}
	if err := form.Close(); err != nil {
		t.Fatalf("close form: %v", err)
	}

	request := httptest.NewRequest(http.MethodPost, apiV1+"/parts/drawings", &body)
	request.Header.Set("Content-Type", form.FormDataContentType())
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)

	var result apiResponse
	if err := json.NewDecoder(recorder.Body).Decode(&result); err != nil {
		t.Fatalf("upload drawing: decode response: %v", err)
	}
	return result
}

// createTestMachine adds a machine of fixture's customer through the machine API.
func createTestMachine(t *testing.T, fixture machineFixture, name string) int64 {
	t.Helper()
	result := sendRequest(t, newMachineRouter(asTestUser), http.MethodPost, apiV1+"/machines", machineBody(fixture, name))
	expectStatus(t, result, http.StatusCreated, "created")
	var machine machineRecord
	decodeData(t, result, &machine)
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(t.Context(), `UPDATE machines SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, machine.ID)
	})
	return machine.ID
}

// partProcessBody is a valid process of the Add Part form.
func partProcessBody(name string, drawing model.PartDrawing, machineIDs []int64) map[string]any {
	return map[string]any{
		"process":           name,
		"programName":       "PRG_" + name,
		"machineIds":        machineIDs,
		"drawingNo":         "DRW-" + name,
		"drawingRevisionNo": "1",
		"drawingFile":       drawing.File,
		"drawingName":       drawing.Name,
		"accuracyId":        1,
		"setupTime":         900,
		"cycleTime":         3600,
		"loadingTime":       900,
		"hasSubProgram":     false,
		"subPrograms":       []any{},
	}
}

func TestPartCRUD(t *testing.T) {
	uploadDir := t.TempDir()
	router := newPartRouter(asTestUser, uploadDir) // superadmin
	suffix := uniqueSuffix()
	fixture := newMachineFixture(t, suffix)
	machineA := createTestMachine(t, fixture, "VMC-A-"+suffix)
	machineB := createTestMachine(t, fixture, "VMC-B-"+suffix)
	otherFixture := newMachineFixture(t, suffix+"9")
	otherMachine := createTestMachine(t, otherFixture, "VMC-O-"+suffix)
	partPath := func(id int64) string { return fmt.Sprintf("%s/parts/%d", apiV1, id) }

	var drawing model.PartDrawing
	t.Run("upload drawing", func(t *testing.T) {
		result := uploadDrawing(t, router, "shaft drawing.pdf", testPDF)
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &drawing)
		if drawing.Name != "shaft drawing.pdf" || len(drawing.File) != 36 {
			t.Fatalf("drawing = %+v; want the uploaded name and a stored file name", drawing)
		}
	})
	if drawing.File == "" {
		t.FailNow()
	}

	t.Run("upload rejects a file that is not a PDF", func(t *testing.T) {
		result := uploadDrawing(t, router, "drawing.pdf", []byte("not a pdf"))
		expectStatus(t, result, http.StatusBadRequest, "Part drawing must be a PDF file")
	})

	withSubPrograms := partProcessBody("P2", drawing, []int64{machineB})
	withSubPrograms["hasSubProgram"] = true
	withSubPrograms["subPrograms"] = []map[string]any{{"name": "S_SUB1", "cycleTime": 1800}, {"name": "S_SUB2", "cycleTime": 900}}
	partBody := map[string]any{
		"userId":    fixture.customerID,
		"partNo":    "PN-" + suffix,
		"partName":  "Shaft " + suffix,
		"partCount": 2,
		"processes": []map[string]any{partProcessBody("P1", drawing, []int64{machineA, machineB, machineA}), withSubPrograms},
	}

	var created model.Part
	t.Run("create", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/parts", partBody)
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		switch {
		case created.ID == 0 || created.UserID != fixture.customerID || created.Status != "active" || created.CustomerName == "":
			t.Fatalf("created = %+v", created)
		case len(created.Processes) != 2 || created.ProcessCount != 2:
			t.Fatalf("processes = %d (count %d); want 2", len(created.Processes), created.ProcessCount)
		// P1: 900 + 3600 + 900; P2: the same + sub programs 1800 + 900
		case created.Processes[0].ProductionTime != 5400 || created.Processes[1].ProductionTime != 8100 || created.ProductionTime != 13500:
			t.Errorf("production times = %d, %d, part %d; want 5400, 8100, 13500",
				created.Processes[0].ProductionTime, created.Processes[1].ProductionTime, created.ProductionTime)
		case len(created.Processes[0].MachineIDs) != 2:
			t.Errorf("P1 machines = %v; want the duplicate dropped", created.Processes[0].MachineIDs)
		case len(created.Processes[1].SubPrograms) != 2 || created.Processes[1].SubPrograms[0].Name != "S_SUB1":
			t.Errorf("P2 sub programs = %+v", created.Processes[1].SubPrograms)
		}
	})
	if created.ID == 0 {
		t.FailNow()
	}

	t.Run("duplicate part no. and name", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/parts", partBody)
		expectStatus(t, result, http.StatusBadRequest, "This customer already has a part with this part no. and name")
	})

	t.Run("machine of another customer", func(t *testing.T) {
		body := with(partBody, "partNo", "PN-X-"+suffix)
		body["processes"] = []map[string]any{partProcessBody("P1", drawing, []int64{otherMachine})}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/parts", body)
		expectStatus(t, result, http.StatusBadRequest, "Please select applicable machines of the selected customer")
	})

	t.Run("drawing that was never uploaded", func(t *testing.T) {
		body := with(partBody, "partNo", "PN-Y-"+suffix)
		body["processes"] = []map[string]any{partProcessBody("P1", model.PartDrawing{File: "0123456789abcdef0123456789abcdef.pdf", Name: "x.pdf"}, []int64{machineA})}
		result := sendRequest(t, router, http.MethodPost, apiV1+"/parts", body)
		expectStatus(t, result, http.StatusBadRequest, "Please upload part drawing (PDF) again")
	})

	t.Run("validation names the process", func(t *testing.T) {
		process := partProcessBody("P1", drawing, []int64{machineA})
		process["cycleTime"] = 0
		result := sendRequest(t, router, http.MethodPost, apiV1+"/parts", with(partBody, "processes", []map[string]any{process}))
		expectStatus(t, result, http.StatusBadRequest, "Please enter valid cycle time (process 1)")
	})

	t.Run("update keeps, drops and adds processes", func(t *testing.T) {
		kept := partProcessBody("P1-edited", drawing, []int64{machineB})
		kept["id"] = created.Processes[0].ID
		body := with(partBody, "partName", "Shaft edited "+suffix)
		body["status"] = "inactive"
		body["processes"] = []map[string]any{kept, partProcessBody("P3", drawing, []int64{machineA})}

		result := sendRequest(t, router, http.MethodPut, partPath(created.ID), body)
		expectStatus(t, result, http.StatusOK, "success")
		var updated model.Part
		decodeData(t, result, &updated)
		switch {
		case updated.PartName != "Shaft edited "+suffix || updated.Status != "inactive":
			t.Errorf("updated = %+v", updated)
		case len(updated.Processes) != 2 || updated.Processes[0].ID != created.Processes[0].ID || updated.Processes[0].Process != "P1-edited":
			t.Fatalf("processes = %+v; want P1 kept (edited) then P3", updated.Processes)
		case updated.Processes[1].Process != "P3" || updated.ProductionTime != 10800:
			t.Errorf("P3 / total = %q / %d; want P3 / 10800", updated.Processes[1].Process, updated.ProductionTime)
		case len(updated.Processes[0].MachineIDs) != 1 || updated.Processes[0].MachineIDs[0] != machineB:
			t.Errorf("P1 machines = %v; want only %d", updated.Processes[0].MachineIDs, machineB)
		}
	})

	t.Run("update rejects a process of another part", func(t *testing.T) {
		process := partProcessBody("P1", drawing, []int64{machineA})
		process["id"] = created.Processes[1].ID // soft-deleted by the update above
		result := sendRequest(t, router, http.MethodPut, partPath(created.ID), with(partBody, "processes", []map[string]any{process}))
		expectStatus(t, result, http.StatusBadRequest, "Please enter valid processes")
	})

	t.Run("drawing download", func(t *testing.T) {
		request := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%s/processes/%d/drawing", partPath(created.ID), created.Processes[0].ID), nil)
		recorder := httptest.NewRecorder()
		router.ServeHTTP(recorder, request)
		content, _ := io.ReadAll(recorder.Body)
		if recorder.Code != http.StatusOK || recorder.Header().Get("Content-Type") != "application/pdf" || !bytes.Equal(content, testPDF) {
			t.Fatalf("drawing = %d %q (%d bytes); want the uploaded PDF", recorder.Code, recorder.Header().Get("Content-Type"), len(content))
		}
	})

	t.Run("company scope", func(t *testing.T) {
		ownCustomer := newPartRouter(actAs(fixture.customerID, model.RoleCustomer), uploadDir)
		otherCustomer := newPartRouter(actAs(otherFixture.customerID, model.RoleCustomer), uploadDir)

		result := sendRequest(t, ownCustomer, http.MethodGet, partPath(created.ID), nil)
		expectStatus(t, result, http.StatusOK, "success")
		result = sendRequest(t, otherCustomer, http.MethodGet, partPath(created.ID), nil)
		expectStatus(t, result, http.StatusNotFound, "Part not found")
		result = sendRequest(t, otherCustomer, http.MethodDelete, partPath(created.ID), nil)
		expectStatus(t, result, http.StatusNotFound, "Part not found")
		result = sendRequest(t, otherCustomer, http.MethodGet, apiV1+"/parts?userId="+fmt.Sprint(fixture.customerID), nil)
		expectStatus(t, result, http.StatusForbidden, "You do not have permission to do this")
	})

	t.Run("operators read but do not write", func(t *testing.T) {
		operator := newPartRouter(actAs(fixture.customerID, model.RoleOperator), uploadDir)
		result := sendRequest(t, operator, http.MethodDelete, partPath(created.ID), nil)
		expectStatus(t, result, http.StatusForbidden, "You do not have permission to do this")
	})

	t.Run("delete", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodDelete, partPath(created.ID), nil)
		expectStatus(t, result, http.StatusOK, "Part deleted")
		result = sendRequest(t, router, http.MethodGet, partPath(created.ID), nil)
		expectStatus(t, result, http.StatusNotFound, "Part not found")
		result = sendRequest(t, router, http.MethodDelete, partPath(created.ID), nil)
		expectStatus(t, result, http.StatusNotFound, "Part not found")
	})
}
