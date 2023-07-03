package visio

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/kdm"
)

type phase struct {
	number int
	open   [16]bool
}

var phases map[int]phase

var ready bool

func load() {
	phases = make(map[int]phase)
	for {
		for !kdm.State.Connect {
			time.Sleep(time.Second)
		}
		ok := true
		for i := 1; i < 33; i++ {
			if i < 17 {
				kdm.RequestChan <- kdm.Request{Start: 0xA00 + uint16(i*14), Lenght: 14}
			} else {
				kdm.RequestChan <- kdm.Request{Start: 0xB00 + uint16((i-17)*14), Lenght: 14}
			}
			rep := <-kdm.ReplayChan
			if rep.Status != nil {
				logger.Error.Printf("%s", rep.Status.Error())
				ok = false
				break
			}
			agHi := rep.Data[9]
			agLo := rep.Data[8]
			if agHi == 0 && agLo == 0 {
				continue
			}
			ph := phase{number: i}
			j := 0
			for i := 0; i < 8; i++ {
				ph.open[j] = false
				if agLo&1 == 1 {
					ph.open[j] = true
				}
				agLo = agLo >> 1
				j++
			}
			agHi = agHi >> 8
			for i := 0; i < 4; i++ {
				ph.open[j] = false
				if agHi&1 == 1 {
					ph.open[j] = true
				}
				agHi = agHi >> 1
				j++
			}
			phases[ph.number] = ph
		}
		if ok {
			break
		}
	}
}
