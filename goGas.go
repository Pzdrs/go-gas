package main

import (
	"github.com/Pzdrs/go-gas/internal/config"
	"github.com/Pzdrs/go-gas/internal/gasstation"
)

func main() {
	config.LoadConfig("config.yaml")

	gasStation := gasstation.NewGasStation(config.GasStationConfiguration)

	gasStation.Inspect()

	gasStation.Setup()

	gasStation.Begin(10)
}
