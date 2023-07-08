package web

import (
	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/kdm"
)

const KDMText = `
GridLayout {
	style = showPage,
	content = [
		GridLayout {
			width = 100%, height = 100%, cell-vertical-align = center, cell-horizontal-align = center,
			content = [
				TextView {
					text-color="red",text-align="center",text-size="24px",
					border = _{ style = solid, width = 1px, color = darkgray },
					text = "РЕГИСТРЫ"
				},
			]
		},
			]
		}
	]
}
`

func KDMShow(session rui.Session) rui.View {
	mutex.Lock()
	defer mutex.Unlock()

	view := rui.CreateViewFromText(session, KDMText)
	if view == nil {
		return nil
	}
	kdm.RequestChan <- kdm.Request{Start: 0, Lenght: 1}
	return view
}
