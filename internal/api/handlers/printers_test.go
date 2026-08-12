package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/printqueue"
)

func TestRetryPrintJobAcceptsFailedJob(t *testing.T) {
	queue := printqueue.NewService(printqueue.NewMemoryRepository())
	failed, err := queue.Queue(context.Background(), printers.Ticket{OrderID: "ORD-410", Destination: "kitchen", Lines: []printers.Line{{Text: "Dosa", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	queue.Failed(context.Background(), printers.Ticket{PrintJobID: failed.ID}, 1, nil)
	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := manager.Register(printers.NewMockDriver("kitchen"), "kitchen"); err != nil {
		t.Fatal(err)
	}

	router := chi.NewRouter()
	router.Post("/queue/{jobID}/retry", RetryPrintJob(manager, queue))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/queue/"+failed.ID+"/retry", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusAccepted, response.Body.String())
	}
	jobs, err := queue.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 2 || jobs[0].Status != printqueue.StatusQueued || !jobs[0].Reprint {
		t.Fatalf("jobs = %#v, want queued retry job first", jobs)
	}
}

func TestRetryPrintJobRejectsNonFailedJob(t *testing.T) {
	queue := printqueue.NewService(printqueue.NewMemoryRepository())
	queued, err := queue.Queue(context.Background(), printers.Ticket{OrderID: "ORD-411", Destination: "kitchen", Lines: []printers.Line{{Text: "Idli", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := manager.Register(printers.NewMockDriver("kitchen"), "kitchen"); err != nil {
		t.Fatal(err)
	}

	router := chi.NewRouter()
	router.Post("/queue/{jobID}/retry", RetryPrintJob(manager, queue))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/queue/"+queued.ID+"/retry", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusConflict, response.Body.String())
	}
}

func TestRequeuePrintJobAcceptsPrintingJob(t *testing.T) {
	queue := printqueue.NewService(printqueue.NewMemoryRepository())
	printing, err := queue.Queue(context.Background(), printers.Ticket{OrderID: "ORD-414", Destination: "kitchen", Lines: []printers.Line{{Text: "Uttapam", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	queue.Printing(context.Background(), printers.Ticket{PrintJobID: printing.ID}, 1)
	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := manager.Register(printers.NewMockDriver("kitchen"), "kitchen"); err != nil {
		t.Fatal(err)
	}

	router := chi.NewRouter()
	router.Post("/queue/{jobID}/requeue", RequeuePrintJob(manager, queue))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/queue/"+printing.ID+"/requeue", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusAccepted {
		t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusAccepted, response.Body.String())
	}
	jobs, err := queue.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 2 || jobs[0].Status != printqueue.StatusQueued || !jobs[0].Reprint {
		t.Fatalf("jobs = %#v, want queued reprint job first", jobs)
	}
}

func TestRequeuePrintJobRejectsNonPrintingJob(t *testing.T) {
	queue := printqueue.NewService(printqueue.NewMemoryRepository())
	failed, err := queue.Queue(context.Background(), printers.Ticket{OrderID: "ORD-415", Destination: "cashier", Lines: []printers.Line{{Text: "Receipt", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	queue.Failed(context.Background(), printers.Ticket{PrintJobID: failed.ID}, 1, nil)
	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := manager.Register(printers.NewMockDriver("cashier"), "cashier"); err != nil {
		t.Fatal(err)
	}

	router := chi.NewRouter()
	router.Post("/queue/{jobID}/requeue", RequeuePrintJob(manager, queue))
	response := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodPost, "/queue/"+failed.ID+"/requeue", nil)
	router.ServeHTTP(response, request)

	if response.Code != http.StatusConflict {
		t.Fatalf("status = %d, want %d: %s", response.Code, http.StatusConflict, response.Body.String())
	}
}
