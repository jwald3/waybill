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

var (
	driversBase   = "/api/v1/drivers"
	driversWithId = "/api/v1/drivers/{id}"

	trucksBase   = "/api/v1/trucks"
	trucksWithId = "/api/v1/trucks/{id}"

	facilitiesBase   = "/api/v1/facilities"
	facilitiesWithId = "/api/v1/facilities/{id}"

	tripsBase   = "/api/v1/trips"
	tripsWithId = "/api/v1/trips/{id}"

	fuelLogsBase   = "/api/v1/fuel-logs"
	fuelLogsWithId = "/api/v1/fuel-logs/{id}"

	incidentReportsBase   = "/api/v1/incident-reports"
	incidentReportsWithId = "/api/v1/incident-reports/{id}"

	maintenanceLogsBase   = "/api/v1/maintenance-logs"
	maintenanceLogsWithId = "/api/v1/maintenance-logs/{id}"
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

	driverRepo := repository.NewDriverRepository(db)
	driverService := service.NewDriverService(db, driverRepo)
	driverHandler := handler.NewDriverHandler(driverService)

	truckRepo := repository.NewTruckRepository(db)
	truckService := service.NewTruckService(db, truckRepo)
	truckHandler := handler.NewTruckHandler(truckService)

	facilityRepo := repository.NewFacilityRepository(db)
	facilityService := service.NewFacilityService(db, facilityRepo)
	facilityHandler := handler.NewFacilityHandler(facilityService)

	tripRepo := repository.NewTripRepository(db)
	tripService := service.NewTripService(db, tripRepo)
	tripHandler := handler.NewTripHandler(tripService)

	fuelLogRepo := repository.NewFuelLogRepository(db)
	fuelLogService := service.NewFuelLogService(db, fuelLogRepo)
	fuelLogHandler := handler.NewFuelLogHandler(fuelLogService)

	incidentReportRepo := repository.NewIncidentReportRepository(db)
	incidentReportService := service.NewIncidentReportService(db, incidentReportRepo)
	incidentReportHandler := handler.NewIncidentReportHandler(incidentReportService)

	maintenanceLogRepo := repository.NewMaintenanceLogRepository(db)
	maintenanceLogService := service.NewMaintenanceLogService(db, maintenanceLogRepo)
	maintenanceLogHandler := handler.NewMaintenanceLogHandler(maintenanceLogService)

	router := mux.NewRouter()

	router.HandleFunc(driversBase, driverHandler.List).Methods(http.MethodGet)
	router.HandleFunc(driversBase, driverHandler.Create).Methods(http.MethodPost)
	router.HandleFunc(driversWithId, driverHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc(driversWithId, driverHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(driversWithId, driverHandler.Delete).Methods(http.MethodDelete)

	router.HandleFunc(trucksBase, truckHandler.List).Methods(http.MethodGet)
	router.HandleFunc(trucksBase, truckHandler.Create).Methods(http.MethodPost)
	router.HandleFunc(trucksWithId, truckHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc(trucksWithId, truckHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(trucksWithId, truckHandler.Delete).Methods(http.MethodDelete)

	router.HandleFunc(facilitiesBase, facilityHandler.List).Methods(http.MethodGet)
	router.HandleFunc(facilitiesBase, facilityHandler.Create).Methods(http.MethodPost)
	router.HandleFunc(facilitiesWithId, facilityHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc(facilitiesWithId, facilityHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(facilitiesWithId, facilityHandler.Delete).Methods(http.MethodDelete)

	router.HandleFunc(tripsBase, tripHandler.List).Methods(http.MethodGet)
	router.HandleFunc(tripsBase, tripHandler.Create).Methods(http.MethodPost)
	router.HandleFunc(tripsWithId, tripHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc(tripsWithId, tripHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(tripsWithId, tripHandler.Delete).Methods(http.MethodDelete)

	router.HandleFunc(fuelLogsBase, fuelLogHandler.List).Methods(http.MethodGet)
	router.HandleFunc(fuelLogsBase, fuelLogHandler.Create).Methods(http.MethodPost)
	router.HandleFunc(fuelLogsWithId, fuelLogHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc(fuelLogsWithId, fuelLogHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(fuelLogsWithId, fuelLogHandler.Delete).Methods(http.MethodDelete)

	router.HandleFunc(incidentReportsBase, incidentReportHandler.List).Methods(http.MethodGet)
	router.HandleFunc(incidentReportsBase, incidentReportHandler.Create).Methods(http.MethodPost)
	router.HandleFunc(incidentReportsWithId, incidentReportHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc(incidentReportsWithId, incidentReportHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(incidentReportsWithId, incidentReportHandler.Delete).Methods(http.MethodDelete)

	router.HandleFunc(maintenanceLogsBase, maintenanceLogHandler.List).Methods(http.MethodGet)
	router.HandleFunc(maintenanceLogsBase, maintenanceLogHandler.Create).Methods(http.MethodPost)
	router.HandleFunc(maintenanceLogsWithId, maintenanceLogHandler.GetById).Methods(http.MethodGet)
	router.HandleFunc(maintenanceLogsWithId, maintenanceLogHandler.Update).Methods(http.MethodPut)
	router.HandleFunc(maintenanceLogsWithId, maintenanceLogHandler.Delete).Methods(http.MethodDelete)

	server := &http.Server{
		Addr:         fmt.Sprintf("%s:%s", cfg.Server.Host, cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
		IdleTimeout:  cfg.Server.IdleTimeout,
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Info("starting server...",
			zap.String("host", cfg.Server.Host),
			zap.String("port", cfg.Server.Port),
		)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("failed to listen and serve", zap.Error(err))
		}
	}()

	<-stopChan

	log.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatal("server forced to shutdown", zap.Error(err))
	}

	log.Info("server gracefully stopped.")
}
