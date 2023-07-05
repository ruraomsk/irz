package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/setup"
	"github.com/ruraomsk/irz/visio"
)

const visioText = `
ListLayout {
	width = 100%, height = 100%, orientation = vertical, padding = 16px,
	content = [
		TextView {
			id=setVisio,semantics="code",
			text = ""
		},
		TextView {
			id=workVisio,semantics="code",
			text = ""
		},
		TableView {cell-horizontal-align = right,
			id="tablePhases"},
  ]
},

`

func VisioShow(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, visioText)
	if view == nil {
		return nil
	}
	if !setup.Set.Visio {
		rui.Set(view, "setVisio", "text", "Отключена индикация")
		return view
	}
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

	var content [][]any
	count := 1
	var header []any
	header = append(header, "Фаза")
	for i := 0; i < 16; i++ {
		header = append(header, fmt.Sprintf("Нап %d", i+1))
	}
	content = append(content, header)

	for _, v := range visio.Phases {
		var line []any
		line = append(line, v.Number)
		for i := 0; i < len(v.Open); i++ {
			line = append(line, v.Open[i])
		}
		content = append(content, line)
		count++
	}

	rui.SetParams(view, "tablePhases", rui.Params{
		rui.Content:    content,
		rui.HeadHeight: count,
	})

	return view
}
