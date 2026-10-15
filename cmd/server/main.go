package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/db"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/logger"
	"7th-sense-server/internal/middleware"
	"7th-sense-server/internal/repository"
	"7th-sense-server/internal/routes"
	"7th-sense-server/internal/service"
)

func main() {
	logger.Init()
	if err := run(); err != nil {
		slog.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

// run starts the server and blocks until Ctrl+C / SIGTERM or a fatal error.
// It returns instead of exiting so every defer (like closing the database) runs.
func run() error {
	appConfig := config.Load()

	// ctx is cancelled on Ctrl+C / SIGTERM.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	dbPool, err := db.Connect(ctx, appConfig.DB)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer dbPool.Close()
	slog.Info("connected to database", "name", appConfig.DB.Name)

	categoryRepository := repository.NewCategoryRepository(dbPool)
	brandRepository := repository.NewBrandRepository(dbPool)
	productRepository := repository.NewProductRepository(dbPool)
	controllerRepository := repository.NewControllerRepository(dbPool)
	machineRepository := repository.NewMachineRepository(dbPool)
	// Organization: each service checks its parent with the service above it
	userCountryService := service.NewUserCountryService(repository.NewUserCountryRepository(dbPool))
	userStateService := service.NewUserStateService(repository.NewUserStateRepository(dbPool), userCountryService)
	userCityService := service.NewUserCityService(repository.NewUserCityRepository(dbPool), userStateService)
	branchService := service.NewBranchService(repository.NewBranchRepository(dbPool), userCityService)
	userLocationService := service.NewUserLocationService(repository.NewUserLocationRepository(dbPool), branchService)
	userShopService := service.NewUserShopService(repository.NewUserShopRepository(dbPool), userLocationService)
	userLineService := service.NewUserLineService(repository.NewUserLineRepository(dbPool), userShopService)
	userCellService := service.NewUserCellService(repository.NewUserCellRepository(dbPool), userLineService)
	userSupervisorService := service.NewUserSupervisorService(repository.NewUserSupervisorRepository(dbPool), userCellService)
	authService := service.NewAuthService(repository.NewAuthRepository(dbPool))
	brandService := service.NewBrandService(brandRepository)
	productService := service.NewProductService(productRepository)
	// Machine catalog: Brand → Controller → Model
	controllerService := service.NewControllerService(controllerRepository, brandService)
	machineModelService := service.NewMachineModelService(repository.NewMachineModelRepository(dbPool), controllerService)
	// A machine checks its customer, brand, controller, model and the customer's line
	machineService := service.NewMachineService(machineRepository, repository.NewUserRepository(dbPool),
		machineModelService, brandService, controllerService, userLineService)
	roleService := service.NewRoleService(repository.NewRoleRepository(dbPool))
	clienteleService := service.NewClienteleService(repository.NewClienteleRepository(dbPool), repository.NewUserRepository(dbPool))
	// A part checks its customer, clientele and applicable machines; drawings are files under UPLOAD_DIR
	partService := service.NewPartService(dbPool, repository.NewUserRepository(dbPool), machineService, clienteleService,
		repository.NewPartDrawingStorage(filepath.Join(appConfig.UploadDir, "part-drawings")))
	router := routes.New(routes.Handlers{
		Auth:         handler.NewAuthHandler(authService, appConfig.CookieSecure),
		User:         handler.NewUserHandler(service.NewUserService(dbPool, branchService, roleService, machineService)),
		Branch:       handler.NewBranchHandler(branchService),
		Category:     handler.NewCategoryHandler(service.NewCategoryService(categoryRepository)),
		Brand:        handler.NewBrandHandler(brandService),
		Product:      handler.NewProductHandler(productService),
		Controller:   handler.NewControllerHandler(controllerService),
		MachineModel: handler.NewMachineModelHandler(machineModelService),
		Machine:      handler.NewMachineHandler(machineService),
		Part:         handler.NewPartHandler(partService),
		Clientele:    handler.NewClienteleHandler(clienteleService),
		Shift:        handler.NewShiftHandler(service.NewShiftService(dbPool, repository.NewUserRepository(dbPool))),
		Break:        handler.NewBreakHandler(service.NewBreakService(dbPool, repository.NewUserRepository(dbPool))),
		Holiday: handler.NewHolidayHandler(service.NewHolidayService(repository.NewHolidayRepository(dbPool),
			repository.NewUserRepository(dbPool))),
		LicenceKey: handler.NewLicenceKeyHandler(service.NewLicenceKeyService(repository.NewLicenceKeyRepository(dbPool),
			repository.NewUserRepository(dbPool))),
		UserCountry:    handler.NewUserCountryHandler(userCountryService),
		UserState:      handler.NewUserStateHandler(userStateService),
		UserCity:       handler.NewUserCityHandler(userCityService),
		UserLocation:   handler.NewUserLocationHandler(userLocationService),
		UserShop:       handler.NewUserShopHandler(userShopService),
		UserLine:       handler.NewUserLineHandler(userLineService),
		UserCell:       handler.NewUserCellHandler(userCellService),
		UserSupervisor: handler.NewUserSupervisorHandler(userSupervisorService),
		Role:           handler.NewRoleHandler(roleService),
	}, middleware.Authenticate(authService))

	server := &http.Server{
		Addr:         ":" + appConfig.Port,
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("server started", "port", appConfig.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// Wait for a startup failure (e.g. port in use) or a shutdown signal.
	select {
	case err := <-serverErrors:
		return fmt.Errorf("listen: %w", err)
	case <-ctx.Done():
	}

	slog.Info("shutting down server")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("shutdown: %w", err)
	}
	slog.Info("server stopped")
	return nil
}
