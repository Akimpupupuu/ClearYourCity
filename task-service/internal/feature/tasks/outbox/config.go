package tasks_outbox

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	BatchLimit int           `envconfig:"LIMIT" default:"20"`
	Interval   time.Duration `envconfig:"INTERVAL" default:"500ms"`
}

func NewConfig() (Config, error) {
	var config Config

	if err := envconfig.Process("TASK_OUTBOX", &config); err != nil {
		return Config{}, fmt.Errorf("process config: %w", err)
	}

	return config, nil
}

func NewConfigMust() Config {
	config, err := NewConfig()
	if err != nil {
		err = fmt.Errorf("get outbox config: %w", err)
		panic(err)
	}

	return config
}
