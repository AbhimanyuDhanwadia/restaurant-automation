package printers

import (
	"context"
	"testing"
	"time"
)

func TestManagerPrintsAndReprints(t *testing.T) {
	driver := NewMockDriver("kitchen")
	manager := NewManager(ESCPosFormatter{}, 2)
	if err := manager.Register(driver, "kitchen"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer manager.Close()
	ticket := Ticket{OrderID: "ORD-300", Destination: "kitchen", Lines: []Line{{Text: "Paneer Tikka", Quantity: 2}}}
	if err := manager.Print(context.Background(), ticket); err != nil {
		t.Fatal(err)
	}
	waitForPrints(t, driver, 1)
	if err := manager.Reprint(context.Background(), ticket.OrderID); err != nil {
		t.Fatal(err)
	}
	waitForPrints(t, driver, 2)
}

func waitForPrints(t *testing.T, driver *MockDriver, count int) {
	t.Helper()
	deadline := time.Now().Add(time.Second)
	for driver.PrintCount() < count && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if driver.PrintCount() != count {
		t.Fatalf("printed %d tickets, want %d", driver.PrintCount(), count)
	}
}

func TestESCPosFormatterIncludesPrinterCommands(t *testing.T) {
	data := (ESCPosFormatter{}).Format(Ticket{OrderID: "ORD-301", Destination: "cashier", Lines: []Line{{Text: "Total"}}})
	if len(data) < 5 || data[0] != 0x1b || data[len(data)-3] != 0x1d {
		t.Fatalf("unexpected ESC/POS payload: %v", data)
	}
}
