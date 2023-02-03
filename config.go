package main

import (
	"encoding/json"
	"fmt"
	"os"
)

// Loads Configuration File
func (c *Config) LoadConfig(confFile string) (err error) {
	conf, err := os.ReadFile(confFile)

	if err != nil {
		return fmt.Errorf("Error while opening configuration: %v", err)
	}

	var m []Metric
	err = json.Unmarshal(conf, &m)

	if err != nil {
		return fmt.Errorf("Error during Unmarshal: %v", err)
	}

	c.Metrics = m

	return nil
}
