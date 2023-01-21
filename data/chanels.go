package data

import "github.com/ruraomsk/ag-server/binding"

type InternalCmd struct {
	Source   int
	Command  int
	Parametr int
}

const (
	Pult = 1
	Server
	Controller
)

var Commands chan InternalCmd
var Arrays chan binding.Arrays

func initChans() {
	Commands = make(chan InternalCmd, 10)
	Arrays = make(chan binding.Arrays)
}
