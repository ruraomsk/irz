package visio

import (
	"time"

	"github.com/goburrow/serial"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/setup"
)

var work bool
var senderChan chan toSend

func sender() {
	mode := serial.Config{Address: setup.Set.VisioDevice.Device,
		BaudRate: setup.Set.VisioDevice.BaudRate,
		Parity:   setup.Set.VisioDevice.Parity,
		RS485:    serial.RS485Config{Enabled: true},
	}
	senderChan = make(chan toSend)
	for {
		work = false
		port, err := serial.Open(&mode)
		if err != nil {
			logger.Error.Printf("%s", err.Error())
			time.Sleep(time.Second)
			continue
		}
		work = true
		for {
			m := <-senderChan
			_, err := port.Write(m.buff160[:])
			if err != nil {
				logger.Error.Printf("%s", err.Error())
				break
			}
			_, err = port.Write(m.buff168[:])
			if err != nil {
				logger.Error.Printf("%s", err.Error())
				break
			}

		}
		work = false
		port.Close()
	}
}
