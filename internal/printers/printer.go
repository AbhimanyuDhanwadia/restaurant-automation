package printers

import (
	"context"
	"errors"
	"sort"
	"sync"
	"sync/atomic"
	"time"
)

type Status string

const (
	StatusReady   Status = "ready"
	StatusOffline Status = "offline"
	StatusError   Status = "error"
)

type Line struct {
	Text     string `json:"text"`
	Quantity int    `json:"quantity,omitempty"`
}
type Ticket struct {
	OrderID     string `json:"order_id"`
	Destination string `json:"destination"`
	Lines       []Line `json:"lines"`
	Reprint     bool   `json:"reprint,omitempty"`
}

type Driver interface {
	Name() string
	Connect() error
	Disconnect() error
	Health() Status
	Print([]byte) error
}

type PrinterHealth struct {
	Name       string `json:"name"`
	Status     Status `json:"status"`
	QueueDepth int    `json:"queue_depth"`
	Printed    uint64 `json:"printed"`
	Failed     uint64 `json:"failed"`
}

var ErrPrinterExists = errors.New("printer is already registered")
var ErrPrinterNotFound = errors.New("printer is not registered")
var ErrTicketNotFound = errors.New("ticket is not available for reprint")

type job struct {
	ticket   Ticket
	attempts int
}
type managedPrinter struct {
	driver  Driver
	queue   chan job
	printed atomic.Uint64
	failed  atomic.Uint64
}

type Manager struct {
	mu         sync.RWMutex
	printers   map[string]*managedPrinter
	routes     map[string]string
	history    map[string]Ticket
	formatter  Formatter
	retryLimit int
	ctx        context.Context
	cancel     context.CancelFunc
	wg         sync.WaitGroup
	started    bool
}

func NewManager(formatter Formatter, retryLimit int) *Manager {
	if retryLimit < 1 {
		retryLimit = 1
	}
	return &Manager{printers: make(map[string]*managedPrinter), routes: make(map[string]string), history: make(map[string]Ticket), formatter: formatter, retryLimit: retryLimit}
}
func (m *Manager) Register(driver Driver, destinations ...string) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	if _, exists := m.printers[driver.Name()]; exists {
		return ErrPrinterExists
	}
	printer := &managedPrinter{driver: driver, queue: make(chan job, 100)}
	m.printers[driver.Name()] = printer
	for _, destination := range destinations {
		m.routes[destination] = driver.Name()
	}
	return nil
}
func (m *Manager) Start(parent context.Context) error {
	m.mu.Lock()
	if m.started {
		m.mu.Unlock()
		return nil
	}
	m.ctx, m.cancel = context.WithCancel(parent)
	m.started = true
	printers := make([]*managedPrinter, 0, len(m.printers))
	for _, printer := range m.printers {
		printers = append(printers, printer)
	}
	m.mu.Unlock()
	for _, printer := range printers {
		if err := printer.driver.Connect(); err != nil {
			return err
		}
		m.wg.Add(1)
		go m.worker(printer)
	}
	return nil
}
func (m *Manager) Print(ctx context.Context, ticket Ticket) error {
	m.mu.Lock()
	printerName, ok := m.routes[ticket.Destination]
	if !ok {
		m.mu.Unlock()
		return ErrPrinterNotFound
	}
	printer := m.printers[printerName]
	m.history[ticket.OrderID] = ticket
	m.mu.Unlock()
	select {
	case printer.queue <- job{ticket: ticket}:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}
func (m *Manager) Reprint(ctx context.Context, orderID string) error {
	m.mu.RLock()
	ticket, ok := m.history[orderID]
	m.mu.RUnlock()
	if !ok {
		return ErrTicketNotFound
	}
	ticket.Reprint = true
	return m.Print(ctx, ticket)
}
func (m *Manager) worker(printer *managedPrinter) {
	defer m.wg.Done()
	for {
		select {
		case <-m.ctx.Done():
			return
		case current := <-printer.queue:
			data := m.formatter.Format(current.ticket)
			var err error
			for attempt := 1; attempt <= m.retryLimit; attempt++ {
				if printer.driver.Health() != StatusReady {
					_ = printer.driver.Connect()
				}
				err = printer.driver.Print(data)
				if err == nil {
					printer.printed.Add(1)
					break
				}
				time.Sleep(time.Duration(attempt) * 10 * time.Millisecond)
			}
			if err != nil {
				printer.failed.Add(1)
				_ = printer.driver.Disconnect()
			}
		}
	}
}
func (m *Manager) Health() []PrinterHealth {
	m.mu.RLock()
	defer m.mu.RUnlock()
	result := make([]PrinterHealth, 0, len(m.printers))
	for name, printer := range m.printers {
		result = append(result, PrinterHealth{Name: name, Status: printer.driver.Health(), QueueDepth: len(printer.queue), Printed: printer.printed.Load(), Failed: printer.failed.Load()})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].Name < result[j].Name })
	return result
}
func (m *Manager) Close() {
	m.mu.Lock()
	if !m.started {
		m.mu.Unlock()
		return
	}
	m.cancel()
	printers := make([]*managedPrinter, 0, len(m.printers))
	for _, printer := range m.printers {
		printers = append(printers, printer)
	}
	m.mu.Unlock()
	m.wg.Wait()
	for _, printer := range printers {
		_ = printer.driver.Disconnect()
	}
}

type MockDriver struct {
	mu     sync.RWMutex
	name   string
	status Status
	prints [][]byte
}

func NewMockDriver(name string) *MockDriver { return &MockDriver{name: name, status: StatusOffline} }
func (d *MockDriver) Name() string          { return d.name }
func (d *MockDriver) Connect() error        { d.mu.Lock(); d.status = StatusReady; d.mu.Unlock(); return nil }
func (d *MockDriver) Disconnect() error {
	d.mu.Lock()
	d.status = StatusOffline
	d.mu.Unlock()
	return nil
}
func (d *MockDriver) Health() Status { d.mu.RLock(); defer d.mu.RUnlock(); return d.status }
func (d *MockDriver) Print(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.status != StatusReady {
		return errors.New("printer is offline")
	}
	d.prints = append(d.prints, append([]byte(nil), data...))
	return nil
}
func (d *MockDriver) PrintCount() int { d.mu.RLock(); defer d.mu.RUnlock(); return len(d.prints) }
