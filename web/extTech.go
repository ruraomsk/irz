package web

import (
	"fmt"

	"github.com/ruraomsk/ag-server/pudge"
)

func getPhaseDU(cc pudge.Controller) string {
	// return fmt.Sprintf("%d %d %d", cc.DK.FDK, cc.DK.FTUDK, cc.DK.FTSDK)
	switch cc.DK.FTUDK {
	case 0:
		return "ЛР"
	case 9:
		return "ПрТакт"
	case 10:
		return "ЖМ"
	case 11:
		return "ОС"
	case 12:
		return "КК"
	case 14:
		return "ЖМ"
	case 15:
		return "ОС"
	}
	return fmt.Sprintf("%d", cc.DK.FTUDK)
}
func getPhaseRU(cc pudge.Controller) string {
	// return fmt.Sprintf("%d %d %d", cc.DK.FDK, cc.DK.FTUDK, cc.DK.FTSDK)
	switch cc.DK.FDK {
	case 0:
		return "ОС"
	case 9:
		return "ПрТакт"
	case 10:
		return "ЖМ"
	case 11:
		return "ОС"
	case 12:
		return "КК"
	case 14:
		return "ЖМ"
	case 15:
		return "ОС"
	}
	return fmt.Sprintf("%d", cc.DK.FDK)
}
func getRezim(cc pudge.Controller) string {
	switch cc.DK.RDK {
	case 1:
		return "РУ"
	case 2:
		return "РП"
	case 3:
		return "ЗУ"
	case 4:
		return "ДУ"
	case 5:
		return "ЛУ"
	case 6:
		return "ЛУ"
	case 7:
		return "РП"
	case 8:
		return "КУ"
	case 9:
		return "КУ"
	}
	return fmt.Sprintf("КОД %d ", cc.DK.RDK)
}
func getBroken(cc pudge.Controller) string {
	switch cc.DK.EDK {
	case 0:
		return "НОРМ"
	case 1:
		return "ПЕРЕХОД"
	case 2:
		return "ОБРЫВ ЛС"
	case 3:
		return "НГ паритет"
	case 4:
		return "Нет кода"
	case 5:
		return "ОС КОНФЛИКТ"
	case 6:
		return "ЖМ перегорание"
	case 7:
		return "Невкл в коорд"
	case 8:
		return "Неподчинение"
	case 9:
		return "Длинный промтакт"
	case 10:
		return "Нет фазы"
	case 11:
		return "Обрыв ЛС с КЗЦ"
	case 12:
		return "Обрыв ЛС с ЭВМ"
	case 13:
		return "Нет информации"
	}
	return fmt.Sprintf("КОД %d ", cc.DK.EDK)
}
func getTechRezim(cc pudge.Controller) string {
	switch cc.TechMode {
	case 1:
		return "ВР-СК"
	case 2:
		return "ВР-НК"
	case 3:
		return "ДУ-СК"
	case 4:
		return "ДУ-НК"
	case 5:
		return "ДУ-ПК"
	case 6:
		return "РП"
	case 7:
		return "КП ИП"
	case 8:
		return "КП С"
	case 9:
		return "ВР"
	case 10:
		return "ПК ХТ"
	case 11:
		return "ПК КТ"
	case 12:
		return "ПЗУ"
	}
	return fmt.Sprintf("КОД %d ", cc.TechMode)
}
