package radar

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/setup"
	"github.com/simonvetter/modbus"
)

var client *modbus.ModbusClient

func pusher() {
	var errClient error
	var regs []uint16
	client, errClient = modbus.NewClient(&modbus.ClientConfiguration{
		URL:     fmt.Sprintf("tcp://127.0.0.1:%d", setup.Set.ModbusRadar.Port),
		Timeout: time.Second,
		Logger:  logger.Info,
	})
	if errClient != nil {
		logger.Error.Printf("Не могу создать клиента %v", errClient)
		return
	}
	client.SetUnitId(1)
	for {
		errClient = client.Open()
		if errClient != nil {
			logger.Error.Printf("Не могу открыть клиента %v", errClient)
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
	for {
		time.Sleep(1 * time.Second)
		regs, errClient = client.ReadRegisters(0, uint16(setup.Set.ModbusRadar.Chanels), modbus.HOLDING_REGISTER)
		if errClient != nil {
			logger.Error.Printf("Не смог прочитать %v ", errClient)
			break
		}
		for i := 0; i < setup.Set.ModbusRadar.Chanels; i++ {
			regs[i] = uint16(rand.Intn(3))

		}
		errClient = client.WriteRegisters(0, regs)
		if errClient != nil {
			logger.Error.Printf("Не смог отправить %v %v", regs, errClient)
			break
		}
	}
	client.Close()
}
