package driver

import (
	"fmt"
	"net"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	libvirtApiClient "github.com/goryszewski/libvirtApi-client/libvirtApiClient"
)

const Name string = "Libvirtapi"

type Driver struct {
	name     string
	endpoint string

	srv     *grpc.Server
	storage *libvirtApiClient.Client
	log     *logrus.Entry
	ready   bool
}

type InputParam struct {
	Name     string
	Endpoint string
}

func NewDriver(params InputParam, conf libvirtApiClient.Config) (*Driver, error) {

	id, err := GetIDNode()
	if err != nil {
		return nil, fmt.Errorf("problem with id: %v", err)
	}
	log := logrus.New().WithFields(logrus.Fields{
		"host_id": id,
	})

	clientStorage, _ := libvirtApiClient.NewClient(conf, &http.Client{Timeout: 10 * time.Second})
	return &Driver{
		name:     params.Name,
		endpoint: params.Endpoint,
		storage:  clientStorage,
		log:      log,
	}, nil
}

func (d *Driver) Run(role *string) error {
	url, err := url.Parse(d.endpoint)
	if err != nil {
		return fmt.Errorf("problem parsing : %s", err.Error())
	}

	if url.Scheme != "unix" {
		return fmt.Errorf("shema is not unix : %s", url.Scheme)
	}

	d.log.Infof("DEBUG: url.Host: %s\n", url.Host)
	d.log.Infof("DEBUG: url.Path: %s\n", url.Path)
	grpcAddress := path.Join(url.Host, filepath.FromSlash(url.Path))
	d.log.Infof("DEBUG: grpcAddress: %s\n", grpcAddress)

	if url.Host == "" {
		grpcAddress = filepath.FromSlash(url.Path)
	}

	if err = os.Remove(grpcAddress); err != nil && os.IsExist(err) {
		return fmt.Errorf("error remove sock: %s", err)
	}

	listener, e := net.Listen(url.Scheme, grpcAddress)
	if e != nil {
		return fmt.Errorf("problem with net.Listen: %s", e)
	}

	d.srv = grpc.NewServer()

	csi.RegisterIdentityServer(d.srv, d)
	if *role == "node" {
		csi.RegisterNodeServer(d.srv, d)
	} else {
		csi.RegisterControllerServer(d.srv, d)
	}
	d.ready = true
	return d.srv.Serve(listener) // blocking call
}
