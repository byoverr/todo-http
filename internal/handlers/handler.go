package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/byoverr/todo-http/internal/models"
	"github.com/byoverr/todo-http/internal/services"
)

type ServerAPI struct {
	svc services.Service
	log *slog.Logger
}

var _ ServerInterface = (*ServerAPI)(nil)

type ServerInterface interface {
	GetTodos(ctx context.Context, w http.ResponseWriter, r *http.Request)
	GetTodoByID(ctx context.Context, w http.ResponseWriter, r *http.Request, id int)
	CreateTodo(ctx context.Context, w http.ResponseWriter, r *http.Request)
	UpdateTodo(ctx context.Context, w http.ResponseWriter, r *http.Request, id int)
	DeleteTodo(ctx context.Context, w http.ResponseWriter, r *http.Request, id int)
}

func New(svc services.Service, log *slog.Logger) *ServerAPI {
	return &ServerAPI{svc: svc, log: log}
}

// GET /todos — получить список всех задач
func (s *ServerAPI) GetTodos(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	const op = "todo.handlers.GetTodos"

	log := s.log.With(slog.String("op", op))
	w.Header().Set("Content-Type", "application/json")

	todos, err := s.svc.GetAll(ctx)
	if err != nil {
		log.Error("failed to get todos", slog.String("error", err.Error()))
		writeJSONError(w, http.StatusInternalServerError, "failed to retrieve todos")
		return
	}

	if err := json.NewEncoder(w).Encode(todos); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
		writeJSONError(w, http.StatusInternalServerError, "failed to encode response")
		return
	}

	log.Info("GetTodos success", slog.Int("count", len(todos)))
}

// GET /todos/{id} — получить задачу по идентификатору
func (s *ServerAPI) GetTodoByID(ctx context.Context, w http.ResponseWriter, r *http.Request, id int) {
	const op = "todo.handlers.GetTodoByID"

	log := s.log.With(slog.String("op", op))
	w.Header().Set("Content-Type", "application/json")

	todo, err := s.svc.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, services.ErrNotFound) {
			log.Warn("todo not found", slog.Int("id", id))
			writeJSONError(w, http.StatusNotFound, "todo not found")
			return
		}

		log.Error("failed to get todo", slog.String("error", err.Error()))
		writeJSONError(w, http.StatusInternalServerError, "failed to retrieve todo")
		return
	}

	if err = json.NewEncoder(w).Encode(todo); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
		writeJSONError(w, http.StatusInternalServerError, "failed to encode response")
		return
	}

	log.Info("GetTodoByID success", slog.Int("id", id))
}

// POST /todos — создать новую задачу
func (s *ServerAPI) CreateTodo(ctx context.Context, w http.ResponseWriter, r *http.Request) {

	const op = "todo.handlers.CreateTodo"

	log := s.log.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Warn("invalid json", slog.String("error", err.Error()))
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	created, err := s.svc.Create(ctx, t)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMissingTitle):
			log.Warn("missing title", slog.String("error", err.Error()))
			writeJSONError(w, http.StatusBadRequest, err.Error())

		default:
			log.Error("failed to create todo", slog.String("error", err.Error()))
			writeJSONError(w, http.StatusInternalServerError, "failed to create todo")
		}
		return
	}

	w.WriteHeader(http.StatusCreated)
	if err = json.NewEncoder(w).Encode(created); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
		writeJSONError(w, http.StatusInternalServerError, "failed to encode response")
		return
	}

	log.Info("CreateTodo success", slog.Int("id", created.ID))
}

// PUT /todos/{id} — обновить задачу по идентификатору
func (s *ServerAPI) UpdateTodo(ctx context.Context, w http.ResponseWriter, r *http.Request, id int) {

	const op = "todo.handlers.UpdateTodo"

	log := s.log.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Warn("invalid json", slog.String("error", err.Error()))
		writeJSONError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	updated, err := s.svc.Update(ctx, id, t)
	if err != nil {
		switch {
		case errors.Is(err, services.ErrMissingTitle):
			log.Warn("missing title", slog.String("error", err.Error()), slog.Int("id", id))
			writeJSONError(w, http.StatusBadRequest, err.Error())

		case errors.Is(err, services.ErrNotFound):
			log.Warn("todo not found", slog.String("error", err.Error()), slog.Int("id", id))
			writeJSONError(w, http.StatusNotFound, "todo not found")

		default:
			log.Error("failed to update todo", slog.String("error", err.Error()), slog.Int("id", id))
			writeJSONError(w, http.StatusInternalServerError, "failed to update todo")
		}
		return
	}

	if err = json.NewEncoder(w).Encode(updated); err != nil {
		log.Error("failed to encode response", slog.String("error", err.Error()))
		writeJSONError(w, http.StatusInternalServerError, "failed to encode response")
		return
	}

	log.Info("UpdateTodo success", slog.Int("id", id))
}

// DELETE /todos/{id} — удалить задачу по идентификатору
func (s *ServerAPI) DeleteTodo(ctx context.Context, w http.ResponseWriter, r *http.Request, id int) {

	const op = "todo.handlers.DeleteTodo"

	log := s.log.With(slog.String("op", op))

	w.Header().Set("Content-Type", "application/json")

	if err := s.svc.Delete(ctx, id); err != nil {
		switch {
		case errors.Is(err, services.ErrNotFound):
			log.Warn("todo not found", slog.String("error", err.Error()), slog.Int("id", id))
			writeJSONError(w, http.StatusNotFound, "todo not found")

		default:
			log.Error("failed to delete todo", slog.String("error", err.Error()), slog.Int("id", id))
			writeJSONError(w, http.StatusInternalServerError, "failed to delete todo")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
	log.Info("DeleteTodo success", slog.Int("id", id))
}

func writeJSONError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{
		"error": message,
	})
}
