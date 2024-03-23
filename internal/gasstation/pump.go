package gasstation

import "time"

func (p *pump) Handle(vehicle *vehicle, station *GasStation) {
	//fmt.Println("vehicle ", vehicle.ID, "is filling up")
	time.Sleep(randomDuration(p.Speed))
	//fmt.Println("vehicle ", vehicle.ID, "is done filling up")
	p.Occupied = false
}

func pumpRoutine(line *line, pump *pump, vehicle *vehicle, station *GasStation) {
	defer func() {
		leastBusyRegister(station).Queue <- vehicle
		line.Wg.Done()
		line.PumpAvailability <- pump
	}()
	pump.Handle(vehicle, station)
}
