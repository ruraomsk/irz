package worker

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
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

var oldDK pudge.DK
var dk = pudge.DK{RDK: 5, FDK: 0, DDK: 4, EDK: 0, PDK: false, EEDK: 0, ODK: false, LDK: 0, FTUDK: 0, TDK: 0, FTSDK: 0, TTCDK: 0}
var endDU time.Timer

// var ctrlDU = time.Duration(60 * time.Second)
var isDU = false
var isPromtakt = false

func Worker() {
	go device.Device()
	data.ToDevice <- 12 //Кругом красный
	dk.RDK = 9
	dk.FDK = 12
	data.DataValue.SetDK(dk)
	oldDK = dk
	data.ToServer <- 0
	time.Sleep(3 * time.Second)
	logger.Info.Print("КК закончили")
	dk.FDK = 0
	dk.RDK = 6
	dk.DDK = data.USDK
	oldDK = dk
	data.DataValue.SetDK(dk)
	if data.DataValue.Controller.Base {
		data.ToDevice <- 0
		data.ToServer <- 0
	} else {
		//Выбираем согласно планов
		// но пока все равно в ЛР
		data.ToDevice <- 0
		data.ToServer <- 0
	}
	tik := time.NewTicker(1 * time.Second)
	logger.Info.Print("Начинаем основной цикл worker")
	for {
		select {
		case cmd := <-data.Commands:
			logger.Info.Printf("Команда %v", cmd)
			dk.DDK = cmd.Source
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
				if isDU {
					if cmd.Parametr != dk.FDK {
						endDU.Stop()
						isDU = false
					} else {
						//Обновляем время удержания
						endDU.Stop()
						endDU = *time.NewTimer(time.Minute)
						continue
					}
				}
				isPromtakt = false
				if cmd.Parametr != 9 {
					if cmd.Parametr == 0 {
						//Перевод в ЛР
						if isDU {
							endDU.Stop()
						}
						isDU = false
						dk.RDK = 6
						dk.FDK = 1
						data.DataValue.SetDK(dk)
						data.ToDevice <- 0
						data.ToServer <- 0
					} else {
						dk.RDK = 4
						dk.FDK = cmd.Parametr
						data.DataValue.SetDK(dk)
						data.ToDevice <- cmd.Parametr
						data.ToServer <- 0
						if cmd.Parametr < 9 {
							//Держим фазу
							isDU = true
							endDU = *time.NewTimer(time.Minute)
						}
					}
				} else {
					//Выключаем ДУ производим выбор нового плана
					// Но пока для  отладки снова врубаем ЛР
					dk.RDK = 6
					dk.FDK = 1
					data.DataValue.SetDK(dk)
					data.ToDevice <- 0
					data.ToServer <- 0
				}
			}
		case <-endDU.C:
			//Перестали удерживать ДУ
			data.Commands <- data.InternalCmd{Source: data.USDK, Command: 9, Parametr: 9}
		case ars := <-data.Arrays:
			logger.Debug.Print("Записываем привязку")
			data.DataValue.SetArrays(ars)
			logger.Debug.Print("Записали привязку")
			// Тут нужно все заново выбрать
		case <-tik.C:
			// dk.TDK += 1
			// dk.TTCDK += 1
			// if isChangeDK() {
			// 	data.DataValue.SetDK(dk)
			// 	oldDK = dk
			// 	data.ToServer <- 0
			// }
		case dev := <-data.FromDevice:
			dk = data.DataValue.Controller.DK
			logger.Info.Printf("От устройства %v", dev)
			dk.FDK = dev.Phase
			dk.TTCDK = dev.TimeTC
			dk.TDK = dev.TimeTU
			dk.FTSDK = dev.PhaseTC
			dk.FTUDK = dev.PhaseTU
			if dk.RDK == 5 || dk.RDK == 6 {
				dk.FTUDK = 0
				dk.TDK = 0
			}
			data.DataValue.SetDK(dk)
			logger.Info.Printf("%v", data.DataValue.Controller.DK)
			data.ToServer <- 0
		}
	}
}
func isChangeDK() bool {
	if oldDK.RDK != dk.RDK {
		return true
	}
	if oldDK.DDK != dk.DDK {
		return true
	}
	if oldDK.EDK != dk.EDK {
		return true
	}
	if oldDK.EEDK != dk.EEDK {
		return true
	}
	if oldDK.FDK != dk.FDK {
		return true
	}
	if oldDK.FTSDK != dk.FTSDK {
		return true
	}
	if oldDK.LDK != dk.LDK {
		return true
	}
	if oldDK.PDK != dk.PDK {
		return true
	}
	return false

}
