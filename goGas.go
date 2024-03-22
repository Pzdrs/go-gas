package main

import (
	"github.com/Pzdrs/go-gas/config"
	"sync"
)

var pumpWg sync.WaitGroup
var registerWg sync.WaitGroup
var finishWg sync.WaitGroup

var gasStation GasStation

func main() {
	config.LoadConfig()

	gasStation = NewGasStation()
	gasStation.Inspect()

	finishWg.Add(1)
	go spawnVehicles()

	for _, pump := range gasStation.Pumps {
		go handlePump(pump)
	}

	for _, register := range gasStation.Registers {
		go handleRegister(register)
	}

	go findPump(gasStation.Pumps)
	go findRegister(gasStation.Registers)

	go aggregate(gasStation.ExitQueue)

	pumpWg.Wait()
	close(gasStation.PaymentQueue)

	registerWg.Wait()
	close(gasStation.ExitQueue)

	finishWg.Wait()
}
