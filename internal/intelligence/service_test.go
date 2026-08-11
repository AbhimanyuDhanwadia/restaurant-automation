package intelligence

import (
	"context"
	"testing"

	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/printers"
)

func TestServiceCreatesInsightsFromFailuresAndOfflinePrinters(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	engine.Start(context.Background())
	defer engine.Close()
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	if err := printerManager.Register(printers.NewMockDriver("kitchen"), "kitchen"); err != nil {
		t.Fatal(err)
	}
	service := NewService(engine, printerManager)
	service.evaluateEvent(automation.NewEvent(automation.EventJobFailed, "ORD-500", nil))
	service.evaluateRuntime()
	insights := service.Insights()
	if len(insights) != 2 {
		t.Fatalf("got %d insights, want 2", len(insights))
	}
	if insights[0].Kind != "printer_anomaly" && insights[1].Kind != "printer_anomaly" {
		t.Fatal("expected printer anomaly")
	}
}

func TestServiceCreatesPrinterInsightFromPrintFailure(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	service := NewService(engine, printerManager)
	service.evaluateEvent(automation.NewEvent(automation.EventPrintFailed, "ORD-501", map[string]string{"destination": "kitchen"}))

	insights := service.Insights()
	if len(insights) != 1 || insights[0].Kind != "printer_anomaly" || insights[0].Resource != "kitchen" {
		t.Fatalf("insights = %#v, want kitchen printer anomaly", insights)
	}
}

func TestServiceRecordsPrintJobNeedingReview(t *testing.T) {
	engine := automation.NewEngine(1, 10, automation.RetryPolicy{MaxAttempts: 1})
	printerManager := printers.NewManager(printers.ESCPosFormatter{}, 1)
	service := NewService(engine, printerManager)
	service.RecordPrintJobNeedsReview("print-1", "cashier")

	insights := service.Insights()
	if len(insights) != 1 || insights[0].Kind != "printer_anomaly" || insights[0].Resource != "cashier" || insights[0].Title != "Print job needs review" {
		t.Fatalf("insights = %#v, want print review anomaly", insights)
	}
}
