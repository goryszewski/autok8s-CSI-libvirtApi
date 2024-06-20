package main

import (
	"flag"
	"fmt"

	"github.com/goryszewski/autok8s-CSI-libvirtApi/pkg/driver"
	libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

func main() {
	fmt.Println("Start 001")
	var (
		endpoint = flag.String("endpoint", "default", "Endpoint gRPC")
		role     = flag.String("role", "controler", "role")
	)
	flag.Parse()

	fmt.Println(*endpoint)

	test := "CSIv"
	URL := "http://10.17.3.1:8050"
	// DOTO FIX static

	conf := libvirtApiClient.Config{Username: &test, Password: &test, Url: &URL}

	drv, err := driver.NewDriver(driver.InputParam{Endpoint: *endpoint, Name: driver.Name}, conf)

	if err != nil {
		fmt.Println("Error load driver", err.Error())
	}

	if err := drv.Run(role); err != nil {
		fmt.Println(err)
	}
}
