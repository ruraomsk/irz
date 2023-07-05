package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/binding"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
)

const CkText = `
TabsLayout { id = tabsCK, width = 100%, height = 100%, tabs = top, tab-close-button = false,
	content = [
		TableView {cell-horizontal-align = right,  title = "СК1", id="ck1"},
		TableView {cell-horizontal-align = right,  title = "СК2", id="ck2"},
		TableView {cell-horizontal-align = right,  title = "СК3", id="ck3"},
		TableView {cell-horizontal-align = right,  title = "СК4", id="ck4"},
		TableView {cell-horizontal-align = right,  title = "СК5", id="ck5"},
		TableView {cell-horizontal-align = right,  title = "СК6", id="ck6"},
		TableView {cell-horizontal-align = right,  title = "СК7", id="ck7"},
		TableView {cell-horizontal-align = right,  title = "СК8", id="ck8"},
		TableView {cell-horizontal-align = right,  title = "СК9", id="ck9"},
		TableView {cell-horizontal-align = right,  title = "СК10", id="ck10"},
		TableView {cell-horizontal-align = right,  title = "СК11", id="ck11"},
		TableView {cell-horizontal-align = right,  title = "СК12", id="ck12"},
	]
}
`

func CKShow(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, CkText)
	if view == nil {
		return nil
	}
	for cl := 1; cl < 13; cl++ {
		var ck = binding.OneDay{Number: 0}
		for _, v := range data.DataValue.Arrays.DaySets.DaySets {
			if v.Number == cl {
				ck = *v
			}
		}
		if ck.Number == 0 {
			logger.Error.Printf("Нет суточного плана %d", cl)
			continue
		}
		var content [][]any
		count := 1
		content = append(content, []any{"Время начала", "План координации"})

		for _, v := range ck.Lines {
			content = append(content, []any{fmt.Sprintf("%02d:%02d", v.Hour, v.Min), v.PKNom})
			count++
		}

		rui.SetParams(view, fmt.Sprintf("ck%d", cl), rui.Params{
			rui.Content:    content,
			rui.HeadHeight: count,
		})

	}
	return view
}
