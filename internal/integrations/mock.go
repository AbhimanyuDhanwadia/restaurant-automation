package integrations

import "sync"

// MockProvider is used by tests and local development. It models an external
// provider without making network requests.
type MockProvider struct {
	name      string
	mu        sync.RWMutex
	status    Status
	orders    chan Order
	closeOnce sync.Once
}

func NewMockProvider(name string, buffer int) *MockProvider {
	if buffer < 1 {
		buffer = 1
	}
	return &MockProvider{name: name, status: StatusDisconnected, orders: make(chan Order, buffer)}
}
func (p *MockProvider) Name() string { return p.name }
func (p *MockProvider) Connect() error {
	p.mu.Lock()
	p.status = StatusConnected
	p.mu.Unlock()
	return nil
}
func (p *MockProvider) Disconnect() error {
	p.mu.Lock()
	p.status = StatusDisconnected
	p.mu.Unlock()
	return nil
}
func (p *MockProvider) Health() Status              { p.mu.RLock(); defer p.mu.RUnlock(); return p.status }
func (p *MockProvider) ReceiveOrders() <-chan Order { return p.orders }
func (p *MockProvider) Submit(order Order) bool {
	p.mu.RLock()
	connected := p.status == StatusConnected
	p.mu.RUnlock()
	if !connected {
		return false
	}
	select {
	case p.orders <- order:
		return true
	default:
		return false
	}
}
func (p *MockProvider) Close() { p.closeOnce.Do(func() { close(p.orders) }) }
