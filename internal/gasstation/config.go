package gasstation

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type StationConfiguration struct {
	VehicleSpawner   VehicleSpawnerConfig  `yaml:"vehicle-spawner"`
	PrometheusConfig PrometheusConfig      `yaml:"prometheus"`
	Pumps            map[string]PumpConfig `yaml:"pumps"`
	Registers        RegisterConfig        `yaml:"registers"`
}

type VehicleSpawnerConfig struct {
	Goal  int             `yaml:"goal"`
	Speed []time.Duration `yaml:"speed"`
}

type PumpConfig struct {
	Name   string          `yaml:"name"`
	Type   string          `yaml:"type"`
	Amount int             `yaml:"amount"`
	Speed  []time.Duration `yaml:"speed"`
}

type RegisterConfig struct {
	Amount int             `yaml:"amount"`
	Speed  []time.Duration `yaml:"speed"`
}

type PrometheusConfig struct {
	PushGateway string `yaml:"push-gateway"`
}

func loadConfig(file string) StationConfiguration {
	var config StationConfiguration

	err := cleanenv.ReadConfig(file, &config)
	if err != nil {
		fmt.Println("Error reading configuration:", err)
	}
	return config
}
