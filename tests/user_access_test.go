package tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

const forbiddenMessage = "You do not have permission to do this"

// newUserHandler builds the user handler on the test database.
func newUserHandler() *handler.UserHandler {
	userCountryService := service.NewUserCountryService(repository.NewUserCountryRepository(testDBPool))
	userStateService := service.NewUserStateService(repository.NewUserStateRepository(testDBPool), userCountryService)
	userCityService := service.NewUserCityService(repository.NewUserCityRepository(testDBPool), userStateService)
	branchService := service.NewBranchService(repository.NewBranchRepository(testDBPool), userCityService)
	roleService := service.NewRoleService(repository.NewRoleRepository(testDBPool))
	return handler.NewUserHandler(service.NewUserService(testDBPool, branchService, roleService, newMachineService()))
}

// newUserAccessRouter serves the user, role and machine routes, acting as
// authenticate's user.
func newUserAccessRouter(authenticate middleware.Middleware) http.Handler {
	return routes.New(routes.Handlers{
		User:    newUserHandler(),
		Role:    handler.NewRoleHandler(service.NewRoleService(repository.NewRoleRepository(testDBPool))),
		Machine: handler.NewMachineHandler(newMachineService()),
	}, authenticate)
}

// actAsMember makes every request act as userID with role, as a member of the
// customer companyID.
func actAsMember(userID int64, role string, companyID int64) middleware.Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			user := model.AuthUser{ID: userID, Role: role, CompanyID: &companyID}
			next.ServeHTTP(writer, request.WithContext(middleware.WithAuthUser(request.Context(), user)))
		})
	}
}

// testCompany is a customer with its organization and machines.
type testCompany struct {
	machineFixture
	machineIDs []int64
}

