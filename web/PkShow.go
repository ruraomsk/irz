package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/binding"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/worker"
)

const PkText = `
	ListLayout {
		width = 100%, height = 100%, orientation = vertical, padding = 16px,
		content = [
			Button{ id = idEdit, content="Изменить выбранный план координации" },
			TabsLayout { id = tabsPK, width = 100%, height = 100%, tabs = top, tab-close-button = false,
				content = [
					TableView {cell-horizontal-align = right,  title = "ПК1", id="pk1"},
					TableView {cell-horizontal-align = right,  title = "ПК2", id="pk2"},
					TableView {cell-horizontal-align = right,  title = "ПК3", id="pk3"},
					TableView {cell-horizontal-align = right,  title = "ПК4", id="pk4"},
					TableView {cell-horizontal-align = right,  title = "ПК5", id="pk5"},
					TableView {cell-horizontal-align = right,  title = "ПК6", id="pk6"},
					TableView {cell-horizontal-align = right,  title = "ПК7", id="pk7"},
					TableView {cell-horizontal-align = right,  title = "ПК8", id="pk8"},
					TableView {cell-horizontal-align = right,  title = "ПК9", id="pk9"},
					TableView {cell-horizontal-align = right,  title = "ПК10", id="pk10"},
					TableView {cell-horizontal-align = right,  title = "ПК11", id="pk11"},
					TableView {cell-horizontal-align = right,  title = "ПК12", id="pk12"},
				]
			}
		]		
	}
`

var updatedPk = false

func makePKShow(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()

	for pl := 1; pl < 13; pl++ {
		var pk = binding.SetPk{Pk: 0}
		for _, v := range data.DataValue.Arrays.SetDK.DK {
			if v.Pk == pl {
				pk = v
			}
		}
		if pk.Pk == 0 {
			logger.Error.Printf("Нет плана координации %d", pl)
			continue
		}
		pk = worker.RepackPlan(pk)
		tp := ""
		// header := 1
		if pk.Tc > 2 {
			tp = "Координированный"
			if pk.TypePU == 1 {
				tp = "Локальный"
			} else {
				tp += fmt.Sprintf(" Сдвиг %d", pk.Shift)
			}
			tp += fmt.Sprintf(" Время цикла %d", pk.Tc)
		} else {
			// header = 0
			if pk.Tc == 0 {
				tp = "Локальное управление"
			} else if pk.Tc == 1 {
				tp = "Желтое мигание"
			} else if pk.Tc == 2 {
				tp = "Отключить светофор"
			}
		}
		var content [][]any
		count := 1
		content = append(content, []any{tp, rui.HorizontalTableJoin{}})
		// if header > 0 {
		content = append(content, []any{"Тип", "Фаза", "Начало", "Конец"})
		count++
		// }

		for _, v := range pk.Stages {
			// if v.Start == 0 && v.Stop == 0 {
			// 	continue
			// }
			// 1 - МГР
			// 2 - 1ТВП
			// 3 - 2ТВП
			// 4 - 1,2ТВП
			// 5 - Зам 1 ТВП
			// 6 - Зам 2 ТВП
			// 7 - Зам
			// 8 - МДК
			// 9 - ВДК

			tf := ""
			switch v.Tf {
			case 0:
				tf = "Простая"
			case 1:
				tf = "МГР"
			case 2:
				tf = "ТВП 1"
			case 3:
				tf = "ТВП 2"
			case 4:
				tf = "ТВП 1,2 "
			case 5:
				tf = "ЗАМ ТВП 1 "
			case 6:
				tf = "ЗАМ ТВП 2 "
			case 7:
				tf = "ЗАМ ТВП 1,2 "
			}
			content = append(content, []any{tf, v.Number, v.Start, v.Stop})
			count++
		}
		// content = append(content)
		rui.SetParams(view, fmt.Sprintf("pk%d", pl), rui.Params{
			rui.Content:     content,
			rui.HeadHeight:  count,
			rui.CellPadding: "4px",
		})

	}

}
func updaterPkShow(view rui.View, session rui.Session) {
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
		if updatedPk {
			makePKShow(view)
			updatedPk = false
		}
	}
}

func PKShow(session rui.Session) rui.View {

	view := rui.CreateViewFromText(session, PkText)
	if view == nil {
		return nil
	}
	makePKShow(view)
	rui.Set(view, "idEdit", rui.ClickEvent, func(rui.View) {
		editPK(session, rui.GetCurrent(view, "tabsPK")+1)
	})
	go updaterPkShow(view, session)
	return view
}
