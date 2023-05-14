package data

import (
	"fmt"

	"github.com/ruraomsk/ag-server/binding"
	"github.com/ruraomsk/ag-server/pudge"
)

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
	Phase   int  //Фаза плана
	PhaseTU int  // Заданная фаза
	PhaseTC int  // Текущая фаза
	Door    bool // Открыта ли дверь
	Lamp    int  // На какой фазе перегорели двери
	TimeTU  int  // Время фазы ТУ
	TimeTC  int  // Время фазы ТС
	Connect bool // Есть ли связь с устройством

}

func (s *StatusDevice) ToString() string {
	return fmt.Sprintf("Фаза %d ТУ %d ТС %d Время ТC %d ТУ %d", s.Phase, s.PhaseTU, s.PhaseTC, s.TimeTC, s.TimeTU)
}

func (s *StatusDevice) NewPhase() {
	s.TimeTC = 0
	s.TimeTU = 0
}

var QAInfo chan QInfo
var AInfo chan bool

var Commands chan InternalCmd
var Arrays chan binding.Arrays
var ToServer chan int
var ToDevice chan int
var FromDevice chan StatusDevice
var Statistics chan pudge.Statistic

func initChans() {
	Commands = make(chan InternalCmd, 100)
	Arrays = make(chan binding.Arrays, 100)
	ToServer = make(chan int)
	ToDevice = make(chan int)
	FromDevice = make(chan StatusDevice, 10)
	Statistics = make(chan pudge.Statistic, 10)
	QAInfo = make(chan QInfo)
	AInfo = make(chan bool)

}
