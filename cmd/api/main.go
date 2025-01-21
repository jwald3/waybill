package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/jwald3/waybill/internal/config"
	"github.com/jwald3/waybill/internal/database"
	"github.com/jwald3/waybill/internal/handler"
	"github.com/jwald3/waybill/internal/logger"
	"github.com/jwald3/waybill/internal/repository"
	"github.com/jwald3/waybill/internal/service"
	"go.uber.org/zap"
)

// feel free to save other routes as variables here. You don't necessarily need to do this, I just didn't like having the warning for multiple routes using the same string literal
var (
	driversBase   = "/api/v1/drivers"
	driversWithId = "/api/v1/drivers/{id}"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.App.LogLevel)
	defer log.Sync()

	db, err := database.NewMongoConnection(*cfg)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	// the bit of dependency injection where you register repos, services, and handlers.
	// If you're using additional or different resources, make sure you include them here as well
	driverRepo := repository.NewDriverRepository(db)
	driverService := service.NewDriverService(db, driverRepo)
	driverHandler := handler.NewDriverHandler(driverService)

	// the gorilla mux router - I went with this dependency to simplify the routing and make handling URL params less of a pain
	router := mux.NewRouter()

	// registering each handler function onto the router (using variables for the route to avoid warnings of overused string literals).
	// Register any additional routes below
	router.HandleFunc(driversBase, driverHandler.List).Methods(http.MethodGet)
	router.HandleFunc(driversBase, driverHandler.Create).Methods(http.MethodPost)
	router.HandleFunc(driversWithId, driverHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc(driversWithId, driverHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(driversWithId, driverHandler.Delete).Methods(http.MethodDelete)

	// start a server with the mux router we just armed with routes and the env variables as loaded into the config.go file
	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	stopChan := make(chan os.Signal, 1)
	// registers os.Interrupt (CTRL+C) and SIGTERM signals to the `stopChan` channel
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	// using a goroutine (labeled with the `go` keyword) ensures that we can run the server on a background thread and handle operations on the main thread
	go func() {
		log.Info("starting server...",
			zap.String("host", cfg.Server.Host),
			zap.String("port", cfg.Server.Port),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to listen and serve", zap.Error(err))
		}
	}()

	// when the server receives that signal and stopChan gets the notification it allows the program to progress to the shutdown logic
	<-stopChan

	log.Info("shutting down server...")

	// the application gives itself 5 seconds to finish using context resources before shutting down
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// attempt the shutdown logic, catching any errors where a forceful shutdown was necessary
	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server gracefully stopped.")
}
