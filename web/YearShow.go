package web

import (
	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/data"
)

const YearText = `
TableView {cell-horizontal-align = right,
	id="yk"}
`

func YearShow(session rui.Session) rui.View {
	mutex.Lock()
	defer mutex.Unlock()

	view := rui.CreateViewFromText(session, YearText)
	if view == nil {
		return nil
	}
	var content [][]any
	count := 1
	var header []any
	header = append(header, "Месяц")
	for i := 1; i < 32; i++ {
		header = append(header, i)
	}
	content = append(content, header)

	for _, v := range data.DataValue.Arrays.MonthSets.MonthSets {
		var line []any
		line = append(line, v.Number)
		for i := 0; i < len(v.Days); i++ {
			line = append(line, v.Days[i])
		}
		content = append(content, line)
		count++
	}

	rui.SetParams(view, "yk", rui.Params{
		rui.Content:     content,
		rui.HeadHeight:  count,
		rui.CellPadding: "2px",
	})

	return view
}
