package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/data"
)

// border = _{ style = solid, width = 1px, color = darkgray },
const statusText = `
ListLayout {
	style = showPage,
	content = [
		GridLayout {
			width = 100%, height = 100%, orientation = vertical, padding = 16px,
			content = [
				TextView {
					text = "Текущее состояние УСДК"
				},
				TextView {
					id=idDevice,
					text = "Номер контроллера"
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
	rui.Set(view, "idDevice", "text", fmt.Sprintf("Номер контроллера %d", data.DataValue.Controller.ID))
	return view
}
