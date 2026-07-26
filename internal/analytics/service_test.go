package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/printers"
)

func TestOverviewAggregatesOrderAndPrinterState(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	if err := engine.SubmitOrder(context.Background(), "ORD-400"); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for len(engine.Events()) < 5 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	driver := printers.NewMockDriver("kitchen")
	if err := printerManager.Register(driver, "kitchen"); err != nil {
		t.Fatal(err)
	}
	if err := printerManager.Start(context.Background()); err != nil {
		t.Fatal(err)
	}
	defer printerManager.Close()
	report := NewService(engine, printerManager).Overview()
	if !report.Orders.Available || report.Orders.Value != 1 {
		t.Fatalf("orders = %+v", report.Orders)
	}
	if !report.PrinterAvailability.Available || report.PrinterAvailability.Value != 100 {
		t.Fatalf("availability = %+v", report.PrinterAvailability)
	}
	if !report.KitchenCompletion.Available || report.KitchenCompletion.Value != 100 {
		t.Fatalf("completion = %+v", report.KitchenCompletion)
	}
}
