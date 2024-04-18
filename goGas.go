package main

import (
	"github.com/Pzdrs/go-gas/internal/gasstation"
)

func main() {
	gasStation := gasstation.NewGasStation("config.yaml")

	gasStation.Inspect()

	gasStation.Setup()

	gasStation.Begin()
}
