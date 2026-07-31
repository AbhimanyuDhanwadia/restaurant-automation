package settings

import (
	"context"
	"testing"
)

func TestServiceSavesSettings(t *testing.T) {
	service := NewService(NewMemoryRepository())
	saved, err := service.Save(context.Background(), Settings{RestaurantName: "The Green Table", Timezone: "Asia/Kolkata", OperationalAlerts: true})
	if err != nil {
		t.Fatal(err)
	}
	current, err := service.Get(context.Background())
	if err != nil || current.RestaurantName != saved.RestaurantName {
		t.Fatalf("settings = %+v, err = %v", current, err)
	}
}
