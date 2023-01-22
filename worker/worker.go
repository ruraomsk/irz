package worker

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
)

func Worker() {
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

		}
	}
}
