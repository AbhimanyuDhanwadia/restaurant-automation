package alerts

import (
	"context"
	"testing"
)

func TestServiceCreatesAndAcknowledgesAlert(t *testing.T) {
	service := NewService(NewMemoryRepository())
	alert, err := service.Create(context.Background(), CreateInput{Label: "Cold room", Detail: "Temperature rising", Severity: "high"})
	if err != nil {
		t.Fatal(err)
	}
	if err := service.Acknowledge(context.Background(), alert.ID); err != nil {
		t.Fatal(err)
	}
	alerts, err := service.ListActive(context.Background())
	if err != nil || len(alerts) != 0 {
		t.Fatalf("alerts = %+v, err = %v", alerts, err)
	}
}
