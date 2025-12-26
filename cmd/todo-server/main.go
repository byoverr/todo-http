package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/byoverr/todo-http/internal/app"
	"github.com/byoverr/todo-http/internal/config"
	loggerConstructor "github.com/byoverr/todo-http/pkg/logger"
)

func main() {

	logger := loggerConstructor.New("info")

	cfg, err := config.MustLoad()
	if err != nil {
		logger.Error("todo.config.MustLoad", "err", err)
	}

	app := app.New(logger, cfg)

	go func() {
		app.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.GracefulStop()
}
