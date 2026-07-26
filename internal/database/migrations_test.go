package database

import "testing"

func TestReadMigrationsRejectsMissingDirectory(t *testing.T) {
	if _, err := readMigrations("missing-migrations-directory"); err == nil {
		t.Fatal("expected missing directory error")
	}
}
