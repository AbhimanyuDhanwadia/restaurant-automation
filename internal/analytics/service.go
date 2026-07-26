package analytics

import (
	"sort"
	"time"

	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/printers"
)

type Metric struct {
	Value     float64 `json:"value"`
	Available bool    `json:"available"`
}

type Hour struct {
	Hour   int `json:"hour"`
	Orders int `json:"orders"`
}

type Report struct {
	GeneratedAt         time.Time `json:"generated_at"`
	Orders              Metric    `json:"orders"`
	Sales               Metric    `json:"sales"`
	AverageTicket       Metric    `json:"average_ticket"`
	KitchenCompletion   Metric    `json:"kitchen_completion"`
	DeliveryTime        Metric    `json:"delivery_time"`
	PrinterAvailability Metric    `json:"printer_availability"`
	PrinterUtilization  Metric    `json:"printer_utilization"`
	StaffProductivity   Metric    `json:"staff_productivity"`
	PeakHours           []Hour    `json:"peak_hours"`
}

type Service struct {
	engine   *automation.Engine
	printers *printers.Manager
}

func NewService(engine *automation.Engine, printerManager *printers.Manager) *Service {
	return &Service{engine: engine, printers: printerManager}
}

// Overview only reports values that can be derived from the current in-memory
// event stream and printer manager. Sales, delivery, and staffing become
// available when persisted order, delivery, and shift data are introduced.
func (s *Service) Overview() Report {
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
	return Report{GeneratedAt: time.Now().UTC(), Orders: Metric{Value: float64(len(orders)), Available: true}, KitchenCompletion: completion, PrinterAvailability: availability, PrinterUtilization: Metric{Value: float64(printed), Available: true}, Sales: Metric{}, AverageTicket: Metric{}, DeliveryTime: Metric{}, StaffProductivity: Metric{}, PeakHours: peakHours}
}
