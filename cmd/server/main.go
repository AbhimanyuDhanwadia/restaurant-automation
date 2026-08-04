package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/restaurantautomation/api/internal/alerts"
	"github.com/restaurantautomation/api/internal/analytics"
	"github.com/restaurantautomation/api/internal/api"
	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/backups"
	"github.com/restaurantautomation/api/internal/config"
	"github.com/restaurantautomation/api/internal/database"
	"github.com/restaurantautomation/api/internal/integrations"
	"github.com/restaurantautomation/api/internal/intelligence"
	"github.com/restaurantautomation/api/internal/inventory"
	"github.com/restaurantautomation/api/internal/logger"
	"github.com/restaurantautomation/api/internal/orders"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/printqueue"
	"github.com/restaurantautomation/api/internal/roles"
	"github.com/restaurantautomation/api/internal/settings"
	"github.com/restaurantautomation/api/internal/staff"
	"github.com/restaurantautomation/api/internal/tables"
	"github.com/restaurantautomation/api/internal/users"
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

	// 3. Open PostgreSQL and apply versioned migrations when configured.
	var eventPersister *database.EventPersister
	var eventStore *database.EventStore
	var databaseInspector *database.Inspector
	var databasePool *pgxpool.Pool
	orderRepository := orders.Repository(orders.NewMemoryRepository())
	tableRepository := tables.Repository(tables.NewMemoryRepository())
	inventoryRepository := inventory.Repository(inventory.NewMemoryRepository())
	staffRepository := staff.Repository(staff.NewMemoryRepository())
	alertRepository := alerts.Repository(alerts.NewMemoryRepository())
	settingsRepository := settings.Repository(settings.NewMemoryRepository())
	printQueueRepository := printqueue.Repository(printqueue.NewMemoryRepository())
	roleRepository := roles.Repository(roles.NewMemoryRepository())
	backupRepository := backups.Repository(backups.NewMemoryRepository())
	userRepository := users.Repository(users.NewMemoryRepository())
	if cfg.Database.DSN != "" {
		startupCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
		pool, err := database.Open(startupCtx, cfg.Database)
		if err == nil {
			err = database.RunMigrations(startupCtx, pool, cfg.Database.MigrationsDir)
		}
		cancel()
		if err != nil {
			log.Fatal().Err(err).Msg("database startup failed")
		}
		defer pool.Close()
		databasePool = pool
		orderRepository = orders.NewPostgresRepository(pool)
		tableRepository = tables.NewPostgresRepository(pool)
		inventoryRepository = inventory.NewPostgresRepository(pool)
		staffRepository = staff.NewPostgresRepository(pool)
		alertRepository = alerts.NewPostgresRepository(pool)
		settingsRepository = settings.NewPostgresRepository(pool)
		printQueueRepository = printqueue.NewPostgresRepository(pool)
		roleRepository = roles.NewPostgresRepository(pool)
		backupRepository = backups.NewPostgresRepository(pool)
		userRepository = users.NewPostgresRepository(pool)
		eventStore = database.NewEventStore(pool)
		databaseInspector = database.NewInspector(pool)
		eventPersister = database.NewEventPersister(eventStore, log)
		log.Info().Msg("database persistence enabled")
	} else {
		log.Warn().Msg("database persistence disabled: DATABASE_URL is not configured")
	}

	// 4. Initialize the automation engine and service boundaries.
	engine := automation.NewEngine(2, 100, automation.RetryPolicy{MaxAttempts: 3})
	engine.Start(context.Background())
	defer engine.Close()
	if eventPersister != nil {
		eventPersister.Start(context.Background(), engine)
		defer eventPersister.Close()
	}
	providers := integrations.NewRegistry()
	providers.Register(integrations.NewMockProvider("mock", 100))
	if cfg.Integrations.WebhookProviderName != "" {
		providers.Register(integrations.NewWebhookProvider(cfg.Integrations.WebhookProviderName, cfg.Integrations.WebhookSecret, 100))
		log.Info().Str("provider", cfg.Integrations.WebhookProviderName).Msg("signed webhook provider enabled")
	}
	if cfg.Integrations.SwiggyWebhookSecret != "" {
		providers.Register(integrations.NewWebhookProvider("swiggy", cfg.Integrations.SwiggyWebhookSecret, 100))
		log.Info().Msg("Swiggy signed webhook provider enabled")
	}
	if err := providers.StartCollectors(context.Background(), func(ctx context.Context, order integrations.Order) error { return engine.SubmitOrder(ctx, order.ID) }); err != nil {
		log.Fatal().Err(err).Msg("start provider collectors")
	}
	defer providers.DisconnectAll()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 3)
	var kitchenDriver printers.Driver = printers.NewMockDriver("kitchen")
	if cfg.Printers.KitchenAddress != "" {
		kitchenDriver = printers.NewTCPDriver("kitchen", cfg.Printers.KitchenAddress, cfg.Printers.ConnectTimeout)
		log.Info().Str("address", cfg.Printers.KitchenAddress).Msg("using TCP kitchen printer")
	}
	printerManager.Register(kitchenDriver, "kitchen")
	printerManager.Start(context.Background())
	defer printerManager.Close()
	orderService := orders.NewService(orderRepository)
	analyticsService := analytics.NewService(engine, printerManager, orderService)
	intelligenceService := intelligence.NewService(engine, printerManager)
	tableService := tables.NewService(tableRepository)
	inventoryService := inventory.NewService(inventoryRepository)
	staffService := staff.NewService(staffRepository)
	alertService := alerts.NewService(alertRepository)
	settingsService := settings.NewService(settingsRepository)
	printQueueService := printqueue.NewService(printQueueRepository)
	roleService := roles.NewService(roleRepository)
	backupService := backups.NewService(backupRepository)
	userService := users.NewService(userRepository, roleService, cfg.Auth.AdminEmails)
	printerManager.SetObserver(printQueueService)
	intelligenceService.Start(context.Background())
	defer intelligenceService.Close()
	router := api.NewRouter(cfg, log, engine, providers, printerManager, analyticsService, intelligenceService, orderService, tableService, inventoryService, staffService, alertService, settingsService, api.Dependencies{Readiness: databasePool, AuditLogReader: eventStore, Database: databaseInspector, PrintQueue: printQueueService, Roles: roleService, Backups: backupService, Users: userService})

	// 5. Setup HTTP Server
	server := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.Port),
		Handler: router,
	}

	// 6. Start Server with Graceful Shutdown
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
