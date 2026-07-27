package handlers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

type failingChecker struct{}

func (failingChecker) Ping(context.Context) error { return errors.New("down") }

func TestReadyWithCheckReturnsServiceUnavailable(t *testing.T) {
	response := httptest.NewRecorder()
	ReadyWithCheck(failingChecker{}).ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ready", nil))
	if response.Code != http.StatusServiceUnavailable {
		t.Fatalf("status = %d", response.Code)
	}
}
