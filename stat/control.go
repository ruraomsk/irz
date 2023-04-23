package stat

import (
	"time"

	"github.com/ruraomsk/irz/data"
)

var ReqCtrl chan interface{}
var ReplCtrl chan ReplayCtrl

// Контроль входов
type ctrl struct {
	code   int
	len    float64
	good   bool
	values []Value
	last   Value
}
type ReplayCtrl struct {
	TVP1 bool
	TVP2 bool
	MGR  bool
}

func (c *ctrl) clearOld() {
	values := make([]Value, 0)
	for _, v := range c.values {
		if time.Now().Sub(v.Time).Seconds() < c.len {
			values = append(values, v)
		}
	}
	c.values = values
}
func (c *ctrl) Update() {
	if c.len == 0 {
		c.good = true
		return
	}
	getInfo <- c.code
	last := <-sendInfo
	if !last.eq(&c.last) {
		c.last = last
		c.values = append(c.values, last)
	}
}
func (c *ctrl) SetGood() {
	if c.len == 0 {
		c.good = true
		return
	}
	if time.Now().Sub(c.last.Time).Seconds() > c.len {
		// logger.Debug.Printf("Обрыв либо долго не жмут")
		c.good = false
		return
	}
	for _, v := range c.values {
		if v.Value == 0 {
			c.good = true
			return
		}
	}
	// logger.Debug.Printf("Зажата кнопка")
	c.good = false
}
func (c *ctrl) setTime(len int) {
	c.len = float64(len * 60)
}

var tvp1 ctrl
var tvp2 ctrl
var mgr ctrl

func CtrlInputs() {
	defaultCtrl()
	ReplCtrl = make(chan ReplayCtrl)
	ReqCtrl = make(chan interface{})
	var replay ReplayCtrl
	if !data.DataValue.Controller.Base {
		setCtrl()
	}
	tik := time.NewTicker(time.Second)
	for {
		select {
		case <-ReqCtrl:

			replay.TVP1 = tvp1.good
			replay.TVP2 = tvp2.good
			replay.MGR = mgr.good
			ReplCtrl <- replay

		case <-tik.C:
			tvp1.clearOld()
			tvp2.clearOld()
			mgr.clearOld()
			tvp1.Update()
			tvp2.Update()
			mgr.Update()
			tvp1.SetGood()
			tvp2.SetGood()
			mgr.SetGood()
			setCtrl()
		case <-ChangeCtrl:
			defaultCtrl()
			if !data.DataValue.Controller.Base {
				setCtrl()
			}
		}
	}
}
func defaultCtrl() {
	tvp1 = ctrl{code: 1, good: true, values: make([]Value, 0), last: Value{Time: time.Unix(0, 0)}, len: 0}
	tvp2 = ctrl{code: 2, good: true, values: make([]Value, 0), last: Value{Time: time.Unix(0, 0)}, len: 0}
	mgr = ctrl{code: 3, good: true, values: make([]Value, 0), last: Value{Time: time.Unix(0, 0)}, len: 0}
}
func setCtrl() {
	if len(data.DataValue.Arrays.SetCtrl.Stage) == 0 {
		return
	}
	t := time.Now().Hour()*60 + time.Now().Minute()
	for _, v := range data.DataValue.Arrays.SetCtrl.Stage {
		l := v.Start.Hour*60 + v.Start.Minute
		r := v.End.Hour*60 + v.End.Minute
		if t <= r && t >= l {
			tvp1.setTime(v.TVPLen)
			tvp2.setTime(v.TVPLen)
			mgr.setTime(v.MGRLen)
			return
		}
	}
	tvp1.setTime(0)
	tvp2.setTime(0)
	mgr.setTime(0)
}
