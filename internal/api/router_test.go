package api_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/restaurantautomation/api/internal/analytics"
	"github.com/restaurantautomation/api/internal/api"
	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/config"
	"github.com/restaurantautomation/api/internal/integrations"
	"github.com/restaurantautomation/api/internal/printers"
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
	router := api.NewRouter(cfg, log, engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager))

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
	router := api.NewRouter(cfg, log, engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager))

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
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/automation/orders", bytes.NewBufferString(`{"order_id":"ORD-100"}`))
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusAccepted)
	}
}

func TestAutomationOrderEndpointValidatesOrderID(t *testing.T) {
	cfg := &config.Config{Server: config.ServerConfig{CORSOrigins: []string{"*"}}}
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	router := api.NewRouter(cfg, zerolog.Nop(), engine, integrations.NewRegistry(), printerManager, analytics.NewService(engine, printerManager))

	req := httptest.NewRequest(http.MethodPost, "/api/v1/automation/orders", bytes.NewBufferString(`{}`))
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)
	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}
}
