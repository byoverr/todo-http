package services

import (
	"context"
	"strings"
	"sync"
	"time"

	"github.com/byoverr/todo-http/internal/models"
)

type Service interface {
	Create(ctx context.Context, t models.Todo) (models.Todo, error)
	GetAll(ctx context.Context) ([]models.Todo, error)
	GetByID(ctx context.Context, id int) (models.Todo, error)
	Update(ctx context.Context, id int, t models.Todo) (models.Todo, error)
	Delete(ctx context.Context, id int) error
}

type Store struct {
	mu     sync.RWMutex
	todos  map[int]models.Todo
	nextID int
}

func NewService() *Store {
	return &Store{
		todos:  make(map[int]models.Todo),
		nextID: 1,
	}
}

var _ Service = (*Store)(nil)

func (s *Store) GetAll(ctx context.Context) ([]models.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	list := make([]models.Todo, 0, len(s.todos))
	for _, t := range s.todos {
		list = append(list, t)
	}
	return list, nil
}

func (s *Store) GetByID(ctx context.Context, id int) (models.Todo, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	t, ok := s.todos[id]
	if !ok {
		return models.Todo{}, ErrNotFound
	}
	return t, nil
}

func (s *Store) Create(ctx context.Context, t models.Todo) (models.Todo, error) {
	if strings.TrimSpace(t.Title) == "" {
		return models.Todo{}, ErrMissingTitle
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	t.ID = s.nextID
	t.Title = strings.TrimSpace(t.Title)
	t.CreatedAt = time.Now()
	t.UpdatedAt = t.CreatedAt

	s.todos[t.ID] = t
	s.nextID++

	return t, nil
}

func (s *Store) Update(ctx context.Context, id int, t models.Todo) (models.Todo, error) {
	if strings.TrimSpace(t.Title) == "" {
		return models.Todo{}, ErrMissingTitle
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	existing, ok := s.todos[id]
	if !ok {
		return models.Todo{}, ErrNotFound
	}

	existing.Title = strings.TrimSpace(t.Title)
	existing.Description = t.Description
	existing.IsCompleted = t.IsCompleted
	existing.UpdatedAt = time.Now()

	s.todos[id] = existing

	return existing, nil
}

func (s *Store) Delete(ctx context.Context, id int) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.todos[id]; !ok {
		return ErrNotFound
	}
	delete(s.todos, id)
	return nil
}
