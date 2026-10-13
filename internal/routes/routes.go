package routes

import (
	"net/http"

	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/response"
)

const apiPrefix = "/api"

// Handlers groups every handler the router needs.
type Handlers struct {
	User       *handler.UserHandler
	Branch     *handler.BranchHandler
	Category   *handler.CategoryHandler
	Brand      *handler.BrandHandler
	Product    *handler.ProductHandler
	Controller *handler.ControllerHandler
	Machine    *handler.MachineHandler
	// Organization (per user): country → state → city → … → supervisor
	UserCountry    *handler.UserCountryHandler
	UserState      *handler.UserStateHandler
	UserCity       *handler.UserCityHandler
	UserBranch     *handler.UserBranchHandler
	UserLocation   *handler.UserLocationHandler
	UserShop       *handler.UserShopHandler
	UserLine       *handler.UserLineHandler
	UserCell       *handler.UserCellHandler
	UserSupervisor *handler.UserSupervisorHandler
}

// New builds the main router with all routes and global middleware.
func New(handlers Handlers) http.Handler {
	router := http.NewServeMux()

	router.HandleFunc("GET /health", func(writer http.ResponseWriter, request *http.Request) {
		response.SendResponse(writer, request, response.Success, nil)
	})

	// Register API versions
	registerV1Routes(router, handlers)

	return middleware.Chain(router, middleware.Logger, middleware.Recoverer, middleware.Authenticate)
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
	if handlers.UserBranch != nil {
		registerUserBranchRoutes(router, v1Prefix, handlers.UserBranch)
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
}
