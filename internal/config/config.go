package config

import (
	"encoding/json"
	"os"
)

type Config struct {
	Simulations []SimulationConfig `json:"simulations"`
}

func Load(path string) (Config, error) {
	cfg := Config{}
	file, err := os.ReadFile(path)

	if err != nil {
		return cfg, err
	}

	err = json.Unmarshal(file, &cfg)

	if err != nil {
		return cfg, err
	}

	return cfg, nil
}
