package staff

import (
	"context"
	"testing"
)

func TestServiceManagesRosterTasksAndHandoff(t *testing.T) {
	service := NewService(NewMemoryRepository())
	member, err := service.CreateMember(context.Background(), CreateMemberInput{Name: "Maya", Role: "Manager", Station: "Front of house"})
	if err != nil || member.Status != "off_shift" {
		t.Fatalf("member = %+v, err = %v", member, err)
	}
	updated, err := service.UpdateMemberStatus(context.Background(), member.ID, "on_shift")
	if err != nil || updated.Status != "on_shift" {
		t.Fatalf("member = %+v, err = %v", updated, err)
	}
	task, err := service.CreateTask(context.Background(), CreateTaskInput{Title: "Approve prep list", Owner: "Maya"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := service.CompleteTask(context.Background(), task.ID); err != nil {
		t.Fatal(err)
	}
	if tasks, err := service.ListTasks(context.Background()); err != nil || len(tasks) != 0 {
		t.Fatalf("tasks = %+v, err = %v", tasks, err)
	}
	if _, err := service.SaveHandoff(context.Background(), "Confirm closing counts."); err != nil {
		t.Fatal(err)
	}
	handoff, err := service.GetHandoff(context.Background())
	if err != nil || handoff.Note != "Confirm closing counts." {
		t.Fatalf("handoff = %+v, err = %v", handoff, err)
	}
}
