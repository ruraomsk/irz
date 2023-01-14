package setup

var (
	Set *Setup
)

type Setup struct {
	LogPath    string     `toml:"logpath"`
	ID         int        `toml:"id"`
	Server     Server     `toml:"server"`
	Modbus     Modbus     `toml:"modbus"`
	SetupPudge SetupPudge `toml:"pudge"`
	WatchDog   WatchDog   `toml:"watchdog"`
	Hardware   Hardware   `toml:"hardware"`
	Netware    Netware    `toml:"netware"`
	Vpu        Vpu        `toml:"vpu"`
	Counter    Counter    `toml:"counter"`
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
type WatchDog struct {
	Step int `toml:"step"`
}
type Hardware struct {
	Step    int    `toml:"step"`
	Connect string `toml:"connect"`
	SPort   int    `toml:"sport"`
	C8count int    `toml:"count"`
	LongKK  int    `toml:"longkk"`
	PinOS   int    `toml:"pinos"`
	PinYB   int    `toml:"pinyb"`
}
type Netware struct {
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Password string `toml:"password"`
}
type Vpu struct {
	Step    int    `toml:"step"`
	Connect string `toml:"connect"`
	SPort   int    `toml:"sport"`
}
type Counter struct {
	Step    int    `toml:"step"`
	Connect string `toml:"connect"`
	SPort   int    `toml:"sport"`
}

func init() {
}
