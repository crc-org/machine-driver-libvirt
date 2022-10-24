package libvirt

import (
	"fmt"

	"github.com/libvirt/libvirt-go"
	"libvirt.org/go/libvirtxml"

	"github.com/crc-org/machine/libmachine/drivers"
)

const macAddress = "52:fd:fc:07:21:82"

func domainXML(d *Driver, machineType string) (string, error) {
	domain := libvirtxml.Domain{
		Type: "kvm",
		Name: d.MachineName,
		Memory: &libvirtxml.DomainMemory{
			Value: uint(d.Memory),
			Unit:  "MiB",
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
						Name: "qemu",
						Type: "qcow2",
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
				{
					Source: &libvirtxml.DomainChardevSource{
						StdIO: &libvirtxml.DomainChardevSourceStdIO{},
					},
				},
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
	if machineType != "" {
		domain.OS.Type.Machine = machineType
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
	if virtiofsSupported(d.conn) == nil && len(d.SharedDirs) != 0 {
		domain.MemoryBacking = &libvirtxml.DomainMemoryBacking{
			MemorySource: &libvirtxml.DomainMemorySource{
				Type: "memfd",
			},
			MemoryAccess: &libvirtxml.DomainMemoryAccess{
				Mode: "shared",
			},
		}
		for _, sharedDir := range d.SharedDirs {
			filesystem := libvirtxml.DomainFilesystem{
				AccessMode: "passthrough",
				Driver: &libvirtxml.DomainFilesystemDriver{
					Type: "virtiofs",
				},
				Source: &libvirtxml.DomainFilesystemSource{
					Mount: &libvirtxml.DomainFilesystemSourceMount{
						Dir: sharedDir.Source,
					},
				},
				Target: &libvirtxml.DomainFilesystemTarget{
					Dir: sharedDir.Tag,
				},
			}
			domain.Devices.Filesystems = append(domain.Devices.Filesystems, filesystem)
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

func virtiofsSupported(conn *libvirt.Connect) error {
	if conn == nil {
		return drivers.ErrNotSupported
	}

	guest, err := getBestGuestFromCaps(conn)
	if err != nil {
		return err
	}

	domainCapsXML, err := conn.GetDomainCapabilities(guest.Arch.Emulator, guest.Arch.Name, getMachineType(guest), "kvm", 0)
	if err != nil {
		return err
	}

	caps := &libvirtxml.DomainCaps{}
	err = caps.Unmarshal(domainCapsXML)
	if err != nil {
		return fmt.Errorf("Error parsing libvirt domain capabilities: %w", err)
	}

	if caps.Devices.FileSystem == nil {
		return drivers.ErrNotSupported
	}

	if caps.Devices.FileSystem.Supported != "yes" {
		return drivers.ErrNotSupported
	}
	for _, enum := range caps.Devices.FileSystem.Enums {
		if enum.Name != "driverType" {
			continue
		}
		for _, val := range enum.Values {
			if val == "virtiofs" {
				return nil
			}
		}
	}

	return drivers.ErrNotSupported
}
