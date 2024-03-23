package main

import (
	"fmt"
	"math/rand/v2"
	"sync"
)

type Pump struct {
	ID       FuelType
	Occupied bool
}

type Vehicle struct {
	ID   int
	Type FuelType
}

type Line struct {
	Type             FuelType
	Pumps            []*Pump
	PumpAvailability chan *Pump
	Queue            chan *Vehicle
	Wg               sync.WaitGroup
}

func countFreePumps(pumps []*Pump) int {
	count := 0
	for _, pump := range pumps {
		if !pump.Occupied {
			count++
		}
	}
	return count
}

var lines = getLines()

type FuelType string

const (
	Gas      FuelType = "gas"
	Diesel   FuelType = "diesel"
	Electric FuelType = "electric"
	Lpg      FuelType = "lpg"
)

var fuelTypes = []FuelType{
	Gas, Diesel, Electric, Lpg,
}

func getLines() []*Line {
	return []*Line{
		{
			Type: "gas",
			Pumps: []*Pump{
				{ID: "gas0", Occupied: false},
				{ID: "gas1", Occupied: false},
				//{ID: "gas2", Occupied: false},
				//{ID: "gas3", Occupied: false},
				//{ID: "gas4", Occupied: false},
			},
			PumpAvailability: make(chan *Pump),
			Queue:            make(chan *Vehicle, 1000),
			// Wait for all the pumps to be finished, only then is the line also finished
			Wg: sync.WaitGroup{},
		},
		{
			Type: "diesel",
			Pumps: []*Pump{
				{ID: "diesel0", Occupied: false},
				{ID: "diesel1", Occupied: false},
			},
			PumpAvailability: make(chan *Pump),
			Queue:            make(chan *Vehicle, 1000),
			// Wait for all the pumps to be finished, only then is the line also finished
			Wg: sync.WaitGroup{},
		},
		{
			Type: "electric",
			Pumps: []*Pump{
				{ID: "electric0", Occupied: false},
			},
			PumpAvailability: make(chan *Pump),
			Queue:            make(chan *Vehicle, 1000),
			// Wait for all the pumps to be finished, only then is the line also finished
			Wg: sync.WaitGroup{},
		},
		{
			Type: "lpg",
			Pumps: []*Pump{
				{ID: "lpg0", Occupied: false},
			},
			PumpAvailability: make(chan *Pump),
			Queue:            make(chan *Vehicle, 1000),
			// Wait for all the pumps to be finished, only then is the line also finished
			Wg: sync.WaitGroup{},
		},
	}
}

type CashRegister struct {
	ID    int
	Queue chan *Vehicle
}

var registers = []*CashRegister{
	{
		ID:    1,
		Queue: make(chan *Vehicle, 1000),
	},
	{
		ID:    2,
		Queue: make(chan *Vehicle, 1000),
	},
}

var exitChannel = make(chan *Vehicle)

var vehiclesLeft = 0
var vehiclesPaid = 0
var vehiclesFilledUp = 0

func main() {
	inspect()
	// Wait for all the lines of pumps to be finished processing cars
	linesWg := sync.WaitGroup{}
	registersWg := sync.WaitGroup{}
	statsWg := sync.WaitGroup{}

	go exitHandler(&statsWg)

	for _, register := range registers {
		registersWg.Add(1)
		go registerHandler(&registersWg, register, &statsWg)
	}

	for _, line := range lines {
		linesWg.Add(1)
		go lineHandler(&linesWg, line, &statsWg)
	}

	go spawnVehicles(1000)

	fmt.Println("Waiting for all lines to be done")
	linesWg.Wait()

	for _, register := range registers {
		close(register.Queue)
	}
	fmt.Println("All lines are done")

	registersWg.Wait()
	fmt.Println("All registers are done")

	// When all the lanes are done (all the line pump goroutines are done) and all the registers are done, there will be no more vehicles leaving because they're all gone already
	close(exitChannel)

	statsWg.Wait()

	fmt.Println(vehiclesLeft)
	fmt.Println(vehiclesPaid)
	fmt.Println(vehiclesFilledUp)
}

func exitHandler(statsWg *sync.WaitGroup) {
	for vehicle := range exitChannel {
		statsWg.Add(1)
		logStatistic(statsWg, func() {
			vehiclesLeft++
		})
		fmt.Println(" --- Vehicle ", vehicle.ID, "is leaving the gas station")
	}
}

func logStatistic(wg *sync.WaitGroup, stat func()) {
	defer wg.Done()
	stat()
}

