package web

import (
	"fmt"
	"strconv"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/radar"
	"github.com/ruraomsk/irz/setup"
	"github.com/ruraomsk/irz/visio"
)

// border = _{ style = solid, width = 1px, color = darkgray },
const statusText = `
		ListLayout {
			width = 100%, height = 100%, orientation = vertical, padding = 16px,
			content = [
				TextView {
					text-color="red",text-align="center",text-size="24px",
					border = _{ style = solid, width = 1px, color = darkgray },
					id=titleStatus,text = ""
				},
				TextView {
					id=idBaseSet,text-color=red,
					text = ""
				},
				TextView {
					id=idConnect,semantics="code",
					text = "server"
				},

				TextView {
					id=idModbus,semantics="code",
					text = "modbus"
				},
				TextView {
					id=setModbusRadar,semantics="code",
					text = "modbusRadar"
				},
				TextView {
					id=setVisio,semantics="code",
					text = ""
				},
				TextView {
					id=workVisio,semantics="code",
					text = ""
				},
				ListLayout {
					orientation = horizontal, list-column-gap=16px,padding = 16px,
					content = [
						TextView {
							text = "<b>Команды управления от Центра</b>"
						},
						TextView {
							id=idSFDK,text-color=red,
							text = ""
						},
						TextView {
							id=idDU,text-color=red,
							text = ""
						},
						TextView {
							id=idPK,text-color=red,
							text = ""
						},
						TextView {
							id=idCK,text-color=red,
							text = ""
						},
						TextView {
							id=idNK,text-color=red,
							text = ""
						},
					]
				},		
				ListLayout {
					orientation = horizontal, list-column-gap=16px,padding = 16px,
					content = [
						TextView {
							text = "<b>Технология</b>"
						},
						TextView {
							id=idTech,
							text = ""
						},
						TextView {
							id=idNPK,text-color=green,
							text = ""
						},
						TextView {
							id=idNCK,text-color=green,
							text = ""
						},
						TextView {
							id=idNNK,text-color=green,
							text = ""
						},
						TextView {
							id=idRezim,text-color=green,
							text = ""
						},
						TextView {
							id=idPhaseDU,text-color=green,
							text = ""
						},
						TextView {
							id=idPhaseRU,text-color=green,
							text = ""
						},
						TextView {
							id=idBroken,text-color=red,
							text = ""
						},


					]
				},		
				ListLayout {
					orientation = horizontal, list-column-gap=16px,padding = 16px,
					border = _{style=solid,width=4px,color=blue },
					content = [
						TextView {
							text = "<b>Изменение настроек</b>"
						},
						Button {
							id=setBase,content="Установить базовую привязку"
						},
						EditView{
							id=idIP,type=text
						},
						NumberPicker {
							id=idPort,type=editor,min=0,max=32000,value=1090
						},
						Button {
							id=setIP,content="Установить сервер связи"
						},
					]
				},		
			]
		}
`

func toString(t time.Time) string {
	t = t.Add(time.Duration(data.DataValue.Arrays.TimeDivice.TimeZone) * time.Hour)
	return fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
func makeViewStatus(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	t := time.Now().Add(time.Duration(data.DataValue.Arrays.TimeDivice.TimeZone) * time.Hour)
	rui.Set(view, "titleStatus", "text", fmt.Sprintf("<b>Текущее состояние УСДК %d </b>%02d:%02d:%02d", data.DataValue.Controller.ID,
		t.Hour(), t.Minute(), t.Second()))

	if data.DataValue.Controller.Base {
		rui.Set(view, "idBaseSet", "text", "<b>БАЗОВАЯ ПРИВЯЗКА</b>")
	} else {
		rui.Set(view, "idBaseSet", "text", "")
	}
	c := fmt.Sprintf("<b>Соединение с сервером %s:%d ", setup.Set.Server.Host, setup.Set.Server.Port)
	if data.DataValue.Controller.StatusConnection {
		c += fmt.Sprintf("установлено %s обмен %s", toString(data.DataValue.Controller.ConnectTime), toString(data.DataValue.Controller.LastOperation))
	} else {
		c += "отсутствует"
	}
	rui.Set(view, "idConnect", "text", c+"</b>")

	c = fmt.Sprintf("Соединение Modbus device %s baud %d parity %s uid %d \t",
		setup.Set.Modbus.Device, setup.Set.Modbus.BaudRate, setup.Set.Modbus.Parity, setup.Set.Modbus.UId)
	if data.DataValue.Connect {
		c += "установлено"
	} else {
		c += "отсутствует"
	}
	rui.Set(view, "idModbus", "text", c)
	if !setup.Set.ModbusRadar.Radar {
		rui.Set(view, "setModbusRadar", "text", "Оключен прием данных от радаров")
	} else {
		rui.Set(view, "setModbusRadar", "text", fmt.Sprintf("От радаров (%s): %s ", radar.GetStatus(), radar.GetValues()))
	}
	if !setup.Set.VisioDevice.Visio {
		rui.Set(view, "setVisio", "text", "Отключена индикация")
	} else {
		if visio.Ready {
			rui.Set(view, "setVisio", "text", "Описание фаз загружено")
		} else {
			rui.Set(view, "setVisio", "text", "Описание фаз не загружено (нет связи по Modbus)")
		}
		vs := fmt.Sprintf("Связь с индикацией device %s baud %d parity %s ", setup.Set.VisioDevice.Device, setup.Set.VisioDevice.BaudRate, setup.Set.VisioDevice.Parity)
		if visio.Work {
			vs += " установлена"
		} else {
			vs += " отсутсвует"
		}
		rui.Set(view, "workVisio", "text", vs)
	}
	showTech(view)
}
func updaterStatus(view rui.View, session rui.Session) {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		if view == nil {
			return
		}
		w, ok := SessionStatus[session.ID()]
		if !ok {
			continue
		}

		if !w {
			continue
		}
		makeViewStatus(view)
	}
}

func statusShow(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, statusText)
	if view == nil {
		return nil
	}
	makeViewStatus(view)
	go updaterStatus(view, session)

	rui.Set(view, "idIP", "text", setup.ExtSet.Server.Host)
	rui.Set(view, "idPort", "value", setup.ExtSet.Server.Port)
	rui.Set(view, "setBase", rui.ClickEvent, func(rui.View) {
		logger.Info.Println("Утановили базовую привязку")
		data.SetBase <- 1
	})
	rui.Set(view, "setIP", rui.ClickEvent, func(rui.View) {
		setup.ExtSet.Server.Host = rui.GetText(view, "idIP")
		a := rui.Get(view, "idPort", "value")
		s, ok := a.(string)
		if ok {
			setup.ExtSet.Server.Port, _ = strconv.Atoi(s)
		}
		f, ok := a.(float64)
		if ok {
			setup.ExtSet.Server.Port = int(f)
		}
		logger.Info.Printf("Изменили центр на %s:%d", setup.ExtSet.Server.Host, setup.ExtSet.Server.Port)
		data.SaveExtSetup <- 1
	})

	return view
}
