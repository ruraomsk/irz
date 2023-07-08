package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/visio"
)

const visioText = `
ListLayout {
	width = 100%, height = 100%, orientation = vertical, padding = 16px,
	content = [
		TableView {cell-horizontal-align = right,
			id="tablePhases"},
  ]
},

`

func VisioShow(session rui.Session) rui.View {
	mutex.Lock()
	defer mutex.Unlock()

	view := rui.CreateViewFromText(session, visioText)
	if view == nil {
		return nil
	}

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
		rui.Content:     content,
		rui.HeadHeight:  count,
		rui.CellPadding: "2px",
	})

	return view
}
