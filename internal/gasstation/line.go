package gasstation

import (
	"strconv"
	"sync"
)

func (l *line) Handle(vehicle *vehicle, station *GasStation) {
	claimedPump := false
	//fmt.Println("line ", l.Type, "received vehicle: ", vehicle.ID)
	vehicle.setWaitingForPump()

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
			vehicle.setDoneWaitingForPump(station)
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
	linesProgress.Add(1)
	//fmt.Println("line", line.Type, "is done and waiting for the pumps to finish")
	line.Wg.Wait()
	//fmt.Println("All pumps are done at line", line.Type)
}
func constructLines(config map[string]PumpConfig) []*line {
	var lines []*line

	for _, pumpConfig := range config {
		if !lineTypeExists(lines, fuelType(pumpConfig.Type)) {
			lines = append(lines, &line{
				Type:             fuelType(pumpConfig.Type),
				Pumps:            []*pump{},
				PumpAvailability: make(chan *pump),
				Queue:            make(chan *vehicle, 1000),
				Wg:               sync.WaitGroup{},
			})
		}
		for _, line := range lines {
			if line.Type == fuelType(pumpConfig.Type) {
				for i := 0; i < pumpConfig.Amount; i++ {
					line.Pumps = append(line.Pumps, &pump{
						ID:       pumpConfig.Type + strconv.Itoa(len(line.Pumps)),
						Name:     pumpConfig.Name,
						Speed:    pumpConfig.Speed,
						Occupied: false,
					})
				}
			}
		}
	}
	return lines
}
func lineTypeExists(lines []*line, lineType fuelType) bool {
	for _, line := range lines {
		if line.Type == lineType {
			return true
		}
	}
	return false
}
