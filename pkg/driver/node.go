package driver

import (
	"context"
	"fmt"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (d *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	// fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.PublishContext] %#+v \n", req.PublishContext)
	// fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.VolumeId] %#+v \n", req.VolumeId)
	// fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.StagingTargetPath] %#+v \n", req.StagingTargetPath)
	// fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.VolumeCapability] %#+v \n", req.VolumeCapability)
	// fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.VolumeContext] %#+v \n", req.VolumeContext)
	// fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest.VolumeCapability.AccessType] %#+v \n", req.VolumeCapability.AccessType)

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

	source_path := ""
	target := req.StagingTargetPath
	if address, ok := req.PublishContext["Address"]; !ok {
		return nil, status.Error(codes.InvalidArgument, "Required address in  req.PublishContext")
	} else {
		source_path = fmt.Sprintf("/dev/disk/by-path/pci-%s", address)
	}

	if req.VolumeContext["encrypt"] == "true" {
		not_encrypted, _ := isNotEncrypt(source_path)

		if not_encrypted {
			err := Encrypter(source_path)
			if err != nil {
				return nil, status.Error(codes.Internal, fmt.Sprintf("Problem Encrypter: %v", err.Error()))
			}
			if req.VolumeContext["tang"] != "false" {
				err := ClevisBind(source_path, req.VolumeContext["tang"])
				if err != nil {
					return nil, status.Error(codes.Internal, fmt.Sprintf("Problem clevis bind: %v", err.Error()))
				}
			}
		}

		err := OpenLuks(source_path, req.VolumeId)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Problem OpenLuks: %v", err.Error()))
		}

		source_path = fmt.Sprintf("/dev/mapper/%v", req.VolumeId)
	}

	not_formated, _ := isNotFormated(source_path)
	if not_formated {
		err := Formater(fsType, source_path)
		if err != nil {
			return nil, status.Error(codes.Internal, fmt.Sprintf("Problem format: %v", err.Error()))
		}
	}

	err := mount(source_path, target, &mnt.MountFlags)
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

	err = CloseLuks(req.VolumeId)
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("Problem CloseLuks: %v", err.Error()))
	}

	return &csi.NodeUnstageVolumeResponse{}, nil
}
func (d *Driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	// fmt.Printf("[DEBUG][NodePublishVolume][StagingTargetPath] %#+v \n", req.StagingTargetPath)
	// fmt.Printf("[DEBUG][NodePublishVolume][TargetPath       ] %#+v \n", req.TargetPath)
	// fmt.Printf("[DEBUG][NodePublishVolume][Secrets          ] %#+v \n", req.Secrets)
	// fmt.Printf("[DEBUG][NodePublishVolume][VolumeContext    ] %#+v \n", req.VolumeContext)
	// fmt.Printf("[DEBUG][NodePublishVolume][Readonly         ] %#+v \n", req.Readonly)
	// fmt.Printf("[DEBUG][NodePublishVolume][VolumeCapability ] %#+v \n", req.VolumeCapability)
	// fmt.Printf("[DEBUG][NodePublishVolume][PublishContext   ] %#+v \n", req.PublishContext)
	// fmt.Printf("[DEBUG][NodePublishVolume][VolumeId         ] %#+v \n", req.VolumeId)

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
	node, err := d.storage.GetNodeByMetadata()
	if err != nil {
		return nil, fmt.Errorf("problem with Metadata %+v \n", err)
	}

	return &csi.NodeGetInfoResponse{
		NodeId:            node.Name,
		MaxVolumesPerNode: 5,
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				"region": "local",
			},
		},
	}, nil
}
