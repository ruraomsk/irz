package web

import (
	"fmt"

	"github.com/anoshenko/rui"
	"github.com/ruraomsk/irz/data"
)

func showTech(view rui.View) {
	if data.DataValue.Controller.StatusCommandDU.IsReqSFDK1 {
		rui.Set(view, "idSFDK", "text", "СФДК")
	} else {
		rui.Set(view, "idSFDK", "text", "")
	}
	if data.DataValue.Controller.StatusCommandDU.IsDUDK1 {
		rui.Set(view, "idDU", "text", "ДУ ДК")
	} else {
		rui.Set(view, "idDU", "text", "")
	}
	if data.DataValue.Controller.StatusCommandDU.IsPK {
		rui.Set(view, "idPK", "text", "ПК")
	} else {
		rui.Set(view, "idPK", "text", "")
	}
	if data.DataValue.Controller.StatusCommandDU.IsCK {
		rui.Set(view, "idCK", "text", "CК")
	} else {
		rui.Set(view, "idCK", "text", "")
	}
	if data.DataValue.Controller.StatusCommandDU.IsNK {
		rui.Set(view, "idNK", "text", "НК")
	} else {
		rui.Set(view, "idNK", "text", "")
	}
	rui.Set(view, "idNPK", "text", fmt.Sprintf("ПК %d", data.DataValue.Controller.PK))
	rui.Set(view, "idNCK", "text", fmt.Sprintf("CК %d", data.DataValue.Controller.CK))
	rui.Set(view, "idNNK", "text", fmt.Sprintf("НК %d", data.DataValue.Controller.NK))
	rui.Set(view, "idTech", "text", getTechRezim(data.DataValue.Controller))
	rui.Set(view, "idRezim", "text", fmt.Sprintf("Режим %s", getRezim(data.DataValue.Controller)))
	rui.Set(view, "idPhaseDU", "text", fmt.Sprintf("Фаза ДУ %s", getPhaseDU(data.DataValue.Controller)))
	rui.Set(view, "idPhaseRU", "text", fmt.Sprintf("Фаза РУ %s", getPhaseRU(data.DataValue.Controller)))
	if data.DataValue.Controller.DK.EDK != 0 {
		rui.Set(view, "idBroken", "text", fmt.Sprintf("ОШИБКА %s", getBroken(data.DataValue.Controller)))
	}
}
