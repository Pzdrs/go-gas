package config

type PumpConfiguration struct {
	Type   string
	Name   string
	Amount int
	Speed  []int
}

type StationConfiguration struct {
	Pumps []PumpConfiguration
}

type VehicleSpawnerConfiguration struct {
	SpawnInterval int   `yaml:"vehicle-spawn-rate"`
	SpawnAmount   []int `yaml:"vehicle-spawn-amount"`
}
