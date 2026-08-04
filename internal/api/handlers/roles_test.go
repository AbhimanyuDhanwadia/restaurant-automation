package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/restaurantautomation/api/internal/roles"
)

func TestRoleHandlersListAndCreateRoles(t *testing.T) {
	service := roles.NewService(roles.NewMemoryRepository())
	listResponse := httptest.NewRecorder()
	ListRoles(service).ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/api/v1/admin/roles", nil))
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d", listResponse.Code)
	}

	createResponse := httptest.NewRecorder()
	body := bytes.NewBufferString(`{"name":"Closer","description":"Closes the shift.","permissions":["orders.manage"]}`)
	CreateRole(service).ServeHTTP(createResponse, httptest.NewRequest(http.MethodPost, "/api/v1/admin/roles", body))
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("create status = %d body=%s", createResponse.Code, createResponse.Body.String())
	}
}
