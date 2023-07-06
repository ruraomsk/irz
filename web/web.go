package web

import (
	"fmt"

	"github.com/anoshenko/rui"
)

const rootViewText = `
GridLayout {
	id = rootLayout, width = 100%, height = 100%, cell-height = "auto, 1fr",
	content = [
		GridLayout {
			id = rootTitle, width = 100%, cell-width = "auto, 1fr", 
			cell-vertical-align = center, background-color = #ffc0ded9, text-color = black,
			content = [
				ImageView { 
					id = rootTitleButton, padding = 8px, src = menu_icon.svg,
					tooltip = "Выберите режим (alt-M)"
				},
				TextView { 
					id = rootTitleText, column = 1, padding-left = 8px, text = "Title",
				}
			],
		},
		StackLayout {
			id = rootViews, row = 1,
		}
	]
}
`

var SessionStatus map[int]bool

type Page struct {
	title   string
	creator func(session rui.Session) rui.View
	view    rui.View
}

type NowSession struct {
	rootView rui.View
	pages    []Page
}

func (d *NowSession) OnStart(session rui.Session) {
	SessionStatus[session.ID()] = true
	rui.DebugLog(fmt.Sprintf("Session start %d", session.ID()))
}

func (d *NowSession) OnFinish(session rui.Session) {
	rui.DebugLog(fmt.Sprintf("Session finish %d", session.ID()))
	_, ok := SessionStatus[session.ID()]
	if !ok {
		rui.DebugLog(fmt.Sprintf("Session not started %v", SessionStatus))
		return
	}
	SessionStatus[session.ID()] = false
}

func (d *NowSession) OnResume(session rui.Session) {
	rui.DebugLog(fmt.Sprintf("Session resume %d", session.ID()))
	_, ok := SessionStatus[session.ID()]
	if !ok {
		rui.DebugLog(fmt.Sprintf("Session not started %v", SessionStatus))
		return
	}
	SessionStatus[session.ID()] = true
}

func (d *NowSession) OnPause(session rui.Session) {
	rui.DebugLog(fmt.Sprintf("Session pause %d", session.ID()))
	_, ok := SessionStatus[session.ID()]
	if !ok {
		rui.DebugLog(fmt.Sprintf("Session not started %v", SessionStatus))
		return
	}
	SessionStatus[session.ID()] = false
}

func (d *NowSession) OnDisconnect(session rui.Session) {
	rui.DebugLog(fmt.Sprintf("Session disconect %d", session.ID()))
	_, ok := SessionStatus[session.ID()]
	if !ok {
		rui.DebugLog(fmt.Sprintf("Session not started %v", SessionStatus))
		return
	}
	SessionStatus[session.ID()] = false
}

func (d *NowSession) OnReconnect(session rui.Session) {
	rui.DebugLog(fmt.Sprintf("Session reconect %d", session.ID()))
	_, ok := SessionStatus[session.ID()]
	if !ok {
		rui.DebugLog(fmt.Sprintf("Session not started %v", SessionStatus))
		return
	}
	SessionStatus[session.ID()] = true
}

func CreateSession(_ rui.Session) rui.SessionContent {
	sessionContent := new(NowSession)
	sessionContent.pages = []Page{
		{"Текущее состояние", statusShow, nil},
		{"Планы координации", PKShow, nil},
		{"Суточные карты", CKShow, nil},
		{"Недельные карты", NKShow, nil},
		{"Годовая карта", YearShow, nil},
		{"Индикатор отсчета", VisioShow, nil},
		// {"Состояние КДМ", KDMShow, nil},
	}

	return sessionContent
}
func (d *NowSession) CreateRootView(session rui.Session) rui.View {
	d.rootView = rui.CreateViewFromText(session, rootViewText)
	if d.rootView == nil {
		return nil
	}

	rui.Set(d.rootView, "rootTitleButton", rui.ClickEvent, d.clickMenuButton)
	session.SetHotKey(rui.KeyM, rui.AltKey, func(session rui.Session) {
		d.clickMenuButton()
	})
	d.showPage(0)
	return d.rootView
}
func (d *NowSession) clickMenuButton() {
	items := make([]string, len(d.pages))
	for i, page := range d.pages {
		items[i] = page.title
	}

	buttonFrame := rui.ViewByID(d.rootView, "rootTitleButton").Frame()

	rui.ShowMenu(d.rootView.Session(), rui.Params{
		rui.Items:           items,
		rui.OutsideClose:    true,
		rui.VerticalAlign:   rui.TopAlign,
		rui.HorizontalAlign: rui.LeftAlign,
		rui.MarginLeft:      rui.Px(buttonFrame.Bottom() / 2),
		rui.Arrow:           rui.LeftArrow,
		rui.ArrowAlign:      rui.LeftAlign,
		rui.ArrowSize:       rui.Px(12),
		rui.ArrowOffset:     rui.Px(buttonFrame.Left + (buttonFrame.Width-12)/2),
		rui.PopupMenuResult: func(n int) {
			d.showPage(n)
		},
	})
}

func (d *NowSession) showPage(index int) {
	if index < 0 || index >= len(d.pages) {
		return
	}

	if stackLayout := rui.StackLayoutByID(d.rootView, "rootViews"); stackLayout != nil {
		if d.pages[index].view == nil {
			d.pages[index].view = d.pages[index].creator(d.rootView.Session())
			stackLayout.Append(d.pages[index].view)
		} else {
			stackLayout.MoveToFront(d.pages[index].view)
		}
		rui.Set(d.rootView, "rootTitleText", rui.Text, d.pages[index].title)
		// d.rootView.Session().SetTitle(d.pages[index].title)
	}
}
func Web() {
	SessionStatus = make(map[int]bool)
	rui.ProtocolInDebugLog = false
	addr := rui.GetLocalIP() + ":8000"
	// addr := "localhost:8000"
	if rui.GetLocalIP() == "192.168.2.100" {
		rui.OpenBrowser("http://" + addr)
	}
	rui.StartApp(addr, CreateSession, rui.AppParams{
		Title:      "Ag-IRZ",
		Icon:       "icon.png",
		TitleColor: rui.Color(0xffc0ded9),
	})

}
