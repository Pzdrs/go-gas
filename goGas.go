package main

import (
	"github.com/Pzdrs/go-gas/internal/gasstation"
)

func main() {
	gasStation := gasstation.NewGasStation()

	gasStation.Inspect()

	gasStation.Setup()

	gasStation.Begin(10000)
}
