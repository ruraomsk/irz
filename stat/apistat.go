package stat

import (
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

func Statistics() {
	InStat = make(chan OneTick, 100)
	ChangeArrays = make(chan interface{})
	ChangeCtrl = make(chan interface{})
	sendStat = make(chan SendStat)
	getInfo = make(chan int)
	sendInfo = make(chan Value)
	statistics = newStatistics()
	statistics.Cron.Start()
	go CtrlInputs()
	// go pusher(18)
	for {
		select {
		case <-ChangeArrays:
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
		}
	}

}
func makerStat(head SendStat) {
	// logger.Debug.Printf("Готовим статистику %d %d", head.Hour, head.Minute)
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
			// logger.Debug.Printf("Послали статистику по %d ", TypeStat)
			data.Statistics <- r
		}
	}
	statistics.ClearInfo()
}
