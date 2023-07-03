package kdm

import (
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
