package driver

import (
	"context"
	"fmt"
	"strconv"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *Driver) CreateVolume(ctx context.Context, req *csi.CreateVolumeRequest) (*csi.CreateVolumeResponse, error) {
	fmt.Printf("[DEBUG][CreateVolume] called")
	// fmt.Printf("req: [%+v] \n", *req)
	// req: [{Name:pvc-5c8397b5-1e68-434a-8ab0-635dad53ce73 CapacityRange:required_bytes:1073741824  VolumeCapabilities:[mount:<> access_mode:<mode:SINGLE_NODE_WRITER > ] Parameters:map[] Secrets:map[] VolumeContentSource:<nil> AccessibilityRequirements:<nil> MutableParameters:map[] XXX_NoUnkeyedLiteral:{} XXX_unrecognized:[] XXX_sizecache:0}]
	if req.Name == "" {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume - Expect name")
	}

	var sizeByte int64 = req.CapacityRange.GetLimitBytes()
	fmt.Printf("[DEBUG][CreateVolume][req.CapacityRange.GetLimitBytes()] %+v \n", sizeByte)
	if req.VolumeCapabilities == nil || len(req.VolumeCapabilities) == 0 {
		return nil, status.Error(codes.InvalidArgument, "CreateVolume - Expect VolumeCapabilities")
	}

	vol, err := d.storage.CreateDisk(8)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed provisioning volume")
	}
	fmt.Printf("[DEBUG][CreateVolume][CreateDisk] %+v \n", vol)

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			CapacityBytes: sizeByte,
			VolumeId:      strconv.Itoa(vol.ID),
		},
	}, nil
}
func (d *Driver) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	fmt.Printf("[DEBUG][DeleteVolume][*csi.DeleteVolumeRequest] %+v \n", req)

	id, _ := strconv.Atoi(req.VolumeId)
	err := d.storage.DeleteDisk(id)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed Delete volume")
	}

	return &csi.DeleteVolumeResponse{}, nil
}
func (d *Driver) ControllerPublishVolume(ctx context.Context, req *csi.ControllerPublishVolumeRequest) (*csi.ControllerPublishVolumeResponse, error) {
	fmt.Printf("[DEBUG][ControllerPublishVolume][*csi.ControllerPublishVolumeRequest] %+v \n", req)
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
	fmt.Printf("[DEBUG][ControllerPublishVolume][DiskBindInfo] %+v \n", DiskBindInfo)
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
	fmt.Printf("[DEBUG][ControllerUnpublishVolume][*csi.ControllerUnpublishVolumeRequest] %+v \n", req)
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
	fmt.Printf("[DEBUG][ControllerGetCapabilities][*csi.ControllerGetCapabilitiesRequest] %+v \n", req)
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
