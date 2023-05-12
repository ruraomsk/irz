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
	flagP := 3 //Флаг переходной фазы
	dk.PDK = true
	data.DataValue.SetDK(dk)

	for {
		if flagP > 0 {
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
		flagP--
		if flagP == 0 {
			stop <- 0
			t := <-getControl
			if t.isGood() {
				flagP = 3
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
	tail := make([]binding.Stage, 0)
	//Находим начальную фазу
	skip := false
	for i, v := range pk.Stages {
		if skip {
			skip = false
			continue
		}
		if v.Start == 0 && v.Stop == 0 {
			continue
		}
		if v.Stop < v.Start {
			continue
		}
		if v.Start < pk.Shift {
			v.Start += pk.Tc - pk.Shift
			v.Stop += pk.Tc - pk.Shift
			tail = append(tail, v)
			continue
		}
		v.Start -= pk.Shift
		v.Stop -= pk.Shift
		var dv binding.Stage
		if isTVP(v.Tf) {
			if i < len(pk.Stages) && !isZAM(pk.Stages[i+1].Tf) {
				logger.Error.Printf("В плане координации %d нарушена структура", pk.Pk)
				return newPk
			}
			//Блок основная замещающая
			skip = true
			dv = pk.Stages[i+1]
			if dv.Trs {
				dv.Start -= pk.Shift
				r := v.Stop - v.Start
				v.Stop = v.Start + dv.Dt + r
				dv.Stop = dv.Start + dv.Dt + r
			}

		} else {
			//обычные фазы
			if v.Trs {
				r := v.Stop - v.Start
				v.Stop = v.Start + v.Dt + r
			}
		}

		newPk.Stages = append(newPk.Stages, v)
		if skip {
			newPk.Stages = append(newPk.Stages, dv)
		}
	}
	newPk.Stages = append(newPk.Stages, tail...)
	logger.Info.Printf("out %s", toSting(newPk))
	return newPk
}
func toSting(pk binding.SetPk) string {
	res := fmt.Sprintf("shift=%d {", pk.Shift)
	for _, v := range pk.Stages {
		if v.Start == 0 && v.Stop == 0 {
			continue
		}
		res += fmt.Sprintf("[%d-%d f=%d t=%d dt=%v pl=%v trs=%v]", v.Start, v.Stop, v.Number, v.Tf, v.Dt, v.Plus, v.Trs)
	}
	res += "}"
	return res

}
func isTVP(tf int) bool {
	// Tf     int `json:"tf"`    //Тип фазы 0 -простая
	// 1 - МГР
	// 2 - 1ТВП
	// 3 - 2ТВП
	// 4 - 1,2ТВП
	// 5 - Зам 1 ТВП
	// 6 - Зам 2 ТВП
	// 7 - Зам
	// 8 - МДК
	// 9 - ВДК

	if tf == 2 || tf == 3 || tf == 4 {
		return true
	}
	return false
}
func isZAM(tf int) bool {
	// Tf     int `json:"tf"`    //Тип фазы 0 -простая
	// 1 - МГР
	// 2 - 1ТВП
	// 3 - 2ТВП
	// 4 - 1,2ТВП
	// 5 - Зам 1 ТВП
	// 6 - Зам 2 ТВП
	// 7 - Зам
	// 8 - МДК
	// 9 - ВДК

	if tf == 5 || tf == 6 || tf == 7 {
		return true
	}
	return false
}
