package kdm

import (
	"fmt"

	"github.com/ruraomsk/irz/data"
	"github.com/simonvetter/modbus"
)

func readData(rq Request) Replay {
	rep := Replay{Request: rq}
	rep.Data, rep.Status = client.ReadRegisters(rq.Start, rq.Lenght, modbus.HOLDING_REGISTER)
	return rep
}
func writeData(wr WriteCmd) Replay {
	rep := Replay{Request: wr.Request}
	rep.Status = client.WriteRegisters(wr.Request.Start, wr.Data)
	return rep
}
func GetData(rq Request) Replay {
	var rep Replay = Replay{Request: rq}
	if !data.DataValue.Connect {
		rep.Status = fmt.Errorf("Нет соединения с Modbus")
		return rep
	}
	RequestChan <- rq
	rep = <-ReplayChan
	return rep
}
func SetData(wr WriteCmd) Replay {
	var rep Replay = Replay{Request: wr.Request}
	if !data.DataValue.Connect {
		rep.Status = fmt.Errorf("Нет соединения с Modbus")
		return rep
	}
	WriteCmdChan <- wr
	rep = <-ReplayChan
	return rep
}
