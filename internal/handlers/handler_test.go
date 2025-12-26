package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/byoverr/todo-http/internal/models"
	"github.com/byoverr/todo-http/internal/services"
)

// Мок
type mockService struct {
	CreateFunc  func(ctx context.Context, t models.Todo) (models.Todo, error)
	GetAllFunc  func(ctx context.Context) ([]models.Todo, error)
	GetByIDFunc func(ctx context.Context, id int) (models.Todo, error)
	UpdateFunc  func(ctx context.Context, id int, t models.Todo) (models.Todo, error)
	DeleteFunc  func(ctx context.Context, id int) error
}

func (m *mockService) Create(ctx context.Context, t models.Todo) (models.Todo, error) {
	return m.CreateFunc(ctx, t)
}
func (m *mockService) GetAll(ctx context.Context) ([]models.Todo, error) {
	return m.GetAllFunc(ctx)
}
func (m *mockService) GetByID(ctx context.Context, id int) (models.Todo, error) {
	return m.GetByIDFunc(ctx, id)
}
func (m *mockService) Update(ctx context.Context, id int, t models.Todo) (models.Todo, error) {
	return m.UpdateFunc(ctx, id, t)
}
func (m *mockService) Delete(ctx context.Context, id int) error {
	return m.DeleteFunc(ctx, id)
}

func TestServerAPI_GetTodos(t *testing.T) {
	mock := &mockService{
		GetAllFunc: func(ctx context.Context) ([]models.Todo, error) {
			return []models.Todo{
				{ID: 1, Title: "Test"},
			}, nil
		},
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	api := New(mock, log)

	req := httptest.NewRequest(http.MethodGet, "/todos", nil)
	w := httptest.NewRecorder()

	api.GetTodos(req.Context(), w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status 200, got %d", resp.StatusCode)
	}

	var todos []models.Todo
	if err := json.NewDecoder(resp.Body).Decode(&todos); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(todos) != 1 || todos[0].ID != 1 {
		t.Fatalf("unexpected todos: %+v", todos)
	}
}

func TestServerAPI_GetTodoByID(t *testing.T) {
	mock := &mockService{
		GetByIDFunc: func(ctx context.Context, id int) (models.Todo, error) {
			if id == 1 {
				return models.Todo{ID: 1, Title: "Test"}, nil
			}
			return models.Todo{}, services.ErrNotFound
		},
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	api := New(mock, log)

	tests := []struct {
		id         int
		wantStatus int
	}{
		{1, http.StatusOK},
		{2, http.StatusNotFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/todos/%d", tt.id), nil)
		w := httptest.NewRecorder()

		api.GetTodoByID(req.Context(), w, req, tt.id)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != tt.wantStatus {
			t.Fatalf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
		}
	}
}

func TestServerAPI_CreateTodo(t *testing.T) {
	mock := &mockService{
		CreateFunc: func(ctx context.Context, t models.Todo) (models.Todo, error) {
			if t.Title == "" {
				return models.Todo{}, services.ErrMissingTitle
			}
			t.ID = 1
			return t, nil
		},
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	api := New(mock, log)

	tests := []struct {
		name       string
		body       string
		wantStatus int
	}{
		{"valid", `{"title":"Test"}`, http.StatusCreated},
		{"empty title", `{"title":""}`, http.StatusBadRequest},
		{"invalid json", `{"title":`, http.StatusBadRequest},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString(tt.body))
			w := httptest.NewRecorder()

			api.CreateTodo(req.Context(), w, req)

			resp := w.Result()
			defer resp.Body.Close()

			if resp.StatusCode != tt.wantStatus {
				t.Fatalf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
			}
		})
	}
}

func TestServerAPI_UpdateTodo(t *testing.T) {
	mock := &mockService{
		UpdateFunc: func(ctx context.Context, id int, t models.Todo) (models.Todo, error) {
			switch {
			case t.Title == "":
				return models.Todo{}, services.ErrMissingTitle
			case id != 1:
				return models.Todo{}, services.ErrNotFound
			default:
				t.ID = id
				return t, nil
			}
		},
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	api := New(mock, log)

	tests := []struct {
		id         int
		body       string
		wantStatus int
	}{
		{1, `{"title":"Updated"}`, http.StatusOK},
		{1, `{"title":""}`, http.StatusBadRequest},
		{2, `{"title":"X"}`, http.StatusNotFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/todos/%d", tt.id), bytes.NewBufferString(tt.body))
		w := httptest.NewRecorder()

		api.UpdateTodo(req.Context(), w, req, tt.id)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != tt.wantStatus {
			t.Fatalf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
		}
	}
}

func TestServerAPI_DeleteTodo(t *testing.T) {
	mock := &mockService{
		DeleteFunc: func(ctx context.Context, id int) error {
			if id != 1 {
				return services.ErrNotFound
			}
			return nil
		},
	}
	log := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	api := New(mock, log)

	tests := []struct {
		id         int
		wantStatus int
	}{
		{1, http.StatusNoContent},
		{2, http.StatusNotFound},
	}

	for _, tt := range tests {
		req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/todos/%d", tt.id), nil)
		w := httptest.NewRecorder()

		api.DeleteTodo(req.Context(), w, req, tt.id)

		resp := w.Result()
		defer resp.Body.Close()

		if resp.StatusCode != tt.wantStatus {
			t.Fatalf("expected status %d, got %d", tt.wantStatus, resp.StatusCode)
		}
	}
}