func inspect() {
	fmt.Println("Inspecting lines")
	for _, line := range lines {
		fmt.Println("Line ", line.Type)
		fmt.Println(" - Pumps:", len(line.Pumps))
	}
}

func randomFuelType() FuelType {
	return fuelTypes[rand.IntN(len(fuelTypes))]
}

func spawnVehicles(goal int) {
	for i := 0; i < goal; i++ {
		vehicle := Vehicle{
			Type: randomFuelType(),
			ID:   i,
		}
		fmt.Println("Spawned vehicle: ", vehicle.ID, " with fuel type: ", vehicle.Type)
		for _, line := range lines {
			if vehicle.Type == line.Type {
				fmt.Println("Sending vehicle ", vehicle.ID, " to line: ", line.Type)
				line.Queue <- &vehicle
				break
			}
		}
		if i < goal-1 {
			//time.Sleep(10 * time.Millisecond)
		}
	}
	for _, line := range lines {
		close(line.Queue)
	}
}

func spawnVehiclesMock() {
	vehicles := []Vehicle{
		{Type: Gas, ID: 1},
		{Type: Gas, ID: 2},
		{Type: Gas, ID: 3},
		{Type: Gas, ID: 4},
		{Type: Gas, ID: 5},
	}

	for _, vehicle := range vehicles {
		fmt.Println("Spawned vehicle: ", vehicle.ID, " with fuel type: ", vehicle.Type)
		for _, line := range lines {
			if vehicle.Type == line.Type {
				fmt.Println("Sending vehicle ", vehicle.ID, " to line: ", line.Type)
				line.Queue <- &vehicle
				break
			}
		}
	}
	for _, line := range lines {
		close(line.Queue)
	}
}

func lineHandler(wg *sync.WaitGroup, line *Line, statsWg *sync.WaitGroup) {
	defer wg.Done()
	for vehicle := range line.Queue {
		claimedPump := false
		fmt.Println("Line ", line.Type, "received vehicle: ", vehicle.ID)

	findPump:
		if countFreePumps(line.Pumps) == 0 {
			fmt.Println("No free pumps available at line ", line.Type)
			freedUpPump := <-line.PumpAvailability
			fmt.Println("Pump ", freedUpPump.ID, "is now free")
		}

		for _, pump := range line.Pumps {
			if !pump.Occupied {
				claimedPump = true
				fmt.Println("Vehicle ", vehicle.ID, "found a free pump: ", pump.ID)
				pump.Occupied = true
				line.Wg.Add(1)
				go pumpHandler(line, pump, vehicle, statsWg)
				break
			}
		}

		// Hack, sometimes even though there are free pumps, the vehicle doesn't find one
		if !claimedPump {
			goto findPump
		}
		fmt.Println("Line ", line.Type, "has finished processing vehicle: ", vehicle.ID)
	}
	fmt.Println("All vehicles for line ", line.Type, "have been processed")
	line.Wg.Wait()
	fmt.Println(" - Line ", line.Type, "is done")
}

func pumpHandler(line *Line, pump *Pump, vehicle *Vehicle, statsWg *sync.WaitGroup) {
	defer func() {
		sendToLeastBusyRegister(vehicle)
		line.Wg.Done()
		line.PumpAvailability <- pump
	}()
	fmt.Println("Vehicle ", vehicle.ID, "is filling up")
	//time.Sleep(100 * time.Millisecond)
	fmt.Println("Vehicle ", vehicle.ID, "is done filling up")
	pump.Occupied = false

	statsWg.Add(1)
	logStatistic(statsWg, func() {
		vehiclesFilledUp++
	})
}

func registerHandler(wg *sync.WaitGroup, reg *CashRegister, statsWg *sync.WaitGroup) {
	defer wg.Done()
	for vehicle := range reg.Queue {
		statsWg.Add(1)
		go logStatistic(statsWg, func() {
			vehiclesPaid++
		})
		fmt.Println("Register ", vehicle.ID, "received vehicle: ", vehicle.ID)
		//time.Sleep(100 * time.Millisecond)
		fmt.Println("Register ", vehicle.ID, "is done processing vehicle: ", vehicle.ID)
		exitChannel <- vehicle
	}
}

func sendToLeastBusyRegister(vehicle *Vehicle) {
	leastBusyRegister := registers[0]
	for _, register := range registers {
		if len(register.Queue) < len(leastBusyRegister.Queue) {
			leastBusyRegister = register
		}
	}
	fmt.Println("Sending vehicle ", vehicle.ID, " to register ", leastBusyRegister.ID)
	leastBusyRegister.Queue <- vehicle
}
