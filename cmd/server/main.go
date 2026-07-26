package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/restaurantautomation/api/internal/api"
	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/config"
	"github.com/restaurantautomation/api/internal/integrations"
	"github.com/restaurantautomation/api/internal/logger"
)

func main() {
	// 1. Load Configuration
	cfg, err := config.Load()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load config: %v\n", err)
		os.Exit(1)
	}

	// 2. Initialize Logger
	log := logger.New(cfg.Log.Level, cfg.Log.Pretty)
	log.Info().Msg("starting restaurant automation api")

	// 3. Initialize the Phase 3 engine and Phase 4 provider registry.
	engine := automation.NewEngine(2, 100, automation.RetryPolicy{MaxAttempts: 3})
	engine.Start(context.Background())
	defer engine.Close()
	providers := integrations.NewRegistry()
	providers.Register(integrations.NewMockProvider("mock", 100))
	router := api.NewRouter(cfg, log, engine, providers)

	// 4. Setup HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// 5. Start Server with Graceful Shutdown
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sig
		log.Info().Msg("shutdown signal received")

		// Shutdown signal with grace period
		shutdownCtx, cancel := context.WithTimeout(serverCtx, cfg.Server.ShutdownTimeout)
		defer cancel()

		go func() {
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal().Msg("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		err := server.Shutdown(shutdownCtx)
		if err != nil {
			log.Fatal().Err(err).Msg("failed to gracefully shutdown server")
		}
		serverStopCtx()
	}()

	// Run the server
	log.Info().Int("port", cfg.Server.Port).Msg("server listening")
	err = server.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		log.Fatal().Err(err).Msg("server failed")
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()
	log.Info().Msg("server exited cleanly")
}
