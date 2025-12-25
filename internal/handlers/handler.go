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

	log := s.log.With(
		slog.String("op", op),
	)

	todos, err := s.svc.GetAll(ctx)
	if err != nil {
		log.Error("err", err.Error())
	}
	err = json.NewEncoder(w).Encode(todos)
	if err != nil {
		log.Error("err", err.Error())
	}
	log.Info("GetTodos success")
}

// GET /todos/{id} — получить задачу по идентификатору
func (s *ServerAPI) GetTodoByID(ctx context.Context, w http.ResponseWriter, r *http.Request, id int) {

	const op = "todo.handlers.GetTodoByID"

	log := s.log.With(
		slog.String("op", op),
	)

	todo, err := s.svc.GetByID(ctx, id)
	if err != nil {
		log.Error("err", err.Error())
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	err = json.NewEncoder(w).Encode(todo)
	if err != nil {
		log.Error("err", err.Error())
	}
	log.Info("GetTodoByID success")
}

// POST /todos — создать новую задачу
func (s *ServerAPI) CreateTodo(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	var t models.Todo

	const op = "todo.handlers.CreateTodo"

	log := s.log.With(
		slog.String("op", op),
	)

	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		log.Error("err", err.Error())
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	created, err := s.svc.Create(ctx, t)
	if errors.Is(err, services.ErrMissingTitle) {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(created)
	if err != nil {
		log.Error("err", err.Error())
	}

	log.Info("CreateTodo success")
}

// PUT /todos/{id} — обновить задачу по идентификатору
func (s *ServerAPI) UpdateTodo(ctx context.Context, w http.ResponseWriter, r *http.Request, id int) {

	const op = "todo.handlers.UpdateTodo"

	log := s.log.With(
		slog.String("op", op),
	)

	var t models.Todo
	if err := json.NewDecoder(r.Body).Decode(&t); err != nil {
		http.Error(w, "invalid json", http.StatusBadRequest)
		return
	}
	updated, err := s.svc.Update(ctx, id, t)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, services.ErrMissingTitle) {
			status = http.StatusBadRequest
		} else if errors.Is(err, services.ErrNotFound) {
			status = http.StatusNotFound
		}
		log.Error("err", err.Error())
		http.Error(w, err.Error(), status)
		return
	}
	err = json.NewEncoder(w).Encode(updated)
	if err != nil {
		log.Error("err", err.Error())
	}
	log.Info("UpdateTodo success")
}

// DELETE /todos/{id} — удалить задачу по идентификатору
func (s *ServerAPI) DeleteTodo(ctx context.Context, w http.ResponseWriter, r *http.Request, id int) {

	const op = "todo.handlers.DeleteTodo"

	log := s.log.With(
		slog.String("op", op),
	)

	err := s.svc.Delete(ctx, id)
	if errors.Is(err, services.ErrNotFound) {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusNoContent)

	log.Info("DeleteTodo success")
}
