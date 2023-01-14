package comm

import (
	"fmt"
	"net"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/transport"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/setup"
)

var socket net.Conn
var err error

func ToServer() {
	data.DataValue.SetConnected(false)
	connected := false
	for !connected {
		socket, err = net.Dial("tcp", fmt.Sprintf("%s:%d", setup.Set.Server.Host, setup.Set.Server.Port))
		if err != nil {
			logger.Error.Printf("Error dial %s %s", setup.Set.Server.Host, err.Error())
			time.Sleep(10 * time.Second)
			continue
		}
		connected = true
	}
	logger.Info.Printf("Начинием работу с %s", socket.RemoteAddr().String())
	fromServer := make(chan transport.HeaderServer, 1)
	toServer := make(chan transport.HeaderDevice, 1)
	errTcp := make(chan net.Conn, 1)
	readTout := time.Duration(int64(data.DataValue.Controller.Status.TObmen*60-30) * int64(time.Second))
	go transport.GetMessagesFromService(socket, fromServer, &readTout, errTcp)
	go transport.SendMessagesToServer(socket, toServer, &readTout, errTcp)

	var hd = transport.CreateHeaderDevice(data.DataValue.Controller.ID, 0, 1, 1)

	mss := make([]transport.SubMessage, 0)
	var ms transport.SubMessage
	ms.Set0x1DDevice(&data.DataValue.Controller)
	mss = append(mss, ms)
	ms.Set0x10Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	ms.Set0x12Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	ms.Set0x1BDevice(&data.DataValue.Controller)
	mss = append(mss, ms)
	ms.Set0x11Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	hd.UpackMessages(mss)
	toServer <- hd
	logger.Info.Printf("Отправили %v", hd.Message)
	for {
		select {
		case in := <-fromServer:
			logger.Debug.Printf("Пришло сообщение %v", in)
		case s := <-errTcp:
			logger.Debug.Printf("Ошибка приема %s", s.RemoteAddr().String())
		}
	}
	// data.DataValue.SetConnected(true)
}
