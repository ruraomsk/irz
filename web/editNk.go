package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/binding"
	"github.com/ruraomsk/irz/data"
)

const NkEdit = `
	GridLayout{
		title="Редактор недельной карты",
		content = [

				TextView{  row=0,column=0,text="День недели",text-size="24px"},
				TextView{  row=0,column=1,text="Суточная карта",text-size="24px"},
				TextView{  row=1,column=0,text="Понедельник"},
				TextView{  row=2,column=0,text="Вторник"},
				TextView{  row=3,column=0,text="Среда"},
				TextView{  row=4,column=0,text="Четверг"},
				TextView{  row=5,column=0,text="Пятница"},
				TextView{  row=6,column=0,text="Суббота"},
				TextView{  row=7,column=0,text="Воскресенье"},


				NumberPicker {
					row=1,column=1,
					id=id1,type=editor,min=0,max=12,value=0
				},
				NumberPicker {
					row=2,column=1,
					id=id2,type=editor,min=0,max=12,value=0
				},
				NumberPicker {
					row=3,column=1,
					id=id3,type=editor,min=0,max=12,value=0
				},
				NumberPicker {
					row=4,column=1,
					id=id4,type=editor,min=0,max=12,value=0
				},
				NumberPicker {
					row=5,column=1,
					id=id5,type=editor,min=0,max=12,value=0
				},
				NumberPicker {
					row=6,column=1,
					id=id6,type=editor,min=0,max=12,value=0
				},
				NumberPicker {
					row=7,column=1,
					id=id7,type=editor,min=0,max=12,value=0
				},

				Button{id=idSave,content="Сохранить", row=20,column=0},
				Button{id=idNo,content="Отменить", row=20,column=1},
		]
	}
`

var edNkPopup rui.Popup

func editNK(session rui.Session, cart int) {
	edit := rui.CreateViewFromText(session, NkEdit)

	var nk = binding.OneWeek{Number: 0}
	for _, v := range data.DataValue.Arrays.WeekSets.WeekSets {
		if v.Number == cart {
			nk = v
		}
	}

	i := 1
	for _, v := range nk.Days {
		rui.Set(edit, fmt.Sprintf("id%d", i), "value", v)
		i++
	}

	rui.Set(edit, "idSave", rui.ClickEvent, func(view rui.View) {

		for i := 0; i < len(nk.Days); i++ {
			nk.Days[i] = getInteger(rui.Get(edit, fmt.Sprintf("id%d", i+1), "value"))

		}
		rui.ShowMessage("Сохранение", fmt.Sprintf("Недельная карта %d сохранена", cart), session)
		edNkPopup.Dismiss()
	})
	rui.Set(edit, "idNo", rui.ClickEvent, func(view rui.View) {
		edNkPopup.Dismiss()
	})
	edNkPopup = rui.ShowPopup(edit, rui.Params{
		rui.Title: "Редактор недельной карты",
	})

}