// newTestCompany creates a customer with machineCount machines (added by the
// superadmin). Everything is soft-deleted when the test ends.
func newTestCompany(t *testing.T, suffix string, machineCount int) testCompany {
	t.Helper()
	company := testCompany{machineFixture: newMachineFixture(t, suffix)}
	superadmin := newUserAccessRouter(asTestUser)
	for index := range machineCount {
		result := sendRequest(t, superadmin, http.MethodPost, apiV1+"/machines", machineBody(company.machineFixture, fmt.Sprintf("UA-%s-%d", suffix, index)))
		expectStatus(t, result, http.StatusCreated, "created")
		var machine machineRecord
		decodeData(t, result, &machine)
		company.machineIDs = append(company.machineIDs, machine.ID)
		t.Cleanup(func() {
			_, _ = testDBPool.Exec(context.Background(), `UPDATE machines SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, machine.ID)
		})
	}
	return company
}

// testRoleID returns the id of the role called name.
func testRoleID(t *testing.T, name string) int64 {
	t.Helper()
	var id int64
	if err := testDBPool.QueryRow(context.Background(), `SELECT id FROM roles WHERE name = $1`, name).Scan(&id); err != nil {
		t.Fatalf("find role %s: %v", name, err)
	}
	return id
}

// userBody is a valid POST /users body: a user named tag with role, on company's branch.
func userBody(t *testing.T, company testCompany, tag, role string, machineIDs []int64) dto.CreateUserRequest {
	t.Helper()
	return dto.CreateUserRequest{
		Username:   "ua_" + tag,
		Email:      "ua_" + tag + "@user-access.example.com",
		Password:   "secret123",
		RoleID:     testRoleID(t, role),
		MachineIDs: machineIDs,
		CustomerDetails: dto.CustomerDetails{
			BranchID: &company.branchID, CompanyName: "Member " + tag, GSTNumber: "24AAACE1234F1Z5", PANNumber: "AAACE1234F",
			AddressLine1: "Plot 12", Country: "India", State: "Gujarat", City: "Rajkot", ZipCode: "360001",
			ContactPerson: "Member", PrimaryPhone: "+919876543210", TimeZone: "Asia/Kolkata", DayStartFrom: "06:00",
		},
	}
}

// updateBody is the PUT /users/{id} body with the fields of body.
func updateBody(body dto.CreateUserRequest) dto.UpdateUserRequest {
	return dto.UpdateUserRequest{
		Username: body.Username, Email: body.Email, RoleID: body.RoleID, CompanyID: body.CompanyID,
		MachineIDs: body.MachineIDs, CustomerDetails: body.CustomerDetails,
	}
}

// userRecord is the part of a user response the access tests check.
type userRecord struct {
	ID         int64   `json:"id"`
	Role       string  `json:"role"`
	CompanyID  *int64  `json:"companyId"`
	MachineIDs []int64 `json:"machineIds"`
}

// createUser posts body as router's user, expects 201 and soft-deletes the user
// when the test ends.
func createUser(t *testing.T, router http.Handler, body dto.CreateUserRequest) userRecord {
	t.Helper()
	result := sendRequest(t, router, http.MethodPost, apiV1+"/users", body)
	expectStatus(t, result, http.StatusCreated, "created")
	var user userRecord
	decodeData(t, result, &user)
	t.Cleanup(func() {
		_, _ = testDBPool.Exec(context.Background(), `UPDATE users SET is_deleted = true, deleted_at = NOW() WHERE id = $1`, user.ID)
	})
	return user
}

// listIDs returns the ids of a list response (the server omits data for an empty list).
func listIDs(t *testing.T, result apiResponse) []int64 {
	t.Helper()
	var records []struct {
		ID int64 `json:"id"`
	}
	if len(result.Data) > 0 {
		decodeData(t, result, &records)
	}
	ids := make([]int64, 0, len(records))
	for _, record := range records {
		ids = append(ids, record.ID)
	}
	slices.Sort(ids)
	return ids
}

// TestUserListAccess checks GET /users: the superadmin, customers and admins of a
// company may list users; an admin without a company and every other role get 403.
func TestUserListAccess(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		authenticate middleware.Middleware
		status       int
		message      string
	}{
		{"superadmin", actAs(testUserID, model.RoleSuperAdmin), http.StatusOK, "success"},
		{"customer", actAs(testUserID, model.RoleCustomer), http.StatusOK, "success"},
		{"admin of a company", actAsMember(testUserID, model.RoleAdmin, testUserID), http.StatusOK, "success"},
		{"admin without a company", actAs(testUserID, model.RoleAdmin), http.StatusForbidden, forbiddenMessage},
		{"operator", actAsMember(testUserID, model.RoleOperator, testUserID), http.StatusForbidden, forbiddenMessage},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			router := newUserAccessRouter(testCase.authenticate)
			expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/users", nil), testCase.status, testCase.message)
		})
	}
}

// TestAssignableRoles checks GET /roles, the role dropdown of the user form.
func TestAssignableRoles(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		authenticate middleware.Middleware
		want         []string
	}{
		{"superadmin", actAs(testUserID, model.RoleSuperAdmin), []string{model.RoleAdmin, model.RoleSupervisor, model.RoleOperator, model.RoleCustomer, model.RoleQC, model.RoleShopHead}},
		{"customer hides customer and superadmin", actAs(testUserID, model.RoleCustomer), []string{model.RoleAdmin, model.RoleSupervisor, model.RoleOperator, model.RoleQC, model.RoleShopHead}},
		{"admin also hides admin", actAsMember(testUserID, model.RoleAdmin, testUserID), []string{model.RoleSupervisor, model.RoleOperator, model.RoleQC, model.RoleShopHead}},
		{"operator gets none", actAsMember(testUserID, model.RoleOperator, testUserID), nil},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			result := sendRequest(t, newUserAccessRouter(testCase.authenticate), http.MethodGet, apiV1+"/roles", nil)
			expectStatus(t, result, http.StatusOK, "success")
			var roles []model.Role
			if len(result.Data) > 0 {
				decodeData(t, result, &roles)
			}
			var names []string
			for _, role := range roles {
				names = append(names, role.Name)
			}
			if !slices.Equal(names, testCase.want) {
				t.Fatalf("roles = %v; want %v", names, testCase.want)
			}
		})
	}
}

// TestNobodyManagesTheSuperadmin checks that the superadmin role cannot be given,
// and that the superadmin cannot be changed or deleted through the user routes.
// Every request is refused, so no data is changed.
func TestNobodyManagesTheSuperadmin(t *testing.T) {
	ctx := context.Background()
	superAdminRoleID := testRoleID(t, model.RoleSuperAdmin)
	var superadminID int64
	err := testDBPool.QueryRow(ctx, `SELECT id FROM users WHERE role_id = $1 AND is_deleted = false ORDER BY id LIMIT 1`, superAdminRoleID).Scan(&superadminID)
	if errors.Is(err, pgx.ErrNoRows) {
		t.Skip("no superadmin user in the database")
	}
	if err != nil {
		t.Fatalf("find superadmin: %v", err)
	}
	superadminPath := fmt.Sprintf("%s/users/%d", apiV1, superadminID)
	company := testCompany{machineFixture: machineFixture{branchID: 1}}

	for _, testCase := range []struct {
		name         string
		authenticate middleware.Middleware
		// status of PUT / DELETE on the superadmin: 404 for users who cannot see it
		status  int
		message string
	}{
		{"superadmin", actAs(testUserID, model.RoleSuperAdmin), http.StatusForbidden, forbiddenMessage},
		{"customer", actAs(testUserID, model.RoleCustomer), http.StatusNotFound, "User not found"},
		{"admin", actAsMember(testUserID, model.RoleAdmin, testUserID), http.StatusNotFound, "User not found"},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			router := newUserAccessRouter(testCase.authenticate)
			body := userBody(t, company, "escalate_"+uniqueSuffix(), model.RoleSuperAdmin, nil)
			expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/users", body), http.StatusForbidden, forbiddenMessage)
			expectStatus(t, sendRequest(t, router, http.MethodPut, superadminPath, updateBody(body)), testCase.status, testCase.message)
			expectStatus(t, sendRequest(t, router, http.MethodDelete, superadminPath, nil), testCase.status, testCase.message)
		})
	}
}

// TestCompanyUserManagement covers user management and machine access across
// two companies end to end: roles, company scoping and assigned machines.
func TestCompanyUserManagement(t *testing.T) {
	suffix := uniqueSuffix()
	companyA := newTestCompany(t, "A"+suffix, 2)
	companyB := newTestCompany(t, "B"+suffix, 1)
	customerA := newUserAccessRouter(actAs(companyA.customerID, model.RoleCustomer))
	customerB := newUserAccessRouter(actAs(companyB.customerID, model.RoleCustomer))
	superadmin := newUserAccessRouter(asTestUser)
	userPath := func(id int64) string { return fmt.Sprintf("%s/users/%d", apiV1, id) }

	admin := createUser(t, customerA, userBody(t, companyA, "adm"+suffix, model.RoleAdmin, companyA.machineIDs))
	adminA := newUserAccessRouter(actAsMember(admin.ID, model.RoleAdmin, companyA.customerID))
	operatorBody := userBody(t, companyA, "op"+suffix, model.RoleOperator, companyA.machineIDs[:1])
	operator := createUser(t, customerA, operatorBody)
	operatorA := newUserAccessRouter(actAsMember(operator.ID, model.RoleOperator, companyA.customerID))

	t.Run("members belong to the customer's company", func(t *testing.T) {
		if admin.Role != model.RoleAdmin || admin.CompanyID == nil || *admin.CompanyID != companyA.customerID {
			t.Fatalf("admin = %+v; want an admin of company %d", admin, companyA.customerID)
		}
		if operator.CompanyID == nil || *operator.CompanyID != companyA.customerID {
			t.Fatalf("operator = %+v; want company %d", operator, companyA.customerID)
		}
	})

	t.Run("admins get every company machine without assignments", func(t *testing.T) {
		if len(admin.MachineIDs) != 0 {
			t.Fatalf("admin machineIds = %v; want none stored", admin.MachineIDs)
		}
		want := slices.Sorted(slices.Values(companyA.machineIDs))
		for name, router := range map[string]http.Handler{"customer": customerA, "admin": adminA} {
			if got := listIDs(t, sendRequest(t, router, http.MethodGet, apiV1+"/machines", nil)); !slices.Equal(got, want) {
				t.Errorf("%s sees machines %v; want %v", name, got, want)
			}
		}
	})

	t.Run("other members see only their assigned machines", func(t *testing.T) {
		if !slices.Equal(operator.MachineIDs, companyA.machineIDs[:1]) {
			t.Fatalf("operator machineIds = %v; want %v", operator.MachineIDs, companyA.machineIDs[:1])
		}
		if got := listIDs(t, sendRequest(t, operatorA, http.MethodGet, apiV1+"/machines", nil)); !slices.Equal(got, companyA.machineIDs[:1]) {
			t.Fatalf("operator sees machines %v; want %v", got, companyA.machineIDs[:1])
		}
		expectStatus(t, sendRequest(t, operatorA, http.MethodGet, fmt.Sprintf("%s/machines/%d", apiV1, companyA.machineIDs[1]), nil),
			http.StatusNotFound, "Machine not found")
		expectStatus(t, sendRequest(t, operatorA, http.MethodGet, fmt.Sprintf("%s/machines?userId=%d", apiV1, companyB.customerID), nil),
			http.StatusForbidden, forbiddenMessage)
	})

	t.Run("machines of another company cannot be assigned", func(t *testing.T) {
		body := userBody(t, companyA, "opx"+suffix, model.RoleOperator, companyB.machineIDs)
		expectStatus(t, sendRequest(t, customerA, http.MethodPost, apiV1+"/users", body), http.StatusBadRequest, "Please select valid machines")
	})

	t.Run("customers and admins give only their allowed roles", func(t *testing.T) {
		for _, role := range []string{model.RoleCustomer, model.RoleSuperAdmin} {
			body := userBody(t, companyA, "deny"+suffix, role, nil)
			expectStatus(t, sendRequest(t, customerA, http.MethodPost, apiV1+"/users", body), http.StatusForbidden, forbiddenMessage)
		}
		body := userBody(t, companyA, "deny"+suffix, model.RoleAdmin, nil)
		expectStatus(t, sendRequest(t, adminA, http.MethodPost, apiV1+"/users", body), http.StatusForbidden, forbiddenMessage)
		// An admin cannot change admins, itself included
		adminBody := userBody(t, companyA, "adm"+suffix, model.RoleAdmin, nil)
		expectStatus(t, sendRequest(t, adminA, http.MethodPut, userPath(admin.ID), updateBody(adminBody)), http.StatusForbidden, forbiddenMessage)

		supervisor := createUser(t, adminA, userBody(t, companyA, "sup"+suffix, model.RoleSupervisor, companyA.machineIDs[1:]))
		if supervisor.CompanyID == nil || *supervisor.CompanyID != companyA.customerID || !slices.Equal(supervisor.MachineIDs, companyA.machineIDs[1:]) {
			t.Fatalf("supervisor = %+v; want company %d with machines %v", supervisor, companyA.customerID, companyA.machineIDs[1:])
		}
	})

	t.Run("other companies cannot see or change the members", func(t *testing.T) {
		expectStatus(t, sendRequest(t, customerB, http.MethodGet, userPath(operator.ID), nil), http.StatusNotFound, "User not found")
		expectStatus(t, sendRequest(t, customerB, http.MethodPut, userPath(operator.ID), updateBody(operatorBody)), http.StatusNotFound, "User not found")
		expectStatus(t, sendRequest(t, customerB, http.MethodDelete, userPath(operator.ID), nil), http.StatusNotFound, "User not found")
		if ids := listIDs(t, sendRequest(t, customerB, http.MethodGet, apiV1+"/users", nil)); slices.Contains(ids, operator.ID) {
			t.Fatalf("company B lists company A's operator: %v", ids)
		}
		expectStatus(t, sendRequest(t, customerA, http.MethodGet, fmt.Sprintf("%s/users?companyId=%d", apiV1, companyB.customerID), nil),
			http.StatusForbidden, forbiddenMessage)
		if ids := listIDs(t, sendRequest(t, customerA, http.MethodGet, apiV1+"/users", nil)); !slices.Contains(ids, operator.ID) || !slices.Contains(ids, admin.ID) {
			t.Fatalf("company A lists %v; want its admin %d and operator %d", ids, admin.ID, operator.ID)
		}
	})

	t.Run("updating replaces the assigned machines", func(t *testing.T) {
		body := operatorBody
		body.MachineIDs = companyA.machineIDs[1:]
		result := sendRequest(t, customerA, http.MethodPut, userPath(operator.ID), updateBody(body))
		expectStatus(t, result, http.StatusOK, "success")
		var updated userRecord
		decodeData(t, result, &updated)
		if !slices.Equal(updated.MachineIDs, companyA.machineIDs[1:]) {
			t.Fatalf("updated machineIds = %v; want %v", updated.MachineIDs, companyA.machineIDs[1:])
		}
		if got := listIDs(t, sendRequest(t, operatorA, http.MethodGet, apiV1+"/machines", nil)); !slices.Equal(got, companyA.machineIDs[1:]) {
			t.Fatalf("operator sees machines %v; want %v", got, companyA.machineIDs[1:])
		}
	})

	t.Run("the superadmin names the company of a member", func(t *testing.T) {
		body := userBody(t, companyB, "sa"+suffix, model.RoleOperator, companyB.machineIDs)
		expectStatus(t, sendRequest(t, superadmin, http.MethodPost, apiV1+"/users", body), http.StatusBadRequest, "Please select valid customer")
		body.CompanyID = &operator.ID // not a customer
		expectStatus(t, sendRequest(t, superadmin, http.MethodPost, apiV1+"/users", body), http.StatusBadRequest, "Please select valid customer")
		body.CompanyID = &companyB.customerID
		member := createUser(t, superadmin, body)
		if member.CompanyID == nil || *member.CompanyID != companyB.customerID || !slices.Equal(member.MachineIDs, companyB.machineIDs) {
			t.Fatalf("member = %+v; want company %d with machines %v", member, companyB.customerID, companyB.machineIDs)
		}
	})

	t.Run("deleting a user removes its machine access", func(t *testing.T) {
		expectStatus(t, sendRequest(t, customerA, http.MethodDelete, userPath(operator.ID), nil), http.StatusOK, "User deleted")
		var count int
		if err := testDBPool.QueryRow(context.Background(), `SELECT COUNT(*) FROM user_machine_access WHERE user_id = $1`, operator.ID).Scan(&count); err != nil {
			t.Fatalf("count machine access: %v", err)
		}
		if count != 0 {
			t.Fatalf("deleted operator still has %d machine access rows", count)
		}
	})
}

// TestSaveOwnTheme checks PUT /me/theme: every logged-in user (here an existing
// user acting as an operator) may save their own theme, and an unknown theme is
// refused. The user's profile row is put back as it was when the test ends.
func TestSaveOwnTheme(t *testing.T) {
	ctx := context.Background()
	var userID int64
	err := testDBPool.QueryRow(ctx, `SELECT id FROM users WHERE is_deleted = false ORDER BY id LIMIT 1`).Scan(&userID)
	if errors.Is(err, pgx.ErrNoRows) {
		t.Skip("no user in the database")
	}
	if err != nil {
		t.Fatalf("find a user: %v", err)
	}

	// The whole profile row as it was, so cleanup can restore it (or remove the
	// row the test added).
	var original struct {
		theme     *string
		isDeleted bool
		deletedAt *time.Time
	}
	err = testDBPool.QueryRow(ctx, `SELECT theme, is_deleted, deleted_at FROM user_profiles WHERE user_id = $1`, userID).
		Scan(&original.theme, &original.isDeleted, &original.deletedAt)
	hadProfile := err == nil
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("read profile of user %d: %v", userID, err)
	}
	t.Cleanup(func() {
		query := `UPDATE user_profiles SET theme = $2, is_deleted = $3, deleted_at = $4 WHERE user_id = $1`
		args := []any{userID, original.theme, original.isDeleted, original.deletedAt}
		if !hadProfile {
			query, args = `DELETE FROM user_profiles WHERE user_id = $1`, []any{userID}
		}
		if _, err := testDBPool.Exec(context.Background(), query, args...); err != nil {
			t.Errorf("restore profile of user %d: %v", userID, err)
		}
	})

	router := routes.New(routes.Handlers{User: newUserHandler()}, actAs(userID, model.RoleOperator))
	themePath := apiV1 + "/me/theme"

	for _, theme := range []string{dto.ThemeLight, dto.ThemeDark} {
		result := sendRequest(t, router, http.MethodPut, themePath, dto.UpdateThemeRequest{Theme: theme})
		expectStatus(t, result, http.StatusOK, "success")
		var saved dto.ThemeResponse
		decodeData(t, result, &saved)

		var stored string
		if err := testDBPool.QueryRow(ctx, `SELECT theme FROM user_profiles WHERE user_id = $1 AND is_deleted = false`, userID).Scan(&stored); err != nil {
			t.Fatalf("read saved theme: %v", err)
		}
		if saved.Theme != theme || stored != theme {
			t.Fatalf("saved %q, stored %q; want %q", saved.Theme, stored, theme)
		}
	}

	expectStatus(t, sendRequest(t, router, http.MethodPut, themePath, dto.UpdateThemeRequest{Theme: "blue"}),
		http.StatusBadRequest, "Please enter valid theme")
}
