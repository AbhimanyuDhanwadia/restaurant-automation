package handlers

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/restaurantautomation/api/internal/backups"
)

func TestBackupHandlersListAndRecordBackup(t *testing.T) {
	service := backups.NewService(backups.NewMemoryRepository())
	createResponse := httptest.NewRecorder()
	RecordBackup(service).ServeHTTP(createResponse, httptest.NewRequest(http.MethodPost, "/api/v1/admin/backups", bytes.NewBufferString(`{"target":"s3://restaurant/backups.sql.gz","size_bytes":1024}`)))
	if createResponse.Code != http.StatusCreated {
		t.Fatalf("record status = %d body=%s", createResponse.Code, createResponse.Body.String())
	}
	listResponse := httptest.NewRecorder()
	ListBackups(service).ServeHTTP(listResponse, httptest.NewRequest(http.MethodGet, "/api/v1/admin/backups", nil))
	if listResponse.Code != http.StatusOK {
		t.Fatalf("list status = %d", listResponse.Code)
	}
}
