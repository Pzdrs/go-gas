package main

import (
	"fmt"
	"github.com/Pzdrs/go-gas/config"
	"time"
)

func (register *CashRegister) Handle(vehicle *Vehicle) {
	fmt.Println(" > REGISTER", register.ID, " -> Handling vehicle", vehicle.ID)
	paymentDuration := getRandomDelayArr(config.GetCashRegisterConfiguration().Speed)
	vehicle.PaymentDuration = paymentDuration
	time.Sleep(paymentDuration)
	fmt.Println(" > REGISTER", register.ID, " -> Finished handling vehicle", vehicle.ID)
}

func handleRegister(register *CashRegister) {
	defer registerWg.Done()
	registerWg.Add(1)
	fmt.Printf("Cash register %d is open\n", register.ID)
	for car := range register.Queue {
		car.TimeSpentInRegisterQueue = time.Duration(time.Since(car.RegisterQueueEnter).Milliseconds())
		register.Handle(car)
		car.CarSync.Done()
		gasStation.ExitQueue <- car
	}
	fmt.Printf("Cash register %d is closed\n", register.ID)
}

// findRegister assigns cars to the shortest register queue
func findRegister(registers []*CashRegister) {
	for car := range gasStation.PaymentQueue {
		var bestRegister *CashRegister
		bestQueueLength := -1
		for _, register := range registers {
			queueLength := len(register.Queue)
			if bestQueueLength == -1 || queueLength < bestQueueLength {
				bestRegister = register
				bestQueueLength = queueLength
			}
		}
		car.RegisterQueueEnter = time.Now()
		bestRegister.Queue <- car
	}
	for _, register := range registers {
		close(register.Queue)
	}
}
