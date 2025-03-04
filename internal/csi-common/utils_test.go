/*
Copyright 2019 ceph-csi authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package csicommon

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/container-storage-interface/spec/lib/go/csi"
	"github.com/csi-addons/spec/lib/go/replication"
	"github.com/stretchr/testify/require"
	mount "k8s.io/mount-utils"
)

var fakeID = "fake-id"

func TestGetReqID(t *testing.T) {
	t.Parallel()
	req := []interface{}{
		&csi.CreateVolumeRequest{
			Name: fakeID,
		},
		&csi.DeleteVolumeRequest{
			VolumeId: fakeID,
		},
		&csi.CreateSnapshotRequest{
			Name: fakeID,
		},
		&csi.DeleteSnapshotRequest{
			SnapshotId: fakeID,
		},

		&csi.ControllerExpandVolumeRequest{
			VolumeId: fakeID,
		},
		&csi.NodeStageVolumeRequest{
			VolumeId: fakeID,
		},
		&csi.NodeUnstageVolumeRequest{
			VolumeId: fakeID,
		},
		&csi.NodePublishVolumeRequest{
			VolumeId: fakeID,
		},
		&csi.NodeUnpublishVolumeRequest{
			VolumeId: fakeID,
		},
		&csi.NodeExpandVolumeRequest{
			VolumeId: fakeID,
		},

		&csi.CreateVolumeGroupSnapshotRequest{
			Name: fakeID,
		},
		&csi.DeleteVolumeGroupSnapshotRequest{
			GroupSnapshotId: fakeID,
		},
		&csi.GetVolumeGroupSnapshotRequest{
			GroupSnapshotId: fakeID,
		},

		&replication.EnableVolumeReplicationRequest{
			VolumeId: fakeID,
		},
		&replication.DisableVolumeReplicationRequest{
			VolumeId: fakeID,
		},
		&replication.PromoteVolumeRequest{
			VolumeId: fakeID,
		},
		&replication.DemoteVolumeRequest{
			VolumeId: fakeID,
		},
		&replication.ResyncVolumeRequest{
			VolumeId: fakeID,
		},
		&replication.GetVolumeReplicationInfoRequest{
			VolumeId: fakeID,
		},

		// volumeId is set in ReplicationSource
		&replication.EnableVolumeReplicationRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volume{
					Volume: &replication.ReplicationSource_VolumeSource{
						VolumeId: fakeID,
					},
				},
			},
		},
		&replication.DisableVolumeReplicationRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volume{
					Volume: &replication.ReplicationSource_VolumeSource{
						VolumeId: fakeID,
					},
				},
			},
		},
		&replication.PromoteVolumeRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volume{
					Volume: &replication.ReplicationSource_VolumeSource{
						VolumeId: fakeID,
					},
				},
			},
		},
		&replication.DemoteVolumeRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volume{
					Volume: &replication.ReplicationSource_VolumeSource{
						VolumeId: fakeID,
					},
				},
			},
		},
		&replication.ResyncVolumeRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volume{
					Volume: &replication.ReplicationSource_VolumeSource{
						VolumeId: fakeID,
					},
				},
			},
		},
		&replication.GetVolumeReplicationInfoRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volume{
					Volume: &replication.ReplicationSource_VolumeSource{
						VolumeId: fakeID,
					},
				},
			},
		},
		// volumeGroupId is set in ReplicationSource
		&replication.EnableVolumeReplicationRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volumegroup{
					Volumegroup: &replication.ReplicationSource_VolumeGroupSource{
						VolumeGroupId: fakeID,
					},
				},
			},
		},
		&replication.DisableVolumeReplicationRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volumegroup{
					Volumegroup: &replication.ReplicationSource_VolumeGroupSource{
						VolumeGroupId: fakeID,
					},
				},
			},
		},
		&replication.PromoteVolumeRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volumegroup{
					Volumegroup: &replication.ReplicationSource_VolumeGroupSource{
						VolumeGroupId: fakeID,
					},
				},
			},
		},
		&replication.DemoteVolumeRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volumegroup{
					Volumegroup: &replication.ReplicationSource_VolumeGroupSource{
						VolumeGroupId: fakeID,
					},
				},
			},
		},
		&replication.ResyncVolumeRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volumegroup{
					Volumegroup: &replication.ReplicationSource_VolumeGroupSource{
						VolumeGroupId: fakeID,
					},
				},
			},
		},
		&replication.GetVolumeReplicationInfoRequest{
			ReplicationSource: &replication.ReplicationSource{
				Type: &replication.ReplicationSource_Volumegroup{
					Volumegroup: &replication.ReplicationSource_VolumeGroupSource{
						VolumeGroupId: fakeID,
					},
				},
			},
		},
	}
	for _, r := range req {
		if got := getReqID(r); got != fakeID {
			t.Errorf("getReqID() = %v, want %v", got, fakeID)
		}
	}

	// test for nil request
	if got := getReqID(nil); got != "" {
		t.Errorf("getReqID() = %v, want empty string", got)
	}

	// test when both volume and group id not set
	if got := getReqID(&replication.EnableVolumeReplicationRequest{}); got != "" {
		t.Errorf("getReqID() = %v, want empty string", got)
	}
}

func TestFilesystemNodeGetVolumeStats(t *testing.T) {
	t.Parallel()

	// ideally this is tested on different filesystems, including CephFS,
	// but we'll settle for the filesystem where the project is checked out
	cwd, err := os.Getwd()
	require.NoError(t, err)

	// retry until a mountpoint is found
	for {
		stats, err := FilesystemNodeGetVolumeStats(context.TODO(), mount.New(""), cwd, true)
		if err != nil && cwd != "/" && strings.HasSuffix(err.Error(), "is not mounted") {
			// try again with the parent directory
			cwd = filepath.Dir(cwd)

			continue
		}

		require.NoError(t, err)
		require.NotEmpty(t, stats.GetUsage())
		for _, usage := range stats.GetUsage() {
			require.NotEqual(t, -1, usage.GetAvailable())
			require.NotEqual(t, -1, usage.GetTotal())
			require.NotEqual(t, -1, usage.GetUsed())
			require.NotEqual(t, 0, usage.GetUnit())
		}

		// tests done, no need to retry again
		break
	}
}

func TestRequirePositive(t *testing.T) {
	t.Parallel()

	require.Equal(t, int64(0), requirePositive(0))
	require.Equal(t, int64(0), requirePositive(-1))
	require.Equal(t, int64(1), requirePositive(1))
}

func TestIsBlockMultiNode(t *testing.T) {
	t.Parallel()

	blockCap := &csi.VolumeCapability{
		AccessType: &csi.VolumeCapability_Block{
			Block: &csi.VolumeCapability_BlockVolume{},
		},
	}

	fsCap := &csi.VolumeCapability{
		AccessType: &csi.VolumeCapability_Mount{
			Mount: &csi.VolumeCapability_MountVolume{},
		},
	}

	multiNodeCap := &csi.VolumeCapability{
		AccessMode: &csi.VolumeCapability_AccessMode{
			Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
		},
	}

	singleNodeCap := &csi.VolumeCapability{
		AccessMode: &csi.VolumeCapability_AccessMode{
			Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_MULTI_WRITER,
		},
	}

	tests := []struct {
		name        string
		caps        []*csi.VolumeCapability
		isBlock     bool
		isMultiNode bool
	}{{
		name:        "block/multi-node",
		caps:        []*csi.VolumeCapability{blockCap, multiNodeCap},
		isBlock:     true,
		isMultiNode: true,
	}, {
		name:        "block/single-node",
		caps:        []*csi.VolumeCapability{blockCap, singleNodeCap},
		isBlock:     true,
		isMultiNode: false,
	}, {
		name:        "filesystem/multi-node",
		caps:        []*csi.VolumeCapability{fsCap, multiNodeCap},
		isBlock:     false,
		isMultiNode: true,
	}, {
		name:        "filesystem/single-node",
		caps:        []*csi.VolumeCapability{fsCap, singleNodeCap},
		isBlock:     false,
		isMultiNode: false,
	}}

	for _, test := range tests {
		isBlock, isMultiNode := IsBlockMultiNode(test.caps)
		require.Equal(t, isBlock, test.isBlock, test.name)
		require.Equal(t, isMultiNode, test.isMultiNode, test.name)
	}
}

func TestIsFileRWO(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name    string
		caps    []*csi.VolumeCapability
		rwoFile bool
	}{
		{
			name: "non valid",
			caps: []*csi.VolumeCapability{
				{
					AccessMode: nil,
					AccessType: nil,
				},
			},
			rwoFile: false,
		},

		{
			name: "single writer FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_SINGLE_WRITER,
					},
				},
			},
			rwoFile: true,
		},
		{
			name: "single node multi writer FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_MULTI_WRITER,
					},
				},
			},
			rwoFile: true,
		},
		{
			name: "multi node multi writer FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
					},
				},
			},
			rwoFile: false,
		},
		{
			name: "multi node multi reader FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
					},
				},
			},
			rwoFile: false,
		},
		{
			name: "single node reader FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
					},
				},
			},
			rwoFile: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			rwoFile := IsFileRWO(tt.caps)
			if rwoFile != tt.rwoFile {
				t.Errorf("IsFileRWO() rwofile = %v, want %v", rwoFile, tt.rwoFile)
			}
		})
	}
}

func TestIsBlockMultiWriter(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name        string
		caps        []*csi.VolumeCapability
		multiWriter bool
		block       bool
	}{
		{
			name: "non valid",
			caps: []*csi.VolumeCapability{
				{
					AccessMode: nil,
					AccessType: nil,
				},
			},
			multiWriter: false,
			block:       false,
		},
		{
			name: "multi node multi writer block mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Block{
						Block: &csi.VolumeCapability_BlockVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
					},
				},
			},
			multiWriter: true,
			block:       true,
		},
		{
			name: "single node multi writer block mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Block{
						Block: &csi.VolumeCapability_BlockVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_MULTI_WRITER,
					},
				},
			},
			multiWriter: true,
			block:       true,
		},
		{
			name: "single writer block mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Block{
						Block: &csi.VolumeCapability_BlockVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_SINGLE_WRITER,
					},
				},
			},
			multiWriter: false,
			block:       true,
		},
		{
			name: "single writer FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_SINGLE_WRITER,
					},
				},
			},
			multiWriter: false,
			block:       false,
		},
		{
			name: "single node multi writer FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_MULTI_WRITER,
					},
				},
			},
			multiWriter: true,
			block:       false,
		},
		{
			name: "multi node multi writer FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
					},
				},
			},
			multiWriter: true,
			block:       false,
		},
		{
			name: "multi node multi reader FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
					},
				},
			},
			multiWriter: false,
			block:       false,
		},
		{
			name: "single node reader FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
					},
				},
			},
			multiWriter: false,
			block:       false,
		},
		{
			name: "multi node reader block mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Block{
						Block: &csi.VolumeCapability_BlockVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
					},
				},
			},
			multiWriter: false,
			block:       true,
		},
		{
			name: "single node reader block mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Block{
						Block: &csi.VolumeCapability_BlockVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
					},
				},
			},
			multiWriter: false,
			block:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			multiWriter, block := IsBlockMultiWriter(tt.caps)
			if multiWriter != tt.multiWriter {
				t.Errorf("IsBlockMultiWriter() multiWriter = %v, want %v", multiWriter, tt.multiWriter)
			}
			if block != tt.block {
				t.Errorf("IsBlockMultiWriter block = %v, want %v", block, tt.block)
			}
		})
	}
}

func TestIsReaderOnly(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name   string
		caps   []*csi.VolumeCapability
		roOnly bool
	}{
		{
			name: "non valid",
			caps: []*csi.VolumeCapability{
				{
					AccessMode: nil,
					AccessType: nil,
				},
			},
			roOnly: false,
		},

		{
			name: "single writer FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_SINGLE_WRITER,
					},
				},
			},
			roOnly: false,
		},
		{
			name: "single node multi writer FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_MULTI_WRITER,
					},
				},
			},
			roOnly: false,
		},
		{
			name: "multi node multi writer FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_MULTI_WRITER,
					},
				},
			},
			roOnly: false,
		},
		{
			name: "multi node multi reader FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
					},
				},
			},
			roOnly: true,
		},
		{
			name: "single node reader FS mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Mount{
						Mount: &csi.VolumeCapability_MountVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
					},
				},
			},
			roOnly: true,
		},
		{
			name: "multi node reader block mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Block{
						Block: &csi.VolumeCapability_BlockVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_MULTI_NODE_READER_ONLY,
					},
				},
			},
			roOnly: true,
		},
		{
			name: "single node reader block mode",
			caps: []*csi.VolumeCapability{
				{
					AccessType: &csi.VolumeCapability_Block{
						Block: &csi.VolumeCapability_BlockVolume{},
					},
					AccessMode: &csi.VolumeCapability_AccessMode{
						Mode: csi.VolumeCapability_AccessMode_SINGLE_NODE_READER_ONLY,
					},
				},
			},
			roOnly: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			roOnly := IsReaderOnly(tt.caps)
			if roOnly != tt.roOnly {
				t.Errorf("isReadOnly() roOnly = %v, want %v", roOnly, tt.roOnly)
			}
		})
	}
}
