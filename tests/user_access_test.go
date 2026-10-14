package tests

import (
	"net/http"
	"testing"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

// TestUserListAccess checks GET /users: the superadmin and admins may list users,
// every other role gets 403.
func TestUserListAccess(t *testing.T) {
	userCountryService := service.NewUserCountryService(repository.NewUserCountryRepository(testDBPool))
	userStateService := service.NewUserStateService(repository.NewUserStateRepository(testDBPool), userCountryService)
	userCityService := service.NewUserCityService(repository.NewUserCityRepository(testDBPool), userStateService)
	branchService := service.NewBranchService(repository.NewBranchRepository(testDBPool), userCityService)
	userHandler := handler.NewUserHandler(service.NewUserService(testDBPool, branchService))

	for _, testCase := range []struct {
		role    string
		status  int
		message string
	}{
		{model.RoleSuperAdmin, http.StatusOK, "success"},
		{model.RoleAdmin, http.StatusOK, "success"},
		{model.RoleCustomer, http.StatusForbidden, "You do not have permission to do this"},
		{model.RoleOperator, http.StatusForbidden, "You do not have permission to do this"},
	} {
		t.Run(testCase.role, func(t *testing.T) {
			router := routes.New(routes.Handlers{User: userHandler}, actAs(testUserID, testCase.role))
			expectStatus(t, sendRequest(t, router, http.MethodGet, apiV1+"/users", nil), testCase.status, testCase.message)
		})
	}
}
