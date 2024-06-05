package iscsi

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
