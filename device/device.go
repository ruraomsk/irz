package device

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
)

// Исполнение команд на устройстве
var promTakt = 5
var phases = []int{10, 20, 30, 40}
var lR = true
var phaseLR = 1
var endPhase time.Timer
var tick = time.NewTicker(1 * time.Second)
var state = data.StatusDevice{Door: false, Lamp: 0, PhaseTC: 0, PhaseTU: 0}

func Device() {
	if lR {
		state.PhaseTC = 1
		state.PhaseTU = 0
		// endPhase = *time.NewTimer(time.Duration(phases[0]) * time.Second)
		endPhase = *time.NewTimer(10 * time.Second)
	}
	logger.Info.Print("Начинаем основной цикл device с ЛР")
	for {
		select {
		case <-tick.C:
			state.TimeTU++
			state.TimeTC++
			// data.FromDevice <- state
		case in := <-data.ToDevice:
			logger.Debug.Printf("from worker %v", in)
			switch in {
			case 0:
				state.PhaseTU = 0
				if lR {
					continue
				}
				if state.PhaseTC != 1 {
					makePromTakt()
				}
				lR = true
				phaseLR = 1
				state.PhaseTC = phaseLR
				// endPhase = *time.NewTimer(time.Duration(phases[phaseLR-1]-promTakt) * time.Second)
				endPhase = *time.NewTimer(10 * time.Second)
			case 10:
				stopLRifNeed()
				state.PhaseTU = 10
				state.PhaseTC = 10
				data.FromDevice <- state
			case 11:
				stopLRifNeed()
				state.PhaseTU = 11
				state.PhaseTC = 11
				data.FromDevice <- state
			case 12:
				stopLRifNeed()
				state.PhaseTU = 12
				state.PhaseTC = 12
				data.FromDevice <- state
			default:
				state.PhaseTU = in
				if in > 0 && in <= len(phases) {
					stopLRifNeed()
					if state.PhaseTC != in {
						makePromTakt()
					} else {
						state.NewPhase()
					}
					state.PhaseTC = in
					data.FromDevice <- state
				} else {
					if !lR {
						if state.PhaseTC != 1 {
							makePromTakt()
						}
						phaseLR = 1
						state.PhaseTC = phaseLR
						// endPhase = *time.NewTimer(time.Duration(phases[phaseLR-1]-promTakt) * time.Second)
						endPhase = *time.NewTimer(10 * time.Second)
						data.FromDevice <- state
					}
				}
			}

		case <-endPhase.C:
			if !lR {
				continue
			}
			data.FromDevice <- state

			phaseLR++
			if phaseLR > len(phases) {
				phaseLR = 1
			}
			logger.Info.Printf("Новая фаза ЛР %d", phaseLR)
			makePromTakt()
			state.PhaseTC = phaseLR
			// endPhase = *time.NewTimer(time.Duration(phases[phaseLR-1]-promTakt) * time.Second)
			endPhase = *time.NewTimer(10 * time.Second)
		}
	}

}
func makePromTakt() {
	tick.Stop()
	state.NewPhase()
	for i := 0; i < promTakt; i++ {
		time.Sleep(time.Second)
	}
	state.PhaseTC = 9
	state.TimeTC += promTakt
	state.TimeTU += promTakt
	data.FromDevice <- state
	tick = time.NewTicker(1 * time.Second)
}
func stopLRifNeed() {
	state.NewPhase()
	if lR {
		endPhase.Stop()
		lR = false
	}
}
