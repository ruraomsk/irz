package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/setup"
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
					text = "Текущее состояние УСДК"
				},
				TextView {
					id=idDevice,
					text = "Номер контроллера"
				},
				TextView {
					id=idConnect,
					text = "server"
				},
				TextView {
					id=idModbus,
					text = "modbus"
				},

			]
		},
	]
}
`

func statusShow(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, statusText)
	if view == nil {
		return nil
	}
	rui.Set(view, "idDevice", "text", fmt.Sprintf("Номер контроллера \t%d", data.DataValue.Controller.ID))

	c := fmt.Sprintf("Соединение с сервером %s:%d \t", data.DataValue.Server.Host, data.DataValue.Server.Port)
	if data.DataValue.Controller.StatusConnection {
		c += "установлено"
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

	return view
}
