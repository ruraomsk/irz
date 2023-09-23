package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/binding"
	"github.com/ruraomsk/irz/data"
)

const CkEdit = `
	GridLayout{
		title="Редактор суточной карты",
		content = [

				TextView{  row=0,column=0,text="Час",text-size="24px"},
				TextView{  row=0,column=1,text="Мин",text-size="24px"},
				TextView{  row=0,column=2,text="План координации",text-size="24px"},

				NumberPicker {
					row=1,column=0,
					id=idh1,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=1,column=1,
					id=idm1,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=1,column=2,
					id=idp1,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=2,column=0,
					id=idh2,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=2,column=1,
					id=idm2,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=2,column=2,
					id=idp2,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=3,column=0,
					id=idh3,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=3,column=1,
					id=idm3,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=3,column=2,
					id=idp3,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=4,column=0,
					id=idh4,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=4,column=1,
					id=idm4,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=4,column=2,
					id=idp4,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=5,column=0,
					id=idh5,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=5,column=1,
					id=idm5,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=5,column=2,
					id=idp5,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=6,column=0,
					id=idh6,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=6,column=1,
					id=idm6,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=6,column=2,
					id=idp6,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=7,column=0,
					id=idh7,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=7,column=1,
					id=idm7,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=7,column=2,
					id=idp7,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=8,column=0,
					id=idh8,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=8,column=1,
					id=idm8,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=8,column=2,
					id=idp8,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=9,column=0,
					id=idh9,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=9,column=1,
					id=idm9,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=9,column=2,
					id=idp9,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=10,column=0,
					id=idh10,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=10,column=1,
					id=idm10,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=10,column=2,
					id=idp10,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=11,column=0,
					id=idh11,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=11,column=1,
					id=idm11,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=11,column=2,
					id=idp11,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=12,column=0,
					id=idh12,type=editor,min=0,max=23,value=0
				},
				NumberPicker {
					row=12,column=1,
					id=idm12,type=editor,min=0,max=59,value=0
				},
				NumberPicker {
					row=12,column=2,
					id=idp12,type=editor,min=0,max=100,value=0
				},

				Button{id=idSave,content="Сохранить", row=20,column=0},
				Button{id=idNo,content="Отменить", row=20,column=2},
		]
	}
`

var edCkPopup rui.Popup

func editCK(session rui.Session, cart int) {
	mutex.Lock()
	defer mutex.Unlock()

	edit := rui.CreateViewFromText(session, CkEdit)

	var ck = binding.OneDay{Number: 0}
	for _, v := range data.DataValue.Arrays.DaySets.DaySets {
		if v.Number == cart {
			ck = *v
		}
	}

	i := 1
	for _, v := range ck.Lines {
		rui.Set(edit, fmt.Sprintf("idh%d", i), "value", v.Hour)
		rui.Set(edit, fmt.Sprintf("idm%d", i), "value", v.Min)
		rui.Set(edit, fmt.Sprintf("idp%d", i), "value", v.PKNom)
		i++
	}

	rui.Set(edit, "idSave", rui.ClickEvent, func(view rui.View) {

		for i := 0; i < len(ck.Lines); i++ {
			ck.Lines[i].Hour = getInteger(rui.Get(edit, fmt.Sprintf("idh%d", i+1), "value"))
			ck.Lines[i].Min = getInteger(rui.Get(edit, fmt.Sprintf("idm%d", i+1), "value"))
			ck.Lines[i].PKNom = getInteger(rui.Get(edit, fmt.Sprintf("idp%d", i+1), "value"))

		}
		for _, v := range data.DataValue.Arrays.DaySets.DaySets {
			if v.Number == cart {
				v = &ck
			}
		}
		rui.ShowMessage("Сохранение", fmt.Sprintf("Суточная карта %d сохранена", cart), session)
		updatedCk = true
		edCkPopup.Dismiss()
	})
	rui.Set(edit, "idNo", rui.ClickEvent, func(view rui.View) {
		edCkPopup.Dismiss()
	})
	edCkPopup = rui.ShowPopup(edit, rui.Params{
		rui.Title: "Редактор суточной карты",
	})

}
