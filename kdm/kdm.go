package kdm

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/setup"
	"github.com/simonvetter/modbus"
)

var client *modbus.ModbusClient
var err error

type KdmStatus struct {
	Phase     int
	PhaseTU   int
	Time      int
	Status    int
	NKPogramm int
	Blink     int
	PhaseRU   int
	BadKeys   bool
	SetKeys   []uint16
	Lamp      int
	Connect   bool
}

var state = data.StatusDevice{Door: false, Lamp: 0, Phase: 0, PhaseTC: 0, PhaseTU: 0, Connect: false}
var sendpromtakt = false
var sendphase = false
var savepromtakt = false
var savestate = state

var statusKdm KdmStatus
var lastcmd = -1

func Kdm() {
	statusKdm.SetKeys = make([]uint16, 32)
	for !state.Connect {
		time.Sleep(time.Second)
		// for an RTU (serial) device/bus
		client, err = modbus.NewClient(&modbus.ClientConfiguration{
			URL:      setup.Set.Modbus.Device,         //"rtu:///dev/ttyUSB0",
			Speed:    uint(setup.Set.Modbus.BaudRate), //19200,                   // default
			DataBits: 8,                               // default, optional
			Parity:   modbus.PARITY_NONE,              // default, optional
			StopBits: 2,                               // default if no parity, optional
			Timeout:  300 * time.Millisecond,
		})

		if err != nil {
			logger.Error.Printf("modbus %v", err.Error())
			continue
			// error out if client creation failed
		}
		client.SetUnitId(uint8(setup.Set.Modbus.UId))
		// now that the client is created and configured, attempt to connect
		err = client.Open()
		if err != nil {
			logger.Error.Printf("modbus open %v", err.Error())
			continue
		}
		data.DataValue.Connect = true
		state.Connect = true
		workModbus()
		state.Connect = false
		data.FromDevice <- state
	}
}
func workModbus() {
	defer client.Close()
	//Начинаем работать
	newSending()
	err = setLocal()
	if err != nil {
		logger.Error.Print(err.Error())
		return
	}

	state.PhaseTC = 0
	state.Phase = 0
	var tick = time.NewTicker(1 * time.Second)
	for {
		select {
		case <-tick.C:
			if state.PhaseTC == 10 || state.PhaseTC == 11 {
				continue
			}

			state.TimeTU++
			state.TimeTC++
			if state.TimeTC > 255 {
				state.TimeTC = 0
			}
			err = getStatus()
			if err != nil {
				logger.Error.Print(err.Error())
				return
			}
			// if state.PhaseTC != 0 {
			// 	if statusKdm.Phase == 9 {
			// 		if !sendpromtakt {
			// 			state.Phase = 9
			// 			state.PhaseTU = statusKdm.PhaseTU
			// 			state.PhaseTC = statusKdm.PhaseTU
			// 			data.FromDevice <- state
			// 			state.Phase = statusKdm.PhaseTU
			// 			state.NewPhase()
			// 			sendpromtakt = true
			// 			sendphase = false
			// 		}
			// 	}
			// 	if statusKdm.Phase != 9 {
			// 		if !sendphase {
			// 			state.Phase = statusKdm.Phase
			// 			state.PhaseTU = statusKdm.PhaseTU
			// 			state.PhaseTC = statusKdm.PhaseTU
			// 			data.FromDevice <- state
			// 			state.TimeTU = 0
			// 			sendphase = true
			// 			sendpromtakt = false
			// 		}
			// 	}
			// 	continue
			// }
			// logger.Info.Printf("От устройства фаза %d %v %v", statusKdm.Phase, sendpromtakt, sendphase)
			if statusKdm.Phase == 9 {
				if !savepromtakt {
					state.Phase = 9
					state.PhaseTU = statusKdm.PhaseTU
					savestate = state
					// data.FromDevice <- state
					state.Phase = statusKdm.PhaseTU
					state.NewPhase()
					savepromtakt = true
					sendphase = false
					sendpromtakt = false
				}
				if statusKdm.PhaseTU != savestate.PhaseTU && !sendpromtakt {
					savestate.PhaseTU = statusKdm.PhaseTU
					sendpromtakt = true
					data.FromDevice <- savestate
				}
			}
			if statusKdm.Phase != 9 {
				savepromtakt = false
				if !sendphase {
					state.Phase = statusKdm.Phase
					state.PhaseTU = statusKdm.PhaseTU
					data.FromDevice <- state
					// state.TimeTU = 0
					state.TimeTC = 0
					sendphase = true
					savepromtakt = false
				}
			}
			if statusKdm.Lamp != state.Lamp || statusKdm.Connect != state.Connect {
				state.Lamp = statusKdm.Lamp
				state.Connect = statusKdm.Connect
				data.FromDevice <- state
			}

		case in := <-data.ToDevice:
			if in == lastcmd {
				continue
			}
			logger.Debug.Printf("from worker %v", in)
			lastcmd = in
			switch in {
			case 0:
				newSending()
				err = setLocal()
				if err != nil {
					logger.Error.Print(err.Error())
					return
				}
				state.PhaseTC = 0
			case 10:
				//ЖМ
				newSending()
				err = setBlink()
				if err != nil {
					logger.Error.Print(err.Error())
					return
				}
				state.Phase = 10
				state.PhaseTC = 10
				state.PhaseTU = 10
				data.FromDevice <- state
				state.NewPhase()
			case 11:
				//OC
				err = setOS()
				if err != nil {
					logger.Error.Print(err.Error())
					return
				}
				newSending()
				state.Phase = 11
				state.PhaseTC = 11
				state.PhaseTU = 11
				state.NewPhase()
				data.FromDevice <- state
			default:
				if in > 0 && in < 9 {
					newSending()
					state.PhaseTC = in
					err = setPhase(in)
					if err != nil {
						logger.Error.Print(err.Error())
						return
					}
				}
			}

		}

	}
}
func stopRUandBlink() error {
	// снимаем команду выключения
	err = client.WriteRegister(0x0a, 0)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	// снимаем команду перейти в РУ
	err = client.WriteRegister(0x0c, 0)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	// снимаем команду перейти в ЖМ
	err = client.WriteRegister(0x0d, 0)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	// снимаем время фазы
	err = client.WriteRegister(0x0c21, 0)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	// снимаем фазу РУ
	err = client.WriteRegister(0x0b, 0)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	return nil
}
func setBlink() error {
	err = stopRUandBlink()
	if err != nil {
		return err
	}
	err = client.WriteRegister(0x0d, 1)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	return nil
}
func setLocal() error {
	err = stopRUandBlink()
	if err != nil {
		return err
	}
	// ставим команду ЛР
	err = client.WriteRegister(4, 1)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	return nil
}
func setPhase(phase int) error {
	// включаем команду перейти в РУ
	err = client.WriteRegister(0x0c, 1)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	// ставим время фазы в РУ
	err = client.WriteRegister(0x0c21, 0xff)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	// Передаем номер фазы
	err = client.WriteRegister(0x0b, uint16(phase))
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	return nil
}
func setOS() error {
	err = stopRUandBlink()
	if err != nil {
		return err
	}
	// включаем даем команду отключения
	err = client.WriteRegister(0x0a, 1)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	return nil
}
func getStatus() error {
	//  читаем первый блок
	var reg16s []uint16
	reg16s, err = client.ReadRegisters(0, 0x16, modbus.HOLDING_REGISTER)
	if err != nil {
		return fmt.Errorf("modbus read %v", err.Error())
	}
	statusKdm.PhaseTU = int(reg16s[3] >> 8)
	// 0 dhtvz это Промтакт
	statusKdm.Time = int(reg16s[3] & 0xff)
	// logger.Debug.Printf("phase %d time %d", statusKdm.PhaseTU, statusKdm.Time)
	if statusKdm.Time == 0 {
		statusKdm.Phase = 9
	} else {
		statusKdm.Phase = statusKdm.PhaseTU
	}
	statusKdm.NKPogramm = int(reg16s[4] >> 8)
	statusKdm.Status = int(reg16s[4] & 0xff)
	switch statusKdm.Status {
	case 0:
		//Отладка контроллера
		statusKdm.Phase = 11
	case 2:
		statusKdm.Connect = false
	case 3:
		statusKdm.Lamp = statusKdm.PhaseTU
	case 4:
		//Выключен контроллер
		statusKdm.Phase = 11
	}
	statusKdm.PhaseRU = int(reg16s[0x0b])

	statusKdm.Blink = int(reg16s[0x0d])

	flag := reg16s[5]
	statusKdm.BadKeys = true
	for i := 0; i < 16; i++ {
		if flag&1 == 0 {
			statusKdm.BadKeys = statusKdm.BadKeys && false
		}
		flag = flag >> 1
	}
	return nil
}
func newSending() {
	sendpromtakt = false
	sendphase = false
	savepromtakt = false
}
