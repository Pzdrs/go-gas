package gasstation

import (
	"github.com/prometheus/client_golang/prometheus"
	"strconv"
	"time"
)

func (reg *cashRegister) Close() {
	close(reg.Queue)
}

func (reg *cashRegister) Serve(vehicle *vehicle, station *GasStation) {
	//fmt.Println("Register ", reg.ID, "received vehicle: ", vehicle.ID)
	vehicle.setDoneWaitingForRegister(station)

	duration := randomDuration(reg.Speed)
	time.Sleep(duration)
	station.logMetric(func() {
		station.Metrics.registerTime.With(prometheus.Labels{"register": strconv.Itoa(reg.ID)}).Observe(float64(duration.Milliseconds()))
	})
	//fmt.Println("Register ", reg.ID, "is done processing vehicle: ", vehicle.ID)
	station.Exit <- vehicle
	if registersProgress != nil {
		registersProgress.Add(1)
	}
}

func registerRoutine(reg *cashRegister, station *GasStation) {
	defer station.RegistersWg.Done()
	for vehicle := range reg.Queue {
		reg.Serve(vehicle, station)
	}
}
func leastBusyRegister(station *GasStation) *cashRegister {
	var leastBusy *cashRegister
	for _, register := range station.Registers {
		if leastBusy == nil || len(leastBusy.Queue) > len(register.Queue) {
			leastBusy = register
		}
	}
	return leastBusy
}
func constructRegisters(config RegisterConfig) []*cashRegister {
	registers := make([]*cashRegister, config.Amount)
	for i := 0; i < config.Amount; i++ {
		registers[i] = &cashRegister{
			ID:    i,
			Queue: make(chan *vehicle, 100),
			Speed: config.Speed,
		}
	}
	return registers
}
