package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/restaurantautomation/api/internal/alerts"
	"github.com/restaurantautomation/api/internal/analytics"
	"github.com/restaurantautomation/api/internal/api"
	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/config"
	"github.com/restaurantautomation/api/internal/integrations"
	"github.com/restaurantautomation/api/internal/intelligence"
	"github.com/restaurantautomation/api/internal/inventory"
	"github.com/restaurantautomation/api/internal/orders"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/roles"
	"github.com/restaurantautomation/api/internal/settings"
	"github.com/restaurantautomation/api/internal/staff"
	"github.com/restaurantautomation/api/internal/tables"
	"github.com/restaurantautomation/api/internal/users"
	"github.com/rs/zerolog"
)

func TestHealthEndpoint(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			CORSOrigins: []string{"*"},
		},
	}
	log := zerolog.Nop()

	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, log, engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	req, err := http.NewRequest("GET", "/health", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	if res["status"] != "ok" {
		t.Errorf("handler returned unexpected body: got %v want %v", res["status"], "ok")
	}
}

func TestReadyEndpoint(t *testing.T) {
	cfg := &config.Config{
		Server: config.ServerConfig{
			CORSOrigins: []string{"*"},
		},
	}
	log := zerolog.Nop()

	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, log, engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	req, err := http.NewRequest("GET", "/ready", nil)
	if err != nil {
		t.Fatal(err)
	}

	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var res map[string]string
	if err := json.NewDecoder(rr.Body).Decode(&res); err != nil {
		t.Fatal(err)
	}

	if res["status"] != "ready" {
		t.Errorf("handler returned unexpected body: got %v want %v", res["status"], "ready")
	}
}

func TestAutomationOrderEndpoint(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/automation/orders", bytes.NewBufferString(`{"order_id":"ORD-100"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusAccepted)
	}
}

