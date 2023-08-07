package radar

import (
	"sync"
	"time"

	"github.com/ruraomsk/irz/setup"
	"github.com/simonvetter/modbus"
)

type handler struct {
	lock    sync.Mutex
	uptime  time.Time
	holding [16]uint16
}

func (h *handler) HandleCoils(req *modbus.CoilsRequest) (res []bool, err error) {
	err = modbus.ErrIllegalFunction
	return
}

func (h *handler) HandleDiscreteInputs(req *modbus.DiscreteInputsRequest) (res []bool, err error) {
	err = modbus.ErrIllegalFunction
	return
}

func (h *handler) HandleHoldingRegisters(req *modbus.HoldingRegistersRequest) (res []uint16, err error) {
	err = nil
	h.lock.Lock()
	defer h.lock.Unlock()

	if req.UnitId != uint8(setup.Set.ModbusRadar.ID) {
		err = modbus.ErrIllegalFunction
		return
	}

	if int(req.Addr)+int(req.Quantity) > len(h.holding) {
		err = modbus.ErrIllegalDataAddress
		return
	}

	if req.IsWrite {
		h.uptime = time.Now()
	}
	for i := 0; i < int(req.Quantity); i++ {
		if req.IsWrite {
			h.holding[int(req.Addr)+i] = req.Args[i]
		}
		res = append(res, h.holding[int(req.Addr)+i])
	}
	return
}

func (h *handler) HandleInputRegisters(req *modbus.InputRegistersRequest) (res []uint16, err error) {
	err = modbus.ErrIllegalFunction
	return
}
