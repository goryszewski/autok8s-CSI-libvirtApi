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
		url        = flag.String("url", "", "")
		user       = flag.String("user", "", "")
		pass       = flag.String("pass", "", "")
	)
	flag.Parse()

	if *configFile == "" && *user == "" && *pass == "" && *url == "" {
		log.Fatal("Config or variables(url,user,pass) Required ")
	}

	var conf libvirtApiClient.Config

	if *configFile != "" {

		data, err := os.ReadFile(*configFile)
		if err != nil {
			log.Fatal("Error read file")
		}

		err = json.Unmarshal(data, &conf)
		if err != nil {
			log.Fatal("error parse config file")
		}
	} else {
		conf = libvirtApiClient.Config{Username: user, Password: pass, Url: url}
	}

	drv, err := driver.NewDriver(driver.InputParam{Endpoint: *endpoint, Name: driver.Name}, conf)

	if err != nil {
		log.Println("Error load driver", err.Error())
	}

	if err := drv.Run(role); err != nil {
		log.Println(err)
	}
}
