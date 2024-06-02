package iscsi

type VolumeCreateRequest struct {
	Name         string
	SizeGigaByte int64
}

func NewIscsi() VolumeCreateRequest {
	return VolumeCreateRequest{}
}
