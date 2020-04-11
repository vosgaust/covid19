package main

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
	"github.com/vosgaust/covid19/entries"
)

type config struct {
	MySQL entries.Config
}

// Get returns the configuration from the environment variables
func getConfig() (config, error) {
	var cfg config
	if err := envconfig.Process("covid19", &cfg); err != nil {
		return cfg, fmt.Errorf("read configuration from env: %v", err)
	}
	return cfg, nil
}
