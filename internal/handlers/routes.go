package handlers

import (
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func Register(mux *http.ServeMux, api *ServerAPI) {
	mux.HandleFunc("/health", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, err := w.Write([]byte("ok"))
		if err != nil {
			api.log.Warn("failed to write health response", slog.String("error", err.Error()))
		}
	})
	mux.HandleFunc("GET /todos", func(w http.ResponseWriter, r *http.Request) {
		api.GetTodos(r.Context(), w, r)
	})

	mux.HandleFunc("POST /todos", func(w http.ResponseWriter, r *http.Request) {
		api.CreateTodo(r.Context(), w, r)
	})

	mux.HandleFunc("GET /todos/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
		id, _ := strconv.Atoi(idStr)
		api.GetTodoByID(r.Context(), w, r, id)
	})

	mux.HandleFunc("PUT /todos/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
		id, _ := strconv.Atoi(idStr)
		api.UpdateTodo(r.Context(), w, r, id)
	})

	mux.HandleFunc("DELETE /todos/", func(w http.ResponseWriter, r *http.Request) {
		idStr := strings.TrimPrefix(r.URL.Path, "/todos/")
		id, _ := strconv.Atoi(idStr)
		api.DeleteTodo(r.Context(), w, r, id)
	})
}
