package data

import (
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/ag-server/pudge"
	"github.com/ruraomsk/irz/setup"
)

var mutex sync.Mutex
var DataValue Common
var pathCommon string

type CommandDU struct {
	PK int
	CK int
	NK int
	DU int
}

func (d *CommandDU) setEmpty() {
	d.PK = 0
	d.CK = 0
	d.NK = 0
	d.DU = 9
}

type Common struct {
	Controller pudge.Controller `json:"controller"`
	change     bool
	CommandDU  CommandDU
}

func (c *Common) setEmpty() {
	c.Controller.ID = setup.Set.ID
	c.Controller.Base = true
	c.Controller.Name = "Нет перекрестка"
	c.Controller.NK = 1
	c.Controller.PK = 1
	c.Controller.CK = 1
	c.Controller.LastOperation = time.Now()
	c.Controller.ConnectTime = time.Unix(0, 0)
	c.Controller.TechMode = 1
	c.Controller.DK.TDK = 1
	c.Controller.Model.VPCPDL = 15
	c.Controller.Model.VPCPDR = 5
	c.Controller.Model.VPBSL = 1
	c.Controller.Model.VPBSR = 1
	c.Controller.Status.TObmen = 5
	c.Controller.Traffic = pudge.Traffic{}
	c.Controller.Arrays = make([]pudge.ArrayPriv, 0)
	c.Controller.LogLines = make([]pudge.LogLine, 0)
	c.CommandDU.setEmpty()
}
func (c *Common) Save() error {
	mutex.Lock()
	defer mutex.Unlock()

	file, err := json.Marshal(DataValue)
	if err != nil {
		logger.Error.Printf("При сборке для записи файла %s %s", pathCommon, err.Error())
		return err
	}
	err = os.WriteFile(pathCommon, file, 0777)
	if err != nil {
		logger.Error.Printf("При записи файла %s %s", pathCommon, err.Error())
		return err
	}
	return nil
}
func (c *Common) SetConnected(status bool) {
	mutex.Lock()
	c.Controller.StatusConnection = status
	mutex.Unlock()
}
func (c *Common) SetSFDK(status bool) {
	mutex.Lock()
	DataValue.Controller.StatusCommandDU.IsReqSFDK1 = status
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetPK(status int) {
	mutex.Lock()
	DataValue.CommandDU.PK = status
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetNK(status int) {
	mutex.Lock()
	DataValue.CommandDU.NK = status
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetCK(status int) {
	mutex.Lock()
	DataValue.CommandDU.CK = status
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetDU(status int) {
	mutex.Lock()
	DataValue.CommandDU.DU = status
	DataValue.change = true
	mutex.Unlock()
}

func LoadAll() {
	pathCommon = setup.Set.SetupPudge.DbPath + "common.json"
	file, err := os.ReadFile(pathCommon)
	if err != nil {
		logger.Info.Printf("Отсутствует настроечный файл %s", pathCommon)
		DataValue.setEmpty()
	} else {
		err = json.Unmarshal(file, &DataValue)
		if err != nil {
			logger.Error.Printf("При разборе файла %s %s", pathCommon, err.Error())
			DataValue.setEmpty()
		}
	}
	file, err = json.Marshal(DataValue)
	if err != nil {
		logger.Error.Printf("При сборке для записи файла %s %s", pathCommon, err.Error())
		DataValue.setEmpty()
	}
	err = os.WriteFile(pathCommon, file, 0777)
	if err != nil {
		logger.Error.Printf("При записи файла %s %s", pathCommon, err.Error())
	}
	DataValue.CommandDU.setEmpty()
	go run()
}
func run() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	saver := time.NewTicker(5 * time.Second)
	for {
		select {
		case <-c:
			DataValue.Save()
			logger.Info.Print("Состояние устройства сохранено")
			logger.Info.Print("Работа Common завершена")
			return
		case <-saver.C:
			DataValue.Save()
			// logger.Info.Print("Состояние устройства сохранено")
		}
	}

}
