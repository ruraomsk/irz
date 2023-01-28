package data

import "github.com/ruraomsk/ag-server/binding"

type InternalCmd struct {
	Source   int
	Command  int
	Parametr int
}

const (
	DKDevice = 1
	VPU
	IngenerUSDK
	USDK
	IngenerDK
	IBM
)

type StatusDevice struct {
	PhaseTU int  // Заданная фаза
	PhaseTC int  // Текущая фаза
	Door    bool // Открыта ди дверь
	Lamp    int  // На какой фазе перегорели двери
	TimeTU  int  // Время фазы ТУ
	TimeTC  int  // Время фазы ТС
}

func (s *StatusDevice) NewPhase() {
	s.TimeTC = 0
	s.TimeTU = 0
}

var Commands chan InternalCmd
var Arrays chan binding.Arrays
var ToServer chan int
var ToDevice chan int
var FromDevice chan StatusDevice

func initChans() {
	Commands = make(chan InternalCmd, 100)
	Arrays = make(chan binding.Arrays, 100)
	ToServer = make(chan int, 100)
	ToDevice = make(chan int, 100)
	FromDevice = make(chan StatusDevice, 100)
}
