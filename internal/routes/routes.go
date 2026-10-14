package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/model"
	"7th-sense-server/internal/response"
)

const apiPrefix = "/api"

// Handlers groups every handler the router needs.
type Handlers struct {
	Auth       *handler.AuthHandler
	User       *handler.UserHandler
	Branch     *handler.BranchHandler
	Category   *handler.CategoryHandler
	Brand      *handler.BrandHandler
	Product    *handler.ProductHandler
	Controller *handler.ControllerHandler
	// MachineModel serves /models (machine models of a controller)
	MachineModel *handler.MachineModelHandler
	Machine      *handler.MachineHandler
	// Organization (per user): country → state → city → … → supervisor
	UserCountry    *handler.UserCountryHandler
	UserState      *handler.UserStateHandler
	UserCity       *handler.UserCityHandler
	UserLocation   *handler.UserLocationHandler
	UserShop       *handler.UserShopHandler
	UserLine       *handler.UserLineHandler
	UserCell       *handler.UserCellHandler
	UserSupervisor *handler.UserSupervisorHandler
	Role           *handler.RoleHandler
}

// New builds the main router with all routes and global middleware.
// authenticate resolves the logged-in user (middleware.Authenticate in the
// server). Only /health, login and logout work without a session; every other
// /api route answers 401 until the user logs in.
func New(handlers Handlers, authenticate middleware.Middleware) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(writer http.ResponseWriter, request *http.Request) {
		response.SendResponse(writer, request, response.Success, nil)
	})

	if handlers.Auth != nil {
		registerAuthRoutes(router, apiPrefix+"/v1", handlers.Auth, authenticate)
	}

	// Register API versions behind authentication
	protectedRouter := http.NewServeMux()
	registerV1Routes(protectedRouter, handlers)
	router.Handle(apiPrefix+"/", authenticate(protectedRouter))

	return middleware.Chain(router, middleware.Logger, middleware.Recoverer, middleware.CrossOrigin)
}

// userManagersOnly lets only the roles that manage users (superadmin, customer,
// admin) call handlerFunc; other users get 403. The service limits each of them
// to the users and roles they may manage.
func userManagersOnly(handlerFunc http.HandlerFunc) http.Handler {
	return middleware.RequireRole(model.UserManagerRoles...)(handlerFunc)
}

// superAdminOnly lets only the superadmin call handlerFunc; everyone else (admins
// included) gets 403.
func superAdminOnly(handlerFunc http.HandlerFunc) http.Handler {
	return middleware.RequireRole(model.RoleSuperAdmin)(handlerFunc)
}

// registerV1Routes registers all v1 API routes.
func registerV1Routes(router *http.ServeMux, handlers Handlers) {
	v1Prefix := apiPrefix + "/v1"

	if handlers.User != nil {
		registerUserRoutes(router, v1Prefix, handlers.User)
	}
	if handlers.Branch != nil {
		registerBranchRoutes(router, v1Prefix, handlers.Branch)
	}
	if handlers.Category != nil {
		registerCategoryRoutes(router, v1Prefix, handlers.Category)
	}
	if handlers.Brand != nil {
		registerBrandRoutes(router, v1Prefix, handlers.Brand)
	}
	if handlers.Product != nil {
		registerProductRoutes(router, v1Prefix, handlers.Product)
	}
	if handlers.Controller != nil {
		registerControllerRoutes(router, v1Prefix, handlers.Controller)
	}
	if handlers.MachineModel != nil {
		registerMachineModelRoutes(router, v1Prefix, handlers.MachineModel)
	}
	if handlers.Machine != nil {
		registerMachineRoutes(router, v1Prefix, handlers.Machine)
	}
	if handlers.UserCountry != nil {
		registerUserCountryRoutes(router, v1Prefix, handlers.UserCountry)
	}
	if handlers.UserState != nil {
		registerUserStateRoutes(router, v1Prefix, handlers.UserState)
	}
	if handlers.UserCity != nil {
		registerUserCityRoutes(router, v1Prefix, handlers.UserCity)
	}
	if handlers.UserLocation != nil {
		registerUserLocationRoutes(router, v1Prefix, handlers.UserLocation)
	}
	if handlers.UserShop != nil {
		registerUserShopRoutes(router, v1Prefix, handlers.UserShop)
	}
	if handlers.UserLine != nil {
		registerUserLineRoutes(router, v1Prefix, handlers.UserLine)
	}
	if handlers.UserCell != nil {
		registerUserCellRoutes(router, v1Prefix, handlers.UserCell)
	}
	if handlers.UserSupervisor != nil {
		registerUserSupervisorRoutes(router, v1Prefix, handlers.UserSupervisor)
	}
	if handlers.Role != nil {
		registerRoleRoutes(router, v1Prefix, handlers.Role)
	}
}
