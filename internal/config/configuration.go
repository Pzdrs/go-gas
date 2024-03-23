package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"time"
)

type GasStationConfig struct {
	VehicleSpawner VehicleSpawnerConfig  `yaml:"vehicle-spawner"`
	Pumps          map[string]PumpConfig `yaml:"pumps"`
	Registers      RegisterConfig        `yaml:"registers"`
}

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

type InfluxConfig struct {
	Url    string `yaml:"influx.url"`
	Token  string `yaml:"influx.token"`
	Org    string `yaml:"influx.org"`
	Bucket string `yaml:"influx.bucket"`
}

var GasStationConfiguration GasStationConfig
var InfluxFluxConfiguration InfluxConfig

func LoadConfig(file string) {
	loadGasStationConfig(file)
	loadInfluxConfig(file)
}

func loadInfluxConfig(file string) {
	err := cleanenv.ReadConfig(file, &InfluxFluxConfiguration)
	if err != nil {
		fmt.Println("Error reading configuration:", err)
	}
}

func loadGasStationConfig(file string) {
	err := cleanenv.ReadConfig(file, &GasStationConfiguration)
	if err != nil {
		fmt.Println("Error reading configuration:", err)
	}
}
