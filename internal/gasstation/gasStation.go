package gasstation

import (
	"context"
	"fmt"
	"github.com/Pzdrs/go-gas/internal/config"
	"github.com/google/uuid"
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
	"log"
	"math/rand/v2"
	"sync"
	"time"
)

var influxClient influxdb2.Client

func (s *GasStation) CollectMetric(getDataPoint func() *write.Point) {
	s.StatsWg.Add(1)
	go func(client influxdb2.Client) {
		defer s.StatsWg.Done()
		//point := getDataPoint()
		point := influxdb2.NewPoint("register-serving", map[string]string{"simulation": s.SimulationID.String()}, map[string]interface{}{"duration": 0}, time.Now())

		writeAPI := client.WriteAPIBlocking("Homelab", "go-gas")
		if err := writeAPI.WritePoint(context.Background(), point); err != nil {
			log.Fatal(err)
		}
		fmt.Println("Metric written")
	}(influxClient)
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

func (s *GasStation) Begin(vehicleGoal int) {
	if s.SimulationRunning {
		panic("The simulation is already running!")
	}
	if s.SimulationComplete {
		panic("The simulation has already been run!")
	}
	s.SimulationRunning = true

	fmt.Println("Starting simulation", s.SimulationID)
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

	fmt.Println("The simulation took:", time.Since(startTime))

	influxClient.Close()

	s.SimulationRunning = false
	s.SimulationComplete = true
}

func (s *GasStation) spawnVehicles(goal int) {
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
		if i < goal-1 {
			//time.Sleep(10 * time.Millisecond)
		}
	}
	s.closeLines()
}

func (s *GasStation) exitHandler() {
	for vehicle := range s.Exit {
		_ = vehicle
		//fmt.Println(" --- vehicle ", vehicle.ID, "is leaving the gas station")
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

func NewGasStation(configuration config.GasStationConfig) *GasStation {
	return &GasStation{
		SimulationID: uuid.New(),
		Lines:        constructLines(config.GasStationConfiguration.Pumps),
		LinesWg:      sync.WaitGroup{},
		Registers:    constructRegisters(configuration.Registers),
		RegistersWg:  sync.WaitGroup{},
		StatsWg:      sync.WaitGroup{},
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
