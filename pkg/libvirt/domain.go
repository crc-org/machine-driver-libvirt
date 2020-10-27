package libvirt

import (
	libvirtxml "github.com/libvirt/libvirt-go-xml"
)

const macAddress = "52:fd:fc:07:21:82"

func domainXML(d *Driver) (string, error) {
	domain := libvirtxml.Domain{
		Type: "kvm",
		Name: d.MachineName,
		Memory: &libvirtxml.DomainMemory{
			Value: uint(d.Memory),
			Unit:  "MB",
		},
		VCPU: &libvirtxml.DomainVCPU{
			Value: uint(d.CPU),
		},
		Features: &libvirtxml.DomainFeatureList{
			ACPI: &libvirtxml.DomainFeature{},
			APIC: &libvirtxml.DomainFeatureAPIC{},
			PAE:  &libvirtxml.DomainFeature{},
		},
		CPU: &libvirtxml.DomainCPU{
			Mode: "host-passthrough",
			// https://bugzilla.redhat.com/show_bug.cgi?id=1806532
			Features: []libvirtxml.DomainCPUFeature{
				{
					Policy: "disable",
					Name:   "rdrand",
				},
			},
		},
		OS: &libvirtxml.DomainOS{
			Type: &libvirtxml.DomainOSType{
				Type: "hvm",
			},
			BootDevices: []libvirtxml.DomainBootDevice{
				{
					Dev: "hd",
				},
			},
			BootMenu: &libvirtxml.DomainBootMenu{
				Enable: "no",
			},
		},
		Clock: &libvirtxml.DomainClock{
			Offset: "utc",
		},
		Devices: &libvirtxml.DomainDeviceList{
			Disks: []libvirtxml.DomainDisk{
				{
					Device: "disk",
					Driver: &libvirtxml.DomainDiskDriver{
						Name:  "qemu",
						Type:  "qcow2",
						Cache: d.CacheMode,
						IO:    d.IOMode,
					},
					Source: &libvirtxml.DomainDiskSource{
						File: &libvirtxml.DomainDiskSourceFile{
							File: d.getDiskImagePath(),
						},
					},
					Target: &libvirtxml.DomainDiskTarget{
						Dev: "vda",
						Bus: "virtio",
					},
				},
			},
			Graphics: []libvirtxml.DomainGraphic{
				{
					VNC: &libvirtxml.DomainGraphicVNC{},
				},
			},
			Consoles: []libvirtxml.DomainConsole{
				{},
			},
			RNGs: []libvirtxml.DomainRNG{
				{
					Model: "virtio",
					Backend: &libvirtxml.DomainRNGBackend{
						Random: &libvirtxml.DomainRNGBackendRandom{
							Device: "/dev/urandom",
						},
					},
				},
			},
			MemBalloon: &libvirtxml.DomainMemBalloon{
				Model: "none",
			},
		},
	}
	if d.Network != "" {
		domain.Devices.Interfaces = []libvirtxml.DomainInterface{
			{
				MAC: &libvirtxml.DomainInterfaceMAC{
					Address: macAddress,
				},
				Source: &libvirtxml.DomainInterfaceSource{
					Network: &libvirtxml.DomainInterfaceSourceNetwork{
						Network: d.Network,
					},
				},
				Model: &libvirtxml.DomainInterfaceModel{
					Type: "virtio",
				},
			},
		}
	}
	if d.VSock {
		domain.Devices.VSock = &libvirtxml.DomainVSock{
			Model: "virtio",
			CID: &libvirtxml.DomainVSockCID{
				Auto: "yes",
			},
		}
	}
	return domain.Marshal()
}
