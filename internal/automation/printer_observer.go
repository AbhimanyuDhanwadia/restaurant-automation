package automation

import (
	"context"

	"github.com/restaurantautomation/api/internal/printers"
)

// PrinterEventObserver adapts printer lifecycle callbacks into automation
// events. The printer package stays independent of the automation engine.
type PrinterEventObserver struct{ engine *Engine }

func NewPrinterEventObserver(engine *Engine) *PrinterEventObserver {
	return &PrinterEventObserver{engine: engine}
}

func (observer *PrinterEventObserver) Printing(ctx context.Context, ticket printers.Ticket, attempt int) {
	if observer.engine != nil {
		_ = observer.engine.recordPrintEvent(ctx, EventPrintStarted, ticket.OrderID, ticket.Destination, attempt, nil)
	}
}

func (observer *PrinterEventObserver) Printed(ctx context.Context, ticket printers.Ticket, attempt int) {
	if observer.engine != nil {
		_ = observer.engine.recordPrintEvent(ctx, EventPrintFinished, ticket.OrderID, ticket.Destination, attempt, nil)
	}
}

func (observer *PrinterEventObserver) Failed(ctx context.Context, ticket printers.Ticket, attempt int, err error) {
	if observer.engine != nil {
		_ = observer.engine.recordPrintEvent(ctx, EventPrintFailed, ticket.OrderID, ticket.Destination, attempt, err)
	}
}
