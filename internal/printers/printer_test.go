package printers

import (
	"context"
	"errors"
	"sync"
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
	if err := manager.Reprint(context.Background(), ticket.OrderID, ""); err != nil {
		t.Fatal(err)
	}
	waitForPrints(t, driver, 2)
}

func TestManagerRoutesTicketsToDedicatedPrinters(t *testing.T) {
	kitchen := NewMockDriver("kitchen")
	cashier := NewMockDriver("cashier")
	manager := NewManager(ESCPosFormatter{}, 1)
	if err := manager.Register(kitchen, "kitchen"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Register(cashier, "cashier"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer manager.Close()
	if err := manager.Print(context.Background(), Ticket{OrderID: "ORD-302", Destination: "kitchen", Lines: []Line{{Text: "Curry"}}}); err != nil {
		t.Fatal(err)
	}
	if err := manager.Print(context.Background(), Ticket{OrderID: "ORD-303", Destination: "cashier", Lines: []Line{{Text: "Receipt"}}}); err != nil {
		t.Fatal(err)
	}
	waitForPrints(t, kitchen, 1)
	waitForPrints(t, cashier, 1)
}

func TestManagerReconnectsOfflinePrinterWithoutQueuedTicket(t *testing.T) {
	driver := &recoveringDriver{name: "kitchen", readyAfter: 2, status: StatusOffline}
	manager := NewManager(ESCPosFormatter{}, 1)
	manager.SetReconnectInterval(time.Millisecond)
	if err := manager.Register(driver, "kitchen"); err != nil {
		t.Fatal(err)
	}
	if err := manager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer manager.Close()

	deadline := time.Now().Add(time.Second)
	for driver.Health() != StatusReady && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	if driver.Health() != StatusReady {
		t.Fatal("printer did not reconnect")
	}
	if driver.ConnectCount() < 2 {
		t.Fatalf("connect attempts = %d, want at least 2", driver.ConnectCount())
	}
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

type recoveringDriver struct {
	mu         sync.RWMutex
	name       string
	status     Status
	connects   int
	readyAfter int
}

func (driver *recoveringDriver) Name() string { return driver.name }
func (driver *recoveringDriver) Connect() error {
	driver.mu.Lock()
	defer driver.mu.Unlock()
	driver.connects++
	if driver.connects < driver.readyAfter {
		driver.status = StatusOffline
		return errors.New("printer unavailable")
	}
	driver.status = StatusReady
	return nil
}
func (driver *recoveringDriver) Disconnect() error {
	driver.mu.Lock()
	driver.status = StatusOffline
	driver.mu.Unlock()
	return nil
}
func (driver *recoveringDriver) Health() Status {
	driver.mu.RLock()
	defer driver.mu.RUnlock()
	return driver.status
}
func (driver *recoveringDriver) Print([]byte) error { return nil }
func (driver *recoveringDriver) ConnectCount() int {
	driver.mu.RLock()
	defer driver.mu.RUnlock()
	return driver.connects
}
