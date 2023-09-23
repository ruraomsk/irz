package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/binding"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
)

const NkText = `
ListLayout {
	width = 100%, height = 100%, orientation = vertical, padding = 16px,
	content = [
		Button{ id = idEdit, content="Изменить выбранную недельную карту" },
		TabsLayout { id = tabsNK, width = 100%, height = 100%, tabs = top, tab-close-button = false,
			content = [
				TableView {cell-horizontal-align = right,  title = "HК1", id="nk1"},
				TableView {cell-horizontal-align = right,  title = "HК2", id="nk2"},
				TableView {cell-horizontal-align = right,  title = "HК3", id="nk3"},
				TableView {cell-horizontal-align = right,  title = "HК4", id="nk4"},
				TableView {cell-horizontal-align = right,  title = "HК5", id="nk5"},
				TableView {cell-horizontal-align = right,  title = "HК6", id="nk6"},
				TableView {cell-horizontal-align = right,  title = "HК7", id="nk7"},
				TableView {cell-horizontal-align = right,  title = "HК8", id="nk8"},
				TableView {cell-horizontal-align = right,  title = "HК9", id="nk9"},
				TableView {cell-horizontal-align = right,  title = "HК10", id="nk10"},
				TableView {cell-horizontal-align = right,  title = "HК11", id="nk11"},
				TableView {cell-horizontal-align = right,  title = "HК12", id="nk12"},
			]
		}
	]		
}
`

var updatedNk = false

func makeNkShow(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	for cl := 1; cl < 13; cl++ {
		var nk binding.OneWeek
		for _, v := range data.DataValue.Arrays.WeekSets.WeekSets {
			if v.Number == cl {
				nk = v
			}
		}
		if nk.Number == 0 {
			logger.Error.Printf("Нет недельного плана %d", cl)
			continue
		}
		var content [][]any
		content = append(content, []any{"пн", "вт", "ср", "чт", "пт", "сб", "вс"})

		content = append(content, []any{nk.Days[0], nk.Days[1], nk.Days[2], nk.Days[3], nk.Days[4], nk.Days[5], nk.Days[6]})

		rui.SetParams(view, fmt.Sprintf("nk%d", cl), rui.Params{
			rui.Content:     content,
			rui.HeadHeight:  2,
			rui.CellPadding: "8px",
		})

	}
}
func updaterNkShow(view rui.View, session rui.Session) {
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
		if updatedNk {
			makeNkShow(view)
			updatedNk = false
		}
	}
}

func NKShow(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, NkText)
	if view == nil {
		return nil
	}
	makeNkShow(view)
	rui.Set(view, "idEdit", rui.ClickEvent, func(rui.View) {
		editNK(session, rui.GetCurrent(view, "tabsNK")+1)
	})
	go updaterNkShow(view, session)
	return view
}
