package driver

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.PublishContext] %#+v \n", req.PublishContext)
	fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.VolumeId] %#+v \n", req.VolumeId)
	fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.StagingTargetPath] %#+v \n", req.StagingTargetPath)
	fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.VolumeCapability] %#+v \n", req.VolumeCapability)
	fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.VolumeContext] %#+v \n", req.VolumeContext)
	fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.VolumeCapability.AccessType] %#+v \n", req.VolumeCapability.AccessType)

	if req.VolumeId == "" {
		return nil, status.Error(codes.InvalidArgument, "Required VolumeID")
	}

	if req.StagingTargetPath == "" {
		return nil, status.Error(codes.InvalidArgument, "Required StagingTargetPath")
	}

	if req.VolumeCapability == nil {
		return nil, status.Error(codes.InvalidArgument, "Required VolumeCapability")
	}

	switch req.VolumeCapability.AccessType.(type) {
	case *csi.VolumeCapability_Block:
		return &csi.NodeStageVolumeResponse{}, nil
	}

	mnt := req.VolumeCapability.GetMount()

	var fsType string = "ext4"
	if mnt.FsType != "" {
		fsType = mnt.FsType
	}

	source := ""
	target := req.StagingTargetPath
	if address, ok := req.PublishContext["Address"]; !ok {
		return nil, status.Error(codes.InvalidArgument, "Required address in  req.PublishContext")
	} else {
		source = fmt.Sprintf("/dev/disk/by-path/pci-%s", address)
	}

	isf, err := isNotFormated(source)
	fmt.Printf("[DEBUG][NodeStageVolume][PRE - Formater] ERROR: %v \n", err)
	if isf {
		fmt.Printf("[DEBUG][NodeStageVolume][REQUIRED - Formater] \n")
		err := Formater(fsType, source)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Problem format: %v", err.Error()))
		}
	}

	err = mount(source, target, &mnt.MountFlags)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Problem Mount: %v", err.Error()))
	}

	return &csi.NodeStageVolumeResponse{}, nil
}
func (d *Driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	err := Umount(req.StagingTargetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Problem Mount: %v", err.Error()))
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}
func (d *Driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	fmt.Printf("[DEBUG][NodePublishVolume][StagingTargetPath] %#+v \n", req.StagingTargetPath)
	fmt.Printf("[DEBUG][NodePublishVolume][TargetPath       ] %#+v \n", req.TargetPath)
	fmt.Printf("[DEBUG][NodePublishVolume][Secrets          ] %#+v \n", req.Secrets)
	fmt.Printf("[DEBUG][NodePublishVolume][VolumeContext    ] %#+v \n", req.VolumeContext)
	fmt.Printf("[DEBUG][NodePublishVolume][Readonly         ] %#+v \n", req.Readonly)
	fmt.Printf("[DEBUG][NodePublishVolume][VolumeCapability ] %#+v \n", req.VolumeCapability)
	fmt.Printf("[DEBUG][NodePublishVolume][PublishContext   ] %#+v \n", req.PublishContext)
	fmt.Printf("[DEBUG][NodePublishVolume][VolumeId         ] %#+v \n", req.VolumeId)

	err := mount(req.StagingTargetPath, req.TargetPath, &[]string{"--bind"})
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Problem Mount: %v", err.Error()))
	}
	return &csi.NodePublishVolumeResponse{}, nil
}
func (d *Driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	err := Umount(req.TargetPath)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Problem Umount: %v", err.Error()))
	}

	return &csi.NodeUnpublishVolumeResponse{}, nil
}
func (d *Driver) NodeGetVolumeStats(ctx context.Context, req *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	fmt.Printf("[DEBUG][NodeGetVolumeStats][*csi.NodeGetVolumeStatsRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) NodeExpandVolume(ctx context.Context, req *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	fmt.Printf("[DEBUG][NodeExpandVolume][*csi.NodeExpandVolumeRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) NodeGetCapabilities(ctx context.Context, req *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
	fmt.Printf("[DEBUG][NodeGetCapabilities][*csi.NodeGetCapabilitiesRequest] %#+v \n", req)

	return &csi.NodeGetCapabilitiesResponse{
		Capabilities: []*csi.NodeServiceCapability{
			&csi.NodeServiceCapability{
				Type: &csi.NodeServiceCapability_Rpc{
					Rpc: &csi.NodeServiceCapability_RPC{
						Type: csi.NodeServiceCapability_RPC_STAGE_UNSTAGE_VOLUME,
					},
				},
			},
		},
	}, nil
}
func (d *Driver) NodeGetInfo(ctx context.Context, req *csi.NodeGetInfoRequest) (*csi.NodeGetInfoResponse, error) {
	fmt.Printf("[DEBUG][NodeGetInfo][*csi.NodeGetInfoRequest] %+v \n", req)
	nodeID, err := os.ReadFile("/id")
	if err != nil {
		fmt.Printf("[DEBUG][NodeGetInfo][os.ReadFile(/id)] %+v \n", err)
	}
	id := strings.ReplaceAll(string(nodeID), "\\n", "")
	id = strings.ReplaceAll(id, "\n", "")
	fmt.Printf("[DEBUG][NodeGetInfo][nodeID] %v \n", id)

	// DOTO endpoint metadata

	return &csi.NodeGetInfoResponse{
		NodeId:            id + ".autok8s.xyz",
		MaxVolumesPerNode: 5,
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				"region": "local",
			},
		},
	}, nil
}
