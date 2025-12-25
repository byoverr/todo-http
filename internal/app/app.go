package app

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/byoverr/todo-http/internal/config"
	"github.com/byoverr/todo-http/internal/handlers"
	"github.com/byoverr/todo-http/internal/handlers/middleware"
	"github.com/byoverr/todo-http/internal/services"
)

type App struct {
	server *http.Server
	log    *slog.Logger
	cfg    *config.Config
}

func New(
	log *slog.Logger,
	cfg *config.Config,
) App {

	mux := http.NewServeMux()

	service := services.NewService()
	api := handlers.New(service, log)

	handlers.Register(mux, api)

	server := &http.Server{
		Addr:         address(cfg.Host, cfg.Port),
		Handler:      middleware.LoggingMiddleware(log, mux),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	return App{
		server: server,
		log:    log,
		cfg:    cfg,
	}
}

func (a App) MustRun() {
	a.log.Info("Server started", slog.String("op", "todo.app.MustRun"), slog.String("port", a.cfg.Port))

	if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		a.log.Info("ListenAndServe error", slog.String("op", "todo.app.MustRun"), slog.String("err", err.Error()))
		panic(err)
	}
}

func (a App) GracefulStop() {
	// здесь может закрываться пул соединений и тд, например у postgres, но у нас его нет
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		a.log.Error("Graceful shutdown failed", slog.String("op", "todo.app.GracefulStop"), slog.String("error", err.Error()))
	} else {
		a.log.Info("Gracefully stopped", slog.String("op", "todo.app.GracefulStop"))
	}
}

func address(host string, port string) string {
	return net.JoinHostPort(host, port)
}
