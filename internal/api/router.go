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
	"github.com/restaurantautomation/api/internal/config"
	"github.com/restaurantautomation/api/internal/integrations"
	"github.com/restaurantautomation/api/internal/intelligence"
	"github.com/restaurantautomation/api/internal/inventory"
	"github.com/restaurantautomation/api/internal/orders"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/settings"
	"github.com/restaurantautomation/api/internal/staff"
	"github.com/restaurantautomation/api/internal/tables"
)

// NewRouter constructs the Chi router with the standard middleware stack
// and registers all API routes.
func NewRouter(cfg *config.Config, log zerolog.Logger, engine *automation.Engine, registry *integrations.Registry, printerManager *printers.Manager, analyticsService *analytics.Service, intelligenceService *intelligence.Service, orderService *orders.Service, tableService *tables.Service, inventoryService *inventory.Service, staffService *staff.Service, alertService *alerts.Service, settingsService *settings.Service, readiness ...handlers.ReadinessChecker) *chi.Mux {
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
	var checker handlers.ReadinessChecker
	if len(readiness) > 0 {
		checker = readiness[0]
	}
	r.Get("/ready", handlers.ReadyWithCheck(checker))
	r.Post("/api/v1/webhooks/{provider}", handlers.ReceiveWebhook(registry))

	// 4. API Routes (Will be authenticated later)
	r.Route("/api/v1", func(r chi.Router) {
		r.Use(appm.Auth(cfg.Auth))
		r.Get("/me", handlers.CurrentUser())
		r.Get("/insights", handlers.Insights(intelligenceService))
		r.Get("/analytics/overview", handlers.AnalyticsOverview(analyticsService))
		r.Get("/integrations", handlers.Integrations(registry))
		r.Get("/printers", handlers.Printers(printerManager))
		r.Post("/printers/tickets", handlers.PrintTicket(printerManager))
		r.Post("/printers/tickets/{orderID}/reprint", handlers.ReprintTicket(printerManager))
		r.Route("/automation", func(r chi.Router) {
			r.Post("/orders", handlers.SubmitAutomationOrder(engine))
			r.Get("/events", handlers.AutomationEvents(engine))
			r.Get("/queue", handlers.AutomationQueue(engine))
		})
		r.Get("/orders", handlers.ListOrders(orderService))
		r.Post("/orders", handlers.CreateOrder(orderService, engine))
		r.Patch("/orders/{orderID}/status", handlers.UpdateOrderStatus(orderService))
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
		r.Get("/settings", handlers.GetSettings(settingsService))
		r.Put("/settings", handlers.SaveSettings(settingsService))
	})

	return r
}
