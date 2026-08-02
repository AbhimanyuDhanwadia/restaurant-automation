package handlers

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/restaurantautomation/api/internal/automation"
)

type auditEventReader struct {
	events []automation.Event
	err    error
}

func (reader auditEventReader) List(context.Context, int) ([]automation.Event, error) {
	return reader.events, reader.err
}

func TestAuditLogsUsesPersistentReader(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	stored := automation.NewEvent(automation.EventOrderStored, "ORD-200", nil)
	response := httptest.NewRecorder()
	AuditLogs(engine, auditEventReader{events: []automation.Event{stored}}).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/admin/audit-logs", nil))
	if response.Code != http.StatusOK {
		t.Fatalf("status = %d", response.Code)
	}
	var payload AuditLogsResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Source != "persistent" || len(payload.Events) != 1 || payload.Events[0].ID != stored.ID {
		t.Fatalf("payload = %#v", payload)
	}
}

func TestAuditLogsFallsBackToRuntimeEvents(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	if err := engine.SubmitOrder(context.Background(), "ORD-201"); err != nil {
		t.Fatal(err)
	}
	response := httptest.NewRecorder()
	AuditLogs(engine, nil).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/api/v1/admin/audit-logs", nil))
	var payload AuditLogsResponse
	if err := json.NewDecoder(response.Body).Decode(&payload); err != nil {
		t.Fatal(err)
	}
	if payload.Source != "runtime" || len(payload.Events) == 0 {
		t.Fatalf("payload = %#v", payload)
	}
}
