package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type VehicleSpawnerConfig struct {
	Goal int             `yaml:"goal"`
	Rate []time.Duration `yaml:"rate"`
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

type GasStationConfig struct {
	VehicleSpawner VehicleSpawnerConfig  `yaml:"vehicle-spawner"`
	Pumps          map[string]PumpConfig `yaml:"pumps"`
	Registers      RegisterConfig        `yaml:"registers"`
}

var GasStationConfiguration GasStationConfig

func LoadConfig(file string) {
	err := cleanenv.ReadConfig(file, &GasStationConfiguration)
	if err != nil {
		fmt.Println("Error reading configuration:", err)
	}
}
