package api

import (
	"bytes"
	"encoding/json"
	"go-final-project/internal/service"
	"go-final-project/internal/types"
	"log/slog"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	srv service.TaskService
}

func NewHandler(srv service.TaskService) *Handler {
	return &Handler{srv: srv}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/nextdate", h.handleNextDate)
	mux.HandleFunc("POST /api/task", h.handlePostTask)
	mux.HandleFunc("GET /api/tasks", h.handleGetTasks)
	mux.HandleFunc("GET /api/task", h.handleGetTaskById)
	mux.HandleFunc("PUT /api/task", h.handleUpdateTask)
}

func (h *Handler) handleNextDate(w http.ResponseWriter, r *http.Request) {
	nowValue := r.FormValue("now")
	dateValue := r.FormValue("date")
	repeat := r.FormValue("repeat")

	if len(repeat) == 0 {
		slog.Debug("empty repeat.")
		http.Error(w, `{"error":"empty repeat"}`, http.StatusBadRequest)
		return
	}

	next, err := h.srv.NextDate(nowValue, dateValue, repeat)
	if err != nil {
		slog.Debug("failed to get next date.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(next))
}

func (h *Handler) handlePostTask(w http.ResponseWriter, r *http.Request) {
	var task types.Task
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		slog.Error("failed to read body.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		slog.Error("failed to unmarshal body.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusBadRequest)
		return
	}

	if task.Date != "today" && task.Date != "" && task.Date != time.Now().Format("20060102") {
		task.Date, err = h.srv.NextDate(time.Now().Format("20060102"), task.Date, task.Repeat)
		if err != nil {
			slog.Debug("failed to get next date.", "err", err)
			http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
			return
		}
	} else {
		task.Date = time.Now().Format("20060102")
	}

	task.ID, err = h.srv.AddNewTask(task)
	if err != nil {
		slog.Error("failed to add new task.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
		return
	}

	stringId := strconv.Itoa(int(task.ID))
	resp, err := json.Marshal(types.ResponseOK{ID: stringId})
	if err != nil {
		slog.Error("failed to marshal id.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		slog.Error("failed to write the response.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
	}
}

func (h *Handler) handleGetTasks(w http.ResponseWriter, r *http.Request) {
	tasks, err := h.srv.GetTasks()
	if err != nil {
		slog.Error("failed to get tasks.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(tasks)
	if err != nil {
		slog.Error("failed to marshal tasks.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		slog.Error("failed to write the response.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
	}
}

func (h *Handler) handleGetTaskById(w http.ResponseWriter, r *http.Request) {
	taskId := r.FormValue("id")
	task, err := h.srv.GetTaskById(taskId)
	if err != nil {
		slog.Error("failed to get task by id.", "id", taskId, "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
		return
	}

	resp, err := json.Marshal(task)
	if err != nil {
		slog.Error("failed to marshal task.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(resp)
	if err != nil {
		slog.Error("failed to write the response.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
	}
}

func (h *Handler) handleUpdateTask(w http.ResponseWriter, r *http.Request) {
	var task types.TaskDTO
	var buf bytes.Buffer

	_, err := buf.ReadFrom(r.Body)
	if err != nil {
		slog.Error("failed to read body.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusBadRequest)
		return
	}

	if err = json.Unmarshal(buf.Bytes(), &task); err != nil {
		slog.Error("failed to unmarshal body.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusBadRequest)
		return
	}

	err = h.srv.UpdateTask(task)
	if err != nil {
		slog.Error("failed to update task.", "err", err)
		http.Error(w, `{"error":"oops, something went wrong"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application-json; charset=UTF-8")
	w.Write([]byte(`{}`))
}
