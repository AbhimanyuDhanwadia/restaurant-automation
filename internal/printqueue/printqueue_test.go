package printqueue

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/restaurantautomation/api/internal/printers"
)

func TestServiceTracksPrintedJobAndReprint(t *testing.T) {
	service := NewService(NewMemoryRepository())
	ticket := printers.Ticket{OrderID: "ORD-400", Destination: "kitchen", Lines: []printers.Line{{Text: "Masala dosa", Quantity: 1}}}
	job, err := service.Queue(context.Background(), ticket)
	if err != nil {
		t.Fatal(err)
	}
	ticket.PrintJobID = job.ID
	service.Printing(context.Background(), ticket, 1)
	service.Printed(context.Background(), ticket, 1)
	jobs, err := service.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if len(jobs) != 1 || jobs[0].Status != StatusPrinted || jobs[0].Attempts != 1 {
		t.Fatalf("jobs = %#v", jobs)
	}
	reprint, err := service.Reprint(context.Background(), "ORD-400")
	if err != nil {
		t.Fatal(err)
	}
	if !reprint.Reprint || reprint.ID == job.ID || len(reprint.Lines) != 1 {
		t.Fatalf("reprint = %#v", reprint)
	}
}

func TestServiceRecoversOnlyQueuedJobs(t *testing.T) {
	service := NewService(NewMemoryRepository())
	queued, err := service.Queue(context.Background(), printers.Ticket{OrderID: "ORD-403", Destination: "kitchen", Lines: []printers.Line{{Text: "Dosa", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	printing, err := service.Queue(context.Background(), printers.Ticket{OrderID: "ORD-404", Destination: "kitchen", Lines: []printers.Line{{Text: "Idli", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	service.Printing(context.Background(), printers.Ticket{PrintJobID: printing.ID}, 1)
	failed, err := service.Queue(context.Background(), printers.Ticket{OrderID: "ORD-405", Destination: "kitchen", Lines: []printers.Line{{Text: "Chai", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	service.Failed(context.Background(), printers.Ticket{PrintJobID: failed.ID}, 1, nil)

	driver := printers.NewMockDriver("kitchen")
	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := manager.Register(driver, "kitchen"); err != nil {
		t.Fatal(err)
	}
	manager.SetObserver(service)
	if err := manager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer manager.Close()

	recovered, err := service.RecoverQueued(context.Background(), manager)
	if err != nil {
		t.Fatal(err)
	}
	if recovered != 1 {
		t.Fatalf("recovered = %d, want 1", recovered)
	}
	waitForDriverPrints(t, driver, 1)

	jobs, err := service.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	statusByID := make(map[string]Status, len(jobs))
	for _, job := range jobs {
		statusByID[job.ID] = job.Status
	}
	if statusByID[queued.ID] != StatusPrinted || statusByID[printing.ID] != StatusPrinting || statusByID[failed.ID] != StatusFailed {
		t.Fatalf("recovery statuses = %#v", statusByID)
	}
}

func TestServiceRecoveryRequiresPrinterManager(t *testing.T) {
	service := NewService(NewMemoryRepository())
	if _, err := service.RecoverQueued(context.Background(), nil); err != ErrPrinterManagerUnavailable {
		t.Fatalf("error = %v, want %v", err, ErrPrinterManagerUnavailable)
	}
}

func TestServiceRecoveryClaimsQueuedJobsOnce(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewService(repository)
	job, err := service.Queue(context.Background(), printers.Ticket{OrderID: "ORD-406", Destination: "kitchen", Lines: []printers.Line{{Text: "Tea", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := manager.Register(printers.NewMockDriver("kitchen"), "kitchen"); err != nil {
		t.Fatal(err)
	}

	recovered, err := service.RecoverQueued(context.Background(), manager)
	if err != nil {
		t.Fatal(err)
	}
	if recovered != 1 {
		t.Fatalf("recovered = %d, want 1", recovered)
	}
	recovered, err = service.RecoverQueued(context.Background(), manager)
	if err != nil {
		t.Fatal(err)
	}
	if recovered != 0 {
		t.Fatalf("second recovery = %d, want 0", recovered)
	}
	jobs, err := service.List(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if jobs[0].ID != job.ID || jobs[0].Status != StatusPrinting || jobs[0].Attempts != 1 {
		t.Fatalf("claimed job = %#v", jobs[0])
	}
}

func TestMemoryRepositoryClaimQueuedSkipsNonQueuedJobs(t *testing.T) {
	repository := NewMemoryRepository()
	service := NewService(repository)
	job, err := service.Queue(context.Background(), printers.Ticket{OrderID: "ORD-407", Destination: "kitchen", Lines: []printers.Line{{Text: "Coffee", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	service.Failed(context.Background(), printers.Ticket{PrintJobID: job.ID}, 2, nil)
	claimed, err := repository.ClaimQueued(context.Background(), job.ID)
	if err != nil {
		t.Fatal(err)
	}
	if claimed {
		t.Fatal("failed job was claimed for recovery")
	}
}

func waitForDriverPrints(t *testing.T, driver *printers.MockDriver, count int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for driver.PrintCount() < count && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if driver.PrintCount() != count {
		t.Fatalf("prints = %d, want %d", driver.PrintCount(), count)
	}
}

func TestManagerNotifiesAdditionalObservers(t *testing.T) {
	service := NewService(NewMemoryRepository())
	driver := printers.NewMockDriver("kitchen")
	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := manager.Register(driver, "kitchen"); err != nil {
		t.Fatal(err)
	}
	manager.SetObserver(service)
	additional := &countingObserver{}
	manager.AddObserver(additional)
	if err := manager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer manager.Close()

	job, err := service.Queue(context.Background(), printers.Ticket{OrderID: "ORD-402", Destination: "kitchen", Lines: []printers.Line{{Text: "Coffee", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.Print(context.Background(), printers.Ticket{OrderID: job.OrderID, Destination: job.Destination, Lines: job.Lines, PrintJobID: job.ID}); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		jobs, err := service.List(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if jobs[0].Status == StatusPrinted {
			if additional.printed.Load() != 1 {
				t.Fatalf("additional observer printed callbacks = %d, want 1", additional.printed.Load())
			}
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("print queue observer did not receive printed callback")
}

type countingObserver struct{ printed atomic.Int32 }

func (*countingObserver) Printing(context.Context, printers.Ticket, int) {}
func (observer *countingObserver) Printed(context.Context, printers.Ticket, int) {
	observer.printed.Add(1)
}
func (*countingObserver) Failed(context.Context, printers.Ticket, int, error) {
}

func TestManagerObserverMarksDurableJobPrinted(t *testing.T) {
	service := NewService(NewMemoryRepository())
	driver := printers.NewMockDriver("kitchen")
	manager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := manager.Register(driver, "kitchen"); err != nil {
		t.Fatal(err)
	}
	manager.SetObserver(service)
	if err := manager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer manager.Close()
	job, err := service.Queue(context.Background(), printers.Ticket{OrderID: "ORD-401", Destination: "kitchen", Lines: []printers.Line{{Text: "Chai", Quantity: 2}}})
	if err != nil {
		t.Fatal(err)
	}
	if err := manager.Print(context.Background(), printers.Ticket{OrderID: job.OrderID, Destination: job.Destination, Lines: job.Lines, PrintJobID: job.ID}); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		jobs, err := service.List(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		if jobs[0].Status == StatusPrinted {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("print job did not reach printed state")
}
