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
	"github.com/jwald3/go_rest_template/internal/config"
	"github.com/jwald3/go_rest_template/internal/database"
	"github.com/jwald3/go_rest_template/internal/handler"
	"github.com/jwald3/go_rest_template/internal/logger"
	"github.com/jwald3/go_rest_template/internal/repository"
	"github.com/jwald3/go_rest_template/internal/service"
	"go.uber.org/zap"
)

func main() {
	cfg := config.Load()

	log := logger.New(cfg.App.LogLevel)
	defer log.Sync()

	db, err := database.NewPostgresConnection(*cfg)
	if err != nil {
		log.Fatal("failed to connect to database", zap.Error(err))
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db.DB)

	userService := service.NewUserService(userRepo)

	userHandler := handler.NewUserHandler(userService)

	router := mux.NewRouter()

	router.HandleFunc("/api/v1/users", userHandler.List).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users", userHandler.Create).Methods(http.MethodPost)
	router.HandleFunc("/api/v1/users/{id}", userHandler.Update).Methods(http.MethodPut)
	router.HandleFunc("/api/v1/users/{id}", userHandler.Get).Methods(http.MethodGet)
	router.HandleFunc("/api/v1/users/{id}", userHandler.Delete).Methods(http.MethodDelete)

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
