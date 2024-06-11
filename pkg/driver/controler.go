package driver

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/goryszewski/autok8s-CSI-libvirtApi/pkg/iscsi"
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

	volReq := iscsi.VolumeCreateRequest{
		Name:         req.Name,
		SizeGigaByte: sizeByte / (1024 * 1024 * 1024),
	}
	fmt.Printf("[DEBUG][CreateVolume][VolumeCreateRequest] %+v \n", volReq)

	vol, err := d.storage.CreateVolume(&volReq)
	if err != nil {
		return nil, status.Error(codes.Internal, "Failed provisioning volume")
	}

	return &csi.CreateVolumeResponse{
		Volume: &csi.Volume{
			CapacityBytes: sizeByte,
			VolumeId:      vol.Id,
		},
	}, nil
}
func (d *Driver) DeleteVolume(ctx context.Context, req *csi.DeleteVolumeRequest) (*csi.DeleteVolumeResponse, error) {
	fmt.Printf("[DEBUG][DeleteVolume][*csi.DeleteVolumeRequest] %+v \n", req)

	volReq := iscsi.VolumeDeleteRequest{
		Id: req.VolumeId,
	}

	fmt.Printf("[DEBUG][DeleteVolume][*iscsi.VolumeDeleteRequest] %+v \n", volReq)

	err := d.storage.DeleteVolume(&volReq)
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

	_, err := d.storage.GetVolume(req.VolumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, "Error get Volume")
	}

	d.storage.Attach(req.VolumeId, req.NodeId)

	return nil, nil
}
func (d *Driver) ControllerUnpublishVolume(ctx context.Context, req *csi.ControllerUnpublishVolumeRequest) (*csi.ControllerUnpublishVolumeResponse, error) {
	fmt.Printf("[DEBUG][ControllerUnpublishVolume][*csi.ControllerUnpublishVolumeRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) ValidateVolumeCapabilities(context.Context, *csi.ValidateVolumeCapabilitiesRequest) (*csi.ValidateVolumeCapabilitiesResponse, error) {
	return nil, nil
}
func (d *Driver) ListVolumes(context.Context, *csi.ListVolumesRequest) (*csi.ListVolumesResponse, error) {
	return nil, nil
}
func (d *Driver) GetCapacity(context.Context, *csi.GetCapacityRequest) (*csi.GetCapacityResponse, error) {
	return nil, nil
}
func (d *Driver) ControllerGetCapabilities(context.Context, *csi.ControllerGetCapabilitiesRequest) (*csi.ControllerGetCapabilitiesResponse, error) {
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
func (d *Driver) CreateSnapshot(context.Context, *csi.CreateSnapshotRequest) (*csi.CreateSnapshotResponse, error) {
	return nil, nil
}
func (d *Driver) DeleteSnapshot(context.Context, *csi.DeleteSnapshotRequest) (*csi.DeleteSnapshotResponse, error) {
	return nil, nil
}
func (d *Driver) ListSnapshots(context.Context, *csi.ListSnapshotsRequest) (*csi.ListSnapshotsResponse, error) {
	return nil, nil
}
func (d *Driver) ControllerExpandVolume(context.Context, *csi.ControllerExpandVolumeRequest) (*csi.ControllerExpandVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) ControllerGetVolume(context.Context, *csi.ControllerGetVolumeRequest) (*csi.ControllerGetVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) ControllerModifyVolume(context.Context, *csi.ControllerModifyVolumeRequest) (*csi.ControllerModifyVolumeResponse, error) {
	return nil, nil
}
