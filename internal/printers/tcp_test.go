package printers

import (
	"io"
	"net"
	"testing"
	"time"
)

func TestTCPDriverPrintsToNetworkSocket(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatal(err)
	}
	defer listener.Close()
	received := make(chan []byte, 1)
	go func() {
		conn, acceptErr := listener.Accept()
		if acceptErr != nil {
			return
		}
		defer conn.Close()
		data := make([]byte, 3)
		if _, readErr := io.ReadFull(conn, data); readErr == nil {
			received <- data
		}
	}()
	driver := NewTCPDriver("kitchen", listener.Addr().String(), time.Second)
	if err := driver.Connect(); err != nil {
		t.Fatal(err)
	}
	defer driver.Disconnect()
	if err := driver.Print([]byte{0x1b, 0x40, 0x0a}); err != nil {
		t.Fatal(err)
	}
	select {
	case data := <-received:
		if len(data) != 3 || data[0] != 0x1b {
			t.Fatalf("received %v", data)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out reading printer bytes")
	}
}
