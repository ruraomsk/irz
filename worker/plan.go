package worker

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/binding"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
)

var workplan = false //
var endPlan chan (interface{})
var notCmdControl = 10 * 10

func stopPlan() {
	if workplan {
		workplan = false
		<-endPlan
	}
}
func exitPlan() {
	endPlan <- 0
}

var state data.StatusDevice

func waitTime(seconds int, phase int) error {
	dk.RDK = 9
	data.DataValue.SetDK(dk)
	tick := time.NewTicker(100 * time.Millisecond)
	endphase := time.NewTimer(time.Duration(seconds) * time.Second)
	repeat := time.NewTicker(time.Minute)
	count := 0
	data.ToDevice <- phase
	data.ToServer <- 0

	logger.Info.Printf("Исполняем фаза %d плана %d длительность %d", phase, nowPlan, seconds)
	for {
		select {
		case <-repeat.C:
			data.ToDevice <- phase
		case <-tick.C:
			if !workplan {
				return fmt.Errorf("end work")
			}
			if state.PhaseTU != state.PhaseTC {
				count++
				if count > notCmdControl {
					return nil
				}
			}

		case <-endphase.C:
			return nil
		case state = <-toPlan:
			logger.Info.Printf("to plan %v", state)
		}
	}
}
func goPlan(pl int) {
	nowPlan = pl

	var pk = binding.SetPk{Pk: 0}
	for _, v := range data.DataValue.Arrays.SetDK.DK {
		if v.Pk == pl {
			pk = v
		}
	}
	if pk.Pk == 0 {
		logger.Error.Printf("Нет плана координации %d", pl)
		return
	}
	//Выполнение простого плана

	dk.EDK = 0
	data.DataValue.SetDK(dk)
	workplan = true
	defer exitPlan()
	logger.Info.Printf("Выполняем план %d", pk.Pk)
	data.DataValue.Controller.PK = pl
	if pk.Tc == 0 {
		//ЛР
		for {
			if waitTime(10000, 0) != nil {
				return
			}
		}
	}
	if pk.Tc == 1 {
		//ЖМ
		for {
			if waitTime(9999999, 10) != nil {
				return
			}
		}

	}
	if pk.Tc == 2 {
		for {
			if waitTime(9999999, 11) != nil {
				return
			}
		}

	}
	pk = repackPlan(pk)
	ctrl := buildControl(pk)
	flagP := true //Флаг переходной фазы
	dk.PDK = true
	data.DataValue.SetDK(dk)

	for {
		if flagP {
			startCycle <- ctrl
		}
		for _, v := range pk.Stages {
			if v.Start == 0 && v.Stop == 0 {
				continue
			}
			if v.Tf == 0 {
				if waitTime(v.Stop-v.Start, v.Number) != nil {

					return
				}
				if v.Number != state.PhaseTU {
					logger.Error.Printf("Неподчинение фазы %d приходит %d", v.Number, state.PhaseTU)
				}
			}
		}
		if flagP {
			stop <- 0
			t := <-getControl
			if t.isGood() {
				flagP = false
				dk.PDK = false
				data.DataValue.SetDK(dk)
			}
		}
	}
}
func buildControl(pk binding.SetPk) control {
	return control{plans: make([]ctrlPlan, 0), lenght: pk.Tc}
}
func repackPlan(pk binding.SetPk) binding.SetPk {
	logger.Info.Printf("in %s", toSting(pk))
	if pk.TypePU == 1 {
		logger.Info.Printf("План локальный ничего не меняем")
		return pk
	}
	if pk.Shift == 0 {
		logger.Info.Printf("План координированный смещение 0")
		return pk
	}
	newPk := pk
	newPk.Stages = make([]binding.Stage, 0)
	//Находим начальную фазу
	last := 0
	for _, v := range pk.Stages {
		if v.Start == 0 && v.Stop == 0 {
			continue
		}
		if v.Start < pk.Shift {
			r := v.Stop - v.Start
			v.Start = last
			v.Stop = last + r
			if v.Trs {
				v.Stop += v.Dt
			}
			if v.Tf == 0 {
				last = v.Stop
			} else {
				last = v.Start
			}
			newPk.Stages = append(newPk.Stages, v)
		} else {
			v.Start -= pk.Shift
			v.Stop -= pk.Shift
			if v.Trs {
				v.Stop += v.Dt
			}
			if v.Tf == 0 {
				last = v.Stop
			} else {
				last = v.Start
			}
			newPk.Stages = append(newPk.Stages, v)
		}
	}
	logger.Info.Printf("out %s", toSting(newPk))
	return newPk
}
func toSting(pk binding.SetPk) string {
	res := fmt.Sprintf("len=%d {", pk.Shift)
	for _, v := range pk.Stages {
		if v.Start == 0 && v.Stop == 0 {
			continue
		}
		res += fmt.Sprintf("[%d-%d f=%d t=%d dt=%v pl=%v trs=%v]", v.Start, v.Stop, v.Number, v.Tf, v.Dt, v.Plus, v.Trs)
	}
	res += "}"
	return res

}
