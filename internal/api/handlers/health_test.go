package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/integrations"
	"github.com/restaurantautomation/api/internal/printers"
)

type failingChecker struct{}

func (failingChecker) Ping(context.Context) error { return errors.New("down") }

func TestReadyWithCheckReturnsServiceUnavailable(t *testing.T) {
	response := httptest.NewRecorder()
	ReadyWithCheck(failingChecker{}).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ready", nil))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", response.Code)
	}
}

func TestSystemHealthReportsComponents(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	registry := integrations.NewRegistry()
	provider := integrations.NewMockProvider("swiggy", 1)
	if err := registry.Register(provider); err != nil {
		t.Fatal(err)
	}
	if err := provider.Connect(); err != nil {
		t.Fatal(err)
	}
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := printerManager.Register(printers.NewMockDriver("kitchen"), "kitchen"); err != nil {
		t.Fatal(err)
	}
	if err := printerManager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer printerManager.Close()

	response := httptest.NewRecorder()
	SystemHealth(engine, registry, printerManager, nil).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/system/health", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	var payload SystemHealthResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Status != "healthy" || len(payload.Components) != 4 {
		t.Fatalf("payload = %#v", payload)
	}
}
