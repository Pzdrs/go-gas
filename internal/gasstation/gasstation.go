package gasstation

import (
	"fmt"
	"math/rand/v2"
	"sync"
	"time"
)

var vehiclesLeft = 0
var vehiclesPaid = 0
var vehiclesFilledUp = 0

type GasStation struct {
	SimulationRunning  bool
	SimulationComplete bool

	Lines   []*Line
	LinesWg sync.WaitGroup

	Registers   []*CashRegister
	RegistersWg sync.WaitGroup

	StatsWg sync.WaitGroup

	Exit chan *Vehicle
}

type Pump struct {
	ID       FuelType
	Occupied bool
}

func (p *Pump) Handle(vehicle *Vehicle, station *GasStation) {
	fmt.Println("Vehicle ", vehicle.ID, "is filling up")
	//time.Sleep(5 * time.Millisecond)
	fmt.Println("Vehicle ", vehicle.ID, "is done filling up")
	p.Occupied = false

	station.LogMetric(func() {
		vehiclesFilledUp++
	})
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

func (l *Line) Handle(vehicle *Vehicle, station *GasStation) {
	claimedPump := false
	fmt.Println("Line ", l.Type, "received vehicle: ", vehicle.ID)

findPump:
	if countFreePumps(l.Pumps) == 0 {
		fmt.Println("No free pumps available at line ", l.Type)
		freedUpPump := <-l.PumpAvailability
		fmt.Println("Pump ", freedUpPump.ID, "is now free")
	}

	for _, pump := range l.Pumps {
		if !pump.Occupied {
			claimedPump = true
			fmt.Println("Vehicle ", vehicle.ID, "found a free pump: ", pump.ID)
			pump.Occupied = true
			l.Wg.Add(1)
			go pumpRoutine(l, pump, vehicle, station)
			break
		}
	}

	// Hack, sometimes even though there are free pumps, the vehicle doesn't find one
	if !claimedPump {
		goto findPump
	}
	fmt.Println("Line ", l.Type, "has finished processing vehicle: ", vehicle.ID)
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

type CashRegister struct {
	ID    int
	Queue chan *Vehicle
}

func (reg *CashRegister) Close() {
	close(reg.Queue)
}

func (reg *CashRegister) Serve(vehicle *Vehicle, station *GasStation) {
	fmt.Println("Register ", reg.ID, "received vehicle: ", vehicle.ID)
	//time.Sleep(10 * time.Millisecond)
	station.LogMetric(func() {
		vehiclesPaid++
	})
	fmt.Println("Register ", vehicle.ID, "is done processing vehicle: ", vehicle.ID)
	station.Exit <- vehicle
}

func (s *GasStation) LogMetric(metric func()) {
	s.StatsWg.Add(1)
	go func() {
		defer s.StatsWg.Done()
		metric()
	}()
}

func (s *GasStation) exitHandler() {
	for vehicle := range s.Exit {
		s.LogMetric(func() {
			vehiclesLeft++
		})
		fmt.Println(" --- Vehicle ", vehicle.ID, "is leaving the gas station")
	}
}

func (s *GasStation) Setup() {
	go s.exitHandler()

	for _, register := range s.Registers {
		s.RegistersWg.Add(1)
		go registerRoutine(register, s)
	}

	for _, line := range s.Lines {
		s.LinesWg.Add(1)
		go lineRoutine(line, s)
	}
}

func (s *GasStation) Inspect() {
	fmt.Println("Inspecting lines")
	for _, line := range s.Lines {
		fmt.Println("Line ", line.Type)
		fmt.Println(" - Pumps:", len(line.Pumps))
	}
}

func (s *GasStation) closeLines() {
	for _, line := range s.Lines {
		close(line.Queue)
	}
}

func (s *GasStation) closeRegisters() {
	for _, register := range s.Registers {
		register.Close()
	}
}

func (s *GasStation) spawnVehicles(goal int) {
	for i := 0; i < goal; i++ {
		vehicle := Vehicle{
			Type: randomFuelType(),
			ID:   i,
		}
		fmt.Println("Spawned vehicle: ", vehicle.ID, " with fuel type: ", vehicle.Type)
		for _, line := range s.Lines {
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
	s.closeLines()
}

func (s *GasStation) Begin(vehicleGoal int) {
	if s.SimulationRunning {
		panic("The simulation is already running!")
	}
	if s.SimulationComplete {
		panic("The simulation has already been run!")
	}
	s.SimulationRunning = true

	startTime := time.Now()

	go s.spawnVehicles(vehicleGoal)

	s.LinesWg.Wait()

	// All lines are done => all pumps are done => all vehicles are in register queues and no more vehicles are coming
	s.closeRegisters()

	fmt.Println("All lines are done")

	s.RegistersWg.Wait()
	fmt.Println("All registers are done")

	s.StatsWg.Wait()
	fmt.Println("All stats are logged")

	// All vehicles are fueled up and paid up and heading out => no more vehicles will be coming through the exit
	s.closeExit()
	fmt.Println("The exit is closed")

	fmt.Println(vehiclesLeft)
	fmt.Println(vehiclesPaid)
	fmt.Println(vehiclesFilledUp)
	fmt.Println("The simulation took: ", time.Since(startTime))

	s.SimulationRunning = false
	s.SimulationComplete = true
}

func (s *GasStation) closeExit() {
	close(s.Exit)
}

func NewGasStation() *GasStation {
	return &GasStation{
		Lines:   getLines(),
		LinesWg: sync.WaitGroup{},
		Registers: []*CashRegister{
			{
				ID:    1,
				Queue: make(chan *Vehicle, 1000),
			},
			{
				ID:    2,
				Queue: make(chan *Vehicle, 1000),
			},
		},
		RegistersWg: sync.WaitGroup{},
		StatsWg:     sync.WaitGroup{},
		Exit:        make(chan *Vehicle),
	}
}

func getLines() []*Line {
	return []*Line{
		{
			Type: "gas",
			Pumps: []*Pump{
				{ID: "gas0", Occupied: false},
				{ID: "gas1", Occupied: false},
				{ID: "gas2", Occupied: false},
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
				{ID: "diesel2", Occupied: false},
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
				{ID: "lpg1", Occupied: false},
			},
			PumpAvailability: make(chan *Pump),
			Queue:            make(chan *Vehicle, 1000),
			// Wait for all the pumps to be finished, only then is the line also finished
			Wg: sync.WaitGroup{},
		},
	}
}

func randomFuelType() FuelType {
	return fuelTypes[rand.IntN(len(fuelTypes))]
}

func lineRoutine(line *Line, station *GasStation) {
	defer station.LinesWg.Done()
	for vehicle := range line.Queue {
		line.Handle(vehicle, station)
	}
	fmt.Println("Line", line.Type, "is done and waiting for the pumps to finish")
	line.Wg.Wait()
	fmt.Println("All pumps are done at line", line.Type)
}

func pumpRoutine(line *Line, pump *Pump, vehicle *Vehicle, station *GasStation) {
	defer func() {
		sendToLeastBusyRegister(vehicle, station)
		line.Wg.Done()
		line.PumpAvailability <- pump
	}()
	pump.Handle(vehicle, station)
}

func registerRoutine(reg *CashRegister, station *GasStation) {
	defer station.RegistersWg.Done()
	for vehicle := range reg.Queue {
		reg.Serve(vehicle, station)
	}
}

func sendToLeastBusyRegister(vehicle *Vehicle, station *GasStation) {
	leastBusyRegister := station.Registers[0]
	for _, register := range station.Registers {
		if len(register.Queue) < len(leastBusyRegister.Queue) {
			leastBusyRegister = register
		}
	}
	fmt.Println("Sending vehicle ", vehicle.ID, " to register ", leastBusyRegister.ID)
	leastBusyRegister.Queue <- vehicle
}
