package driver

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func (d *Driver) NodeStageVolume(ctx context.Context, req *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	fmt.Printf("[DEBUG][NodeStageVolume][*csi.NodeStageVolumeRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) NodeUnstageVolume(ctx context.Context, req *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	fmt.Printf("[DEBUG][NodeUnstageVolume][*csi.NodeUnstageVolumeRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) NodePublishVolume(ctx context.Context, req *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	fmt.Printf("[DEBUG][NodePublishVolume][*csi.NodePublishVolumeRequest] %+v \n", req)
	return nil, nil
}
func (d *Driver) NodeUnpublishVolume(ctx context.Context, req *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	fmt.Printf("[DEBUG][NodeUnpublishVolume][*csi.NodeUnpublishVolumeRequest] %+v \n", req)
	return nil, nil
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
	fmt.Printf("[DEBUG][NodeGetCapabilities][*csi.NodeGetCapabilitiesRequest] %+v \n", req)
	return nil, nil
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
