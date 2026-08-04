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

type DeliveryMetric struct {
	Value           float64 `json:"value"`
	Available       bool    `json:"available"`
	DeliveredOrders int     `json:"delivered_orders"`
	ActiveOrders    int     `json:"active_orders"`
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
	DeliveryTime        DeliveryMetric `json:"delivery_time"`
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
// totals and delivery timestamps. Staffing metrics remain unavailable until
// their source records are introduced.
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
	durableOrders, durableOrdersAvailable := s.durableOrders(ctx)
	sales, averageTicket := sales(durableOrders, durableOrdersAvailable)
	deliveryTime := deliveryTime(durableOrders, durableOrdersAvailable)
	return Report{GeneratedAt: time.Now().UTC(), Orders: Metric{Value: float64(len(orders)), Available: true}, KitchenCompletion: completion, PrinterAvailability: availability, PrinterUtilization: Metric{Value: float64(printed), Available: true}, Sales: sales, AverageTicket: averageTicket, DeliveryTime: deliveryTime, StaffProductivity: Metric{}, PeakHours: peakHours}
}

func (s *Service) durableOrders(ctx context.Context) ([]orders.Order, bool) {
	if s.orders == nil {
		return nil, false
	}
	allOrders, err := s.orders.List(ctx)
	if err != nil {
		return nil, false
	}
	return allOrders, true
}

func sales(allOrders []orders.Order, available bool) (MonetaryMetric, MonetaryMetric) {
	if !available {
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

func deliveryTime(allOrders []orders.Order, available bool) DeliveryMetric {
	if !available {
		return DeliveryMetric{}
	}
	result := DeliveryMetric{}
	var totalMinutes float64
	for _, order := range allOrders {
		if order.DeliveryPartner == "" || order.Status == "cancelled" {
			continue
		}
		if order.DeliveredAt == nil {
			result.ActiveOrders++
			continue
		}
		result.DeliveredOrders++
		if deliveredMinutes := order.DeliveredAt.Sub(order.CreatedAt).Minutes(); deliveredMinutes > 0 {
			totalMinutes += deliveredMinutes
		}
	}
	if result.DeliveredOrders > 0 {
		result.Value = totalMinutes / float64(result.DeliveredOrders)
		result.Available = true
	}
	return result
}
