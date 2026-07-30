package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/restaurantautomation/api/internal/staff"
)

func ListStaff(service *staff.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		members, err := service.ListMembers(r.Context())
		if err != nil {
			http.Error(w, "staff unavailable", 500)
			return
		}
		writeJSON(w, 200, members)
	}
}
func CreateStaffMember(service *staff.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input staff.CreateMemberInput
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			http.Error(w, "invalid staff member", 400)
			return
		}
		member, err := service.CreateMember(r.Context(), input)
		if err != nil {
			http.Error(w, "invalid staff member", 400)
			return
		}
		writeJSON(w, 201, member)
	}
}
func UpdateStaffStatus(service *staff.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Status string `json:"status"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			http.Error(w, "invalid staff status", 400)
			return
		}
		member, err := service.UpdateMemberStatus(r.Context(), chi.URLParam(r, "staffID"), body.Status)
		if errors.Is(err, staff.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "invalid staff status", 400)
			return
		}
		writeJSON(w, 200, member)
	}
}
func UpdateStaffHandoff(service *staff.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Handoff string `json:"handoff"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			http.Error(w, "invalid handoff", 400)
			return
		}
		member, err := service.UpdateMemberHandoff(r.Context(), chi.URLParam(r, "staffID"), body.Handoff)
		if errors.Is(err, staff.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "invalid handoff", 400)
			return
		}
		writeJSON(w, 200, member)
	}
}
func ListShiftTasks(service *staff.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tasks, err := service.ListTasks(r.Context())
		if err != nil {
			http.Error(w, "tasks unavailable", 500)
			return
		}
		writeJSON(w, 200, tasks)
	}
}
func CreateShiftTask(service *staff.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var input staff.CreateTaskInput
		if json.NewDecoder(r.Body).Decode(&input) != nil {
			http.Error(w, "invalid shift task", 400)
			return
		}
		task, err := service.CreateTask(r.Context(), input)
		if err != nil {
			http.Error(w, "invalid shift task", 400)
			return
		}
		writeJSON(w, 201, task)
	}
}
func CompleteShiftTask(service *staff.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		task, err := service.CompleteTask(r.Context(), chi.URLParam(r, "taskID"))
		if errors.Is(err, staff.ErrNotFound) {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "task unavailable", 500)
			return
		}
		writeJSON(w, 200, task)
	}
}
func GetShiftHandoff(service *staff.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		handoff, err := service.GetHandoff(r.Context())
		if err != nil {
			http.Error(w, "handoff unavailable", 500)
			return
		}
		writeJSON(w, 200, handoff)
	}
}
func SaveShiftHandoff(service *staff.Service) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var body struct {
			Note string `json:"note"`
		}
		if json.NewDecoder(r.Body).Decode(&body) != nil {
			http.Error(w, "invalid handoff", 400)
			return
		}
		handoff, err := service.SaveHandoff(r.Context(), body.Note)
		if err != nil {
			http.Error(w, "handoff unavailable", 500)
			return
		}
		writeJSON(w, 200, handoff)
	}
}
