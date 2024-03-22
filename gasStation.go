package main

import (
	"fmt"
	"github.com/Pzdrs/go-gas/config"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func NewGasStation() GasStation {
	gasStation := GasStation{
		ArrivalQueue: make(chan *Vehicle, 20),
		ExitQueue:    make(chan *Vehicle),
		PaymentQueue: make(chan *Vehicle, 10),
	}

	for _, pumpConfig := range config.GetPumpConfiguration() {
		for i := range pumpConfig.Amount {
			gasStation.AddPump(&Pump{
				ID:    strings.ToLower(pumpConfig.Type) + strconv.Itoa(i),
				Name:  pumpConfig.Name,
				Type:  pumpConfig.Type,
				Speed: pumpConfig.Speed,
				Queue: make(chan *Vehicle, 2),
			})
		}
	}

	for i := range config.GetCashRegisterConfiguration().Amount {
		gasStation.AddRegister(&CashRegister{
			ID:    i,
			Queue: make(chan *Vehicle, 3),
		})
	}

	return gasStation
}

func (gasStation *GasStation) Inspect() {
	fmt.Println("Gas Station")
	fmt.Println("===========")
	fmt.Println("Pumps:")
	for _, pump := range gasStation.Pumps {
		fmt.Println(" - Pump", pump.ID, "(", pump.Name, ")")
		fmt.Println("   - Type:", pump.Type)
	}

	fmt.Println("Registers:")
	for _, register := range gasStation.Registers {
		fmt.Println(" - Register", register.ID)
	}

	fmt.Println()
}

func (gasStation *GasStation) AddPump(pump *Pump) {
	gasStation.Pumps = append(gasStation.Pumps, pump)
}

func (gasStation *GasStation) AddRegister(register *CashRegister) {
	gasStation.Registers = append(gasStation.Registers, register)
}

func getRandomDelay(low time.Duration, high time.Duration) time.Duration {
	speedMinNano := int64(low)
	speedMaxNano := int64(high)

	randomNano := speedMinNano + rand.Int63n(speedMaxNano-speedMinNano+1)
	randomDuration := time.Duration(randomNano)

	return randomDuration
}

func getRandomDelayArr(durationRange []time.Duration) time.Duration {
	return getRandomDelay(durationRange[0], durationRange[1])
}