func TestAutomationQueueEndpoint(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(2, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	request := httptest.NewRequest(http.MethodGet, "/api/v1/automation/queue", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var stats automation.Stats
	if err := json.NewDecoder(recorder.Body).Decode(&stats); err != nil {
		t.Fatal(err)
	}
	if stats.Workers != 2 || stats.QueueCapacity != 10 {
		t.Fatalf("stats = %#v, want two workers and capacity ten", stats)
	}
}

func TestAutomationEventsEndpoint(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	if err := engine.SubmitOrder(context.Background(), "ORD-101"); err != nil {
		t.Fatal(err)
	}
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	request := httptest.NewRequest(http.MethodGet, "/api/v1/automation/events", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var events []automation.Event
	if err := json.NewDecoder(recorder.Body).Decode(&events); err != nil {
		t.Fatal(err)
	}
	if len(events) == 0 || events[0].OrderID != "ORD-101" {
		t.Fatalf("events = %#v, want event for ORD-101", events)
	}
}

func TestIntegrationsEndpoint(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	registry := integrations.NewRegistry()
	provider := integrations.NewMockProvider("swiggy", 1)
	if err := registry.Register(provider); err != nil {
		t.Fatal(err)
	}
	if err := provider.Connect(); err != nil {
		t.Fatal(err)
	}
	router := api.NewRouter(cfg, zerolog.Nop(), engine, registry, printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	request := httptest.NewRequest(http.MethodGet, "/api/v1/integrations", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var providers []integrations.Health
	if err := json.NewDecoder(recorder.Body).Decode(&providers); err != nil {
		t.Fatal(err)
	}
	if len(providers) != 1 || providers[0].Name != "swiggy" || providers[0].Status != integrations.StatusConnected {
		t.Fatalf("providers = %#v, want connected swiggy", providers)
	}
}

func TestPrintersEndpoint(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := printerManager.Register(printers.NewMockDriver("kitchen"), "kitchen"); err != nil {
		t.Fatal(err)
	}
	if err := printerManager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer printerManager.Close()
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	request := httptest.NewRequest(http.MethodGet, "/api/v1/printers", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	var fleet []printers.PrinterHealth
	if err := json.NewDecoder(recorder.Body).Decode(&fleet); err != nil {
		t.Fatal(err)
	}
	if len(fleet) != 1 || fleet[0].Name != "kitchen" || fleet[0].Status != printers.StatusReady {
		t.Fatalf("fleet = %#v, want ready kitchen printer", fleet)
	}
}

func TestAutomationOrderEndpointValidatesOrderID(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/automation/orders", bytes.NewBufferString(`{}`))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}

func TestTableEndpoints(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	create := httptest.NewRequest(http.MethodPost, "/api/v1/tables", bytes.NewBufferString(`{"id":"T1","seats":4}`))
	create.Header.Set("Content-Type", "application/json")
	created := httptest.NewRecorder()
	router.ServeHTTP(created, create)
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", created.Code, http.StatusCreated)
	}

	list := httptest.NewRequest(http.MethodGet, "/api/v1/tables", nil)
	listed := httptest.NewRecorder()
	router.ServeHTTP(listed, list)
	if listed.Code != http.StatusOK {
		t.Fatalf("list status = %d, want %d", listed.Code, http.StatusOK)
	}
}

func TestInventoryEndpoints(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	create := httptest.NewRequest(http.MethodPost, "/api/v1/inventory", bytes.NewBufferString(`{"name":"Tomatoes","category":"Produce","on_hand":2,"par_level":4,"unit":"kg"}`))
	create.Header.Set("Content-Type", "application/json")
	created := httptest.NewRecorder()
	router.ServeHTTP(created, create)
	if created.Code != http.StatusCreated {
		t.Fatalf("create status = %d, want %d", created.Code, http.StatusCreated)
	}
}

func TestStaffEndpoints(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/staff", bytes.NewBufferString(`{"name":"Maya","role":"Manager","station":"Floor"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
}

func TestAlertEndpoints(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))
	request := httptest.NewRequest(http.MethodPost, "/api/v1/alerts", bytes.NewBufferString(`{"label":"Cold room","detail":"Temperature rising","severity":"high"}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}
}

func TestSettingsEndpoints(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()))

	update := httptest.NewRequest(http.MethodPut, "/api/v1/settings", bytes.NewBufferString(`{"restaurant_name":"The Green Table","timezone":"Asia/Kolkata","operational_alerts":true,"auto_advance_tickets":true}`))
	update.Header.Set("Content-Type", "application/json")
	updated := httptest.NewRecorder()
	router.ServeHTTP(updated, update)
	if updated.Code != http.StatusOK {
		t.Fatalf("update status = %d, want %d", updated.Code, http.StatusOK)
	}

	request := httptest.NewRequest(http.MethodGet, "/api/v1/settings", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusOK {
		t.Fatalf("get status = %d, want %d", recorder.Code, http.StatusOK)
	}
}

func TestUsersRouteRequiresAssignedPermission(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}, Auth: config.AuthConfig{Required: true, SupabaseJWTSecret: "secret", AdminEmails: []string{"manager@example.com"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	roleService := roles.NewService(roles.NewMemoryRepository())
	userService := users.NewService(users.NewMemoryRepository(), roleService, cfg.Auth.AdminEmails)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager), intelligence.NewService(engine, printerManager), orders.NewService(orders.NewMemoryRepository()), tables.NewService(tables.NewMemoryRepository()), inventory.NewService(inventory.NewMemoryRepository()), staff.NewService(staff.NewMemoryRepository()), alerts.NewService(alerts.NewMemoryRepository()), settings.NewService(settings.NewMemoryRepository()), api.Dependencies{Roles: roleService, Users: userService})

	operatorRequest := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	operatorRequest.Header.Set("Authorization", "Bearer "+signedAccessToken(t, "operator", "operator@example.com"))
	operatorResponse := httptest.NewRecorder()
	router.ServeHTTP(operatorResponse, operatorRequest)
	if operatorResponse.Code != http.StatusForbidden {
		t.Fatalf("operator status = %d, want %d", operatorResponse.Code, http.StatusForbidden)
	}

	adminRequest := httptest.NewRequest(http.MethodGet, "/api/v1/admin/users", nil)
	adminRequest.Header.Set("Authorization", "Bearer "+signedAccessToken(t, "manager", "manager@example.com"))
	adminResponse := httptest.NewRecorder()
	router.ServeHTTP(adminResponse, adminRequest)
	if adminResponse.Code != http.StatusOK {
		t.Fatalf("administrator status = %d, want %d body=%s", adminResponse.Code, http.StatusOK, adminResponse.Body.String())
	}
}

func signedAccessToken(t *testing.T, subject, email string) string {
	t.Helper()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"sub": subject, "email": email, "exp": time.Now().Add(time.Minute).Unix()})
	signed, err := token.SignedString([]byte("secret"))
	if err != nil {
		t.Fatal(err)
	}
	return signed
}
