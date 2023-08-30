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
	"github.com/ruraomsk/irz/stat"
)

var socket net.Conn
var err error
var priv bool
var lastTechmode int
var areaPriv []uint8

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
				slp := time.NewTimer(10 * time.Second)
				to := true
				for to {
					select {
					case cmd := <-data.ToServer:
						logger.Info.Printf("Пропущена команда %d", cmd)
					case <-slp.C:
						to = false
					case s := <-data.Statistics:
						logger.Info.Printf("Пропущена отправка статистики %v", s)
					}
				}
				continue
			}
			connected = true
		}
		logger.Info.Printf("Начинаем работу с %s", socket.RemoteAddr().String())
		readTout := time.Duration(int64(data.DataValue.Controller.Status.TObmen*60-30) * int64(time.Second))
		longTout := time.Duration(24 * time.Hour)
		go transport.GetMessagesFromServer(socket, fromServer, &longTout, errTcp)
		go transport.SendMessagesToServer(socket, toServer, &readTout, errTcp)
		priv = false
		toServer <- makeHeaderForConnect()
		ticker := time.NewTicker(1 * time.Second)
		work := true
		for work {
			select {
			case stat := <-data.Statistics:
				toServer <- makeStatistics(stat)
			case cmd := <-data.ToServer:
				// logger.Info.Print("Просят по СФДК отправить")
				// logger.Debug.Printf("Пришла команда %d", cmd)
				if cmd == 1 { //Что то изменилось и нам нужно сообщить об этом
					data.DataValue.SetLastOperation()
					toServer <- makeStatus()
				}
				if cmd == 0 && data.DataValue.Controller.StatusCommandDU.IsReqSFDK1 {
					// Изменине слать если есть контроль со стороны сервера
					data.DataValue.SetLastOperation()
					toServer <- makeSFDKreplay()
				}

			case <-ticker.C:
				data.DataValue.SetNowTime()
				if data.DataValue.Controller.IsConnected() {
					if time.Since((data.DataValue.Controller.LastOperation)) > readTout {
						data.DataValue.SetLastOperation()
						toServer <- makeStatus()
					}
				}
			case in := <-fromServer:
				logger.Debug.Printf("Пришло сообщение %v", in)
				data.DataValue.SetLastOperation()
				if !data.DataValue.Controller.IsConnected() {
					//Нам разрешили работать
					data.DataValue.SetConnectTime()
					logger.Info.Printf("Произошло подключение к серверу %s", socket.RemoteAddr().String())
					data.DataValue.SetConnected(true)
					continue
				}
				replay, need := makeReplay(in)
				if need {
					data.DataValue.SetLastOperation()
					toServer <- replay
				}
			case s := <-errTcp:
				logger.Debug.Printf("Ошибка приема %s", s.RemoteAddr().String())
				data.DataValue.SetConnected(false)
				socket.Close()
				connected = false
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

func makeReplay(in transport.HeaderServer) (transport.HeaderDevice, bool) {
	var hd = transport.CreateHeaderDevice(data.DataValue.Controller.ID, 0, 0, 1)
	mss := make([]transport.SubMessage, 0)
	var ms transport.SubMessage
	need := false
	if in.Number != 0 {
		ms.Set0x01Device(int(in.Number), time.Now().Minute(), time.Now().Second(), 0, 0)
		mss = append(mss, ms)
		need = true
	}
	if priv {
		if in.Message[0] == 0 && in.Message[1] == 2 && in.Message[2] == 2 && in.Message[3] == 2 {
			logger.Info.Print("Включить управление")
			priv = false
			data.DataValue.Controller.TechMode = lastTechmode
			if moveArrasIsGood() {
				logger.Info.Print("Привязки хорошие")
				stat.ChangeArrays <- 1
				data.DataValue.SetBase(false)
			} else {
				logger.Info.Print("Привязки плохие")
			}
			ms.Set0x12Device(&data.DataValue.Controller)
			mss = append(mss, ms)
			ms.Set0x11Device(&data.DataValue.Controller)
			mss = append(mss, ms)
			need = true
		} else {
			if in.Message[0] != 0 {
				areaPriv = append(areaPriv, in.Message...)
			}
		}
		if need {
			hd.UpackMessages(mss)
		}
		return hd, need
	}
	dmess := in.ParseMessage()
	sendStatus := false
	for _, mes := range dmess {
		switch mes.GetCodeCommandServer() {
		case 1:
			logger.Info.Printf("Подтвердили сообощение %d", mes.Get0x01Server())
		case 2:
			if mes.Get0x02Server() {
				logger.Info.Print("Включить управление передано перед отключить")
			} else {
				logger.Info.Print("Выключить управление")
				priv = true
				areaPriv = make([]uint8, 0)
				lastTechmode = data.DataValue.Controller.TechMode
				data.DataValue.SetTechMode(8)
				sendStatus = true
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
			data.Commands <- data.InternalCmd{Source: data.IBM, Command: 5, Parametr: mes.Get0x05Server()}
		case 6:
			//Смена НК
			data.Commands <- data.InternalCmd{Source: data.IBM, Command: 6, Parametr: mes.Get0x06Server()}
		case 7:
			//Смена CК
			data.Commands <- data.InternalCmd{Source: data.IBM, Command: 7, Parametr: mes.Get0x07Server()}
		case 9:
			//Смена ДУ
			data.Commands <- data.InternalCmd{Source: data.IBM, Command: 9, Parametr: mes.Get0x09Server()}
		}
	}
	if need {
		if sendStatus {
			ms.Set0x12Device(&data.DataValue.Controller)
			mss = append(mss, ms)
			ms.Set0x11Device(&data.DataValue.Controller)
			mss = append(mss, ms)
		}
		hd.UpackMessages(mss)
	}
	return hd, need
}
