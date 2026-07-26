package printers

import (
	"errors"
	"net"
	"sync"
	"time"
)

// TCPDriver sends already-rendered ESC/POS bytes to a network printer, which
// commonly exposes a raw TCP socket on port 9100.
type TCPDriver struct {
	mu      sync.RWMutex
	name    string
	address string
	timeout time.Duration
	conn    net.Conn
	status  Status
}

func NewTCPDriver(name, address string, timeout time.Duration) *TCPDriver {
	if timeout <= 0 {
		timeout = 3 * time.Second
	}
	return &TCPDriver{name: name, address: address, timeout: timeout, status: StatusOffline}
}

func (d *TCPDriver) Name() string { return d.name }

func (d *TCPDriver) Connect() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.conn != nil {
		d.status = StatusReady
		return nil
	}
	conn, err := net.DialTimeout("tcp", d.address, d.timeout)
	if err != nil {
		d.status = StatusOffline
		return err
	}
	d.conn = conn
	d.status = StatusReady
	return nil
}

func (d *TCPDriver) Disconnect() error {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.status = StatusOffline
	if d.conn == nil {
		return nil
	}
	err := d.conn.Close()
	d.conn = nil
	return err
}

func (d *TCPDriver) Health() Status { d.mu.RLock(); defer d.mu.RUnlock(); return d.status }

func (d *TCPDriver) Print(data []byte) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if d.conn == nil || d.status != StatusReady {
		return errors.New("printer is offline")
	}
	d.conn.SetWriteDeadline(time.Now().Add(d.timeout))
	for len(data) > 0 {
		written, err := d.conn.Write(data)
		if err != nil {
			d.conn.Close()
			d.conn = nil
			d.status = StatusOffline
			return err
		}
		data = data[written:]
	}
	return nil
}
