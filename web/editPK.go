package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/binding"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/worker"
)

// import (
// 	"fmt"

// 	"github.com/anoshenko/rui"
// 	"github.com/ruraomsk/ag-server/binding"
// 	"github.com/ruraomsk/ag-server/logger"
// 	"github.com/ruraomsk/irz/data"
// 	"github.com/ruraomsk/irz/worker"
// )

// const PkText = `
//
//					TabsLayout { id = tabsPK, width = 100%, height = 100%, tabs = top, tab-close-button = false,
//						content = [
//							TableView {cell-horizontal-align = right,  title = "ПК1", id="pk1"},
//							TableView {cell-horizontal-align = right,  title = "ПК2", id="pk2"},
//							TableView {cell-horizontal-align = right,  title = "ПК3", id="pk3"},
//							TableView {cell-horizontal-align = right,  title = "ПК4", id="pk4"},
//							TableView {cell-horizontal-align = right,  title = "ПК5", id="pk5"},
//							TableView {cell-horizontal-align = right,  title = "ПК6", id="pk6"},
//							TableView {cell-horizontal-align = right,  title = "ПК7", id="pk7"},
//							TableView {cell-horizontal-align = right,  title = "ПК8", id="pk8"},
//							TableView {cell-horizontal-align = right,  title = "ПК9", id="pk9"},
//							TableView {cell-horizontal-align = right,  title = "ПК10", id="pk10"},
//							TableView {cell-horizontal-align = right,  title = "ПК11", id="pk11"},
//							TableView {cell-horizontal-align = right,  title = "ПК12", id="pk12"},
//						]
//	}
//
// `
const PkEdit = `
	GridLayout{
		title="Редактор плана координации",
		content = [
				TextView{ row=0, text="Тип плана"},
				DropDownList {
					row=0,column=1,
					id=idType, width = 100%, height = 100%, orientation = vertical,current=0,
					items = ["Координированный план", "Локальный план","Локальное управление","Желтое мигание","Отключить светофор"],
				},
				TextView{ row=1, text="Время цикла"},
				NumberPicker {
					row=1,column=1,
					id=idTc,type=editor,min=0,max=1000,value=0
				},
				TextView{  row=2, text="Сдвиг "},
				NumberPicker {
					row=2,column=1,
					id=idShift,type=editor,min=0,max=1000,value=0
				},

				TextView{  row=3,column=0,text="Тип",text-size="24px"},
				TextView{  row=3,column=1,text="Номер фазы",text-size="24px"},
				TextView{  row=3,column=2,text="Длительность",text-size="24px"},

				DropDownList {
					row=4,column=0,
					id=idt1,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=4,column=1,
					id=idf1,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=4,column=2,
					id=idl1,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=5,column=0,
					id=idt2,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=5,column=1,
					id=idf2,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=5,column=2,
					id=idl2,type=editor,min=0,max=1000,value=0
				},
				DropDownList {
					row=6,column=0,
					id=idt3,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=6,column=1,
					id=idf3,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=6,column=2,
					id=idl3,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=7,column=0,
					id=idt4,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=7,column=1,
					id=idf4,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=7,column=2,
					id=idl4,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=8,column=0,
					id=idt5,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=8,column=1,
					id=idf5,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=8,column=2,
					id=idl5,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=9,column=0,
					id=idt6,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=9,column=1,
					id=idf6,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=9,column=2,
					id=idl6,type=editor,min=0,max=1000,value=0
				},
				DropDownList {
					row=10,column=0,
					id=idt7,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=10,column=1,
					id=idf7,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=10,column=2,
					id=idl7,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=11,column=0,
					id=idt8,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=11,column=1,
					id=idf8,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=11,column=2,
					id=idl8,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=12,column=0,
					id=idt9,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=12,column=1,
					id=idf9,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=12,column=2,
					id=idl9,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=13,column=0,
					id=idt10,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=13,column=1,
					id=idf10,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=13,column=2,
					id=idl10,type=editor,min=0,max=1000,value=0
				},
				DropDownList {
					row=14,column=0,
					id=idt11,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=14,column=1,
					id=idf11,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=14,column=2,
					id=idl11,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=15,column=0,
					id=idt12,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 1 ТВП","Зам"],
				},
				NumberPicker {
					row=15,column=1,
					id=idf12,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=15,column=2,
					id=idl12,type=editor,min=0,max=1000,value=0
				},
		
				Button{id=idSave,content="Сохранить", row=20,column=1},
				Button{id=idNo,content="Отменить", row=20,column=2},
		]
	}
`

var edPkPopup rui.Popup

func editPK(session rui.Session, plan int) {
	edit := rui.CreateViewFromText(session, PkEdit)
	var pk = binding.SetPk{Pk: 0}
	for _, v := range data.DataValue.Arrays.SetDK.DK {
		if v.Pk == plan {
			pk = v
		}
	}
	if pk.Pk == 0 {
		logger.Error.Printf("Нет плана координации %d", plan)
		rui.ShowMessage("Ошибка", fmt.Sprintf("Нет плана координации %d", plan), session)
		return
	}
	pk = worker.RepackPlan(pk)
	if pk.Tc < 3 {
		switch pk.Tc {
		case 0:
			rui.Set(edit, "idType", "current", 2)
		case 1:
			rui.Set(edit, "idType", "current", 3)
		case 2:
			rui.Set(edit, "idTtype", "current", 4)
		}
	} else {
		if pk.TypePU == 1 {
			rui.Set(edit, "idType", "current", 1)
		} else {
			rui.Set(edit, "idType", "current", 0)
		}
	}
	for i := 0; i < len(pk.Stages); i++ {

	}
	rui.Set(edit, "idTc", "value", pk.Tc)
	rui.Set(edit, "idShift", "value", pk.Shift)
	rui.Set(edit, "idSave", rui.ClickEvent, func(view rui.View) {
		rui.ShowMessage("Сохранение", fmt.Sprintf("План координации %d сохранен", plan), session)
		edPkPopup.Dismiss()
	})
	rui.Set(edit, "idNo", rui.ClickEvent, func(view rui.View) {
		rui.ShowMessage("Отмена", fmt.Sprintf("План координации %d не сохранен", plan), session)
		edPkPopup.Dismiss()
	})
	edPkPopup = rui.ShowPopup(edit, rui.Params{
		rui.Title: "Редактор плана координации",
	})
	// rui.ShowQuestion(title string, text string, session rui.Session, onYes func(), onNo func())
}
