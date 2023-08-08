package worker

import (
	"time"

	"github.com/ruraomsk/ag-server/binding"
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
type nowLocal struct {
	mes  int
	day  int
	nday int
	hour int
	min  int
	sec  int
}

var dk = pudge.DK{RDK: 5, FDK: 1, DDK: data.USDK, EDK: 0, PDK: false, EEDK: 0, ODK: false, LDK: 0, FTUDK: 0, TDK: 0, FTSDK: 0, TTCDK: 0}
var endDUPhase time.Timer
var toPlan chan data.StatusDevice
var nowPlan = 0
var prom = 0

// var ctrlDU = time.Duration(60 * time.Second)
var isDUPhase = false

func Worker() {
	toPlan = make(chan data.StatusDevice, 100)
	endPlan = make(chan interface{})
	data.DataValue.SetDK(dk)
	for !data.DataValue.Connect {
		dk = data.DataValue.GetDK()
		dk.EDK = 11
		data.DataValue.SetDK(dk)
		time.Sleep(1 * time.Second)
	}

	if data.DataValue.Controller.Base {
		data.ToDevice <- 0
		data.ToServer <- 0
	}
	tik := time.NewTicker(1 * time.Second)
	go controllerPlans()
	logger.Info.Print("Начинаем основной цикл worker")
	for {
		select {
		case cmd := <-data.Commands:
			if !data.DataValue.Connect {
				dk = data.DataValue.GetDK()
				dk.EDK = 11
				dk.DDK = 8
				data.DataValue.SetDK(dk)
				continue
			}
			// logger.Info.Printf("Команда %v", cmd)
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
				stopPlan()
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
					if isDUPhase {
						endDUPhase.Stop()
					}
					isDUPhase = false
					dk.RDK = 8
					data.DataValue.SetDK(dk)
					// data.ToServer <- 0
					choicePlan()
				}
			}
		case <-endDUPhase.C:
			if !data.DataValue.Connect {
				dk = data.DataValue.GetDK()
				dk.EDK = 11
				dk.DDK = 8
				data.DataValue.SetDK(dk)
				continue
			}
			//Перестали удерживать ДУ
			data.Commands <- data.InternalCmd{Source: data.USDK, Command: 9, Parametr: 9}
		case ars := <-data.Arrays:
			data.DataValue.SetArrays(ars)
			if !data.DataValue.Connect {
				dk = data.DataValue.GetDK()
				dk.EDK = 11
				dk.DDK = 8
				data.DataValue.SetDK(dk)
				continue
			}
			stopPlan()
			// Тут нужно все заново выбрать
			choicePlan()
		case <-tik.C:
			//Выбираем план
			if !data.DataValue.Connect {
				dk = data.DataValue.GetDK()
				dk.EDK = 11
				dk.DDK = 8
				data.DataValue.SetDK(dk)
				continue
			}
			choicePlan()
		case dev := <-data.FromDevice:
			// logger.Info.Printf("От устройства %s", dev.ToString())
			// data.DataValue.Connect = dev.Connect
			dk = data.DataValue.GetDK()
			if workplan {
				toPlan <- dev
			}
			dk.FDK = dev.Phase
			if dev.Phase == 9 {
				dk.TDK = 0
			} else {
				prom = dev.TimeTC
				dk.TDK = dev.TimeTU
			}
			dk.FTSDK = dev.PhaseTC
			dk.TTCDK = dev.TimeTC
			if dev.Phase == 9 {
				dk.TTCDK = dev.TimeTU - prom
			}
			dk.FTUDK = dev.PhaseTU
			if dk.RDK == 5 || dk.RDK == 6 {
				dk.FTUDK = 0
				dk.TDK = 0
			}
			data.DataValue.SetDK(dk)
			// logger.Info.Printf("%v", data.DataValue.Controller.DK)
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
		logger.Debug.Printf("Назначен план %d", data.DataValue.CommandDU.PK)
		var p = binding.SetPk{Pk: 0}
		for _, v := range data.DataValue.Arrays.SetDK.DK {
			if v.Pk == data.DataValue.CommandDU.PK {
				p = v
			}
		}
		if p.TypePU == 1 || p.Tc < 3 {
			stopPlan()
			data.DataValue.Controller.PK = data.DataValue.CommandDU.PK
			dk.RDK = 8
			data.DataValue.SetDK(dk)
			go goPlan(data.DataValue.Controller.PK)
			time.Sleep(time.Second)
		} else {
			r := getLocalTime()
			w := (r.hour*3600 + r.min*60 + r.sec) % p.Tc
			if w != 0 {
				w = p.Tc - w
			}
			w += p.Shift
			if w > 0 {
				for w != 1 {
					// logger.Debug.Printf("До старта плана %d %d %d %d:%d:%d", p.Pk, p.Tc, w, hour, min, sec)
					time.Sleep(time.Second)
					w--
				}
			}
			stopPlan()
			dk.RDK = 8
			data.DataValue.Controller.PK = data.DataValue.CommandDU.PK
			data.DataValue.SetDK(dk)
			go goPlan(data.DataValue.Controller.PK)
			time.Sleep(time.Second)
		}

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
	r := getLocalTime()
	if data.DataValue.Controller.NK == 0 {
		mk := 0
		for _, v := range data.DataValue.Arrays.MonthSets.MonthSets {
			if v.Number == int(r.mes) {
				mk = v.Days[r.day-1]
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
				ck = v.Days[r.nday-1]
			}
		}
		// logger.Debug.Printf("find CK %d", ck)
		data.DataValue.Controller.CK = ck
	}
	pk := 0
	for _, v := range data.DataValue.Arrays.DaySets.DaySets {
		if v.Number == int(data.DataValue.Controller.CK) {
			for _, v := range v.Lines {
				if r.hour < v.Hour {
					pk = v.PKNom
					break
				} else {
					if r.hour == v.Hour && r.min <= v.Min {
						pk = v.PKNom
						break
					}
				}
			}
			break
		}
	}
	// logger.Debug.Printf("Выбран PK %d", pk)

	if pk == 0 {
		//Все плохо свалимся в ЛР
		logger.Error.Println("плохо свалимся в ЛР")
		stopPlan()
		data.DataValue.Controller.PK = pk
		dk.FDK = 1
		dk.RDK = 6
		dk.EDK = 4
		dk.DDK = data.USDK
		data.DataValue.SetDK(dk)
		data.ToDevice <- 0
		data.ToServer <- 0
		return
	}
	if workplan && nowPlan == pk {
		// logger.Debug.Printf("Выбран PK %d и он исполняется", pk)

		return
	}

	var p = binding.SetPk{Pk: 0}
	for _, v := range data.DataValue.Arrays.SetDK.DK {
		if v.Pk == pk {
			p = v
		}
	}
	logger.Debug.Printf("Переходим на план %d", pk)
	if p.TypePU == 1 || p.Tc < 3 {
		stopPlan()
		data.DataValue.Controller.PK = pk
		dk.RDK = 8
		data.DataValue.SetDK(dk)
		go goPlan(data.DataValue.Controller.PK)
		time.Sleep(time.Second)
	} else {
		r := getLocalTime()
		w := (r.hour*3600 + r.min*60 + r.sec) % p.Tc
		if w != 0 {
			w = p.Tc - w
		}
		data.DataValue.Controller.DK.PDK = true
		w += p.Shift
		if w > 0 {
			for w != 1 {
				// logger.Debug.Printf("До старта плана %d %d %d %d:%d:%d", p.Pk, p.Tc, w, hour, min, sec)
				time.Sleep(time.Second)
				w--
			}
		}
		stopPlan()
		data.DataValue.Controller.PK = pk
		dk.RDK = 8
		data.DataValue.SetDK(dk)
		go goPlan(data.DataValue.Controller.PK)
		time.Sleep(time.Second)
	}
}
func getLocalTime() nowLocal {
	var r nowLocal
	t := time.Now().Add(time.Duration(data.DataValue.Arrays.TimeDivice.TimeZone) * time.Hour)
	r.hour = t.Hour()
	r.min = t.Minute()
	r.sec = t.Second()
	r.mes = int(t.Month())
	r.day = t.Day()
	r.nday = int(t.Weekday())
	if r.nday == 0 {
		r.nday = 7
	}
	// logger.Info.Printf("now %d:%d:%d", r.hour, r.min, r.sec)
	return r
}
