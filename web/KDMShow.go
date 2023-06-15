package web

import "github.com/anoshenko/rui"

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
	view := rui.CreateViewFromText(session, KDMText)
	if view == nil {
		return nil
	}

	return view
}
