// Command server wires config -> repository -> service -> controller and serves HTTP.
package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/starman011/urbanism/backend/internal/config"
	ctrl "github.com/starman011/urbanism/backend/internal/controller/http"
	"github.com/starman011/urbanism/backend/internal/repository/postgres"
	"github.com/starman011/urbanism/backend/internal/service"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg, err := config.Load()
	if err != nil {
		log.Error("config", "err", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("database", "err", err)
		os.Exit(1)
	}
	defer db.Close()
	if err := db.Migrate(ctx); err != nil {
		log.Error("migrate", "err", err)
		os.Exit(1)
	}
	if err := db.Seed(ctx); err != nil {
		log.Error("seed", "err", err)
		os.Exit(1)
	}

	svc := service.New(db.Repositories(), time.Now)
	auth := &service.Auth{Repo: db.Accounts(), Secret: os.Getenv("AUTH_JWT_SECRET"), Bootstrap: os.Getenv("AUTH_BOOTSTRAP_TOKEN"), Now: time.Now}
	if !auth.Enabled() {
		log.Warn("Board sign-in disabled: configure AUTH_JWT_SECRET (at least 32 characters)")
	}
	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           ctrl.NewRouter(svc, log, cfg.AllowedOrigins, auth),
		ReadHeaderTimeout: 10 * time.Second,
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}

	go func() {
		log.Info("listening", "port", cfg.Port, "env", cfg.Env)
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("serve", "err", err)
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = srv.Shutdown(shutdownCtx)
	log.Info("stopped")
}
