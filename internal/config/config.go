package config

import (
	"bufio"
	"os"
	"strings"
)

type Config struct {
	Port string
	DB   string
}

func MustLoadConfig() (*Config, error) {

	err := loadEnv(".env")
	if err != nil {
		return nil, err
	}

	return &Config{
		Port: getEnv("PORT", "8080"),
	}, nil
}

func getEnv(key, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultVal
}

func loadEnv(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if _, exists := os.LookupEnv(key); !exists {

			err = os.Setenv(key, value)
			if err != nil {
				return err
			}
		}
	}
	return scanner.Err()
}
