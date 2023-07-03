package visio

import (
	"github.com/ruraomsk/irz/kdm"
	"github.com/ruraomsk/irz/setup"
)

func Visio() {
	if !setup.Set.Visio {
		for {
			select {
			case <-kdm.InfoChan:
			}
		}
	}
	load()
	for {

	}
}
