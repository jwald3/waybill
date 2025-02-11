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
	"github.com/jwald3/waybill/internal/middleware"
	"github.com/jwald3/waybill/internal/repository"
	"github.com/jwald3/waybill/internal/service"
	"go.uber.org/zap"
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

	handlers := initializeHandlers(db)

	router := mux.NewRouter()
	router.Use(middleware.Logging(log))
	router.Use(middleware.Recovery(log))
	router.Use(middleware.CORS)

	router.HandleFunc("/health", handler.HealthCheck).Methods(http.MethodGet)

	v1 := router.PathPrefix("/api/v1").Subrouter()
	registerDriverRoutes(v1, handlers.driver)
	registerFacilityRoutes(v1, handlers.facility)
	registerFuelLogRoutes(v1, handlers.fuelLog)
	registerIncidentReportRoutes(v1, handlers.incidentReport)
	registerMaintenanceLogRoutes(v1, handlers.maintenanceLog)
	registerTripRoutes(v1, handlers.trip)
	registerTruckRoutes(v1, handlers.truck)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	go func() {
		log.Info("starting server...",
			zap.String("host", cfg.Server.Host),
			zap.String("port", cfg.Server.Port),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to listen and serve", zap.Error(err))
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server gracefully stopped.")
}

type handlers struct {
	driver         *handler.DriverHandler
	facility       *handler.FacilityHandler
	fuelLog        *handler.FuelLogHandler
	incidentReport *handler.IncidentReportHandler
	maintenanceLog *handler.MaintenanceLogHandler
	trip           *handler.TripHandler
	truck          *handler.TruckHandler
}

func initializeHandlers(db *database.MongoDB) *handlers {
	// Initialize repositories
	driverRepo := repository.NewDriverRepository(db)
	facilityRepo := repository.NewFacilityRepository(db)
	fuelLogRepo := repository.NewFuelLogRepository(db)
	incidentReportRepo := repository.NewIncidentReportRepository(db)
	maintenanceLogRepo := repository.NewMaintenanceLogRepository(db)
	tripRepo := repository.NewTripRepository(db)
	truckRepo := repository.NewTruckRepository(db)

	// Initialize services
	driverService := service.NewDriverService(db, driverRepo)
	facilityService := service.NewFacilityService(db, facilityRepo)
	fuelLogService := service.NewFuelLogService(db, fuelLogRepo)
	incidentReportService := service.NewIncidentReportService(db, incidentReportRepo)
	maintenanceLogService := service.NewMaintenanceLogService(db, maintenanceLogRepo)
	tripService := service.NewTripService(db, tripRepo)
	truckService := service.NewTruckService(db, truckRepo)

	// Initialize handlers
	return &handlers{
		driver:         handler.NewDriverHandler(driverService),
		facility:       handler.NewFacilityHandler(facilityService),
		fuelLog:        handler.NewFuelLogHandler(fuelLogService),
		incidentReport: handler.NewIncidentReportHandler(incidentReportService),
		maintenanceLog: handler.NewMaintenanceLogHandler(maintenanceLogService),
		trip:           handler.NewTripHandler(tripService),
		truck:          handler.NewTruckHandler(truckService),
	}
}

func registerDriverRoutes(r *mux.Router, h *handler.DriverHandler) {
	r.HandleFunc("/drivers", h.List).Methods(http.MethodGet)
	r.HandleFunc("/drivers", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/drivers/{id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/drivers/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/drivers/{id}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/drivers/{id}/employment-status", h.UpdateEmploymentStatus).Methods(http.MethodPatch)
}

func registerFacilityRoutes(r *mux.Router, h *handler.FacilityHandler) {
	r.HandleFunc("/facilities", h.List).Methods(http.MethodGet)
	r.HandleFunc("/facilities", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/facilities/{id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/facilities/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/facilities/{id}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/facilities/{id}/services", h.UpdateAvailableFacilityServices).Methods(http.MethodPatch)
}

func registerFuelLogRoutes(r *mux.Router, h *handler.FuelLogHandler) {
	r.HandleFunc("/fuel-logs", h.List).Methods(http.MethodGet)
	r.HandleFunc("/fuel-logs", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/fuel-logs/{id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/fuel-logs/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/fuel-logs/{id}", h.Delete).Methods(http.MethodDelete)
}

func registerIncidentReportRoutes(r *mux.Router, h *handler.IncidentReportHandler) {
	r.HandleFunc("/incident-reports", h.List).Methods(http.MethodGet)
	r.HandleFunc("/incident-reports", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/incident-reports/{id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/incident-reports/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/incident-reports/{id}", h.Delete).Methods(http.MethodDelete)
}

func registerMaintenanceLogRoutes(r *mux.Router, h *handler.MaintenanceLogHandler) {
	r.HandleFunc("/maintenance-logs", h.List).Methods(http.MethodGet)
	r.HandleFunc("/maintenance-logs", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/maintenance-logs/{id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/maintenance-logs/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/maintenance-logs/{id}", h.Delete).Methods(http.MethodDelete)
}

func registerTripRoutes(r *mux.Router, h *handler.TripHandler) {
	r.HandleFunc("/trips", h.List).Methods(http.MethodGet)
	r.HandleFunc("/trips", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/trips/{id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/trips/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/trips/{id}", h.Delete).Methods(http.MethodDelete)
}

func registerTruckRoutes(r *mux.Router, h *handler.TruckHandler) {
	r.HandleFunc("/trucks", h.List).Methods(http.MethodGet)
	r.HandleFunc("/trucks", h.Create).Methods(http.MethodPost)
	r.HandleFunc("/trucks/{id}", h.GetById).Methods(http.MethodGet)
	r.HandleFunc("/trucks/{id}", h.Update).Methods(http.MethodPut)
	r.HandleFunc("/trucks/{id}", h.Delete).Methods(http.MethodDelete)
	r.HandleFunc("/trucks/{id}/status", h.UpdateTruckStatus).Methods(http.MethodPatch)
	r.HandleFunc("/trucks/{id}/mileage", h.UpdateTruckMileage).Methods(http.MethodPatch)
	r.HandleFunc("/trucks/{id}/maintenance", h.UpdateTruckLastMaintenance).Methods(http.MethodPatch)
}
