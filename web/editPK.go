package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/binding"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/worker"
)

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
				TextView{  row=3,column=2,text="Старт",text-size="24px"},
				TextView{  row=3,column=3,text="Стоп",text-size="24px"},

				DropDownList {
					row=4,column=0,
					id=idt1,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=4,column=1,
					id=idf1,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=4,column=2,
					id=ids1,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=4,column=3,
					id=ide1,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=5,column=0,
					id=idt2,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=5,column=1,
					id=idf2,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=5,column=2,
					id=ids2,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=5,column=3,
					id=ide2,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=6,column=0,
					id=idt3,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=6,column=1,
					id=idf3,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=6,column=2,
					id=ids3,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=6,column=3,
					id=ide3,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=7,column=0,
					id=idt4,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=7,column=1,
					id=idf4,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=7,column=2,
					id=ids4,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=7,column=3,
					id=ide4,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=8,column=0,
					id=idt5,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=8,column=1,
					id=idf5,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=8,column=2,
					id=ids5,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=8,column=3,
					id=ide5,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=9,column=0,
					id=idt6,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=9,column=1,
					id=idf6,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=9,column=2,
					id=ids6,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=9,column=3,
					id=ide6,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=10,column=0,
					id=idt7,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=10,column=1,
					id=idf7,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=10,column=2,
					id=ids7,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=10,column=3,
					id=ide7,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=11,column=0,
					id=idt8,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=11,column=1,
					id=idf8,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=11,column=2,
					id=ids8,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=11,column=3,
					id=ide8,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=12,column=0,
					id=idt9,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=12,column=1,
					id=idf9,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=12,column=2,
					id=ids9,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=12,column=3,
					id=ide9,type=editor,min=0,max=100,value=0
				},

				DropDownList {
					row=13,column=0,
					id=idt10,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=13,column=1,
					id=idf10,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=13,column=2,
					id=ids10,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=13,column=3,
					id=ide10,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=14,column=0,
					id=idt11,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=14,column=1,
					id=idf11,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=14,column=2,
					id=ids11,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=14,column=3,
					id=ide11,type=editor,min=0,max=1000,value=0
				},

				DropDownList {
					row=15,column=0,
					id=idt12,  orientation = vertical,current=0,
					items = ["Простая", "МГР","1ТВП","2ТВП","1,2ТВП","Зам 1 ТВП","Зам 2 ТВП","Зам 1,2"],
				},
				NumberPicker {
					row=15,column=1,
					id=idf12,type=editor,min=0,max=100,value=0
				},
				NumberPicker {
					row=15,column=2,
					id=ids12,type=editor,min=0,max=1000,value=0
				},
				NumberPicker {
					row=15,column=3,
					id=ide12,type=editor,min=0,max=1000,value=0
				},
		
				Button{id=idSave,content="Сохранить", row=20,column=1},
				Button{id=idNo,content="Отменить", row=20,column=2},
		]
	}
`

var edPkPopup rui.Popup

func editPK(session rui.Session, plan int) {
	mutex.Lock()
	defer mutex.Unlock()
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
			rui.Set(edit, "idType", "current", 4)
		}
	} else {
		if pk.TypePU == 1 {
			rui.Set(edit, "idType", "current", 1)
		} else {
			rui.Set(edit, "idType", "current", 0)
		}
	}
	i := 1
	for _, v := range pk.Stages {

		tf := fmt.Sprintf("idt%d", i)
		// fmt.Printf("%s %d \n", tf, v.Nline)

		switch v.Tf {
		case 0:
			rui.Set(edit, tf, "current", 0)
		case 1:
			rui.Set(edit, tf, "current", 1)
		case 2:
			rui.Set(edit, tf, "current", 2)
		case 3:
			rui.Set(edit, tf, "current", 3)
		case 4:
			rui.Set(edit, tf, "current", 4)
		case 5:
			rui.Set(edit, tf, "current", 5)
		case 6:
			rui.Set(edit, tf, "current", 6)
		case 7:
			rui.Set(edit, tf, "current", 7)
		}
		rui.Set(edit, fmt.Sprintf("idf%d", i), "value", v.Number)
		rui.Set(edit, fmt.Sprintf("ids%d", i), "value", v.Start)
		rui.Set(edit, fmt.Sprintf("ide%d", i), "value", v.Stop)
		i++
	}
	rui.Set(edit, "idTc", "value", pk.Tc)
	rui.Set(edit, "idShift", "value", pk.Shift)
	rui.Set(edit, "idSave", rui.ClickEvent, func(view rui.View) {
		j := rui.GetCurrent(edit, "idType")
		// fmt.Println(j)
		if j > 1 {
			pk.Tc = j - 2
			pk.TypePU = 1
			for i := 0; i < len(pk.Stages); i++ {
				pk.Stages[i].Number = 0
				pk.Stages[i].Tf = 0
				pk.Stages[i].Start = 0
				pk.Stages[i].Stop = 0
			}
		} else {
			if j == 0 {
				pk.TypePU = 10
			} else {
				pk.TypePU = 1
			}
			pk.Tc = getInteger(rui.Get(edit, "idTc", "value"))
			pk.Shift = getInteger(rui.Get(edit, "idShift", "value"))
			for i := 0; i < len(pk.Stages); i++ {
				pk.Stages[i].Number = getInteger(rui.Get(edit, fmt.Sprintf("idf%d", i+1), "value"))
				pk.Stages[i].Start = getInteger(rui.Get(edit, fmt.Sprintf("ids%d", i+1), "value"))
				pk.Stages[i].Stop = getInteger(rui.Get(edit, fmt.Sprintf("ide%d", i+1), "value"))
				pk.Stages[i].Tf = rui.GetCurrent(edit, fmt.Sprintf("idt%d", i+1))
			}
		}
		for i := 0; i < len(data.DataValue.Arrays.SetDK.DK); i++ {
			if data.DataValue.Arrays.SetDK.DK[i].Pk != pk.Pk {
				continue
			}
			data.DataValue.Arrays.SetDK.DK[i] = pk
			break
		}
		rui.ShowMessage("Сохранение", fmt.Sprintf("План координации %d сохранен", plan), session)
		updatedPk = true
		edPkPopup.Dismiss()
	})
	rui.Set(edit, "idNo", rui.ClickEvent, func(view rui.View) {
		// rui.ShowMessage("Отмена", fmt.Sprintf("План координации %d не сохранен", plan), session)
		edPkPopup.Dismiss()
	})
	edPkPopup = rui.ShowPopup(edit, rui.Params{
		rui.Title: "Редактор плана координации",
	})
}
