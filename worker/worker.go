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

var dk = pudge.DK{RDK: 5, FDK: 0, DDK: 4, EDK: 0, PDK: false, EEDK: 0, ODK: false, LDK: 0, FTUDK: 0, TDK: 0, FTSDK: 0, TTCDK: 0}
var endDU time.Timer
var ctrlDU = time.Duration(60 * time.Second)
var isDU = false
var isPromtakt = false

func Worker() {
	go device.Device()
	data.ToDevice <- 12 //Кругом красный
	dk.FDK = 12
	data.DataValue.SetDK(dk)
	data.ToServer <- 0
	time.Sleep(3 * time.Second)
	dk.FDK = 0
	dk.DDK = data.USDK
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
						endDU = *time.NewTimer(ctrlDU)
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
						dk.FDK = 0
						dk.TDK = 0
						dk.TTCDK = 0
						data.DataValue.SetDK(dk)
						data.ToDevice <- 0
						data.ToServer <- 0
					} else {
						isDU = true
						dk.FDK = cmd.Parametr
						dk.TDK = 0
						dk.TTCDK = 0
						data.DataValue.SetDK(dk)
						data.ToDevice <- cmd.Parametr
						data.ToServer <- 0
					}
				} else {
					//Выключаем ДУ производим выбор нового плана
					// Но пока для  отладки снова врубаем ЛР
					dk.FDK = 0
					dk.TDK = 0
					dk.TTCDK = 0
					data.DataValue.SetDK(dk)
					data.ToDevice <- 0
					data.ToServer <- 0
				}
			}
		case <-endDU.C:
			//Перестали удерживать ДУ
			data.Commands <- data.InternalCmd{Source: data.USDK, Command: 9, Parametr: 9}
		case ars := <-data.Arrays:
			data.DataValue.SetArrays(ars)
			// Тут нужно все заново выбрать
		case <-tik.C:
			dk.TDK++
			dk.TTCDK++
		case dev := <-data.FromDevice:
			logger.Info.Printf("От устройства %v", dev)
			if dev.Phase == 9 && !isPromtakt {

				dk.TDK = 0
				dk.TTCDK = 0
				isPromtakt = true
			}
			if dev.Phase != 9 {
				isPromtakt = false
			}
			dk.FDK = dev.Need
			dk.FTSDK = dev.Need
			dk.FTSDK = dev.Phase

		}
	}
}
