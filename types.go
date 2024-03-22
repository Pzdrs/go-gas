package main

import (
	"sync"
	"time"
)

type GasStation struct {
	Pumps        []*Pump
	Registers    []*CashRegister
	ArrivalQueue chan *Vehicle
	ExitQueue    chan *Vehicle
	PaymentQueue chan *Vehicle
}

type Pump struct {
	ID    string
	Name  string
	Type  string
	Speed []time.Duration
	Queue chan *Vehicle
}

type CashRegister struct {
	ID    int
	Queue chan *Vehicle
}

const (
	Gas      = "gas"
	Diesel   = "diesel"
	LPG      = "lpg"
	Electric = "electric"
)

type Vehicle struct {
	ID                       int
	Fuel                     string
	PumpQueueEnter           time.Time
	RegisterQueueEnter       time.Time
	TimeSpentInPumpQueue     time.Duration
	TimeSpentInRegisterQueue time.Duration
	FuelingDuration          time.Duration
	PaymentDuration          time.Duration
	TotalTime                time.Duration
	CarSync                  *sync.WaitGroup
}
