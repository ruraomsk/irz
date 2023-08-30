package comm

import (
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/ag-server/transport"
	"github.com/ruraomsk/irz/data"
)

func makeStatistics(s pudge.Statistic) transport.HeaderDevice {
	var hd = transport.CreateHeaderDevice(data.DataValue.Controller.ID, 0, 0, 1)
	mss := make([]transport.SubMessage, 0)
	var ms transport.SubMessage
	ms.Set0x09Device(&s)
	mss = append(mss, ms)
	ms.Set0x0ADevice(&s)
	// logger.Debug.Printf("message %v", ms)
	mss = append(mss, ms)
	hd.UpackMessages(mss)
	return hd
}
func makeStatus() transport.HeaderDevice {
	var hd = transport.CreateHeaderDevice(data.DataValue.Controller.ID, 0, 0, 1)
	mss := make([]transport.SubMessage, 0)
	var ms transport.SubMessage
	ms.Set0x12Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	ms.Set0x11Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	ms.Set0x10Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	hd.UpackMessages(mss)
	return hd
}
func makeSFDKreplay() transport.HeaderDevice {
	var hd = transport.CreateHeaderDevice(data.DataValue.Controller.ID, 0, 0, 1)
	mss := make([]transport.SubMessage, 0)
	var ms transport.SubMessage
	ms.Set0x12Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	hd.UpackMessages(mss)
	return hd
}
func makeHeaderForConnect() transport.HeaderDevice {
	var hd = transport.CreateHeaderDevice(data.DataValue.Controller.ID, 0, 1, 0x7f)

	mss := make([]transport.SubMessage, 0)
	var ms transport.SubMessage
	ms.Set0x1DDevice(&data.DataValue.Controller)
	mss = append(mss, ms)
	ms.Set0x10Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	ms.Set0x12Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	// ms.Set0x1BDevice(&data.DataValue.Controller)
	// mss = append(mss, ms)
	// ms.Set0x11Device(&data.DataValue.Controller)
	// mss = append(mss, ms)
	hd.UpackMessages(mss)
	return hd
}
