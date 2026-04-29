package core_logger

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level  string `envconfig:"LEVEL" required:"true"`
	Folder string `envconfig:"FOLDER" required:"true"`
}

func NewLoggerConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process logger config: %w", err)
	}
	return config, nil
}

func NewConfigMust() Config {
	config, err := NewLoggerConfig()
	if err != nil {
		err := (fmt.Sprintf("failed to load logger config: %v", err))
		panic(err)
	}
	return config
}
