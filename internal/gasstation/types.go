package gasstation

import (
	"github.com/google/uuid"
	"sync"
	"time"
)

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
	SimulationID uuid.UUID

	Metrics metrics

	SimulationRunning  bool
	SimulationComplete bool

	Lines   []*line
	LinesWg sync.WaitGroup

	Registers   []*cashRegister
	RegistersWg sync.WaitGroup

	MetricsWg sync.WaitGroup

	Exit chan *vehicle
}

type cashRegister struct {
	ID    int
	Queue chan *vehicle
	Speed []time.Duration
}

type line struct {
	Type             fuelType
	Pumps            []*pump
	PumpAvailability chan *pump
	Queue            chan *vehicle
	// Wait for all the pumps to be finished, only then is the line also finished
	Wg sync.WaitGroup
}

type pump struct {
	ID       string
	Name     string
	Speed    []time.Duration
	Occupied bool
}

type vehicle struct {
	ID                   int
	Type                 fuelType
	LineQueueArrivalTime time.Time
	LineWaitTime         time.Duration

	RegisterQueueArrivalTime time.Time
	RegisterWaitTime         time.Duration
}
