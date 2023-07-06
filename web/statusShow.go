package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/setup"
	"github.com/ruraomsk/irz/visio"
)

// border = _{ style = solid, width = 1px, color = darkgray },
const statusText = `
ListLayout {
	style = showPage,
	content = [
		ListLayout {
			width = 100%, height = 100%, orientation = vertical, padding = 16px,
			content = [
				TextView {
					text-color="red",text-align="center",text-size="24px",
					border = _{ style = solid, width = 1px, color = darkgray },
					id=titleStatus,text = ""
				},
				TextView {
					id=idDevice,semantics="code",
					text = "Номер контроллера"
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
					id=setVisio,semantics="code",
					text = ""
				},
				TextView {
					id=workVisio,semantics="code",
					text = ""
				},
			]
		},
	]
}
`

func toString(t time.Time) string {
	return fmt.Sprintf("%04d/%02d/%02d %02d:%02d:%02d", t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second())
}
func makeViewStatus(view rui.View) {
	rui.Set(view, "titleStatus", "text", fmt.Sprintf("Текущее состояние УСДК %02d:%02d:%02d",
		time.Now().Hour(), time.Now().Minute(), time.Now().Second()))
	rui.Set(view, "idDevice", "text", fmt.Sprintf("Номер контроллера \t%d", data.DataValue.Controller.ID))

	c := fmt.Sprintf("Соединение с сервером %s:%d \t", data.DataValue.Server.Host, data.DataValue.Server.Port)
	if data.DataValue.Controller.StatusConnection {
		c += fmt.Sprintf("установлено %s обмен %s", toString(data.DataValue.Controller.ConnectTime), toString(data.DataValue.Controller.LastOperation))
	} else {
		c += "отсутствует"
	}
	rui.Set(view, "idConnect", "text", c)

	c = fmt.Sprintf("Соединение Modbus device %s baud %d parity %s uid %d \t",
		setup.Set.Modbus.Device, setup.Set.Modbus.BaudRate, setup.Set.Modbus.Parity, setup.Set.Modbus.UId)
	if data.DataValue.Connect {
		c += "установлено"
	} else {
		c += "отсутствует"
	}
	rui.Set(view, "idModbus", "text", c)
	if !setup.Set.Visio {
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

}
func updaterStatus(view rui.View, session rui.Session) {
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		if view == nil {
			return
		}
		if !SessionStatus[session.ID()] {
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
	return view
}
