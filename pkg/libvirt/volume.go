package libvirt

import (
	"libvirt.org/go/libvirtxml"
)

func volumeXML(name string, size uint64) (string, error) {
	volume := libvirtxml.StorageVolume{
		Type: "file",
		Name: name,
		Capacity: &libvirtxml.StorageVolumeSize{
			Unit:  "bytes",
			Value: size,
		},
		Target: &libvirtxml.StorageVolumeTarget{
			Format: &libvirtxml.StorageVolumeTargetFormat{
				Type: "qcow2",
			},
			Features: []libvirtxml.StorageVolumeTargetFeature{
				{
					LazyRefcounts: &struct{}{},
				},
			},
		},
	}
	return volume.Marshal()
}
