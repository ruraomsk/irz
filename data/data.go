package data

import (
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ruraomsk/ag-server/binding"
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
type Server struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}

type Common struct {
	Controller pudge.Controller `json:"controller"`
	Arrays     binding.Arrays
	change     bool
	CommandDU  CommandDU
	Server     Server
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
	c.Controller.DK.RDK = 8
	c.Controller.DK.FDK = 1
	c.Controller.DK.DDK = 2
	c.Controller.DK.EDK = 0
	c.Controller.DK.TDK = 1
	c.Controller.DK.ODK = false
	c.Controller.DK.LDK = 0
	c.Controller.Model.VPCPDL = 15
	c.Controller.Model.VPCPDR = 5
	c.Controller.Model.VPBSL = 1
	c.Controller.Model.VPBSR = 1
	c.Controller.Model.DKA = true
	c.Controller.Status.TObmen = 5
	c.Controller.Traffic = pudge.Traffic{}
	c.Controller.Arrays = make([]pudge.ArrayPriv, 0)
	c.Controller.LogLines = make([]pudge.LogLine, 0)
	c.Arrays = *binding.NewArrays()
	c.Server.Host = setup.Set.Server.Host
	c.Server.Port = setup.Set.Server.Port
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
	if status == 0 {
		c.Controller.StatusCommandDU.IsPK = false
	} else {
		c.Controller.StatusCommandDU.IsPK = true
	}
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetNK(status int) {
	mutex.Lock()
	DataValue.CommandDU.NK = status
	if status == 0 {
		c.Controller.StatusCommandDU.IsNK = false
	} else {
		c.Controller.StatusCommandDU.IsNK = true
	}
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetCK(status int) {
	mutex.Lock()
	DataValue.CommandDU.CK = status
	if status == 0 {
		c.Controller.StatusCommandDU.IsCK = false
	} else {
		c.Controller.StatusCommandDU.IsCK = true
	}
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetDU(status int) {
	mutex.Lock()
	DataValue.CommandDU.DU = status
	if status == 9 {
		c.Controller.StatusCommandDU.IsDUDK1 = false
	} else {
		c.Controller.StatusCommandDU.IsDUDK1 = true
	}

	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetLastOperation() {
	mutex.Lock()
	c.Controller.LastOperation = time.Now()
	mutex.Unlock()
}
func (c *Common) SetConnectTime() {
	mutex.Lock()
	c.Controller.ConnectTime = time.Now()
	mutex.Unlock()
}
func (c *Common) SetNowTime() {
	mutex.Lock()
	c.Controller.TimeDevice = time.Now()
	mutex.Unlock()
}

func (c *Common) SetBase(status bool) {
	mutex.Lock()
	DataValue.Controller.Base = status
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetTechMode(status int) {
	mutex.Lock()
	DataValue.Controller.TechMode = status
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetArrays(arrays binding.Arrays) {
	mutex.Lock()
	DataValue.Arrays = arrays
	DataValue.change = true
	mutex.Unlock()
}
func (c *Common) SetDK(dk pudge.DK) {
	mutex.Lock()
	DataValue.Controller.DK = dk
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

	DataValue.SetNK(0)
	DataValue.SetPK(0)
	DataValue.SetCK(0)
	DataValue.SetDU(9)
	DataValue.SetTechMode(1)
	DataValue.SetSFDK(false)
	DataValue.SetConnected(false)
	DataValue.Controller.GPS.Ok = true
	DataValue.Controller.Status.Ethernet = true

	initChans()

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
