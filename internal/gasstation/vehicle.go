package gasstation

import "time"

func (v *vehicle) setWaitingForPump() {
	v.LineQueueArrivalTime = time.Now()
}

func (v *vehicle) setDoneWaitingForPump(station *GasStation) {
	waitTime := time.Since(v.LineQueueArrivalTime)
	station.logMetric(func() {
		station.Metrics.lineQueueTime.Observe(float64(waitTime.Milliseconds()))
	})
}

func (v *vehicle) setWaitingForRegister() {
	v.RegisterQueueArrivalTime = time.Now()
}

func (v *vehicle) setDoneWaitingForRegister(station *GasStation) {
	waitTime := time.Since(v.RegisterQueueArrivalTime)
	station.logMetric(func() {
		station.Metrics.registerQueueTime.Observe(float64(waitTime.Milliseconds()))
	})
}
