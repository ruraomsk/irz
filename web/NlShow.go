package web

import (
	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/data"
)

const NkText = `
		TableView {cell-horizontal-align = right,
			  id="nk"}
`

func NKShow(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, NkText)
	if view == nil {
		return nil
	}
	var content [][]any
	count := 1
	content = append(content, []any{"Номер карты", "пн", "вт", "ср", "чт", "пт", "сб", "вс"})

	for _, v := range data.DataValue.Arrays.WeekSets.WeekSets {
		content = append(content, []any{v.Number, v.Days[0], v.Days[1], v.Days[2], v.Days[3], v.Days[4], v.Days[5], v.Days[6]})
		count++
	}

	rui.SetParams(view, "nk", rui.Params{
		rui.Content:     content,
		rui.HeadHeight:  count,
		rui.CellPadding: "4px",
	})

	return view
}
