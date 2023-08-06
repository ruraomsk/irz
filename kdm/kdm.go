package kdm

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/setup"
	"github.com/ruraomsk/irz/stat"
	"github.com/simonvetter/modbus"
)

var client *modbus.ModbusClient
var err error

type KdmStatus struct {
	Phase      int
	PhaseTU    int
	Time       int
	Status     int
	NKPogramm  int
	Blink      int
	PhaseRU    int
	BadKeys    bool
	SetKeys    []uint16
	Lamp       int
	Connect    bool
	isBlink    bool
	isOS       bool
	isSetPhase bool
}
type WriteCmd struct {
	Request Request
	Data    []uint16
}
type Request struct {
	Start  uint16
	Lenght uint16
}
type Replay struct {
	Request Request
	Status  error
	Data    []uint16
}
type Info struct {
	Phase  int
	Lenght int
}

var State = data.StatusDevice{Door: false, Lamp: 0, Phase: 0, PhaseTC: 0, PhaseTU: 0, Connect: false}
var sendpromtakt = false
var sendphase = false
var savepromtakt = false
var savestate = State

var statusKdm KdmStatus
var lastcmd = -1
var RequestChan chan Request
var ReplayChan chan Replay
var WriteCmdChan chan WriteCmd
var InfoChan chan Info

