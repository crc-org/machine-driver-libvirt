package libvirt

import (
	"testing"

	"github.com/crc-org/machine/drivers/libvirt"
	"github.com/crc-org/machine/libmachine/drivers"
	"github.com/stretchr/testify/assert"
)

func TestTemplating(t *testing.T) {
	xml, err := domainXML(&Driver{
		Driver: &libvirt.Driver{
			VMDriver: &drivers.VMDriver{
				BaseDriver: &drivers.BaseDriver{
					MachineName: "domain",
				},
				ImageSourcePath: "disk_path",
				ImageFormat:     "test",
				Memory:          4096,
				CPU:             4,
			},
			Network:   "network",
			CacheMode: "default",
			IOMode:    "threads",
			VSock:     false,
		},
	}, "q35")

	assert.NoError(t, err)
	assert.Equal(t, `<domain type="kvm">
  <name>domain</name>
  <memory unit="MiB">4096</memory>
  <vcpu>4</vcpu>
  <os>
    <type machine="q35">hvm</type>
    <boot dev="hd"></boot>
    <bootmenu enable="no"></bootmenu>
  </os>
  <features>
    <pae></pae>
    <acpi></acpi>
    <apic></apic>
  </features>
  <cpu mode="host-passthrough"></cpu>
  <clock offset="utc"></clock>
  <devices>
    <disk type="file" device="disk">
      <driver name="qemu" type="qcow2"></driver>
      <source file="machines/domain/domain.test"></source>
      <target dev="vda" bus="virtio"></target>
    </disk>
    <interface type="network">
      <mac address="52:fd:fc:07:21:82"></mac>
      <source network="network"></source>
      <model type="virtio"></model>
    </interface>
    <console type="stdio"></console>
    <graphics type="vnc"></graphics>
    <memballoon model="none"></memballoon>
    <rng model="virtio">
      <backend model="random">/dev/urandom</backend>
    </rng>
  </devices>
</domain>`, xml)
}

func TestVSockTemplating(t *testing.T) {
	xml, err := domainXML(&Driver{
		Driver: &libvirt.Driver{
			VMDriver: &drivers.VMDriver{
				BaseDriver: &drivers.BaseDriver{
					MachineName: "domain",
				},
				ImageSourcePath: "disk_path",
				ImageFormat:     "test",
				Memory:          4096,
				CPU:             4,
			},
			Network:   "crc",
			CacheMode: "default",
			IOMode:    "threads",
			VSock:     true,
		},
	}, "")
	assert.NoError(t, err)
	assert.Regexp(t, `(?s)<devices>(.*?)<vsock model="virtio">\s*<cid auto="yes">\s*</cid>\s*</vsock>(.*?)</devices>`, xml)
	assert.Regexp(t, `(?s)<os>(.*?)<type>hvm</type>(.*?)</os>`, xml)
}

func TestNetworkTemplating(t *testing.T) {
	xml, err := domainXML(&Driver{
		Driver: &libvirt.Driver{
			VMDriver: &drivers.VMDriver{
				BaseDriver: &drivers.BaseDriver{
					MachineName: "domain",
				},
				ImageSourcePath: "disk_path",
				ImageFormat:     "test",
				Memory:          4096,
				CPU:             4,
			},
			Network:   "crc",
			CacheMode: "default",
			IOMode:    "threads",
			VSock:     true,
		},
	}, "q35")
	assert.NoError(t, err)
	assert.Contains(t, xml, `<interface type="network">
      <mac address="52:fd:fc:07:21:82"></mac>
      <source network="crc"></source>
      <model type="virtio"></model>
    </interface>`)
}
