package worker

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/device"
)

type NowState struct {
	Source int
	PK     int
	NK     int
	CK     int
	DU     int
	Phase  int
}

func Worker() {
	go device.Device()
	if data.DataValue.Controller.Base {
		data.ToDevice <- 0
	}
	tik := time.NewTicker(1 * time.Second)
	for {
		select {
		case cmd := <-data.Commands:
			logger.Info.Printf("Команда %v", cmd)
			switch cmd.Command {
			case 5:
				//Смена плана ПК
				data.DataValue.SetPK(cmd.Parametr)
			case 6:
				//Смена НК
				data.DataValue.SetNK(cmd.Parametr)
			case 7:
				//Смена CК
				data.DataValue.SetCK(cmd.Parametr)
			case 9:
				//Смена ДУ
				data.DataValue.SetDU(cmd.Parametr)

			}
		case ars := <-data.Arrays:
			data.DataValue.SetArrays(ars)
		case <-tik.C:
		case dev := <-data.FromDevice:
			logger.Info.Printf("От устройства %v", dev)

		}
	}
}