func Kdm() {

	RequestChan = make(chan Request)
	ReplayChan = make(chan Replay)
	WriteCmdChan = make(chan WriteCmd)
	InfoChan = make(chan Info)

	statusKdm.SetKeys = make([]uint16, 32)

	for !State.Connect {
		data.DataValue.Connect = false
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
		workModbus()
		logger.Error.Printf("Завершили обмен с ModBus")
		data.DataValue.Connect = false
		State.Connect = false
		data.FromDevice <- State
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

	data.DataValue.Connect = true
	State.Connect = true
	data.FromDevice <- State

	State.PhaseTC = 0
	State.PhaseTU = 0
	State.Phase = 0
	var tick = time.NewTicker(1 * time.Second)
	for {
		select {
		case req := <-RequestChan:
			ReplayChan <- readData(req)
		case wr := <-WriteCmdChan:
			ReplayChan <- writeData(wr)
		case <-tick.C:
			err = getStatus()
			if err != nil {
				logger.Error.Print(err.Error())
				return
			}
			if State.PhaseTC == 10 || State.PhaseTC == 11 {
				continue
			}

			State.TimeTU++
			State.TimeTC++
			if State.TimeTC > 255 {
				State.TimeTC = 0
			}
			if State.TimeTU > 255 {
				State.TimeTU = 0
			}
			err = getStatus()
			if err != nil {
				logger.Error.Print(err.Error())
				return
			}
			if statusKdm.Phase == 9 {
				if !savepromtakt {
					State.Phase = 9
					State.PhaseTU = statusKdm.PhaseTU
					savestate = State
					// data.FromDevice <- state
					State.Phase = statusKdm.PhaseTU
					savepromtakt = true
					sendphase = false
					sendpromtakt = false
				}
			}
			if statusKdm.PhaseTU != savestate.PhaseTU && !sendpromtakt {
				savestate.PhaseTU = statusKdm.PhaseTU
				sendpromtakt = true
				savestate.TimeTC = State.TimeTC - savestate.TimeTC
				data.FromDevice <- savestate
				State.TimeTU = State.TimeTU - savestate.TimeTU
				State.TimeTC = 0

			}
			if statusKdm.Phase != 9 {
				if !sendphase {
					State.Phase = statusKdm.Phase
					State.PhaseTU = statusKdm.PhaseTU
					data.FromDevice <- State
					sendphase = true
					savepromtakt = false
				}
			}
			if statusKdm.Lamp != State.Lamp || statusKdm.Connect != State.Connect {
				State.Lamp = statusKdm.Lamp
				State.Connect = statusKdm.Connect
				data.FromDevice <- State
			}

		case in := <-data.ToDevice:
			// logger.Debug.Printf("set phase %v", in&0xff)
			if in&0xff == lastcmd {
				continue
			}
			InfoChan <- Info{Phase: in & 0xff, Lenght: in >> 8}
			lastcmd = in & 0xff
			switch in {
			case 0:
				newSending()
				err = setLocal()
				if err != nil {
					logger.Error.Print(err.Error())
					return
				}
				State.PhaseTC = 0
				State.PhaseTU = 0
			case 10:
				//ЖМ
				newSending()
				err = setBlink()
				if err != nil {
					logger.Error.Print(err.Error())
					return
				}
				State.Phase = 10
				State.PhaseTC = 10
				State.PhaseTU = 10
				data.FromDevice <- State
				State.NewPhase()
			case 11:
				//OC
				err = setOS()
				if err != nil {
					logger.Error.Print(err.Error())
					return
				}
				newSending()
				State.Phase = 11
				State.PhaseTC = 11
				State.PhaseTU = 11
				State.NewPhase()
				data.FromDevice <- State
			default:
				lenght := in >> 8
				if lenght == 0 {
					lenght = 0xff
				}
				in = in & 0xff
				if in > 0 && in < 9 {
					// newSending()
					State.PhaseTC = in
					err = setPhase(in, lenght)
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
	// // снимаем команду выключения
	if statusKdm.isBlink {
		// снимаем команду перейти в ЖМ
		err = client.WriteRegister(0x0d, 0)
		if err != nil {
			return fmt.Errorf("modbus write %v", err.Error())
		}
	}
	// // снимаем команду перейти в РУ
	// err = client.WriteRegister(0x0c, 0)
	// if err != nil {
	// 	return fmt.Errorf("modbus write %v", err.Error())
	// }
	// // снимаем время фазы
	// err = client.WriteRegister(0x0c21, 0)
	// if err != nil {
	// 	return fmt.Errorf("modbus write %v", err.Error())
	// }
	// // снимаем фазу РУ
	// err = client.WriteRegister(0x0b, 0)
	// if err != nil {
	// 	return fmt.Errorf("modbus write %v", err.Error())
	// }
	return nil
}
func setBlink() error {
	if statusKdm.isOS {
		// снимаем команду отключения
		err = client.WriteRegister(0x0a, 0)
		if err != nil {
			return fmt.Errorf("modbus write %v", err.Error())
		}
	}
	if statusKdm.isSetPhase {
		err = client.WriteRegister(0x0c, 0)
		if err != nil {
			return fmt.Errorf("modbus write %v", err.Error())
		}
	}
	err = client.WriteRegister(0x0c21, 0xff)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	err = client.WriteRegister(0x0d, 1)
	if err != nil {
		return fmt.Errorf("modbus write %v", err.Error())
	}
	return nil
}
func setLocal() error {
	if statusKdm.isOS {
		// снимаем команду отключения
		err = client.WriteRegister(0x0a, 0)
		if err != nil {
			return fmt.Errorf("modbus write Снять ОС setLocal %v", err.Error())
		}
	}

	if statusKdm.isBlink {
		// снимаем команду перейти в ЖМ
		err = client.WriteRegister(0x0d, 0)
		if err != nil {
			return fmt.Errorf("modbus write Снять ЖМ setLocal %v", err.Error())
		}
		return nil
	}
	if statusKdm.isSetPhase {
		err = client.WriteRegister(0x0c, 0)
		if err != nil {
			return fmt.Errorf("modbus write Снять Фазжу setLocal %v", err.Error())
		}
	}
	err = client.WriteRegister(0x08, 0)
	if err != nil {
		return fmt.Errorf("modbus write Снять 0x08 setLocal %v", err.Error())
	}
	return nil
}
func setPhase(phase int, lenght int) error {
	// включаем команду перейти в РУ
	err = client.WriteRegister(0x0c, 1)
	if err != nil {
		return fmt.Errorf("modbus write РУ %v", err.Error())
	}
	// ставим время фазы в РУ
	err = client.WriteRegister(0x0c21, uint16(lenght))
	if err != nil {
		return fmt.Errorf("modbus write время %v", err.Error())
	}
	// Передаем номер фазы
	err = client.WriteRegister(0x0b, uint16(phase))
	if err != nil {
		return fmt.Errorf("modbus write Фаза %v", err.Error())
	}
	if statusKdm.isBlink {
		// снимаем команду перейти в ЖМ
		err = client.WriteRegister(0x0d, 0)
		if err != nil {
			return fmt.Errorf("modbus write Снять ЖМ %v", err.Error())
		}
	}
	if statusKdm.isOS {
		// снимаем команду отключения
		err = client.WriteRegister(0x0a, 0)
		if err != nil {
			return fmt.Errorf("modbus write Снять ОС %v", err.Error())
		}
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
		return fmt.Errorf("modbus read Статус %v", err.Error())
	}
	statusKdm.isBlink = false
	if reg16s[0x0d] != 0 {
		statusKdm.isBlink = true
	}
	statusKdm.isOS = false
	if reg16s[0x0a] != 0 {
		statusKdm.isOS = true
	}
	statusKdm.PhaseTU = int(reg16s[3] >> 8)
	// 0 dhtvz это Промтакт
	statusKdm.Time = int(reg16s[3] & 0xff)
	// logger.Debug.Printf("phase %d time %d", statusKdm.PhaseTU, statusKdm.Time)
	if statusKdm.Time == 0 {
		statusKdm.Phase = 9
	} else {
		statusKdm.Phase = statusKdm.PhaseTU
		if statusKdm.PhaseTU == 0 {
			statusKdm.Phase = 12
		}
	}
	tvp1 := int(reg16s[2] & 1)
	tvp2 := int(reg16s[2] & 2)
	if tvp2 != 0 {
		tvp2 = 1
	}
	stat.InStat <- stat.OneTick{Nomber: 1, Value: stat.Value{Time: time.Now(), Status: 0, Value: tvp1}}
	stat.InStat <- stat.OneTick{Nomber: 2, Value: stat.Value{Time: time.Now(), Status: 0, Value: tvp2}}

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
	statusKdm.isSetPhase = false
	if reg16s[0x0c] == 1 {
		statusKdm.isSetPhase = true
	}
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
