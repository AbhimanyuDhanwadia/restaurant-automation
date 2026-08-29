package analytics

import (
	"context"
	"testing"
	"time"

	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/orders"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/staff"
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
	orderService := orders.NewService(orders.NewMemoryRepository())
	total := int64(12550)
	if _, err := orderService.Create(context.Background(), orders.CreateInput{Channel: "counter", TotalMinor: &total, Currency: "INR", Items: []orders.Item{{Name: "Paneer", Quantity: 1}}}); err != nil {
		t.Fatal(err)
	}
	report := NewService(engine, printerManager, orderService).Overview(context.Background())
	if !report.Orders.Available || report.Orders.Value != 1 {
		t.Fatalf("orders = %+v", report.Orders)
	}
	if !report.PrinterAvailability.Available || report.PrinterAvailability.Value != 100 {
		t.Fatalf("availability = %+v", report.PrinterAvailability)
	}
	if !report.KitchenCompletion.Available || report.KitchenCompletion.Value != 0 || report.KitchenSource != "durable" {
		t.Fatalf("completion = %+v", report.KitchenCompletion)
	}
	if !report.Sales.Available || report.Sales.Value != 125.5 || report.Sales.Currency != "INR" {
		t.Fatalf("sales = %+v", report.Sales)
	}
	if !report.AverageTicket.Available || report.AverageTicket.Value != 125.5 {
		t.Fatalf("average ticket = %+v", report.AverageTicket)
	}
}

func TestOverviewDoesNotCombineCurrencies(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	orderService := orders.NewService(orders.NewMemoryRepository())
	inrTotal := int64(10000)
	usdTotal := int64(10000)
	for _, input := range []orders.CreateInput{
		{Channel: "counter", TotalMinor: &inrTotal, Currency: "INR", Items: []orders.Item{{Name: "Paneer", Quantity: 1}}},
		{Channel: "counter", TotalMinor: &usdTotal, Currency: "USD", Items: []orders.Item{{Name: "Paneer", Quantity: 1}}},
	} {
		if _, err := orderService.Create(context.Background(), input); err != nil {
			t.Fatal(err)
		}
	}
	report := NewService(engine, printerManager, orderService).Overview(context.Background())
	if report.Sales.Available || report.AverageTicket.Available || report.Sales.IncludedOrders != 2 {
		t.Fatalf("sales = %+v average = %+v", report.Sales, report.AverageTicket)
	}
}

func TestOverviewUsesDurableOrdersForVolumeAndPeakHours(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	orderService := orders.NewService(orders.NewMemoryRepository())
	for _, item := range []string{"Dosa", "Idli"} {
		if _, err := orderService.Create(context.Background(), orders.CreateInput{Channel: "counter", Items: []orders.Item{{Name: item, Quantity: 1}}}); err != nil {
			t.Fatal(err)
		}
	}

	report := NewService(engine, printerManager, orderService).Overview(context.Background())
	if report.OrderVolumeSource != "durable" || report.Orders.Value != 2 {
		t.Fatalf("order volume = %+v, source = %q", report.Orders, report.OrderVolumeSource)
	}
	if len(report.PeakHours) != 1 || report.PeakHours[0].Orders != 2 {
		t.Fatalf("peak hours = %+v", report.PeakHours)
	}
}

func TestOverviewFallsBackToRuntimeOrderVolume(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	if err := engine.SubmitOrder(context.Background(), "ORD-RUNTIME"); err != nil {
		t.Fatal(err)
	}
	deadline := time.Now().Add(time.Second)
	for len(engine.Events()) < 5 && time.Now().Before(deadline) {
		time.Sleep(time.Millisecond)
	}

	report := NewService(engine, printers.NewManager(printers.ESCPosFormatter{}, 1)).Overview(context.Background())
	if report.OrderVolumeSource != "runtime" || report.Orders.Value != 1 {
		t.Fatalf("order volume = %+v, source = %q", report.Orders, report.OrderVolumeSource)
	}
	if report.KitchenSource != "runtime" || report.KitchenCompletion.Value != 100 {
		t.Fatalf("kitchen completion = %+v, source = %q", report.KitchenCompletion, report.KitchenSource)
	}
}

func TestOverviewUsesDurableKitchenCompletion(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	orderService := orders.NewService(orders.NewMemoryRepository())
	statuses := []string{"ready", "delivered", "preparing", "cancelled"}
	for index, status := range statuses {
		order, err := orderService.Create(context.Background(), orders.CreateInput{Channel: "counter", Items: []orders.Item{{Name: "Item", Quantity: index + 1}}})
		if err != nil {
			t.Fatal(err)
		}
		if _, err := orderService.UpdateStatus(context.Background(), order.ID, status); err != nil {
			t.Fatal(err)
		}
	}

	report := NewService(engine, printerManager, orderService).Overview(context.Background())
	if report.KitchenSource != "durable" || !report.KitchenCompletion.Available || report.KitchenCompletion.Value != 66.7 {
		t.Fatalf("kitchen completion = %+v, source = %q", report.KitchenCompletion, report.KitchenSource)
	}
	if report.KitchenCompletion.CompletedOrders != 2 || report.KitchenCompletion.EligibleOrders != 3 {
		t.Fatalf("kitchen completion counts = %+v", report.KitchenCompletion)
	}
}

func TestOverviewAggregatesDeliveryTime(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	orderService := orders.NewService(orders.NewMemoryRepository())
	order, err := orderService.Create(context.Background(), orders.CreateInput{Channel: "swiggy", DeliveryPartner: "Swiggy", Items: []orders.Item{{Name: "Paneer", Quantity: 1}}})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := orderService.UpdateStatus(context.Background(), order.ID, "delivered"); err != nil {
		t.Fatal(err)
	}
	report := NewService(engine, printerManager, orderService).Overview(context.Background())
	if !report.DeliveryTime.Available || report.DeliveryTime.DeliveredOrders != 1 || report.DeliveryTime.ActiveOrders != 0 {
		t.Fatalf("delivery = %+v", report.DeliveryTime)
	}
}

func TestOverviewAggregatesStaffTaskCompletion(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	staffService := staff.NewService(staff.NewMemoryRepository())
	for index, title := range []string{"Open station", "Check stock", "Close station"} {
		task, err := staffService.CreateTask(context.Background(), staff.CreateTaskInput{Title: title, Owner: "Maya"})
		if err != nil {
			t.Fatal(err)
		}
		if index < 2 {
			if _, err := staffService.CompleteTask(context.Background(), task.ID); err != nil {
				t.Fatal(err)
			}
		}
	}

	report := NewService(engine, printerManager).WithStaff(staffService).Overview(context.Background())
	if !report.StaffProductivity.Available || report.StaffProductivity.Value != 66.7 {
		t.Fatalf("staff productivity = %+v", report.StaffProductivity)
	}
	if report.StaffProductivity.CompletedTasks != 2 || report.StaffProductivity.TotalTasks != 3 {
		t.Fatalf("staff productivity counts = %+v", report.StaffProductivity)
	}
}
