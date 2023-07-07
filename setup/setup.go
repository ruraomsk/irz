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
	Visio       bool        `toml:"visio"`
	VisioDevice VisioDevice `toml:"visiodevice"`
}
type ExtSetup struct {
	Server      Server      `toml:"server" json:"server"`
	Modbus      Modbus      `toml:"modbus" json:"modbus"`
	Visio       bool        `toml:"visio" json:"visio"`
	VisioDevice VisioDevice `toml:"visiodevice" json:"visiodevice"`
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
type VisioDevice struct {
	Device   string `toml:"device" json:"device"`
	BaudRate int    `toml:"baudrate" json:"baudrate"`
	Parity   string `toml:"parity" json:"parity"`
}
type SetupPudge struct {
	DbPath string `toml:"dbpath"`
}
