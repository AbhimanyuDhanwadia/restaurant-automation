package printqueue

import (
	"context"
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
