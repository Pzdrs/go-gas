package gasstation

func (p *pump) Handle(vehicle *vehicle, station *GasStation) {
	//fmt.Println("vehicle ", vehicle.ID, "is filling up")
	//time.Sleep(5 * time.Millisecond)
	//fmt.Println("vehicle ", vehicle.ID, "is done filling up")
	p.Occupied = false

	station.CollectMetric(func() {
		vehiclesFilledUpMutex.Lock()
		vehiclesFilledUp++
		vehiclesFilledUpMutex.Unlock()
	})
}

func pumpRoutine(line *line, pump *pump, vehicle *vehicle, station *GasStation) {
	defer func() {
		queueLeastBusyRegister(vehicle, station)
		line.Wg.Done()
		line.PumpAvailability <- pump
	}()
	pump.Handle(vehicle, station)
}
