package worker

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/irz/data"
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
var endDUPhase time.Timer
var toPlan chan data.StatusDevice
var nowPlan = 0

// var ctrlDU = time.Duration(60 * time.Second)
var isDUPhase = false

func Worker() {
	toPlan = make(chan data.StatusDevice, 100)
	endPlan = make(chan interface{})
	for !data.DataValue.Connect {
		dk.EDK = 11
		dk.DDK = 8
		data.DataValue.SetDK(dk)
		time.Sleep(1 * time.Second)
	}
	dk.EDK = 0
	dk.FDK = 1
	dk.DDK = data.USDK
	data.DataValue.SetDK(dk)
	if data.DataValue.Controller.Base {
		data.ToDevice <- 0
		data.ToServer <- 0
	} else {
		//Выбираем согласно планов
		// но пока все равно в ЛР
		choicePlan()
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
				choicePlan()
			case 6:
				//Смена НК
				data.DataValue.SetNK(cmd.Parametr)
				choicePlan()
			case 7:
				//Смена CК
				data.DataValue.SetCK(cmd.Parametr)
				choicePlan()
			case 9:
				//Смена ДУ
				data.DataValue.SetDU(cmd.Parametr)
				dk = data.DataValue.GetDK()
				if isDUPhase {
					if cmd.Parametr != dk.FDK {
						endDUPhase.Stop()
						isDUPhase = false
					} else {
						//Обновляем время удержания
						endDUPhase.Stop()
						endDUPhase = *time.NewTimer(time.Minute)
						continue
					}
				}
				if workplan {
					stopPlan()
				}
				if cmd.Parametr != 9 {
					if cmd.Parametr == 0 {
						//Перевод в ЛР
						if isDUPhase {
							endDUPhase.Stop()
						}
						isDUPhase = false
						dk.RDK = 5
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
							isDUPhase = true
							endDUPhase = *time.NewTimer(time.Minute)
						}
					}
				} else {
					//Выключаем ДУ производим выбор нового плана
					choicePlan()
				}
			}
		case <-endDUPhase.C:
			//Перестали удерживать ДУ
			data.Commands <- data.InternalCmd{Source: data.USDK, Command: 9, Parametr: 9}
		case ars := <-data.Arrays:
			data.DataValue.SetArrays(ars)

			if workplan {
				stopPlan()
			}
			// Тут нужно все заново выбрать
			choicePlan()
		case <-tik.C:
			//Выбираем план
			for !data.DataValue.Connect {
				for !data.DataValue.Connect {
					dk.EDK = 11
					dk.DDK = 8
					data.DataValue.SetDK(dk)
					data.ToDevice <- 0
					time.Sleep(1 * time.Second)
				}
			}
			choicePlan()
		case dev := <-data.FromDevice:
			logger.Info.Printf("От устройства %s", dev.ToString())
			// data.DataValue.Connect = dev.Connect
			dk = data.DataValue.GetDK()
			if workplan {
				toPlan <- dev
			}
			dk.FDK = dev.Phase
			dk.TTCDK = dev.TimeTC
			// dk.TDK = dev.TimeTU
			dk.TDK = 0
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
func choicePlan() {
	if data.DataValue.Controller.StatusCommandDU.IsDUDK1 {
		return
	}
	if data.DataValue.Controller.StatusCommandDU.IsPK {
		if workplan && nowPlan == data.DataValue.CommandDU.PK {
			return
		}
		if workplan {
			stopPlan()
		}

		go goPlan(data.DataValue.CommandDU.PK)

		return
	}
	data.DataValue.CommandDU.PK = 0
	if data.DataValue.Controller.StatusCommandDU.IsCK {
		data.DataValue.Controller.CK = data.DataValue.CommandDU.CK
	} else {
		data.DataValue.Controller.CK = 0
	}
	if data.DataValue.Controller.StatusCommandDU.IsNK {
		data.DataValue.Controller.NK = data.DataValue.CommandDU.NK
	} else {
		data.DataValue.Controller.NK = 0
	}
	mes := time.Now().Month()
	day := time.Now().Day()
	nday := time.Now().Weekday()
	hour := time.Now().Hour()
	min := time.Now().Minute()
	if data.DataValue.Controller.NK == 0 {
		mk := 0
		for _, v := range data.DataValue.Arrays.MonthSets.MonthSets {
			if v.Number == int(mes) {
				mk = v.Days[day-1]
				break
			}
		}
		// logger.Debug.Printf("find NK %d", mk)
		data.DataValue.Controller.NK = mk
	}
	if data.DataValue.Controller.CK == 0 {
		ck := 0
		for _, v := range data.DataValue.Arrays.WeekSets.WeekSets {
			if v.Number == int(data.DataValue.Controller.NK) {
				ck = v.Days[nday-1]
			}
		}
		// logger.Debug.Printf("find CK %d", ck)
		data.DataValue.Controller.CK = ck
	}
	if data.DataValue.Controller.PK == 0 {
		pk := 0
		for _, v := range data.DataValue.Arrays.DaySets.DaySets {
			if v.Number == int(data.DataValue.Controller.CK) {
				for _, v := range v.Lines {
					if hour < v.Hour {
						pk = v.PKNom
						break
					} else {
						if hour == v.Hour && min <= v.Min {
							pk = v.PKNom
							break
						}
					}
				}
				break
			}
		}
		// logger.Debug.Printf("find PK %d", pk)

		data.DataValue.Controller.PK = pk
	}
	if data.DataValue.Controller.PK == 0 {
		//Все плохо свалимся в ЛР
		if workplan {
			stopPlan()
		}
		dk.FDK = 1
		dk.RDK = 6
		dk.EDK = 4
		dk.DDK = data.USDK
		data.DataValue.SetDK(dk)
		data.ToDevice <- 0
		data.ToServer <- 0
		return
	}
	if workplan && nowPlan == data.DataValue.Controller.PK {
		return
	}
	if workplan {
		stopPlan()
	}
	go goPlan(data.DataValue.Controller.PK)
	// logger.Info.Printf("%d %d %d:%d %d", mes, day, hour, min, nday)
	// data.DataValue.Controller.NK=0
}
