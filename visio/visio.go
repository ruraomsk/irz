package visio

import (
	"fmt"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/kdm"
	"github.com/ruraomsk/irz/setup"
)

var status vdata

type toSend struct {
	buff160 [29]byte
	buff168 [29]byte
}

func (t *toSend) maker(d vdata) {
	t.buff160 = d.makeBuffer(160)
	t.buff168 = d.makeBuffer(168)
}
func (t *toSend) makerOS(d vdata) {
	t.buff160 = d.makeSpecial(160, 0)
	t.buff168 = d.makeSpecial(168, 0)
}
func (t *toSend) makerYell(d vdata) {
	t.buff160 = d.makeSpecial(160, 2)
	t.buff168 = d.makeSpecial(168, 2)
}

var toTo toSend

func Visio() {
	if !setup.Set.Visio {
		for {
			select {
			case <-kdm.InfoChan:
			}
		}
	}
	load()
	status.init()
	toTo.makerOS(status)
	for {
		select {
		case in := <-kdm.InfoChan:
			switch in.Phase {
			case 0:
				// ЛР
				status.init()
				toTo.makerOS(status)
			case 10:
				//ЖМ
				status.init()
				toTo.makerYell(status)
			case 11:
				//OC
				status.init()
				toTo.makerOS(status)
			default:
				p, ok := phases[in.Phase]
				if !ok {
					logger.Error.Printf("Нет такой фазы %d", in.Phase)
					continue
				}
				status.set(in.Lenght, p.open)
				toTo.maker(status)
			}
			//Передаем
			if work {
				senderChan <- toTo
			} else {
				fmt.Printf("Передача\n%v\n%v\n", toTo.buff160, toTo.buff168)
			}
		}
	}
}
