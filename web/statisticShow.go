package web

import (
	"fmt"
	"time"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/stat"
)

const statisticText = `
		ListLayout {
			width = 100%, height = 100%, orientation = vertical, padding = 16px,
			content = [
				TextView {
					text-align="center",text-size="24px",
					border = _{ style = solid, width = 1px},
					id=titleStatistic,text = "Статистика"
				},
				TableView {cell-horizontal-align = right,
					id="stat"}
	  
			]
		}
`

var names = []any{"Значение", "Статус"}

func updaterStatistic(view rui.View, session rui.Session) {
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
		makeViewStatistic(view)
	}
}
func makeViewStatistic(view rui.View) {
	mutex.Lock()
	defer mutex.Unlock()
	stat.GetLast <- 1
	last := <-stat.LastSending

	var result string
	if last.LastTime == time.Unix(0, 0) {
		result = "Еще нет отправленных данных статистики"
	} else {
		result = fmt.Sprintf("Отправлено %s ", toString(last.LastTime))
		result += fmt.Sprintf(" %02d:%02d ", last.Last.Hour, last.Last.Min)
		if last.Last.Type == 1 {
			result += " интенсивность "
		} else {
			result += " скорость "
		}
		result += fmt.Sprintf(" интервал минут %d ", last.Last.TLen)
	}
	rui.Set(view, "titleStatistic", "text", result)
	var content [][]any
	count := 1
	content = append(content, []any{"Тип", "1 кан", "2 кан", "3 кан", "4 кан", "5 кан", "6 кан", "7 кан", "8 кан", "9 кан", "10 кан", "11 кан", "12 кан", "13 кан", "14 кан", "15 кан", "16 кан", "17 кан", "18 кан"})
	for i := 0; i < 2; i++ {
		var line []any
		line = append(line, names[i])
		for j := 0; j < 18; j++ {
			line = append(line, "")
		}
		if i == 0 {
			for _, v := range last.Last.Datas {
				switch last.Last.Type {
				case 1:
					line[v.Chanel] = v.Intensiv
				case 2:
					line[v.Chanel] = v.Speed
				}
			}
		}
		if i == 1 {
			for _, v := range last.Last.Datas {
				line[v.Chanel] = v.Status
			}
		}

		count++
		// fmt.Printf("%v", line)
		content = append(content, line)
	}

	rui.SetParams(view, "stat", rui.Params{
		rui.Content:     content,
		rui.HeadHeight:  count,
		rui.CellPadding: "4px",
	})
}
func statisticShow(session rui.Session) rui.View {

	view := rui.CreateViewFromText(session, statisticText)
	if view == nil {
		return nil
	}
	makeViewStatistic(view)
	go updaterStatistic(view, session)
	return view
}
