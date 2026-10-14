package tests

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/jackc/pgx/v5"

	"7th-sense-server/internal/dto"
	"7th-sense-server/internal/handler"
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
	return handler.NewUserHandler(service.NewUserService(testDBPool, branchService))
}

// TestUserListAccess checks GET /users: the superadmin and admins may list users,
// every other role gets 403.
func TestUserListAccess(t *testing.T) {
	userHandler := newUserHandler()

	for _, testCase := range []struct {
		role    string
		status  int
		message string
	}{
		{model.RoleSuperAdmin, http.StatusOK, "success"},
		{model.RoleAdmin, http.StatusOK, "success"},
		{model.RoleCustomer, http.StatusForbidden, forbiddenMessage},
		{model.RoleOperator, http.StatusForbidden, forbiddenMessage},
	} {
		t.Run(testCase.role, func(t *testing.T) {
			router := routes.New(routes.Handlers{User: userHandler}, actAs(testUserID, testCase.role))
			expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/users", nil), testCase.status, testCase.message)
		})
	}
}

// TestAdminCannotManageSuperadmin checks that an admin may not give the
// superadmin role, nor change or delete the superadmin. Every request is
// refused, so no data is changed.
func TestAdminCannotManageSuperadmin(t *testing.T) {
	ctx := context.Background()
	var superAdminRoleID int64
	err := testDBPool.QueryRow(ctx, `SELECT id FROM roles WHERE name = $1`, model.RoleSuperAdmin).Scan(&superAdminRoleID)
	if errors.Is(err, pgx.ErrNoRows) {
		t.Skip("no superadmin role in the database")
	}
	if err != nil {
		t.Fatalf("find superadmin role: %v", err)
	}

	router := routes.New(routes.Handlers{User: newUserHandler()}, actAs(testUserID, model.RoleAdmin))
	branchID := int64(1)
	details := dto.CustomerDetails{
		BranchID: &branchID, CompanyName: "Escalation Test", GSTNumber: "24AAACE1234F1Z5", PANNumber: "AAACE1234F",
		AddressLine1: "Plot 12", Country: "India", State: "Gujarat", City: "Rajkot", ZipCode: "360001",
		ContactPerson: "Test", PrimaryPhone: "+919876543210", TimeZone: "Asia/Kolkata", DayStartFrom: "06:00",
	}

	t.Run("create with superadmin role", func(t *testing.T) {
		suffix := uniqueSuffix()
		body := dto.CreateUserRequest{
			Username: "escalate_" + suffix, Email: fmt.Sprintf("escalate_%s@example.com", suffix), Password: "secret123",
			RoleID: superAdminRoleID, CustomerDetails: details,
		}
		expectStatus(t, sendRequest(t, router, http.MethodPost, apiV1+"/users", body), http.StatusForbidden, forbiddenMessage)
	})

	var superadminID int64
	err = testDBPool.QueryRow(ctx, `SELECT id FROM users WHERE role_id = $1 AND is_deleted = false ORDER BY id LIMIT 1`, superAdminRoleID).Scan(&superadminID)
	if errors.Is(err, pgx.ErrNoRows) {
		t.Skip("no superadmin user in the database")
	}
	if err != nil {
		t.Fatalf("find superadmin: %v", err)
	}
	superadminPath := fmt.Sprintf("%s/users/%d", apiV1, superadminID)

	t.Run("update superadmin", func(t *testing.T) {
		body := dto.UpdateUserRequest{
			Username: "superadmin", Email: "superadmin@example.com", RoleID: superAdminRoleID, CustomerDetails: details,
		}
		expectStatus(t, sendRequest(t, router, http.MethodPut, superadminPath, body), http.StatusForbidden, forbiddenMessage)
	})

	t.Run("delete superadmin", func(t *testing.T) {
		expectStatus(t, sendRequest(t, router, http.MethodDelete, superadminPath, nil), http.StatusForbidden, forbiddenMessage)
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
