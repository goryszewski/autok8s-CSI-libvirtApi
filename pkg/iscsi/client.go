package iscsi

type VolumeCreateRequest struct {
	Name         string
	SizeGigaByte int64
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
