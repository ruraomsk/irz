package web

import "github.com/anoshenko/rui"

const YearText = `
GridLayout {
	style = showPage,
	content = [
		GridLayout {
			width = 100%, height = 100%, cell-vertical-align = center, cell-horizontal-align = center,
			content = [
				TextView {
					id = textStyleText, padding = 16px, max-width = 80%, 
					border = _{ style = solid, width = 1px, color = darkgray },
					text = "Здесь должен быть Годовая карта"
				}
			]
		},
			]
		}
	]
}
`

func YearShow(session rui.Session) rui.View {
	view := rui.CreateViewFromText(session, YearText)
	if view == nil {
		return nil
	}

	return view
}
