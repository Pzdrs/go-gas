package gasstation

import "sync"

func (l *line) Handle(vehicle *vehicle, station *GasStation) {
	claimedPump := false
	//fmt.Println("line ", l.Type, "received vehicle: ", vehicle.ID)

findPump:
	if !l.hasUnoccupiedPumps() {
		//fmt.Println("No free pumps available at line ", l.Type)
		freedUpPump := <-l.PumpAvailability
		_ = freedUpPump
		//fmt.Println("pump ", freedUpPump.ID, "is now free")
	}

	for _, pump := range l.Pumps {
		if !pump.Occupied {
			claimedPump = true
			//fmt.Println("vehicle ", vehicle.ID, "found a free pump: ", pump.ID)
			pump.Occupied = true
			l.Wg.Add(1)
			go pumpRoutine(l, pump, vehicle, station)
			break
		}
	}

	// Hack, sometimes even though there are free pumps, the vehicle doesn't find one
	if !claimedPump {
		goto findPump
	}
	//fmt.Println("line ", l.Type, "has finished processing vehicle: ", vehicle.ID)
}

func (l *line) hasUnoccupiedPumps() bool {
	count := 0
	for _, pump := range l.Pumps {
		if !pump.Occupied {
			count++
		}
	}
	return count > 0
}

func lineRoutine(line *line, station *GasStation) {
	defer station.LinesWg.Done()
	for vehicle := range line.Queue {
		line.Handle(vehicle, station)
	}
	//fmt.Println("line", line.Type, "is done and waiting for the pumps to finish")
	line.Wg.Wait()
	//fmt.Println("All pumps are done at line", line.Type)
}

func getLines() []*line {
	return []*line{
		{
			Type: "gas",
			Pumps: []*pump{
				{ID: "gas0", Occupied: false},
				{ID: "gas1", Occupied: false},
				{ID: "gas2", Occupied: false},
			},
			PumpAvailability: make(chan *pump),
			Queue:            make(chan *vehicle, 1000),
			// Wait for all the pumps to be finished, only then is the line also finished
			Wg: sync.WaitGroup{},
		},
		{
			Type: "diesel",
			Pumps: []*pump{
				{ID: "diesel0", Occupied: false},
				{ID: "diesel1", Occupied: false},
				{ID: "diesel2", Occupied: false},
			},
			PumpAvailability: make(chan *pump),
			Queue:            make(chan *vehicle, 1000),
			// Wait for all the pumps to be finished, only then is the line also finished
			Wg: sync.WaitGroup{},
		},
		{
			Type: "electric",
			Pumps: []*pump{
				{ID: "electric0", Occupied: false},
			},
			PumpAvailability: make(chan *pump),
			Queue:            make(chan *vehicle, 1000),
			// Wait for all the pumps to be finished, only then is the line also finished
			Wg: sync.WaitGroup{},
		},
		{
			Type: "lpg",
			Pumps: []*pump{
				{ID: "lpg0", Occupied: false},
				{ID: "lpg1", Occupied: false},
			},
			PumpAvailability: make(chan *pump),
			Queue:            make(chan *vehicle, 1000),
			// Wait for all the pumps to be finished, only then is the line also finished
			Wg: sync.WaitGroup{},
		},
	}
}
