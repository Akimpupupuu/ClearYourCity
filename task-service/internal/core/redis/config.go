package core_redis

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Addr     string `envconfig:"ADDR" default:"localhost:6379"`
	Password string `envconfig:"PASSWORD" default:""`
	DB       int    `envconfig:"DB" default:"0"`
}

func NewConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("TASK_REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("process redis config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("failed to get redis config: %w", err)
		panic(err)
	}

	return config
}
