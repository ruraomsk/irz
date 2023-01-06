package main

import (
	"embed"
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/setup"
)

var (
	//go:embed config
	config embed.FS
)

func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFS(config, "config/config.toml", &setup.Set); err != nil {
		fmt.Println("Dissmis config.toml")
		os.Exit(-1)
		return
	}

	os.MkdirAll(setup.Set.LogPath, 0777)
	os.MkdirAll(setup.Set.SetupBrams.DbPath, 0777)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := logger.Init(setup.Set.LogPath); err != nil {
		log.Panic("Error logger system", err.Error())
		return
	}
	fmt.Println("iRZ start")
	logger.Info.Println("iRZ start")
	fmt.Println("iRZ stop")
	logger.Info.Println("iRZ stop")

}
