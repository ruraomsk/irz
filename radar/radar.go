package radar

import (
	"fmt"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/setup"
	"github.com/ruraomsk/irz/stat"
	"github.com/simonvetter/modbus"
)

var eh *handler
var work = false

func GetValues() string {
	if !work {
		return "Не запущен пока"
	}
	eh.lock.Lock()
	defer eh.lock.Unlock()
	return fmt.Sprintf("%s %v", eh.uptime.String(), eh.holding)
}
func Radar() {
	if !setup.Set.ModbusRadar.Radar {
		return
	}
	eh = &handler{uptime: time.Unix(0, 0)}
	go modbusServer()
	go pusher()
	work = true
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
		eh.lock.Lock()
		var send []stat.OneTick
		if time.Now().Sub(eh.uptime).Seconds() > 2 {
			send = badStatistics()
		} else {
			send = goodStatistics()
		}
		eh.lock.Unlock()
		for _, v := range send {
			stat.InStat <- v
		}
	}

}
func badStatistics() []stat.OneTick {
	r := make([]stat.OneTick, 0)
	t := time.Now()
	for i := 0; i < setup.Set.ModbusRadar.Chanels; i++ {
		r = append(r, stat.OneTick{Nomber: i + 3, Value: stat.Value{Status: 1, Time: t, Value: 0}})
	}
	return r
}
func goodStatistics() []stat.OneTick {
	r := make([]stat.OneTick, 0)
	t := time.Now()
	for i := 0; i < setup.Set.ModbusRadar.Chanels; i++ {
		r = append(r, stat.OneTick{Nomber: i + 3, Value: stat.Value{Status: 0, Time: t, Value: int(eh.holding[i])}})
	}
	return r
}

var server *modbus.ModbusServer
var err error

func modbusServer() {
	server, err = modbus.NewServer(&modbus.ServerConfiguration{
		URL:        fmt.Sprintf("tcp://0.0.0.0:%d", setup.Set.ModbusRadar.Port),
		Timeout:    30 * time.Second,
		MaxClients: 5,
	}, eh)
	if err != nil {
		logger.Error.Printf("Не могу создать сервер %v", err)
		return
	}
	err = server.Start()
	if err != nil {
		logger.Error.Printf("Не могу запустить сервер %v", err)
		return
	}
	ticker := time.NewTicker(time.Second)
	for {
		<-ticker.C
	}

}
