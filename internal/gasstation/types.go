package gasstation

import "sync"

type fuelType string

const (
	gas      fuelType = "gas"
	diesel   fuelType = "diesel"
	electric fuelType = "electric"
	lpg      fuelType = "lpg"
)

var fuelTypes = []fuelType{
	gas, diesel, electric, lpg,
}

type GasStation struct {
	SimulationRunning  bool
	SimulationComplete bool

	DebugLogging bool

	Lines   []*line
	LinesWg sync.WaitGroup

	Registers   []*cashRegister
	RegistersWg sync.WaitGroup

	StatsWg sync.WaitGroup

	Exit chan *vehicle
}

type cashRegister struct {
	ID    int
	Queue chan *vehicle
}

type line struct {
	Type             fuelType
	Pumps            []*pump
	PumpAvailability chan *pump
	Queue            chan *vehicle
	Wg               sync.WaitGroup
}

type pump struct {
	ID       fuelType
	Occupied bool
}

type vehicle struct {
	ID   int
	Type fuelType
}
