package main

import (
	"flag"
	"fmt"

	"github.com/goryszewski/autok8s-CSI-libvirtApi/pkg/driver"
)

func main() {

	var (
		endpoint = flag.String("endpoint", "default", "Endpoint gRPC")
	)
	flag.Parse()

	fmt.Println(*endpoint)

	drv, err := driver.NewDriver(driver.InputParam{Endpoint: *endpoint, Name: driver.Name})

	if err != nil {
		fmt.Println("Error load driver", err.Error())
	}

	if err := drv.Run(); err != nil {
		fmt.Println(err)
	}
}
