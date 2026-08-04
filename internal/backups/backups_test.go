package backups

import (
	"context"
	"errors"
	"testing"
)

func TestServiceRecordsExternalBackup(t *testing.T) {
	service := NewService(NewMemoryRepository())
	record, err := service.Record(context.Background(), RecordInput{Target: "s3://restaurant-backups/2026-08-04.sql.gz", SizeBytes: 1024})
	if err != nil {
		t.Fatal(err)
	}
	if record.Status != StatusRecorded || record.SizeBytes != 1024 {
		t.Fatalf("record = %+v", record)
	}
}

func TestServiceRejectsInvalidBackupRecord(t *testing.T) {
	service := NewService(NewMemoryRepository())
	if _, err := service.Record(context.Background(), RecordInput{Target: " ", SizeBytes: -1}); !errors.Is(err, ErrInvalidRecord) {
		t.Fatalf("error = %v", err)
	}
}
