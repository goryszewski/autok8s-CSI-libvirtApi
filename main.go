package main

import (
	"encoding/json"
	"flag"
	"os"

	"log"

	"github.com/goryszewski/autok8s-CSI-libvirtApi/pkg/driver"
	"github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

func main() {
	var (
		endpoint   = flag.String("endpoint", "default", "Endpoint gRPC")
		role       = flag.String("role", "controler", "role")
		configFile = flag.String("config", "", "config")
	)
	flag.Parse()

	if *configFile == "" {
		log.Fatal("Config Required")
	}

	var conf libvirtApiClient.Config

	data, err := os.ReadFile(*configFile)

	if err != nil {
		log.Fatal("Error read file")
	}

	err = json.Unmarshal(data, &conf)
	if err != nil {
		log.Fatal("error parse config file")
	}

	drv, err := driver.NewDriver(driver.InputParam{Endpoint: *endpoint, Name: driver.Name}, conf)

	if err != nil {
		log.Println("Error load driver", err.Error())
	}

	if err := drv.Run(role); err != nil {
		log.Println(err)
	}
}
