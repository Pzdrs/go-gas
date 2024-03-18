package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type PumpConfiguration struct {
	Type   string
	Name   string
	Amount int
	Speed  []int
}

type Configuration struct {
	vehicleSpawnerConfiguration VehicleSpawnerConfig
}

type VehicleSpawnerConfig struct {
	SpawnInterval int   `yaml:"vehicle-spawn-rate"`
	SpawnAmount   []int `yaml:"vehicle-spawn-amount"`
}

type RegisterConfig struct {
	Amount int   `yaml:"registers.amount"`
	Speed  []int `yaml:"registers.speed"`
}

var RegisterConfiguration RegisterConfig
var VehicleSpawnerConfiguration VehicleSpawnerConfig

func LoadConfig() {
	loadRegisterConfig()
	loadVehicleSpawnerConfig()
}

func loadRegisterConfig() {
	err := cleanenv.ReadConfig("config.yaml", &RegisterConfiguration)
	if err != nil {
		fmt.Println("Error reading configuration:", err)
	}
}

func loadVehicleSpawnerConfig() {
	err := cleanenv.ReadConfig("config.yaml", &VehicleSpawnerConfiguration)
	if err != nil {
		fmt.Println("Error reading configuration:", err)
	}
}
