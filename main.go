package main

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/BurntSushi/toml"
	"github.com/anoshenko/rui"
	"github.com/ruraomsk/ag-server/logger"
	"github.com/ruraomsk/irz/comm"
	"github.com/ruraomsk/irz/data"
	"github.com/ruraomsk/irz/device"
	"github.com/ruraomsk/irz/kdm"
	"github.com/ruraomsk/irz/setup"
	"github.com/ruraomsk/irz/stat"
	"github.com/ruraomsk/irz/web"
	"github.com/ruraomsk/irz/worker"
)

func init() {
	setup.Set = new(setup.Setup)
	if _, err := toml.DecodeFS(resources, "config/base.toml", &setup.Set); err != nil {
		fmt.Println("Dissmis base.toml")
		os.Exit(-1)
		return
	}
	if _, err := os.Stat("config.toml"); err == nil {
		if _, err := toml.DecodeFile("config.toml", &setup.ExtSet); err != nil {
			fmt.Println("Dissmis config.toml")
			os.Exit(-1)
			return
		}
		setup.Set.Modbus = setup.ExtSet.Modbus
		setup.Set.Server = setup.ExtSet.Server
	}
	os.MkdirAll(setup.Set.LogPath, 0777)
	os.MkdirAll(setup.Set.SetupPudge.DbPath, 0777)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	if err := logger.Init(setup.Set.LogPath); err != nil {
		log.Panic("Error logger system", err.Error())
		return
	}
	fmt.Println("iRZ start")
	logger.Info.Println("iRZ start")
	data.LoadAll()
	go worker.Worker()
	// go kdm.KdmTest()
	go comm.ToServer()
	if setup.Set.Immitator {
		go device.Device()
	} else {
		go kdm.Kdm()
	}
	go stat.Statistics()

	// c := make(chan os.Signal, 1)
	// signal.Notify(c, os.Interrupt)
	// <-c
	// time.Sleep(5 * time.Second)
	// fmt.Println("iRZ stop")
	// logger.Info.Println("iRZ stop")
	rui.ProtocolInDebugLog = true
	/*
		rui.SetDebugLog(func(text string) {
			if len(text) > 120 {
				text = text[:120] + "..."
			}
			log.Println(text)
		})
	*/
	rui.AddEmbedResources(&resources)

	//addr := rui.GetLocalIP() + ":8080"
	addr := "localhost:8000"
	fmt.Print(addr)
	rui.OpenBrowser("http://" + addr)
	rui.StartApp(addr, web.CreateSession, rui.AppParams{
		Title:      "Ag-IRZ",
		Icon:       "icon.svg",
		TitleColor: rui.Color(0xffc0ded9),
	})
}
