package tests

import (
	"context"
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

// newOperatorTaskRouter builds the real router with the maintenance task and operator task
// handlers wired to testDBPool; authenticate decides who the requests act as.
func newOperatorTaskRouter(authenticate middleware.Middleware) http.Handler {
	brandService := service.NewBrandService(repository.NewBrandRepository(testDBPool))
	controllerService := service.NewControllerService(repository.NewControllerRepository(testDBPool), brandService)
	machineModelService := service.NewMachineModelService(repository.NewMachineModelRepository(testDBPool), controllerService)
	maintenanceTaskService := service.NewMaintenanceTaskService(repository.NewMaintenanceTaskRepository(testDBPool))
	operatorTaskService := service.NewOperatorTaskService(repository.NewOperatorTaskRepository(testDBPool),
		repository.NewUserRepository(testDBPool), machineModelService, maintenanceTaskService)
	return routes.New(routes.Handlers{
		MaintenanceTask: handler.NewMaintenanceTaskHandler(maintenanceTaskService),
		OperatorTask:    handler.NewOperatorTaskHandler(operatorTaskService),
	}, authenticate)
}

// TestOperatorTasks covers the maintenance task list (seeded, quick add, duplicates) and operator
// task CRUD: validation, the once-per-customer-and-machine-type rule, superadmin-only writes and
// company-scoped reads.
func TestOperatorTasks(t *testing.T) {
	router := newOperatorTaskRouter(asTestUser) // superadmin
	suffix := uniqueSuffix()
	fixture := newMachineFixture(t, suffix)
	customer := fixture.customerID
	ctx := context.Background()

	var seeded []model.MaintenanceTask
	t.Run("the seeded task list", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, apiV1+"/maintenance-tasks", nil)
		expectStatus(t, result, http.StatusOK, "success")
		decodeData(t, result, &seeded)
		if len(seeded) < 15 {
			t.Fatalf("maintenance tasks = %d; want the 15 seeded at least", len(seeded))
		}
	})

	var added model.MaintenanceTask
	t.Run("add a task to the list", func(t *testing.T) {
		description := "Check coolant level " + suffix
		result := sendRequest(t, router, http.MethodPost, apiV1+"/maintenance-tasks", map[string]any{"description": "  Check   coolant level " + suffix + " "})
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &added)
		if added.ID == 0 || added.Description != description || added.Status != "active" {
			t.Fatalf("added = %+v; want %q, active", added, description)
		}
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/maintenance-tasks", map[string]any{"description": "CHECK COOLANT LEVEL " + suffix}),
			http.StatusBadRequest, "This task already exists")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/maintenance-tasks", map[string]any{"description": "   "}),
			http.StatusBadRequest, "Please enter task")
	})
	if added.ID == 0 {
		t.FailNow()
	}
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(ctx, `UPDATE maintenance_tasks SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, added.ID)
	})

	body := map[string]any{"userId": customer, "modelId": fixture.modelID, "maintenanceTaskId": added.ID}
	var created model.OperatorTask
	t.Run("create", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPost, apiV1+"/operator-tasks", body)
		expectStatus(t, result, http.StatusCreated, "created")
		decodeData(t, result, &created)
		switch {
		case created.ID == 0 || created.UserID != customer || created.ModelID == nil || *created.ModelID != fixture.modelID:
			t.Fatalf("created = %+v", created)
		case created.Task != added.Description || created.ModelName == "" || created.ControllerName == "" || created.BrandName == "" || created.CustomerName == "":
			t.Errorf("names = %+v", created)
		}
	})
	if created.ID == 0 {
		t.FailNow()
	}
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(ctx, `UPDATE operator_tasks SET is_deleted = true, deleted_at = NOW() WHERE user_id = $1`, customer)
	})
	taskPath := fmt.Sprintf("%s/operator-tasks/%d", apiV1, created.ID)

	t.Run("validation and duplicates", func(t *testing.T) {
		operatorID := createCompanyMember(t, customer, "otask_op_"+suffix, model.RoleOperator)
		cases := []struct {
			name    string
			body    map[string]any
			message string
		}{
			{"same task, customer and machine type", body, "This task is already added for this customer and machine type"},
			{"no customer", with(body, "userId", 0), "Please select customer"},
			{"not a customer", with(body, "userId", operatorID), "Please select valid customer"},
			{"unknown model", with(body, "modelId", 999999999), "Please select valid machine type"},
			{"no task", with(body, "maintenanceTaskId", 0), "Please select task"},
			{"unknown task", with(body, "maintenanceTaskId", 999999999), "Please select valid task"},
		}
		for _, testCase := range cases {
			t.Run(testCase.name, func(t *testing.T) {
				expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/operator-tasks", testCase.body), http.StatusBadRequest, testCase.message)
			})
		}
		// No machine type is its own type: the same task for every machine is allowed once
		allMachines := with(body, "modelId", nil)
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/operator-tasks", allMachines), http.StatusCreated, "created")
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/operator-tasks", allMachines), http.StatusBadRequest,
			"This task is already added for this customer and machine type")
	})

	t.Run("update", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPut, taskPath, with(with(body, "modelId", fixture.otherModelID), "maintenanceTaskId", seeded[0].ID))
		expectStatus(t, result, http.StatusOK, "success")
		var updated model.OperatorTask
		decodeData(t, result, &updated)
		if updated.ModelID == nil || *updated.ModelID != fixture.otherModelID || updated.MaintenanceTaskID != seeded[0].ID || updated.Task != seeded[0].Description {
			t.Fatalf("updated = %+v", updated)
		}
		expectStatus(t, sendRequest(t, router, http.MethodPut, apiV1+"/operator-tasks/999999999", body), http.StatusNotFound, "Operator task not found")
	})

	t.Run("only the superadmin writes; members read their company's", func(t *testing.T) {
		member := newOperatorTaskRouter(actAsMember(customer, model.RoleCustomer, customer))
		expectStatus(t, sendRequest(t, member, http.MethodPost, apiV1+"/operator-tasks", body), http.StatusForbidden, "You do not have permission to do this")
		expectStatus(t, sendRequest(t, member, http.MethodPost, apiV1+"/maintenance-tasks", map[string]any{"description": "x"}), http.StatusForbidden, "You do not have permission to do this")
		expectStatus(t, sendRequest(t, member, http.MethodDelete, taskPath, nil), http.StatusForbidden, "You do not have permission to do this")

		result := sendRequest(t, member, http.MethodGet, apiV1+"/operator-tasks", nil)
		expectStatus(t, result, http.StatusOK, "success")
		var tasks []model.OperatorTask
		decodeData(t, result, &tasks)
		if len(tasks) != 2 {
			t.Fatalf("customer sees %d operator tasks; want its 2", len(tasks))
		}

		other := newMachineFixture(t, suffix+"5")
		outsider := newOperatorTaskRouter(actAsMember(other.customerID, model.RoleCustomer, other.customerID))
		expectStatus(t, sendRequest(t, outsider, http.MethodGet, taskPath, nil), http.StatusNotFound, "Operator task not found")
		expectStatus(t, sendRequest(t, outsider, http.MethodGet, fmt.Sprintf("%s/operator-tasks?userId=%d", apiV1, customer), nil),
			http.StatusForbidden, "You do not have permission to do this")
	})

	t.Run("list and delete", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, fmt.Sprintf("%s/operator-tasks?userId=%d", apiV1, customer), nil)
		expectStatus(t, result, http.StatusOK, "success")
		var tasks []model.OperatorTask
		decodeData(t, result, &tasks)
		if len(tasks) != 2 {
			t.Fatalf("operator tasks = %d; want 2", len(tasks))
		}
		expectStatus(t, sendRequest(t, router, http.MethodDelete, taskPath, nil), http.StatusOK, "Operator task deleted")
		expectStatus(t, sendRequest(t, router, http.MethodGet, taskPath, nil), http.StatusNotFound, "Operator task not found")
		expectStatus(t, sendRequest(t, router, http.MethodDelete, taskPath, nil), http.StatusNotFound, "Operator task not found")
	})
}

// TestMaintenanceTaskList covers the Task List screen's endpoints: edit (description, status),
// the usage count, an inactive task refused for new operator tasks, and delete refused while in use.
func TestMaintenanceTaskList(t *testing.T) {
	router := newOperatorTaskRouter(asTestUser) // superadmin
	suffix := uniqueSuffix()
	fixture := newMachineFixture(t, suffix)
	ctx := context.Background()

	createTask := func(description string) model.MaintenanceTask {
		t.Helper()
		result := sendRequest(t, router, http.MethodPost, apiV1+"/maintenance-tasks", map[string]any{"description": description})
		expectStatus(t, result, http.StatusCreated, "created")
		var task model.MaintenanceTask
		decodeData(t, result, &task)
		t.Cleanup(func() {
			_, _ = testDBPool.Exec(ctx, `UPDATE maintenance_tasks SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, task.ID)
		})
		return task
	}
	used := createTask("Coolant not working and need to fix " + suffix)
	unused := createTask("Check door seal " + suffix)
	taskPath := func(id int64) string { return fmt.Sprintf("%s/maintenance-tasks/%d", apiV1, id) }

	result := sendRequest(t, router, http.MethodPost, apiV1+"/operator-tasks",
		map[string]any{"userId": fixture.customerID, "maintenanceTaskId": used.ID})
	expectStatus(t, result, http.StatusCreated, "created")
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(ctx, `UPDATE operator_tasks SET is_deleted = true, deleted_at = NOW() WHERE user_id = $1`, fixture.customerID)
	})

	t.Run("usage count", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodGet, apiV1+"/maintenance-tasks", nil)
		expectStatus(t, result, http.StatusOK, "success")
		var tasks []model.MaintenanceTask
		decodeData(t, result, &tasks)
		counts := map[int64]int{}
		for _, task := range tasks {
			counts[task.ID] = task.UsageCount
		}
		if counts[used.ID] != 1 || counts[unused.ID] != 0 {
			t.Fatalf("usage counts = %d, %d; want 1, 0", counts[used.ID], counts[unused.ID])
		}
	})

	t.Run("edit description and status", func(t *testing.T) {
		result := sendRequest(t, router, http.MethodPut, taskPath(unused.ID), map[string]any{"description": " Check door seals " + suffix + " ", "status": "inactive"})
		expectStatus(t, result, http.StatusOK, "success")
		var updated model.MaintenanceTask
		decodeData(t, result, &updated)
		if updated.Description != "Check door seals "+suffix || updated.Status != "inactive" {
			t.Fatalf("updated = %+v", updated)
		}
		expectStatus(t, sendRequest(t, router, http.MethodPut, taskPath(unused.ID), map[string]any{"description": used.Description}),
			http.StatusBadRequest, "This task already exists")
		expectStatus(t, sendRequest(t, router, http.MethodPut, taskPath(unused.ID), map[string]any{"description": "x", "status": "maybe"}),
			http.StatusBadRequest, "Please enter valid status")
		expectStatus(t, sendRequest(t, router, http.MethodPut, taskPath(999999999), map[string]any{"description": "x"}),
			http.StatusNotFound, "Task not found")
	})

	t.Run("an inactive task cannot be given to a customer", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/operator-tasks",
			map[string]any{"userId": fixture.customerID, "maintenanceTaskId": unused.ID}), http.StatusBadRequest, "Please select valid task")
	})

	t.Run("delete", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, taskPath(used.ID), nil), http.StatusBadRequest,
			"This task is used by 1 operator task(s). Remove it from them first, or make it inactive")
		expectStatus(t, sendRequest(t, router, http.MethodDelete, taskPath(unused.ID), nil), http.StatusOK, "Task deleted")
		expectStatus(t, sendRequest(t, router, http.MethodDelete, taskPath(unused.ID), nil), http.StatusNotFound, "Task not found")
	})

	t.Run("only the superadmin edits and deletes", func(t *testing.T) {
		member := newOperatorTaskRouter(actAsMember(fixture.customerID, model.RoleCustomer, fixture.customerID))
		expectStatus(t, sendRequest(t, member, http.MethodPut, taskPath(used.ID), map[string]any{"description": "x"}),
			http.StatusForbidden, "You do not have permission to do this")
		expectStatus(t, sendRequest(t, member, http.MethodDelete, taskPath(used.ID), nil),
			http.StatusForbidden, "You do not have permission to do this")
	})
}
