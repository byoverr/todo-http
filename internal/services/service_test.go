package services

import (
	"context"
	"errors"
	"testing"

	"github.com/byoverr/todo-http/internal/models"
)

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	svc := NewService()

	tests := []struct {
		name      string
		input     models.Todo
		wantErr   error
		wantEmpty bool
	}{
		{
			name:    "valid todo",
			input:   models.Todo{Title: "Example"},
			wantErr: nil,
		},
		{
			name:      "missing title",
			input:     models.Todo{Title: ""},
			wantErr:   ErrMissingTitle,
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			todo, err := svc.Create(ctx, tt.input)

			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected error %v, got %v", tt.wantErr, err)
			}
			if err == nil && todo.Title != tt.input.Title {
				t.Errorf("expected title %q, got %q", tt.input.Title, todo.Title)
			}
			if tt.wantEmpty && todo.ID != 0 {
				t.Errorf("expected empty result, got ID=%d", todo.ID)
			}
		})
	}
}

func TestService_GetAll(t *testing.T) {
	ctx := context.Background()

	svc := NewService()

	// тестировали create раньше, доверяем ему
	_, _ = svc.Create(ctx, models.Todo{Title: "Task 1"})
	_, _ = svc.Create(ctx, models.Todo{Title: "Task 2"})

	todos, err := svc.GetAll(ctx)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(todos) != 2 {
		t.Fatalf("expected 2 todos, got %d", len(todos))
	}

}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	svc := NewService()

	created, _ := svc.Create(ctx, models.Todo{Title: "Test task"})

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name:    "found",
			id:      created.ID,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      123,
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := svc.GetByID(ctx, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()

	svc := NewService()

	created, _ := svc.Create(ctx, models.Todo{Title: "Original"})

	tests := []struct {
		name    string
		id      int
		update  models.Todo
		wantErr error
	}{
		{
			name:    "success",
			id:      created.ID,
			update:  models.Todo{Title: "Updated title"},
			wantErr: nil,
		},
		{
			name:    "missing title",
			id:      created.ID,
			update:  models.Todo{Title: ""},
			wantErr: ErrMissingTitle,
		},
		{
			name:    "not found",
			id:      123,
			update:  models.Todo{Title: "Not exist"},
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			todo, err := svc.Update(ctx, tt.id, tt.update)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
			if err == nil && todo.Title != tt.update.Title {
				t.Errorf("expected title %q, got %q", tt.update.Title, todo.Title)
			}
		})
	}
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()

	svc := NewService()

	created, _ := svc.Create(ctx, models.Todo{Title: "Delete task"})

	tests := []struct {
		name    string
		id      int
		wantErr error
	}{
		{
			name:    "success",
			id:      created.ID,
			wantErr: nil,
		},
		{
			name:    "not found",
			id:      123,
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := svc.Delete(ctx, tt.id)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("expected %v, got %v", tt.wantErr, err)
			}
		})
	}
}
