package stat

import (
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/irz/data"
)

var InStat chan OneTick
var ChangeArrays chan interface{}
var ChangeCtrl chan interface{}
var getInfo chan int
var sendInfo chan Value

var sendStat chan SendStat
var statistics Chanels

var GetLast chan interface{}
var LastSending chan LastSend

type LastSend struct {
	LastTime time.Time
	Last     pudge.Statistic
}

var lastSend = LastSend{LastTime: time.Unix(0, 0), Last: pudge.Statistic{}}

func Statistics() {
	InStat = make(chan OneTick, 100)
	ChangeArrays = make(chan interface{})
	ChangeCtrl = make(chan interface{})
	sendStat = make(chan SendStat)
	getInfo = make(chan int)
	sendInfo = make(chan Value)

	GetLast = make(chan interface{})
	LastSending = make(chan LastSend)

	statistics = newStatistics()
	statistics.Cron.Start()
	go CtrlInputs()
	// go pusher(18)
	for {
		select {
		case <-ChangeArrays:
			logger.Debug.Println("Смена массивов привязки")
			statistics.Cron.Stop()
			statistics = newStatistics()
			statistics.Cron.Start()
			ChangeCtrl <- 0
		case in := <-InStat:
			statistics.AddInfo(in)
		case qa := <-data.QAInfo:
			data.AInfo <- statistics.GetInfo(qa.TypeDev, qa.Interval)
		case head := <-sendStat:
			makerStat(head)
		case d := <-getInfo:
			sendInfo <- statistics.GetLast(d)
		case <-GetLast:
			LastSending <- lastSend
		}
	}

}
func makerStat(head SendStat) {
	// logger.Debug.Printf("Готовим статистику %2d %2d", head.Hour, head.Minute)
	if data.DataValue.Controller.Base {
		statistics.ClearInfo()
		return
	}

	for TypeStat := 1; TypeStat <= 2; TypeStat++ {
		r := pudge.Statistic{Period: head.Period, Type: TypeStat, TLen: head.Interval, Hour: head.Hour, Min: head.Minute, Datas: make([]pudge.DataStat, 0)}
		for _, v := range statistics.Chanels {
			if v.TypeStat == TypeStat {
				rez := 0
				count := 0
				status := 0
				for _, d := range v.Values {
					rez += d.Value
					count++
					status = status | d.Status
				}
				switch TypeStat {
				case 1:
					r.Datas = append(r.Datas, pudge.DataStat{Chanel: v.Number, Status: status, Intensiv: rez})
				case 2:
					rez = rez / count
					r.Datas = append(r.Datas, pudge.DataStat{Chanel: v.Number, Status: status, Speed: rez})

				}
			}
		}
		if len(r.Datas) > 0 {
			// logger.Debug.Printf("Послали статистику %v", r)
			lastSend = LastSend{LastTime: time.Now(), Last: r}
			data.Statistics <- r
		}
	}
	statistics.ClearInfo()
}
