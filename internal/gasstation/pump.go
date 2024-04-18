package gasstation

import (
	"github.com/prometheus/client_golang/prometheus"
	"time"
)

func (p *pump) Handle(vehicle *vehicle, station *GasStation) {
	duration := randomDuration(p.Speed)
	time.Sleep(duration)
	station.logMetric(func() {
		station.Metrics.fuelTime.With(prometheus.Labels{"fuel_type": string(vehicle.Type)}).Observe(float64(duration.Milliseconds()))
	})
	p.Occupied = false
}

func pumpRoutine(line *line, pump *pump, vehicle *vehicle, station *GasStation) {
	defer func() {
		leastBusyRegister(station).Queue <- vehicle
		vehicle.setWaitingForRegister()
		line.Wg.Done()
		line.PumpAvailability <- pump
	}()
	pump.Handle(vehicle, station)
}
