package driver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	_   = iota
	kiB = 1 << (10 * iota)
	miB
	giB
	tiB
)

const defaultVolumeSizeInBytes int64 = 16 * giB

func (d *Driver) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume - Expect name")
	}

	requiredBytes := req.CapacityRange.GetRequiredBytes()
	requiredSet := 0 < requiredBytes
	limitBytes := req.CapacityRange.GetLimitBytes()
	limitSet := 0 < limitBytes
	size := requiredBytes

	if !requiredSet && !limitSet {
		size = defaultVolumeSizeInBytes
	}

	if req.VolumeCapabilities == nil || len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume - Expect VolumeCapabilities")
	}

	d.log.Infof(" requiredBytes:%v | limitBytes:%v | size:%v | size big:%v", requiredBytes, limitBytes, size, int(size/giB))
	vol, err := d.storage.CreateDisk(int(size / giB))
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed provisioning volume")
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			CapacityBytes: size,
			VolumeId:      strconv.Itoa(vol.ID),
		},
	}, nil
}
func (d *Driver) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	id, _ := strconv.Atoi(req.VolumeId)

	err := d.storage.DeleteDisk(id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed Delete volume")
	}

	return &csi.DeleteVolumeResponse{}, nil
}
func (d *Driver) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "VolumeID is mandatory")
	}

	if req.NodeId == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeId is mandatory")
	}

	id, _ := strconv.Atoi(req.VolumeId)

	DiskBindInfo, err := d.storage.BindDisk(id, req.NodeId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Error Internal: %v", err))
	}

	if (DiskBindInfo.ID == "") || (DiskBindInfo.Address == "") {
		return nil, status.Error(codes.InvalidArgument, "Bad Api Response ")
	}

	return &csi.ControllerPublishVolumeResponse{
		PublishContext: map[string]string{
			"libvirtCSI": req.VolumeId,
			"ID":         DiskBindInfo.ID,
			"Path":       DiskBindInfo.Path,
			"Target":     DiskBindInfo.Target,
			"Address":    DiskBindInfo.Address,
		},
	}, nil
}
func (d *Driver) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "VolumeID is mandatory")
	}

	if req.NodeId == "" {
		return nil, status.Error(codes.InvalidArgument, "NodeId is mandatory")
	}

	id, _ := strconv.Atoi(req.VolumeId)

	err := d.storage.UnBindDisk(id, req.NodeId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Error Internal: %v", err))
	}

	return &csi.ControllerUnpublishVolumeResponse{}, nil
}
func (d *Driver) ValidateVolumeCapabilities(ctx context.Context, req *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	fmt.Printf("[DEBUG][ValidateVolumeCapabilities][*csi.ValidateVolumeCapabilitiesRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) ListVolumes(ctx context.Context, req *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	fmt.Printf("[DEBUG][ListVolumes][*csi.ListVolumesRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) GetCapacity(ctx context.Context, req *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	fmt.Printf("[DEBUG][GetCapacity][*csi.GetCapacityRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) ControllerGetCapabilities(ctx context.Context, req *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
	caps := []*csi.ControllerServiceCapability{}

	for _, c := range []csi.ControllerServiceCapability_RPC_Type{
		csi.ControllerServiceCapability_RPC_CREATE_DELETE_VOLUME,
		csi.ControllerServiceCapability_RPC_PUBLISH_UNPUBLISH_VOLUME,
	} {
		caps = append(caps, &csi.ControllerServiceCapability{
			Type: &csi.ControllerServiceCapability_Rpc{
				Rpc: &csi.ControllerServiceCapability_RPC{
					Type: c,
				},
			},
		})
	}

	return &csi.ControllerGetCapabilitiesResponse{
		Capabilities: caps,
	}, nil
}
func (d *Driver) CreateSnapshot(ctx context.Context, req *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	fmt.Printf("[DEBUG][CreateSnapshot][*csi.CreateSnapshotRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) DeleteSnapshot(ctx context.Context, req *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	fmt.Printf("[DEBUG][DeleteSnapshot][*csi.DeleteSnapshotRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) ListSnapshots(ctx context.Context, req *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	fmt.Printf("[DEBUG][ListSnapshots][*csi.ListSnapshotsRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) ControllerExpandVolume(ctx context.Context, req *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	fmt.Printf("[DEBUG][ControllerExpandVolume][*csi.ControllerExpandVolumeRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) ControllerGetVolume(ctx context.Context, req *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	fmt.Printf("[DEBUG][ControllerGetVolume][*csi.ControllerGetVolumeRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) ControllerModifyVolume(ctx context.Context, req *csi.ControllerModifyVolumeRequest) (*csi.ControllerModifyVolumeResponse, error) {
	fmt.Printf("[DEBUG][ControllerModifyVolume][*csi.ControllerModifyVolumeRequest] %+v \n", req)
	return nil, nil
}
