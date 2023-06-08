package setup

var (
	Set    *Setup
	ExtSet *ExtSetup
)

type Setup struct {
	LogPath    string     `toml:"logpath"`
	ID         int        `toml:"id"`
	Immitator  bool       `toml:"immitator"`
	Server     Server     `toml:"server"`
	Modbus     Modbus     `toml:"modbus"`
	SetupPudge SetupPudge `toml:"pudge"`
}
type ExtSetup struct {
	Server Server `toml:"server"`
	Modbus Modbus `toml:"modbus"`
}
type Server struct {
	Host string `toml:"host"`
	Port int    `toml:"port"`
}
type Modbus struct {
	Device   string `toml:"device"`
	BaudRate int    `toml:"baudrate"`
	Parity   string `toml:"parity"`
	UId      int    `toml:"uid"`
}
type SetupPudge struct {
	DbPath string `toml:"dbpath"`
}
