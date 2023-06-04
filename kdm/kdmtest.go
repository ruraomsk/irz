package kdm

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/simonvetter/modbus"
)

func KdmTest() {

	// for an RTU (serial) device/bus
	client, err = modbus.NewClient(&modbus.ClientConfiguration{
		URL:      "rtu:///dev/ttyUSB0",
		Speed:    19200,              // default
		DataBits: 8,                  // default, optional
		Parity:   modbus.PARITY_NONE, // default, optional
		StopBits: 2,                  // default if no parity, optional
		Timeout:  300 * time.Millisecond,
	})

	if err != nil {
		logger.Error.Printf("modbus %v", err.Error())
		return
		// error out if client creation failed
	}
	client.SetUnitId(247)
	statusKdm.SetKeys = make([]uint16, 32)
	// now that the client is created and configured, attempt to connect
	err = client.Open()
	if err != nil {
		logger.Error.Printf("modbus open %v", err.Error())
		return
	}
	defer client.Close()
	//Начинаем работать
	err = client.WriteRegister(4, 1)
	if err != nil {
		logger.Error.Printf("modbus write %v", err.Error())
		return
	}
	for {
		stopRUandBlink()
		for t := 0; t < 90; t++ {
			getStatus()
			logger.Info.Printf("Фаза %d Время %d ", statusKdm.Phase, statusKdm.Time)
			time.Sleep(1 * time.Second)
		}
		stopRUandBlink()
		logger.Info.Print("Переходим в ЖМ")
		setBlink()
		for t := 0; t < 30; t++ {
			getStatus()
			logger.Info.Printf("Фаза %d Время %d ЖМ %d ", statusKdm.Phase, statusKdm.Time, statusKdm.Blink)
			time.Sleep(1 * time.Second)
		}
		stopRUandBlink()
		for phase := 1; phase <= 3; phase++ {
			setPhase(phase, 255)
			for t := 0; t < 30; t++ {
				getStatus()
				logger.Info.Printf("Фаза %d Время %d Ру %d ", statusKdm.Phase, statusKdm.Time, statusKdm.PhaseRU)
				time.Sleep(1 * time.Second)
			}
		}
	}
}

// func stopRUandBlink() {
// 	err = client.WriteRegister(0x0c, 0)
// 	if err != nil {
// 		logger.Error.Printf("modbus write %v", err.Error())
// 		return
// 	}
// 	err = client.WriteRegister(0x0d, 0)
// 	if err != nil {
// 		logger.Error.Printf("modbus write %v", err.Error())
// 		return
// 	}
// 	err = client.WriteRegister(0x0c21, 0)
// 	if err != nil {
// 		logger.Error.Printf("modbus write %v", err.Error())
// 		return
// 	}

// 	err = client.WriteRegister(0x0b, 0)
// 	if err != nil {
// 		logger.Error.Printf("modbus write %v", err.Error())
// 		return
// 	}

// }
// func setBlink() {
// 	err = client.WriteRegister(0x0d, 1)
// 	if err != nil {
// 		logger.Error.Printf("modbus write %v", err.Error())
// 		return
// 	}
// }
// func setPhase(phase int, timeru int) {
// 	err = client.WriteRegister(0x0c, 1)
// 	if err != nil {
// 		logger.Error.Printf("modbus write %v", err.Error())
// 		return
// 	}
// 	err = client.WriteRegister(0x0c21, uint16(timeru))
// 	if err != nil {
// 		logger.Error.Printf("modbus write %v", err.Error())
// 		return
// 	}

// 	err = client.WriteRegister(0x0b, uint16(phase))
// 	if err != nil {
// 		logger.Error.Printf("modbus write %v", err.Error())
// 		return
// 	}

// }
// func getKeys() {
// 	var reg16s []uint16
// 	reg16s, err = client.ReadRegisters(0x400, 32, modbus.HOLDING_REGISTER)
// 	if err != nil {
// 		logger.Error.Printf("modbus read %v", err.Error())
// 		return
// 	}
// 	Status.SetKeys = reg16s
// 	// logger.Info.Printf("keys %v", reg16s)

// }
// func getStatus() {
// 	//  читаем первый блок
// 	var reg16s []uint16
// 	reg16s, err = client.ReadRegisters(0, 0x16, modbus.HOLDING_REGISTER)
// 	if err != nil {
// 		logger.Error.Printf("modbus read %v", err.Error())
// 		return
// 	}
// 	Status.Phase = int(reg16s[3] >> 8)
// 	Status.Time = int(reg16s[3] & 0xff)
// 	Status.NKPogramm = int(reg16s[4] >> 8)
// 	Status.Status = int(reg16s[4] & 0xff)
// 	Status.PhaseRU = int(reg16s[0x0b])
// 	Status.Blink = int(reg16s[0x0d])
// 	flag := reg16s[5]
// 	Status.BadKeys = make([]bool, 0)
// 	for i := 0; i < 16; i++ {
// 		if flag&1 == 1 {
// 			Status.BadKeys = append(Status.BadKeys, true)
// 		} else {
// 			Status.BadKeys = append(Status.BadKeys, false)
// 		}
// 		flag = flag >> 1
// 	}

// }
