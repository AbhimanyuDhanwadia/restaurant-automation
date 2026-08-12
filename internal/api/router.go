package api

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/rs/zerolog"

	"github.com/restaurantautomation/api/internal/alerts"
	"github.com/restaurantautomation/api/internal/analytics"
	"github.com/restaurantautomation/api/internal/api/handlers"
	appm "github.com/restaurantautomation/api/internal/api/middleware"
	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/backups"
	"github.com/restaurantautomation/api/internal/config"
	"github.com/restaurantautomation/api/internal/integrations"
	"github.com/restaurantautomation/api/internal/intelligence"
	"github.com/restaurantautomation/api/internal/inventory"
	"github.com/restaurantautomation/api/internal/orders"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/printqueue"
	"github.com/restaurantautomation/api/internal/roles"
	"github.com/restaurantautomation/api/internal/settings"
	"github.com/restaurantautomation/api/internal/staff"
	"github.com/restaurantautomation/api/internal/tables"
	"github.com/restaurantautomation/api/internal/users"
)

// NewRouter constructs the Chi router with the standard middleware stack
// and registers all API routes.
type Dependencies struct {
	Readiness      handlers.ReadinessChecker
	AuditLogReader handlers.OperationalEventReader
	Database       handlers.DatabaseInspector
	PrintQueue     *printqueue.Service
	Roles          *roles.Service
	Backups        *backups.Service
	Users          *users.Service
}

func NewRouter(cfg *config.Config, log zerolog.Logger, engine *automation.Engine, registry *integrations.Registry, printerManager *printers.Manager, analyticsService *analytics.Service, intelligenceService *intelligence.Service, orderService *orders.Service, tableService *tables.Service, inventoryService *inventory.Service, staffService *staff.Service, alertService *alerts.Service, settingsService *settings.Service, dependencies ...Dependencies) *chi.Mux {
	r := chi.NewRouter()

	// 1. Basic Middleware
	r.Use(middleware.RealIP)
	r.Use(appm.RequestID)
	r.Use(appm.Logger(log))
	r.Use(middleware.Recoverer)

	// 2. CORS
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.Server.CORSOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	// 3. System Routes (Unauthenticated)
	r.Get("/health", handlers.Health())
	var dependenciesConfig Dependencies
	if len(dependencies) > 0 {
		dependenciesConfig = dependencies[0]
	}
	checker := dependenciesConfig.Readiness
	r.Get("/ready", handlers.ReadyWithCheck(checker))
	r.Post("/api/v1/webhooks/{provider}", handlers.ReceiveWebhook(registry))

	// 4. API Routes (Will be authenticated later)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(appm.Auth(cfg.Auth))
		r.Use(appm.TrackUser(dependenciesConfig.Users))
		r.Get("/me", handlers.CurrentUser(dependenciesConfig.Users))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "automation.view")).Get("/system/health", handlers.SystemHealth(engine, registry, printerManager, checker))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "audit.view")).Get("/admin/audit-logs", handlers.AuditLogs(engine, dependenciesConfig.AuditLogReader))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "audit.view")).Get("/admin/access-audit", handlers.ListRoleAssignments(dependenciesConfig.Users))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "database.view")).Get("/admin/database", handlers.DatabaseStatus(dependenciesConfig.Database))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "roles.manage")).Get("/admin/roles", handlers.ListRoles(dependenciesConfig.Roles))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "roles.manage")).Post("/admin/roles", handlers.CreateRole(dependenciesConfig.Roles))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "backups.manage")).Get("/admin/backups", handlers.ListBackups(dependenciesConfig.Backups))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "backups.manage")).Post("/admin/backups", handlers.RecordBackup(dependenciesConfig.Backups))
		r.Route("/admin/users", func(r chi.Router) {
			r.Use(appm.RequirePermission(dependenciesConfig.Users, "users.manage"))
			r.Get("/", handlers.ListUsers(dependenciesConfig.Users))
			r.Patch("/{subject}/role", handlers.UpdateUserRole(dependenciesConfig.Users))
		})
		r.With(appm.RequirePermission(dependenciesConfig.Users, "automation.view")).Get("/insights", handlers.Insights(intelligenceService))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "analytics.view")).Get("/analytics/overview", handlers.AnalyticsOverview(analyticsService))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "integrations.manage")).Get("/integrations", handlers.Integrations(registry))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "printers.manage")).Get("/printers", handlers.Printers(printerManager))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "printers.manage")).Get("/printers/queue", handlers.PrintQueue(dependenciesConfig.PrintQueue))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "printers.manage")).Post("/printers/queue/{jobID}/requeue", handlers.RequeuePrintJob(printerManager, dependenciesConfig.PrintQueue))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "printers.manage")).Post("/printers/queue/{jobID}/retry", handlers.RetryPrintJob(printerManager, dependenciesConfig.PrintQueue))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "printers.manage")).Post("/printers/tickets", handlers.PrintTicket(printerManager, dependenciesConfig.PrintQueue))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "printers.manage")).Post("/printers/tickets/{orderID}/reprint", handlers.ReprintTicket(printerManager, dependenciesConfig.PrintQueue))
		r.Route("/automation", func(r chi.Router) {
			r.Use(appm.RequirePermission(dependenciesConfig.Users, "automation.view"))
			r.Post("/orders", handlers.SubmitAutomationOrder(engine))
			r.Get("/events", handlers.AutomationEvents(engine))
			r.Get("/queue", handlers.AutomationQueue(engine))
		})
		r.Get("/orders", handlers.ListOrders(orderService))
		r.Post("/orders", handlers.CreateOrder(orderService, engine, printerManager, dependenciesConfig.PrintQueue))
		r.Patch("/orders/{orderID}/status", handlers.UpdateOrderStatus(orderService, engine))
		r.Get("/tables", handlers.ListTables(tableService))
		r.Post("/tables", handlers.CreateTable(tableService))
		r.Patch("/tables/{tableID}/status", handlers.UpdateTableStatus(tableService))
		r.Get("/inventory", handlers.ListInventory(inventoryService))
		r.Post("/inventory", handlers.CreateInventoryItem(inventoryService))
		r.Patch("/inventory/{itemID}/stock", handlers.UpdateInventoryStock(inventoryService))
		r.Patch("/inventory/{itemID}/status", handlers.UpdateInventoryStatus(inventoryService))
		r.Get("/staff", handlers.ListStaff(staffService))
		r.Post("/staff", handlers.CreateStaffMember(staffService))
		r.Patch("/staff/{staffID}/status", handlers.UpdateStaffStatus(staffService))
		r.Patch("/staff/{staffID}/handoff", handlers.UpdateStaffHandoff(staffService))
		r.Get("/staff/tasks", handlers.ListShiftTasks(staffService))
		r.Post("/staff/tasks", handlers.CreateShiftTask(staffService))
		r.Patch("/staff/tasks/{taskID}/complete", handlers.CompleteShiftTask(staffService))
		r.Get("/staff/handoff", handlers.GetShiftHandoff(staffService))
		r.Put("/staff/handoff", handlers.SaveShiftHandoff(staffService))
		r.Get("/alerts", handlers.ListAlerts(alertService))
		r.Post("/alerts", handlers.CreateAlert(alertService))
		r.Patch("/alerts/{alertID}/acknowledge", handlers.AcknowledgeAlert(alertService))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "settings.manage")).Get("/settings", handlers.GetSettings(settingsService))
		r.With(appm.RequirePermission(dependenciesConfig.Users, "settings.manage")).Put("/settings", handlers.SaveSettings(settingsService))
	})

	return r
}
