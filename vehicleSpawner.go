package main

import (
	"fmt"
	"github.com/Pzdrs/go-gas/config"
	"math/rand"
	"sync"
	"time"
)

func spawnVehicles() {
	goal := config.GetVehicleSpawnerConfiguration().Goal

	fmt.Println(" >> VEHICLE SPAWNER -> Started spawning vehicles")
	for i := 0; i < goal; i++ {
		vehicle := Vehicle{ID: i, Fuel: genFuelType(), CarSync: &sync.WaitGroup{}, PumpQueueEnter: time.Now()}
		gasStation.ArrivalQueue <- &vehicle
		time.Sleep(getRandomDelayArr(config.GetVehicleSpawnerConfiguration().Rate))
	}
	fmt.Println(" >> VEHICLE SPAWNER -> Finished spawning", goal, "vehicles")
	close(gasStation.ArrivalQueue)
}

// genFuelType returns a random fuel type.
func genFuelType() string {
	fuelTypes := []string{Gas, Diesel, Electric, LPG}
	randomIndex := rand.Intn(len(fuelTypes))
	return fuelTypes[randomIndex]
}
