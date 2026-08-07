package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/orders"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/printqueue"
)

func TestCreateOrderQueuesAndPrintsKitchenTicket(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()

	driver := printers.NewMockDriver("kitchen")
	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := manager.Register(driver, "kitchen"); err != nil {
		t.Fatal(err)
	}
	queue := printqueue.NewService(printqueue.NewMemoryRepository())
	manager.SetObserver(queue)
	if err := manager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer manager.Close()

	handler := CreateOrder(orders.NewService(orders.NewMemoryRepository()), engine, manager, queue)
	request := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(`{"channel":"dine_in","items":[{"name":"Paneer Tikka","quantity":2}]}`))
	request.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}

	var order orders.Order
	if err := json.NewDecoder(recorder.Body).Decode(&order); err != nil {
		t.Fatal(err)
	}

	deadline := time.Now().Add(time.Second)
	for driver.PrintCount() != 1 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if driver.PrintCount() != 1 {
		t.Fatalf("prints = %d, want 1", driver.PrintCount())
	}

	jobs, err := queue.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 {
		t.Fatalf("jobs = %#v, want one job", jobs)
	}
	if jobs[0].OrderID != order.ID || jobs[0].Destination != "kitchen" || jobs[0].Status != printqueue.StatusPrinted {
		t.Fatalf("job = %#v, want printed kitchen job for %q", jobs[0], order.ID)
	}
	if len(jobs[0].Lines) != 1 || jobs[0].Lines[0].Text != "Paneer Tikka" || jobs[0].Lines[0].Quantity != 2 {
		t.Fatalf("ticket lines = %#v, want Paneer Tikka x2", jobs[0].Lines)
	}
}

func TestCreateOrderRecordsFailedKitchenDispatch(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()

	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	queue := printqueue.NewService(printqueue.NewMemoryRepository())
	handler := CreateOrder(orders.NewService(orders.NewMemoryRepository()), engine, manager, queue)
	request := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewBufferString(`{"channel":"takeaway","items":[{"name":"Masala Chai","quantity":1}]}`))
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusCreated {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusCreated)
	}

	jobs, err := queue.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || jobs[0].Status != printqueue.StatusFailed || jobs[0].Attempts != 0 {
		t.Fatalf("jobs = %#v, want one failed undispatched job", jobs)
	}
}

func TestUpdateOrderStatusRecordsOneLifecycleEvent(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()

	service := orders.NewService(orders.NewMemoryRepository())
	order, err := service.Create(context.Background(), orders.CreateInput{Channel: "dine_in", Items: []orders.Item{{Name: "Veg Biryani", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	handler := UpdateOrderStatus(service, engine)
	for range 2 {
		request := httptest.NewRequest(http.MethodPatch, "/orders/"+order.ID+"/status", bytes.NewBufferString(`{"status":"preparing"}`))
		routeContext := chi.NewRouteContext()
		routeContext.URLParams.Add("orderID", order.ID)
		request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeContext))
		recorder := httptest.NewRecorder()
		handler.ServeHTTP(recorder, request)
		if recorder.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
		}
	}

	events := engine.Events()
	if len(events) != 1 || events[0].Type != automation.EventKitchenAccepted || events[0].OrderID != order.ID {
		t.Fatalf("events = %#v, want one kitchen acceptance event for %q", events, order.ID)
	}
}
