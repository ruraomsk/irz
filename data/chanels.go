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

var Commands chan InternalCmd
var Arrays chan binding.Arrays
var ToServer chan int

func initChans() {
	Commands = make(chan InternalCmd, 10)
	Arrays = make(chan binding.Arrays)
	ToServer = make(chan int)
}
