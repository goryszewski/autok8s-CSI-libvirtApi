package iscsi

import "fmt"

type VolumeCreateRequest struct {
	Name         string
	SizeGigaByte int64
}

type VolumeDeleteRequest struct {
	Id string
}

type Volume struct {
	Id string
}

type StorageService struct{}

func NewIscsi() *StorageService {
	return &StorageService{}
}

func (d *StorageService) CreateVolume(*VolumeCreateRequest) (*Volume, error) {
	return &Volume{
		Id: "1",
	}, nil
}
func (d *StorageService) DeleteVolume(*VolumeDeleteRequest) error {
	return nil
}

func (d *StorageService) GetVolume(volumeId string) (*Volume, error) {
	return &Volume{}, nil
}

func (d *StorageService) Attach(volID string, nodeid string) {
	fmt.Printf("{DEBUG}{MOCK} [VOL_ID: %v] [nodeID: %v]", volID, nodeid)
}
