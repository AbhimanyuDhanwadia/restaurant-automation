package intelligence

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/restaurantautomation/api/internal/automation"
	"github.com/restaurantautomation/api/internal/printers"
)

type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityWarning  Severity = "warning"
	SeverityCritical Severity = "critical"
)

type Insight struct {
	ID        string    `json:"id"`
	Kind      string    `json:"kind"`
	Severity  Severity  `json:"severity"`
	Resource  string    `json:"resource"`
	Title     string    `json:"title"`
	Detail    string    `json:"detail"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Service is a separate, read-only consumer of operational events. It makes
// explainable rule-based observations now and can be replaced by an AI model
// consumer later without modifying the automation engine.
type Service struct {
	engine   *automation.Engine
	printers *printers.Manager
	mu       sync.RWMutex
	insights map[string]Insight
	cancel   context.CancelFunc
	stopSub  func()
	wg       sync.WaitGroup
	started  bool
}

func NewService(engine *automation.Engine, printerManager *printers.Manager) *Service {
	return &Service{engine: engine, printers: printerManager, insights: make(map[string]Insight)}
}

func (s *Service) Start(parent context.Context) {
	s.mu.Lock()
	if s.started {
		s.mu.Unlock()
		return
	}
	ctx, cancel := context.WithCancel(parent)
	s.cancel = cancel
	events, stopSub := s.engine.Subscribe(64)
	s.stopSub = stopSub
	s.started = true
	s.mu.Unlock()
	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		ticker := time.NewTicker(10 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case event, ok := <-events:
				if !ok {
					return
				}
				s.evaluateEvent(event)
			case <-ticker.C:
				s.evaluateRuntime()
			}
		}
	}()
	s.evaluateRuntime()
}

func (s *Service) evaluateEvent(event automation.Event) {
	switch event.Type {
	case automation.EventJobRetrying:
		s.record("printer-anomaly:"+event.OrderID, Insight{Kind: "printer_anomaly", Severity: SeverityWarning, Resource: event.OrderID, Title: "Print job is retrying", Detail: "The automation engine retried a job for this order."})
	case automation.EventJobFailed:
		s.record("delayed-order:"+event.OrderID, Insight{Kind: "delayed_order", Severity: SeverityCritical, Resource: event.OrderID, Title: "Order requires attention", Detail: "The automation engine could not finish processing this order."})
	case automation.EventPrintFailed:
		resource := event.Payload["destination"]
		if resource == "" {
			resource = event.OrderID
		}
		s.record("printer-anomaly:"+resource, Insight{Kind: "printer_anomaly", Severity: SeverityWarning, Resource: resource, Title: "Print job failed", Detail: "A printer could not complete a ticket and the job is available for reprint."})
	}
}

func (s *Service) evaluateRuntime() {
	stats := s.engine.Stats()
	if stats.QueueCapacity > 0 && stats.QueueDepth*100 >= stats.QueueCapacity*80 {
		s.record("kitchen-load:queue", Insight{Kind: "kitchen_load", Severity: SeverityWarning, Resource: "automation-queue", Title: "Automation queue is under pressure", Detail: "More than 80% of queue capacity is occupied."})
	}
	for _, printer := range s.printers.Health() {
		if printer.Status != printers.StatusReady {
			s.record("printer-anomaly:"+printer.Name, Insight{Kind: "printer_anomaly", Severity: SeverityWarning, Resource: printer.Name, Title: "Printer is unavailable", Detail: "The printer is not ready and queued tickets may be delayed."})
		}
	}
}

func (s *Service) RecordPrintJobNeedsReview(jobID, destination string) {
	resource := destination
	if resource == "" {
		resource = jobID
	}
	s.record("print-review:"+jobID, Insight{Kind: "printer_anomaly", Severity: SeverityWarning, Resource: resource, Title: "Print job needs review", Detail: "A ticket was already printing before startup and was not replayed automatically."})
}

func (s *Service) record(key string, insight Insight) {
	now := time.Now().UTC()
	s.mu.Lock()
	defer s.mu.Unlock()
	if existing, ok := s.insights[key]; ok {
		insight.ID = existing.ID
		insight.CreatedAt = existing.CreatedAt
	} else {
		insight.ID = key
		insight.CreatedAt = now
	}
	insight.UpdatedAt = now
	s.insights[key] = insight
}

func (s *Service) Insights() []Insight {
	s.mu.RLock()
	result := make([]Insight, 0, len(s.insights))
	for _, insight := range s.insights {
		result = append(result, insight)
	}
	s.mu.RUnlock()
	sort.Slice(result, func(i, j int) bool { return result[i].UpdatedAt.After(result[j].UpdatedAt) })
	return result
}

func (s *Service) Close() {
	s.mu.Lock()
	if !s.started {
		s.mu.Unlock()
		return
	}
	s.cancel()
	s.stopSub()
	s.mu.Unlock()
	s.wg.Wait()
}
