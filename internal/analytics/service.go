package analytics

import (
	"context"
	"math"
	"sort"
	"time"

	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/orders"
	"github.com/restaurantautomation/api/internal/printers"
	"github.com/restaurantautomation/api/internal/staff"
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

type StaffMetric struct {
	Metric
	CompletedTasks int `json:"completed_tasks"`
	TotalTasks     int `json:"total_tasks"`
}

type KitchenMetric struct {
	Metric
	CompletedOrders int `json:"completed_orders"`
	EligibleOrders  int `json:"eligible_orders"`
}

type StaffTaskSummaryReader interface {
	TaskSummary(context.Context) (staff.TaskSummary, error)
}

type Hour struct {
	Hour   int `json:"hour"`
	Orders int `json:"orders"`
}

type Report struct {
	GeneratedAt         time.Time      `json:"generated_at"`
	Orders              Metric         `json:"orders"`
	OrderVolumeSource   string         `json:"order_volume_source"`
	Sales               MonetaryMetric `json:"sales"`
	AverageTicket       MonetaryMetric `json:"average_ticket"`
	KitchenCompletion   KitchenMetric  `json:"kitchen_completion"`
	KitchenSource       string         `json:"kitchen_completion_source"`
	DeliveryTime        DeliveryMetric `json:"delivery_time"`
	PrinterAvailability Metric         `json:"printer_availability"`
	PrinterUtilization  Metric         `json:"printer_utilization"`
	StaffProductivity   StaffMetric    `json:"staff_productivity"`
	PeakHours           []Hour         `json:"peak_hours"`
}

type Service struct {
	engine   *automation.Engine
	printers *printers.Manager
	orders   *orders.Service
	staff    StaffTaskSummaryReader
}

func NewService(engine *automation.Engine, printerManager *printers.Manager, orderServices ...*orders.Service) *Service {
	service := &Service{engine: engine, printers: printerManager}
	if len(orderServices) > 0 {
		service.orders = orderServices[0]
	}
	return service
}

func (s *Service) WithStaff(staffReader StaffTaskSummaryReader) *Service {
	s.staff = staffReader
	return s
}

// Overview combines runtime automation and printer telemetry with durable
// order, delivery, and staff-task records.
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
	durableOrders, durableOrdersAvailable := s.durableOrders(ctx)
	orderCount := len(orders)
	orderVolumeSource := "runtime"
	if durableOrdersAvailable {
		orderCount = len(durableOrders)
		orderVolumeSource = "durable"
		hours = make(map[int]int)
		for _, order := range durableOrders {
			hours[order.CreatedAt.Hour()]++
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
	completion, kitchenSource := kitchenCompletion(orders, queued, durableOrders, durableOrdersAvailable)
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
	sales, averageTicket := sales(durableOrders, durableOrdersAvailable)
	deliveryTime := deliveryTime(durableOrders, durableOrdersAvailable)
	return Report{GeneratedAt: time.Now().UTC(), Orders: Metric{Value: float64(orderCount), Available: true}, OrderVolumeSource: orderVolumeSource, KitchenCompletion: completion, KitchenSource: kitchenSource, PrinterAvailability: availability, PrinterUtilization: Metric{Value: float64(printed), Available: true}, Sales: sales, AverageTicket: averageTicket, DeliveryTime: deliveryTime, StaffProductivity: s.staffProductivity(ctx), PeakHours: peakHours}
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

func (s *Service) staffProductivity(ctx context.Context) StaffMetric {
	if s.staff == nil {
		return StaffMetric{}
	}
	summary, err := s.staff.TaskSummary(ctx)
	if err != nil {
		return StaffMetric{}
	}
	metric := StaffMetric{CompletedTasks: summary.Completed, TotalTasks: summary.Total}
	if summary.Total > 0 {
		metric.Value = math.Round(float64(summary.Completed)/float64(summary.Total)*1000) / 10
		metric.Available = true
	}
	return metric
}

func kitchenCompletion(runtimeOrders, runtimeQueued map[string]struct{}, durableOrders []orders.Order, durableAvailable bool) (KitchenMetric, string) {
	metric := KitchenMetric{}
	source := "runtime"
	if durableAvailable {
		source = "durable"
		for _, order := range durableOrders {
			if order.Status == "cancelled" {
				continue
			}
			metric.EligibleOrders++
			if order.Status == "ready" || order.Status == "delivered" {
				metric.CompletedOrders++
			}
		}
	} else {
		metric.EligibleOrders = len(runtimeOrders)
		metric.CompletedOrders = len(runtimeQueued)
	}
	if metric.EligibleOrders > 0 {
		metric.Value = math.Round(float64(metric.CompletedOrders)/float64(metric.EligibleOrders)*1000) / 10
		metric.Available = true
	}
	return metric, source
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
