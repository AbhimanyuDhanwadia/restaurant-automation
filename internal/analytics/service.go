package analytics

import (
	"context"
	"sort"
	"time"

	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/orders"
	"github.com/restaurantautomation/api/internal/printers"
)

type Metric struct {
	Value     float64 `json:"value"`
	Available bool    `json:"available"`
}

type MonetaryMetric struct {
	Value          float64 `json:"value"`
	Available      bool    `json:"available"`
	Currency       string  `json:"currency,omitempty"`
	IncludedOrders int     `json:"included_orders"`
	ExcludedOrders int     `json:"excluded_orders"`
}

type Hour struct {
	Hour   int `json:"hour"`
	Orders int `json:"orders"`
}

type Report struct {
	GeneratedAt         time.Time      `json:"generated_at"`
	Orders              Metric         `json:"orders"`
	Sales               MonetaryMetric `json:"sales"`
	AverageTicket       MonetaryMetric `json:"average_ticket"`
	KitchenCompletion   Metric         `json:"kitchen_completion"`
	DeliveryTime        Metric         `json:"delivery_time"`
	PrinterAvailability Metric         `json:"printer_availability"`
	PrinterUtilization  Metric         `json:"printer_utilization"`
	StaffProductivity   Metric         `json:"staff_productivity"`
	PeakHours           []Hour         `json:"peak_hours"`
}

type Service struct {
	engine   *automation.Engine
	printers *printers.Manager
	orders   *orders.Service
}

func NewService(engine *automation.Engine, printerManager *printers.Manager, orderServices ...*orders.Service) *Service {
	service := &Service{engine: engine, printers: printerManager}
	if len(orderServices) > 0 {
		service.orders = orderServices[0]
	}
	return service
}

// Overview combines runtime automation and printer telemetry with durable order
// totals. Delivery and staffing metrics remain unavailable until their source
// records are introduced.
func (s *Service) Overview(ctx context.Context) Report {
	events := s.engine.Events()
	orders := make(map[string]struct{})
	queued := make(map[string]struct{})
	hours := make(map[int]int)
	for _, event := range events {
		if event.Type == automation.EventOrderReceived {
			orders[event.OrderID] = struct{}{}
			hours[event.CreatedAt.Hour()]++
		}
		if event.Type == automation.EventOrderQueued {
			queued[event.OrderID] = struct{}{}
		}
	}
	printerHealth := s.printers.Health()
	ready := 0
	printed := uint64(0)
	for _, printer := range printerHealth {
		if printer.Status == printers.StatusReady {
			ready++
		}
		printed += printer.Printed
	}
	availability := Metric{}
	if len(printerHealth) > 0 {
		availability = Metric{Value: float64(ready) / float64(len(printerHealth)) * 100, Available: true}
	}
	completion := Metric{}
	if len(orders) > 0 {
		completion = Metric{Value: float64(len(queued)) / float64(len(orders)) * 100, Available: true}
	}
	peakHours := make([]Hour, 0, len(hours))
	for hour, count := range hours {
		peakHours = append(peakHours, Hour{Hour: hour, Orders: count})
	}
	sort.Slice(peakHours, func(i, j int) bool {
		if peakHours[i].Orders == peakHours[j].Orders {
			return peakHours[i].Hour < peakHours[j].Hour
		}
		return peakHours[i].Orders > peakHours[j].Orders
	})
	sales, averageTicket := s.sales(ctx)
	return Report{GeneratedAt: time.Now().UTC(), Orders: Metric{Value: float64(len(orders)), Available: true}, KitchenCompletion: completion, PrinterAvailability: availability, PrinterUtilization: Metric{Value: float64(printed), Available: true}, Sales: sales, AverageTicket: averageTicket, DeliveryTime: Metric{}, StaffProductivity: Metric{}, PeakHours: peakHours}
}

func (s *Service) sales(ctx context.Context) (MonetaryMetric, MonetaryMetric) {
	if s.orders == nil {
		return MonetaryMetric{}, MonetaryMetric{}
	}
	allOrders, err := s.orders.List(ctx)
	if err != nil {
		return MonetaryMetric{}, MonetaryMetric{}
	}
	var sum int64
	currencies := make(map[string]struct{})
	included := 0
	excluded := 0
	for _, order := range allOrders {
		if order.Status == "cancelled" {
			continue
		}
		if order.TotalMinor == nil {
			excluded++
			continue
		}
		included++
		sum += *order.TotalMinor
		currencies[order.Currency] = struct{}{}
	}
	result := MonetaryMetric{IncludedOrders: included, ExcludedOrders: excluded}
	if included == 0 || len(currencies) != 1 {
		return result, result
	}
	for currency := range currencies {
		result.Currency = currency
	}
	result.Value = float64(sum) / 100
	result.Available = true
	average := result
	average.Value /= float64(included)
	return result, average
}
