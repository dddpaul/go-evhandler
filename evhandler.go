package main

import (
	"fmt"
	"github.com/dddpaul/golang-evdev/evdev"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/exec"
	"strings"
)

const (
	usage = "usage: evhandler [device name]"
)

func initConfig(fn string) {
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	viper.SetConfigName(fn)
	viper.AddConfigPath(pwd)
	viper.ReadInConfig()
	log.Printf("Configuration from %s was loaded\n", viper.ConfigFileUsed())
}

func worker(ch <-chan string) {
	for cmd := range ch {
		parts := strings.Fields(cmd)
		head := parts[0]
		parts = parts[1:len(parts)]
		out, err := exec.Command(head, parts...).CombinedOutput()
		if err != nil {
			log.Println(err)
		}
		log.Println(string(out))
	}
}

func main() {
	initConfig("evhandler")
	params := viper.GetStringMapString("params")
	actions := viper.GetStringMapString("actions")

	var devName string
	var dev *evdev.InputDevice
	var events []evdev.InputEvent
	var err error

	switch len(os.Args) {
	case 1:
		var ok bool
		if devName, ok = params["device"]; !ok {
			log.Fatalf("Specify input device please\n")
		}
	case 2:
		devName = os.Args[1]
	default:
		fmt.Printf(usage + "\n")
		os.Exit(1)
	}

	dev, err = evdev.Open(devName)
	if err != nil {
		log.Fatalf("unable to open input device: %s\n", devName)
	}

	info := fmt.Sprintf("bus 0x%04x, vendor 0x%04x, product 0x%04x, version 0x%04x",
		dev.Bustype, dev.Vendor, dev.Product, dev.Version)
	log.Printf("Device name: %s\n", dev.Name)
	log.Printf("Device info: %s\n", info)
	log.Printf("Listening for events ...\n")

	ch := make(chan string)
	go worker(ch)
	for {
		events, err = dev.Read()
		for _, ev := range events {
			var codeName string
			code := int(ev.Code)
			evType := int(ev.Type)
			if m, ok := evdev.ByEventType[evType]; ok {
				codeName = m[code]
			}
			if evType == evdev.EV_KEY && ev.Value == 1 {
				if cmd, ok := actions[codeName]; ok {
					log.Printf("%s was pressed, command = '%s'\n", codeName, cmd)
					ch <- cmd
				}
			}
		}
	}
}
