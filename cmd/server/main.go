package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"7th-sense-server/internal/config"
	"7th-sense-server/internal/db"
	"7th-sense-server/internal/handler"
	"7th-sense-server/internal/logger"
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

	branchRepository := repository.NewBranchRepository(dbPool)
	branchService := service.NewBranchService(branchRepository)
	categoryRepository := repository.NewCategoryRepository(dbPool)
	brandRepository := repository.NewBrandRepository(dbPool)
	productRepository := repository.NewProductRepository(dbPool)
	controllerRepository := repository.NewControllerRepository(dbPool)
	machineRepository := repository.NewMachineRepository(dbPool)
	router := routes.New(routes.Handlers{
		User:       handler.NewUserHandler(service.NewUserService(dbPool, branchService)),
		Branch:     handler.NewBranchHandler(branchService),
		Category:   handler.NewCategoryHandler(service.NewCategoryService(categoryRepository)),
		Brand:      handler.NewBrandHandler(service.NewBrandService(brandRepository)),
		Product:    handler.NewProductHandler(service.NewProductService(productRepository)),
		Controller: handler.NewControllerHandler(service.NewControllerService(controllerRepository)),
		Machine:    handler.NewMachineHandler(service.NewMachineService(machineRepository)),
	})

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
