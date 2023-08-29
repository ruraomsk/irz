package setup

var (
	Set    *Setup
	ExtSet *ExtSetup
)

type Setup struct {
	LogPath     string      `toml:"logpath"`
	ID          int         `toml:"id"`
	Immitator   bool        `toml:"immitator"`
	Server      Server      `toml:"server"`
	Modbus      Modbus      `toml:"modbus"`
	SetupPudge  SetupPudge  `toml:"pudge"`
	VisioDevice VisioDevice `toml:"visiodevice"`
	ModbusRadar ModbusRadar `toml:"modbusradar"`
	TrafficData TrafficData `toml:"trafficdata"`
}
type ExtSetup struct {
	Server      Server      `toml:"server" json:"server"`
	Modbus      Modbus      `toml:"modbus" json:"modbus"`
	VisioDevice VisioDevice `toml:"visiodevice" json:"visiodevice"`
	ModbusRadar ModbusRadar `toml:"modbusradar" json:"modbusradar"`
	TrafficData TrafficData `toml:"trafficdata" json:"trafficdata"`
}
type Server struct {
	Host string `toml:"host" json:"host"`
	Port int    `toml:"port" json:"port"`
}
type Modbus struct {
	Device   string `toml:"device" json:"device"`
	BaudRate int    `toml:"baudrate" json:"baudrate"`
	Parity   string `toml:"parity" json:"parity"`
	UId      int    `toml:"uid" json:"uid"`
}
type ModbusRadar struct {
	Radar   bool   `toml:"radar" json:"radar"`
	Master  bool   `toml:"master" json:"master"`
	Debug   bool   `toml:"debug" json:"debug"`
	Host    string `toml:"host" json:"host"`
	Port    int    `toml:"port" json:"port"`
	ID      int    `toml:"id" json:"id"`
	Chanels int    `toml:"chanels" json:"chanels"`
}

type VisioDevice struct {
	Visio    bool   `toml:"visio"`
	Device   string `toml:"device" json:"device"`
	BaudRate int    `toml:"baudrate" json:"baudrate"`
	Parity   string `toml:"parity" json:"parity"`
}
type SetupPudge struct {
	DbPath string `toml:"dbpath"`
}
type TrafficData struct {
	Work    bool   `toml:"work" json:"work"`
	Debug   bool   `toml:"debug" json:"debug"`
	Host    string `toml:"host" json:"host"`
	Port    int    `toml:"port" json:"port"`
	Listen  int    `toml:"listen" json:"listen"`
	Chanels int    `toml:"chanels" json:"chanels"`
}

func (s *Setup) Update(es ExtSetup) {
	s.Server = es.Server
	s.Modbus = es.Modbus
	s.VisioDevice = es.VisioDevice
	s.ModbusRadar = es.ModbusRadar
	s.TrafficData = es.TrafficData
}
func (es *ExtSetup) Update(s Setup) {
	es.Server = s.Server
	es.Modbus = s.Modbus
	es.VisioDevice = s.VisioDevice
	es.ModbusRadar = s.ModbusRadar
	es.TrafficData = s.TrafficData
}
