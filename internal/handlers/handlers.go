package handlers

import (
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"welcomeProject/internal/models"
	"welcomeProject/internal/repository"
	"welcomeProject/pkg/utils"
)

type Handlers struct {
	repo *repository.Repository
}

func NewHandlers(repo *repository.Repository) *Handlers {
	return &Handlers{repo: repo}
}

func (h *Handlers) HelloHandler(w http.ResponseWriter, r *http.Request) {
	utils.JSONResponse(w, http.StatusOK, map[string]string{"hello": "world"})
}

func (h *Handlers) CreateTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if task.Description == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Description is required")
		return
	}

	if task.Status == "" {
		task.Status = "new"
	}

	if !isValidStatus(task.Status) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid status")
		return
	}

	if err := h.repo.CreateTask(&task); err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusCreated, task)
}

func (h *Handlers) GetTasksHandler(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.repo.GetAllTasks()
	if err != nil {
		utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		return
	}

	utils.JSONResponse(w, http.StatusOK, tasks)
}

func (h *Handlers) GetTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	task, err := h.repo.GetTaskByID(id)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, "Task not found")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.JSONResponse(w, http.StatusOK, task)
}

func (h *Handlers) UpdateTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var task models.Task
	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	task.ID = id

	if task.Description == "" {
		utils.ErrorResponse(w, http.StatusBadRequest, "Description is required")
		return
	}

	if task.Status != "" && !isValidStatus(task.Status) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid status")
		return
	}

	if err := h.repo.UpdateTask(&task); err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, "Task not found")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.JSONResponse(w, http.StatusOK, task)
}

func (h *Handlers) DeleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	if err := h.repo.DeleteTask(id); err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, "Task not found")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *Handlers) UpdateTaskStatusHandler(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid task ID")
		return
	}

	var statusUpdate models.StatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&statusUpdate); err != nil {
		utils.ErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	if !isValidStatus(statusUpdate.Status) {
		utils.ErrorResponse(w, http.StatusBadRequest, "Invalid status")
		return
	}

	task, err := h.repo.UpdateTaskStatus(id, statusUpdate.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			utils.ErrorResponse(w, http.StatusNotFound, "Task not found")
		} else {
			utils.ErrorResponse(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	utils.JSONResponse(w, http.StatusOK, task)
}

func isValidStatus(status string) bool {
	return status == "new" || status == "in progress" || status == "done"
}
