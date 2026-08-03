package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/restaurantautomation/api/internal/database"
)

type databaseInspector struct {
	migrations []database.Migration
	err        error
}

func (inspector databaseInspector) Migrations(context.Context) ([]database.Migration, error) {
	return inspector.migrations, inspector.err
}

func TestDatabaseStatusReturnsMigrationMetadata(t *testing.T) {
	appliedAt := time.Date(2026, time.August, 3, 10, 0, 0, 0, time.UTC)
	response := httptest.NewRecorder()
	DatabaseStatus(databaseInspector{migrations: []database.Migration{{Version: 8, AppliedAt: appliedAt}}}).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/admin/database", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	var payload DatabaseStatusResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Status != "available" || len(payload.Migrations) != 1 || payload.Migrations[0].Version != 8 {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestDatabaseStatusReportsNotConfigured(t *testing.T) {
	response := httptest.NewRecorder()
	DatabaseStatus(nil).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/admin/database", nil))
	var payload DatabaseStatusResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Status != "not_configured" || payload.Migrations == nil {
		t.Fatalf("payload = %#v", payload)
	}
}
