package device

import (
	"time"

	"github.com/ruraomsk/irz/data"
)

// Исполнение команд на устройстве
var promTakt = 5
var phases = []int{10, 20, 30, 40}
var lR = true
var phaseLR = 1
var endPhase time.Timer
var tick = time.NewTicker(1 * time.Second)
var state = data.StatusDevice{Door: false, Lamp: 0, Phase: 0, Need: 0}

func Device() {
	if lR {
		state.Phase = 1
		state.Need = 0
		endPhase = *time.NewTimer(time.Duration(phases[0]) * time.Second)
	}
	for {
		select {
		case <-tick.C:
			data.FromDevice <- state
		case in := <-data.ToDevice:
			switch in {
			case 0:
				if lR {
					continue
				}
				state.Need = 0
				if state.Phase != 1 {
					makePromTakt()
				}
				phaseLR = 1
				state.Phase = phaseLR
				endPhase = *time.NewTimer(time.Duration(phases[phaseLR-1]-promTakt) * time.Second)
				data.FromDevice <- state
			case 0x0a:
				stopLRifNeed()
				state.Need = 10
				state.Phase = 10
				data.FromDevice <- state
			case 0x0b:
				stopLRifNeed()
				state.Need = 11
				state.Phase = 11
				data.FromDevice <- state
			default:
				state.Need = in
				if in > 0 && in <= len(phases) {
					stopLRifNeed()
					if state.Phase != in {
						makePromTakt()
					}
					state.Phase = in
					data.FromDevice <- state
				} else {
					if !lR {
						if state.Phase != 1 {
							makePromTakt()
						}
						phaseLR = 1
						state.Phase = phaseLR
						endPhase = *time.NewTimer(time.Duration(phases[phaseLR-1]-promTakt) * time.Second)
						data.FromDevice <- state
					}
				}
			}

		case <-endPhase.C:
			phaseLR++
			if phaseLR > len(phases) {
				phaseLR = 1
			}
			makePromTakt()
			state.Phase = phaseLR
			endPhase = *time.NewTimer(time.Duration(phases[phaseLR-1]-promTakt) * time.Second)
			data.FromDevice <- state
		}
	}

}
func makePromTakt() {
	tick.Stop()
	state.Phase = 9
	for i := 0; i < promTakt; i++ {
		data.FromDevice <- state
		time.Sleep(time.Second)
	}
	tick = time.NewTicker(1 * time.Second)
}
func stopLRifNeed() {
	if lR {
		endPhase.Stop()
		lR = false
	}
}
