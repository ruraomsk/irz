package stat

import (
	"strconv"
	"strings"
	"time"

	"github.com/robfig/cron/v3"
	"github.com/ruraomsk/irz/data"
)

type SendStat struct {
	Period   int
	Interval int
	Hour     int
	Minute   int
}
type OneTick struct {
	Nomber int
	Value  Value
}
type Chanels struct {
	Cron     *cron.Cron
	Interval int
	Chanels  map[int]*OneChanel
}
type OneChanel struct {
	Number   int
	TypeStat int // 0 - нет статистики 1 - интенсивность 2 - скорость
	TypeDev  int // 0 для сбора статистики 1 - твп 1 2 ТПВ2 3 МГР 4 ВПУ
	Phases   []int
	InterTE  float32 //Интервал
	Values   []Value
	Last     Value
}
type Value struct {
	Time   time.Time
	Status int
	Value  int
}

func (c *Value) eq(cc *Value) bool {
	return c.Time.Equal(cc.Time) && c.Status == cc.Status && c.Value == cc.Value
}

// GetInfo возвращает былв информация за последние interval секунд по данному каналу
func (Ch *Chanels) GetInfo(TypeDev int, interval int) bool {
	if data.DataValue.Controller.Base {
		return false
	}
	for _, v := range Ch.Chanels {
		if v.TypeDev == TypeDev {
			return v.GetInfo(interval)
		}
	}
	return false
}
func (Ch *Chanels) GetLast(TypeDev int) Value {
	if data.DataValue.Controller.Base {
		return Value{Time: time.Unix(0, 0)}
	}
	for _, v := range Ch.Chanels {
		if v.TypeDev == TypeDev {
			return v.Last
		}
	}
	return Value{Time: time.Unix(0, 0)}
}

func (Ch *Chanels) AddInfo(tick OneTick) {
	if data.DataValue.Controller.Base {
		return
	}
	oc, is := Ch.Chanels[tick.Nomber]
	if !is {
		return
	}
	oc.AddInfo(tick)
}
func (Ch *Chanels) ClearInfo() {
	for _, v := range Ch.Chanels {
		v.ClearInfo()
	}
}
func (ch *OneChanel) AddInfo(tick OneTick) {
	ch.Last = tick.Value
	ch.Values = append(ch.Values, tick.Value)
	// logger.Debug.Printf("%v", tick)
}
func (ch *OneChanel) GetInfo(interval int) bool {
	if ch.Last.Status != 0 {
		return true
	}
	if time.Since(ch.Last.Time).Seconds() <= float64(interval) {
		return true
	}
	return false
}
func (ch *OneChanel) ClearInfo() {
	ch.Values = make([]Value, 0)
}
func newStatistics() Chanels {
	result := Chanels{Chanels: make(map[int]*OneChanel), Cron: cron.New()}
	if data.DataValue.Controller.Base {
		return result
	}
	result.Interval = data.DataValue.Arrays.StatDefine.Levels[0].Period
	for i := 0; i < len(data.DataValue.Arrays.PointSet.Points); i++ {
		ch := new(OneChanel)
		ch.Number = data.DataValue.Arrays.PointSet.Points[i].NumPoint
		ch.Last = Value{Time: time.Unix(0, 0), Status: 0, Value: 0}
		ch.Phases = make([]int, 0)
		ch.Values = make([]Value, 0)
		ch.TypeStat = data.DataValue.Arrays.PointSet.Points[i].TypeSt
		ch.TypeDev = data.DataValue.Arrays.SetTimeUse.Uses[i].Tvps
		ch.Phases = make([]int, 0)
		ch.InterTE = data.DataValue.Arrays.SetTimeUse.Uses[i].Long
		ls := strings.Split(data.DataValue.Arrays.SetTimeUse.Uses[i].Fazes, ",")
		for _, v := range ls {
			phase, _ := strconv.Atoi(v)
			ch.Phases = append(ch.Phases, int(phase))
		}
		result.Chanels[ch.Number] = ch
	}
	// init := "0 "
	// for i := 1; i < 60; i++ {
	// 	if i%result.Interval == 0 {
	// 		init += strconv.Itoa(i) + " "
	// 	}
	// }
	init := "*/" + strconv.Itoa(result.Interval)
	result.Cron.AddFunc(init+" * * * *", func() {
		// logger.Info.Print("Cron")
		m := 0
		if time.Now().Hour() == 0 && time.Now().Minute() == 0 {
			m = (24 * 60) / statistics.Interval
		} else {
			m = (time.Now().Hour()*60 + time.Now().Minute()) / statistics.Interval
		}
		sendStat <- SendStat{Interval: statistics.Interval, Period: m, Hour: time.Now().Hour(), Minute: time.Now().Minute()}
	})
	return result
}
