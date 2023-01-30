package device

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
)

// Исполнение команд на устройстве
var promTakt = 5
var phases = []int{10, 20, 30}
var lR = true
var phaseLR = 1
var endPhase time.Timer
var tick = time.NewTicker(1 * time.Second)
var state = data.StatusDevice{Door: false, Lamp: 0, Phase: 0, PhaseTC: 0, PhaseTU: 0}

func Device() {
	if lR {
		state.Phase = 1
		state.PhaseTC = 1
		state.PhaseTU = 1
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
				if lR {
					continue
				}
				phaseLR = 1
				if state.PhaseTC != 1 {
					makePromTakt()
				}
				lR = true

				state.PhaseTC = phaseLR
				state.PhaseTU = state.PhaseTC
				state.Phase = phaseLR
				endPhase = *time.NewTimer(time.Duration(phases[phaseLR-1]-promTakt) * time.Second)
				// endPhase = *time.NewTimer(20 * time.Second)
			case 10:
				stopLRifNeed()
				state.Phase = 10
				state.PhaseTU = 10
				state.PhaseTC = 10
				data.FromDevice <- state
			case 11:
				stopLRifNeed()
				state.Phase = 11
				state.PhaseTU = 11
				state.PhaseTC = 11
				data.FromDevice <- state
			case 12:
				stopLRifNeed()
				state.Phase = 12
				state.PhaseTU = 12
				state.PhaseTC = 12
				data.FromDevice <- state
			default:
				data.FromDevice <- state
				state.PhaseTU = in
				if in > 0 && in <= len(phases) {

					state.PhaseTU = in
					stopLRifNeed()
					if state.PhaseTC != in {
						phaseLR = in
						makePromTakt()
					} else {
						state.NewPhase()
					}
					state.PhaseTC = in
					state.Phase = in
					data.FromDevice <- state
				} else {
					if !lR {
						state.PhaseTU = phaseLR
						if state.PhaseTC != 1 {
							makePromTakt()
						}
						phaseLR = 1
						state.PhaseTC = phaseLR
						state.Phase = phaseLR
						endPhase = *time.NewTimer(time.Duration(phases[phaseLR-1]-promTakt) * time.Second)
						// endPhase = *time.NewTimer(10 * time.Second)
						data.FromDevice <- state
					}
				}
			}

		case <-endPhase.C:
			if !lR {
				continue
			}
			phaseLR++
			if phaseLR > len(phases) {
				phaseLR = 1
			}
			makePromTakt()
			logger.Info.Printf("Новая фаза ЛР %d", phaseLR)
			endPhase = *time.NewTimer(time.Duration(phases[phaseLR-1]-promTakt) * time.Second)
			// endPhase = *time.NewTimer(20 * time.Second)
		}
	}

}
func makePromTakt() {
	tick.Stop()
	state.Phase = 9
	data.FromDevice <- state
	time.Sleep(time.Second)
	state.NewPhase()
	state.PhaseTC = 9
	state.PhaseTU = phaseLR
	state.Phase = phaseLR
	for i := 0; i < promTakt; i++ {
		time.Sleep(time.Second)
	}
	state.TimeTC = promTakt
	state.TimeTU = promTakt
	data.FromDevice <- state
	state.TimeTC -= promTakt
	state.TimeTU -= promTakt
	state.PhaseTC = phaseLR
	// state.TimeTC = promTakt
	// state.TimeTU += promTakt
	// state.TimeTC = promTakt
	// state.TimeTC = 0
	tick = time.NewTicker(1 * time.Second)
}
func stopLRifNeed() {
	state.NewPhase()
	if lR {
		endPhase.Stop()
		lR = false
	}
}
