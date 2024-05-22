package driver

import (
	"fmt"
	"net"
	"net/url"
	"path"
	"path/filepath"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc"
)

const Name string = "Libvirtapi"

type Driver struct {
	name     string
	endpoint string

	srv *grpc.Server
}

type InputParam struct {
	Name     string
	Endpoint string
}

func NewDriver(params InputParam) *Driver {
	return &Driver{
		name:     params.Name,
		endpoint: params.Endpoint,
	}
}

func (d *Driver) Run() error {
	url, err := url.Parse(d.endpoint)
	if err != nil {
		return fmt.Errorf("problem parsing : %s", err.Error())
	}

	if url.Scheme != "unix" {
		return fmt.Errorf("shema is not unix : %s", url.Scheme)
	}

	fmt.Printf("DEBUG: url.Host: %s\n", url.Host)
	fmt.Printf("DEBUG: url.Path: %s\n", url.Path)
	grpcAddress := path.Join(url.Host, filepath.FromSlash(url.Path))
	fmt.Printf("DEBUG: grpcAddress: %s\n", grpcAddress)
	if url.Host == "" {
		grpcAddress = filepath.FromSlash(url.Path)
	}

	listener, e := net.Listen(url.Scheme, grpcAddress)
	if e != nil {
		return fmt.Errorf("problem with net.Listen: %s", e)
	}
	fmt.Println(listener)
	d.srv = grpc.NewServer()

	csi.RegisterNodeServer(d.srv, d)

	return nil
}
