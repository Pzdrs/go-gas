package main

import (
	"github.com/Pzdrs/go-gas/internal/config"
	"time"
)

var station_types = []string{"gas", "diesel", "electric", "lpg"}

type Pump struct {
	occupied bool
}

type Line struct {
	fuelType string
	pumps    []Pump
}

type Vehicle struct {
}

func (pump *Pump) Fill(vehicle *Vehicle) {
	pump.occupied = true
	time.Sleep(1 * time.Second)
	pump.occupied = false
}

func main() {
	config.LoadConfig()

}
