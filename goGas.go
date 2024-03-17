package main

import (
	"fmt"
	"github.com/Pzdrs/go-gas/internal/config"
	vehiclespawner "github.com/Pzdrs/go-gas/internal/vehicle-spawner"
	"github.com/ilyakaznacheev/cleanenv"
)

func main() {
	var vehicleSpawnerConfiguration config.VehicleSpawnerConfiguration

	err := cleanenv.ReadConfig("config.yaml", &vehicleSpawnerConfiguration)
	if err != nil {
		fmt.Println("Error reading configuration:", err)
	}

	vehicleSpawner := vehiclespawner.NewVehicleSpawner(vehicleSpawnerConfiguration)
	vehicleSpawner.Begin()
	vehicleSpawner.Stop()
}
