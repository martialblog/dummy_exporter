package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Metrics []Metric
}

// LoadConfig loads the configuration file
func (c *Config) LoadConfig(confFile string) (err error) {
	conf, err := os.ReadFile(confFile)
	if err != nil {
		return fmt.Errorf("error while opening configuration: %w", err)
	}

	var m []Metric

	err = json.Unmarshal(conf, &m)
	if err != nil {
		return fmt.Errorf("error during Unmarshal: %w", err)
	}

	c.Metrics = m

	return nil
}
