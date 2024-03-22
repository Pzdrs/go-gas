package main

import (
	"fmt"
	"time"
)

type Stats struct {
	TotalCars        int
	TotalTime        time.Duration
	TotalQueueTime   time.Duration
	MaxQueueTime     time.Duration
	AverageQueueTime time.Duration
}

func (stats *Stats) Print() {
	fmt.Println("  Total cars: ", stats.TotalCars)
	fmt.Println("  Total time: ", stats.TotalTime)
	fmt.Println("  Total queue time: ", stats.TotalQueueTime)
	fmt.Println("  Max queue time: ", stats.MaxQueueTime)
	fmt.Println("  Average queue time: ", stats.AverageQueueTime)
}

func aggregate(exitQueue chan *Vehicle) {
	registerStats := Stats{}
	pumpStats := map[string]*Stats{
		Gas:      {},
		Diesel:   {},
		LPG:      {},
		Electric: {},
	}

	for car := range exitQueue {
		// Register stats collection
		registerStats.TotalCars++
		registerStats.TotalTime += car.PaymentDuration
		registerStats.TotalQueueTime += car.TimeSpentInRegisterQueue
		if car.TimeSpentInRegisterQueue > registerStats.MaxQueueTime {
			registerStats.MaxQueueTime = car.TimeSpentInRegisterQueue
		}
		registerStats.AverageQueueTime += car.TimeSpentInRegisterQueue

		// Pump stats collection
		stats := pumpStats[car.Fuel]
		stats.TotalCars++
		stats.TotalTime += car.FuelingDuration
		stats.TotalQueueTime += car.TimeSpentInPumpQueue
		if car.TimeSpentInPumpQueue > stats.MaxQueueTime {
			stats.MaxQueueTime = car.TimeSpentInPumpQueue
		}
		stats.AverageQueueTime += car.TimeSpentInPumpQueue
	}

	// Average calculations
	registerStats.AverageQueueTime = time.Duration(registerStats.AverageQueueTime.Milliseconds() / int64(registerStats.TotalCars))

	for _, stats := range pumpStats {
		if stats.TotalCars == 0 {
			continue
		}
		stats.AverageQueueTime = time.Duration(stats.TotalQueueTime.Milliseconds() / int64(stats.TotalCars))
	}

	// Print statistics
	fmt.Println("Final statistics")

	fmt.Println("Registers:")
	registerStats.Print()

	fmt.Println("Pumps:")
	for fuelType, stats := range pumpStats {
		fmt.Printf("%s:\n", fuelType)
		stats.Print()
	}

	finishWg.Done()
}
