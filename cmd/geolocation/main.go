// Service for FindHotel Coding Challenge.
//
// # Description of the REST API of the service for working with location data.
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Schemes: http, https
// Host: localhost:8087
// Version: 1.0.0
//
// swagger:meta
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/lib/pq"

	"vio/internal/api"
	"vio/internal/config"
	"vio/internal/database"
	"vio/internal/lib/logger/sl"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.MustLoad("geolocation")
	log := sl.SetupLogger(cfg.Env)
	log.Info(
		"starting geolocation server",
		slog.String("env", cfg.Env),
		slog.String("version", cfg.Version),
	)

	if err := run(log, cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %s\n", err)
		os.Exit(2)
	}
}

func run(log *slog.Logger, cfg *config.Config) error {
	log.Debug("starting geolocation server")

	log.Debug("starting db connect ", "connect", cfg.DBConnect)
	db, err := database.GetDB(cfg.DBConnect)
	if err != nil {
		return err
	}
	defer db.Close()

	log.Debug("db connected successfully")

	router := mux.NewRouter()

	router.HandleFunc("/api/geolocation/{ip_address}", api.GetGeoLocation(db)).Methods("GET")

	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	cfgAddress := fmt.Sprintf(":%d", cfg.Port)

	srv := &http.Server{
		Addr:    cfgAddress,
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			if !errors.Is(err, http.ErrServerClosed) {
				log.Error("failed to start server", "error", err)
			}
		}
	}()

	log.Info("started geolocation server", "port", cfgAddress)

	<-done
	log.Info("stopping geolocation server")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("failed to stop server", "error", err)

		return err
	}

	log.Info("geolocation server stopped")

	return nil
}
