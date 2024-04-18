package gasstation

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/schollz/progressbar/v3"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math/rand/v2"
	"sync"
	"time"
)

var configuration StationConfiguration
var vehicleSpawnerProgress, registersProgress, linesProgress, metricsProgress *progressbar.ProgressBar

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
	fmt.Println(" === GAS STATION DETAILS === ")
	fmt.Println("Lines: ", len(s.Lines))
	for i, line := range s.Lines {
		fmt.Println("Line", i, ":", line.Type)
		for _, pump := range line.Pumps {
			fmt.Printf(" - %s (%s)\n", pump.ID, pump.Name)
		}
	}

	fmt.Println("Registers: ", len(s.Registers))
	fmt.Println("=============================")
}

func (s *GasStation) Begin() {
	if s.SimulationRunning {
		panic("The simulation is already running!")
	}
	if s.SimulationComplete {
		panic("The simulation has already been run!")
	}
	s.SimulationRunning = true

	fmt.Println(fmt.Sprintf(">> Starting simulation %s", s.SimulationID))

	startTime := time.Now()

	go s.spawnVehicles()

	s.awaitLines()
	s.awaitRegisters()

	// All vehicles are fueled up and paid up and heading out => no more vehicles will be coming through the exit
	s.closeExit()

	simulationTime := time.Since(startTime)

	s.logMetric(func() {
		s.Metrics.simulationTime.Set(float64(simulationTime.Milliseconds()))
	})

	s.pushMetrics()
	s.awaitMetrics()

	p := message.NewPrinter(language.English)

	p.Printf(">> Done simulating %d vehicles. Took: %s <<\n", configuration.VehicleSpawner.Goal, simulationTime)

	s.SimulationRunning = false
	s.SimulationComplete = true
}

func (s *GasStation) spawnVehicles() {
	goal := configuration.VehicleSpawner.Goal
	vehicleSpawnerProgress = progressbar.Default(int64(goal), "Spawning vehicles")
	for i := 0; i < goal; i++ {
		vehicle := vehicle{
			Type: randomFuelType(),
			ID:   i,
		}
		for _, line := range s.Lines {
			if vehicle.Type == line.Type {
				//fmt.Println("Sending vehicle ", vehicle.ID, " to line: ", line.Type)
				line.Queue <- &vehicle
				break
			}
		}
		vehicleSpawnerProgress.Add(1)
		if i < goal-1 {
			time.Sleep(randomDuration(configuration.VehicleSpawner.Speed))
		}
	}
	fmt.Println(" 🗹 All vehicles have been spawned")
	linesProgress = progressbar.Default(int64(len(s.Lines)), "Processing lines")
	s.closeLines()
}

func (s *GasStation) awaitLines() {
	// Wait for all the lines to be done (a line is done when all its pumps are done)
	s.LinesWg.Wait()
	// All lines are done => all pumps are done => all vehicles are in register queues and no more vehicles are coming
	s.closeRegisters()
	fmt.Println(" 🗹 All lines are done")
}

func (s *GasStation) awaitRegisters() {
	registersProgress = progressbar.Default(-1, "Processing registers")
	// Wait for all the registers to be done
	s.RegistersWg.Wait()
	registersProgress.Finish()
	fmt.Println(" 🗹 All registers are done")
}

func (s *GasStation) awaitMetrics() {
	metricsProgress = progressbar.Default(-1, "Collecting metrics")
	// Wait for all the metrics to be collected
	s.MetricsWg.Wait()
	metricsProgress.Finish()
	fmt.Println(" 🗹 All metrics have been collected")
}

func (s *GasStation) exitHandler() {
	for vehicle := range s.Exit {
		_ = vehicle
		s.logMetric(func() {
			s.Metrics.carsProcessedTotal.Inc()
		})
	}
}

func (s *GasStation) closeRegisters() {
	for _, register := range s.Registers {
		register.Close()
	}
}

func (s *GasStation) closeLines() {
	for _, line := range s.Lines {
		close(line.Queue)
	}
}

func (s *GasStation) closeExit() {
	close(s.Exit)
}

func NewGasStation(configFile string) *GasStation {
	configuration = loadConfig(configFile)

	return &GasStation{
		SimulationID: uuid.New(),
		Metrics:      *registerMetrics(),
		Lines:        constructLines(configuration.Pumps),
		LinesWg:      sync.WaitGroup{},
		Registers:    constructRegisters(configuration.Registers),
		RegistersWg:  sync.WaitGroup{},
		MetricsWg:    sync.WaitGroup{},
		Exit:         make(chan *vehicle),
	}
}

func randomFuelType() fuelType {
	return fuelTypes[rand.IntN(len(fuelTypes))]
}
func randomDuration(durationRange []time.Duration) time.Duration {
	if len(durationRange) == 1 {
		return durationRange[0]
	}
	if len(durationRange) != 2 {
		panic("Invalid duration range")
	}

	speedMinNano := int64(durationRange[0])
	speedMaxNano := int64(durationRange[1])

	randomNano := speedMinNano + rand.Int64N(speedMaxNano-speedMinNano+1)

	return time.Duration(randomNano)
}
