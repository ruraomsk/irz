package data

import "github.com/ruraomsk/ag-server/binding"

type InternalCmd struct {
	Source   int
	Command  int
	Parametr int
}

const (
	Pult = 1
	IBM
	Controller
)

type StatusDevice struct {
	Need  int  // Заданная фаза
	Phase int  // Текущая фаза
	Door  bool // Открыта ди дверь
	Lamp  int  // На какой фазе перегорели двери
}

var Commands chan InternalCmd
var Arrays chan binding.Arrays
var ToServer chan int
var ToDevice chan int
var FromDevice chan StatusDevice

func initChans() {
	Commands = make(chan InternalCmd, 10)
	Arrays = make(chan binding.Arrays)
	ToServer = make(chan int)
	ToDevice = make(chan int)
	FromDevice = make(chan StatusDevice)
}
