package vehicle_spawner

import (
	"fmt"
	"github.com/Pzdrs/go-gas/internal/config"
	"time"
)

type Spawner interface {
	Begin()
	Stop()
}

func (vs *vehicleSpawner) Begin() {
	fmt.Println("Vehicle spawner started")
}

func (vs *vehicleSpawner) Stop() {
	fmt.Println("Vehicle spawner stopped")
}

type vehicleSpawner struct {
	config config.VehicleSpawnerConfiguration
	ticker *time.Ticker
}

func NewVehicleSpawner(configuration config.VehicleSpawnerConfiguration) Spawner {
	return &vehicleSpawner{
		config: configuration,
		ticker: time.NewTicker(time.Duration(configuration.SpawnInterval) * time.Millisecond),
	}
}
