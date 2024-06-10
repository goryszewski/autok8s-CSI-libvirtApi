package driver

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/container-storage-interface/spec/lib/go/csi"
)

func (d *Driver) NodeStageVolume(context.Context, *csi.NodeStageVolumeRequest) (*csi.NodeStageVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) NodeUnstageVolume(context.Context, *csi.NodeUnstageVolumeRequest) (*csi.NodeUnstageVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) NodePublishVolume(context.Context, *csi.NodePublishVolumeRequest) (*csi.NodePublishVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) NodeUnpublishVolume(context.Context, *csi.NodeUnpublishVolumeRequest) (*csi.NodeUnpublishVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) NodeGetVolumeStats(context.Context, *csi.NodeGetVolumeStatsRequest) (*csi.NodeGetVolumeStatsResponse, error) {
	return nil, nil
}
func (d *Driver) NodeExpandVolume(context.Context, *csi.NodeExpandVolumeRequest) (*csi.NodeExpandVolumeResponse, error) {
	return nil, nil
}
func (d *Driver) NodeGetCapabilities(context.Context, *csi.NodeGetCapabilitiesRequest) (*csi.NodeGetCapabilitiesResponse, error) {
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

	return &csi.NodeGetInfoResponse{
		NodeId:            id,
		MaxVolumesPerNode: 5,
		AccessibleTopology: &csi.Topology{
			Segments: map[string]string{
				"region": "local",
			},
		},
	}, nil
}
