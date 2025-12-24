package main

import (
	"github.com/byoverr/todo-http/internal/config"
	loggerConstructor "github.com/byoverr/todo-http/pkg/logger"
)

func main() {

	logger := loggerConstructor.New("info")

	conf, err := config.MustLoadConfig()
	if err != nil {
		logger.Error("todo.config", "err", err)
	}

	logger.Info("starting server on port", conf.Port)
}
