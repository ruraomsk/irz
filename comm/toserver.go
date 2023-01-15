package comm

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/transport"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/setup"
)

var socket net.Conn
var err error

func ToServer() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	data.DataValue.SetConnected(false)
	connected := false
	fromServer := make(chan transport.HeaderServer, 1)
	toServer := make(chan transport.HeaderDevice, 1)
	errTcp := make(chan net.Conn, 1)
	for {
		for !connected {
			socket, err = net.Dial("tcp", fmt.Sprintf("%s:%d", setup.Set.Server.Host, setup.Set.Server.Port))
			if err != nil {
				logger.Error.Printf("Error dial %s %s", setup.Set.Server.Host, err.Error())
				time.Sleep(10 * time.Second)
				continue
			}
			connected = true
		}
		logger.Info.Printf("Начинаем работу с %s", socket.RemoteAddr().String())
		readTout := time.Duration(int64(data.DataValue.Controller.Status.TObmen*60-30) * int64(time.Second))
		go transport.GetMessagesFromServer(socket, fromServer, &readTout, errTcp)
		go transport.SendMessagesToServer(socket, toServer, &readTout, errTcp)

		toServer <- makeHeaderForConnect()
		ticker := time.NewTicker(readTout)
		work := true
		for work {
			select {
			case <-ticker.C:
				toServer <- makeStatus()
			case in := <-fromServer:
				logger.Debug.Printf("Пришло сообщение %v", in)
				if !data.DataValue.Controller.IsConnected() {
					//Нам разрешили работать
					logger.Info.Printf("Произошло подключение к серверу %s", socket.RemoteAddr().String())
					data.DataValue.SetConnected(true)
					continue
				}
				replay, need := makeReplay(in)
				if need {
					toServer <- replay
				}
			case s := <-errTcp:
				logger.Debug.Printf("Ошибка приема %s", s.RemoteAddr().String())
				data.DataValue.SetConnected(false)
				socket.Close()
				time.Sleep(readTout)
				work = false
			case <-stop:
				logger.Info.Print("Останов системы")
				socket.Close()
				return
			}
		}
	}

	// data.DataValue.SetConnected(true)
}
func makeStatus() transport.HeaderDevice {
	var hd = transport.CreateHeaderDevice(data.DataValue.Controller.ID, 0, 0, 1)
	mss := make([]transport.SubMessage, 0)
	var ms transport.SubMessage
	ms.Set0x12Device(&data.DataValue.Controller)
	mss = append(mss, ms)
	hd.UpackMessages(mss)
	return hd
}
func makeReplay(in transport.HeaderServer) (transport.HeaderDevice, bool) {
	var hd = transport.CreateHeaderDevice(data.DataValue.Controller.ID, 0, 0, 1)
	mss := make([]transport.SubMessage, 0)
	var ms transport.SubMessage
	dmess := in.ParseMessage()
	need := false
	sendStatus := false
	if in.Number != 0 {
		ms.Set0x01Device(int(in.Number), time.Now().Minute(), time.Now().Second(), 0, 0)
		mss = append(mss, ms)
	}
	for _, mes := range dmess {
		switch mes.Type {
		case 1:
			logger.Info.Printf("Подтвердили сообощение %d", mes.Get0x01Server())
		case 2:
			if mes.Get0x02Server() {
				logger.Info.Print("Включить управление")
			} else {
				logger.Info.Print("Выключить управление")
			}
		case 3:
			//Запрос состояния устройства
			sendStatus = true
			need = true
		case 4:
			//Включение выключение СФДК
			data.DataValue.SetSFDK(mes.Get0x04Server()[0])
			sendStatus = true
			need = true
		case 5:
			//Смена плана ПК
			data.DataValue.SetPK(mes.Get0x05Server())
		case 6:
			//Смена НК
			data.DataValue.SetNK(mes.Get0x06Server())
		case 7:
			//Смена CК
			data.DataValue.SetCK(mes.Get0x07Server())
		case 9:
			//Смена ДУ
			data.DataValue.SetDU(mes.Get0x09Server())
		}
	}
	if need {
		if sendStatus {
			ms.Set0x12Device(&data.DataValue.Controller)
			mss = append(mss, ms)
		}
		hd.UpackMessages(mss)
	}
	return hd, need
}
func makeHeaderForConnect() transport.HeaderDevice {
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
	return hd
}
