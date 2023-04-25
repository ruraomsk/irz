package worker

import (
	"time"

	"github.com/ruraomsk/irz/data"
)

// Контроль исполнения плана
type control struct {
	isCtrl bool
	lenght int
	tc     int
	tfact  ctrlPlan
	facts  []ctrlPlan
	plans  []ctrlPlan
}
type ctrlPlan struct {
	start int
	stop  int
	phase int
}

func (c control) isGood() bool {

	return true
}

var nowCtrl = control{isCtrl: false, plans: make([]ctrlPlan, 0)}
var nowStage = ctrlPlan{}
var chanCtrl chan data.StatusDevice
var startCycle chan control
var stop chan interface{}
var getControl chan control
var isGood chan bool

func cintrollerPlans() {
	chanCtrl = make(chan data.StatusDevice)
	startCycle = make(chan control)
	stop = make(chan interface{})
	getControl = make(chan control)
	ticker := time.NewTicker(time.Second)
	for {
		select {
		case <-ticker.C:
			if nowCtrl.isCtrl {
				nowCtrl.tc++
				if nowCtrl.tc > nowCtrl.lenght {
					nowCtrl.tfact.stop = nowCtrl.tc
					nowCtrl.facts = append(nowCtrl.facts, nowCtrl.tfact)
					nowCtrl.isCtrl = false
				}
			}
		case d := <-chanCtrl:
			if !nowCtrl.isCtrl {
				continue
			}
			if !d.Connect {
				continue
			}
			// тут напишем все остальное
		case nowCtrl = <-startCycle:
			nowCtrl.isCtrl = true
			nowCtrl.facts = make([]ctrlPlan, 0)
			nowCtrl.tc = 0
			nowCtrl.tfact.start = 0
		case <-stop:
			nowCtrl.tfact.stop = nowCtrl.tc
			nowCtrl.facts = append(nowCtrl.facts, nowCtrl.tfact)
			getControl <- nowCtrl
			nowCtrl.isCtrl = false
		}
	}
}
