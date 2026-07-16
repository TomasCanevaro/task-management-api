package handlers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"task-management-api/internal/middleware"
	"task-management-api/internal/models"
	"task-management-api/internal/services"

	"github.com/gorilla/mux"
)

type TaskHandler struct {
	service *services.TaskService
}

func NewTaskHandler(service *services.TaskService) *TaskHandler {
	return &TaskHandler{
		service: service,
	}
}

// Helper functions

func writeJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{
		"error": message,
	})
}

// Request structs

type CreateTaskRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type AssignTaskRequest struct {
	WorkerID int `json:"worker_id"`
}

type UpdateStatusRequest struct {
	Status models.TaskStatus `json:"status"`
}

func (h *TaskHandler) CreateTask(w http.ResponseWriter, r *http.Request) {

	user := middleware.GetUser(r)

	var req CreateTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	if req.Title == "" {
		writeError(w, http.StatusBadRequest, "title is required")
		return
	}

	task, err := h.service.CreateTask(
		req.Title,
		req.Description,
		user,
	)

	if err != nil {

		if errors.Is(err, services.ErrForbidden) {
			writeError(w, http.StatusForbidden, err.Error())
			return
		}

		writeError(w, http.StatusInternalServerError, err.Error())
		return
	}

	writeJSON(w, http.StatusCreated, task)
}

func (h *TaskHandler) ListTasks(w http.ResponseWriter, r *http.Request) {

	user := middleware.GetUser(r)

	tasks := h.service.ListTasks(user)

	writeJSON(w, http.StatusOK, tasks)
}

func (h *TaskHandler) GetTask(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	task, err := h.service.GetTaskByID(id, user)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, services.ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())
		default:
			writeError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	writeJSON(w, http.StatusOK, task)
}

func (h *TaskHandler) AssignTask(w http.ResponseWriter, r *http.Request) {

	user := middleware.GetUser(r)

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req AssignTaskRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.service.AssignTask(id, req.WorkerID, user)

	if err != nil {

		switch {

		case errors.Is(err, services.ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())

		case errors.Is(err, services.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		case errors.Is(err, services.ErrUserNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "task assigned",
	})
}

func (h *TaskHandler) UpdateStatus(w http.ResponseWriter, r *http.Request) {

	user := middleware.GetUser(r)

	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		writeError(w, http.StatusBadRequest, "invalid task id")
		return
	}

	var req UpdateStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	err = h.service.UpdateStatus(id, req.Status, user)

	if err != nil {

		switch {

		case errors.Is(err, services.ErrForbidden):
			writeError(w, http.StatusForbidden, err.Error())

		case errors.Is(err, services.ErrTaskNotFound):
			writeError(w, http.StatusNotFound, err.Error())

		default:
			writeError(w, http.StatusBadRequest, err.Error())
		}

		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"message": "status updated",
	})
}

func (h *TaskHandler) GetNotifications(w http.ResponseWriter, r *http.Request) {
	user := middleware.GetUser(r)

	notifications := h.service.GetNotifications(user)

	writeJSON(w, http.StatusOK, notifications)
}
