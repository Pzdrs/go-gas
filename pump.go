package main

import (
	"fmt"
	"time"
)

func (pump *Pump) Handle(vehicle *Vehicle) {
	fmt.Println(" > PUMP", pump.ID, " -> Handling vehicle", vehicle.ID)
	fuelingDuration := getRandomDelayArr(pump.Speed)
	vehicle.FuelingDuration = fuelingDuration
	time.Sleep(fuelingDuration)
	fmt.Println(" > PUMP", pump.ID, " -> Finished handling vehicle", vehicle.ID)
}

func handlePump(pump *Pump) {
	defer pumpWg.Done()
	pumpWg.Add(1)
	fmt.Println("Pump", pump.ID, "is open")
	for car := range pump.Queue {
		car.TimeSpentInPumpQueue = time.Duration(time.Since(car.PumpQueueEnter).Milliseconds())
		pump.Handle(car)
		car.CarSync.Add(1)
		gasStation.PaymentQueue <- car
		car.CarSync.Wait()
	}
	fmt.Println("Pump", pump.ID, "is closed")
}

// findPump assigns cars to the shortest queue of the correct fuel type
func findPump(pumps []*Pump) {
	for car := range gasStation.ArrivalQueue {

		var bestStand *Pump
		bestQueueLength := -1

		for _, stand := range pumps {
			if stand.Type == car.Fuel {
				queueLength := len(stand.Queue)
				if bestQueueLength == -1 || queueLength < bestQueueLength {
					bestStand = stand
					bestQueueLength = queueLength
				}
			}
		}
		bestStand.Queue <- car
	}
	for _, stand := range pumps {
		close(stand.Queue)
	}
}
