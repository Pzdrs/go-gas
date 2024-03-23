package gasstation

func (reg *cashRegister) Close() {
	close(reg.Queue)
}

func (reg *cashRegister) Serve(vehicle *vehicle, station *GasStation) {
	//fmt.Println("Register ", reg.ID, "received vehicle: ", vehicle.ID)
	//time.Sleep(10 * time.Millisecond)
	station.CollectMetric(func() {
		vehiclesPaidMutex.Lock()
		vehiclesPaid++
		vehiclesPaidMutex.Unlock()
	})
	//fmt.Println("Register ", vehicle.ID, "is done processing vehicle: ", vehicle.ID)
	station.Exit <- vehicle
}

func registerRoutine(reg *cashRegister, station *GasStation) {
	defer station.RegistersWg.Done()
	for vehicle := range reg.Queue {
		reg.Serve(vehicle, station)
	}
}

func queueLeastBusyRegister(vehicle *vehicle, station *GasStation) {
	leastBusyRegister := station.Registers[0]
	for _, register := range station.Registers {
		if len(register.Queue) < len(leastBusyRegister.Queue) {
			leastBusyRegister = register
		}
	}
	//fmt.Println("Sending vehicle ", vehicle.ID, " to register ", leastBusyRegister.ID)
	leastBusyRegister.Queue <- vehicle
}
