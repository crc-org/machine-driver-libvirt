package libvirt

import (
	"testing"

	"github.com/code-ready/machine/drivers/libvirt"
	"github.com/code-ready/machine/libmachine/drivers"
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
	})

	assert.NoError(t, err)
	assert.Equal(t, `<domain type="kvm">
  <name>domain</name>
  <memory unit="MB">4096</memory>
  <vcpu placement="static">4</vcpu>
  <os>
    <type arch="x86_64">hvm</type>
    <boot dev="hd"></boot>
    <bootmenu enable="no"></bootmenu>
  </os>
  <features>
    <pae></pae>
    <acpi></acpi>
    <apic></apic>
  </features>
  <cpu mode="host-passthrough">
    <feature policy="disable" name="rdrand"></feature>
  </cpu>
  <clock offset="utc"></clock>
  <devices>
    <disk type="file" device="disk">
      <driver name="qemu" type="qcow2" cache="default" io="threads"></driver>
      <source file="machines/domain/domain.test"></source>
      <target dev="vda" bus="virtio"></target>
    </disk>
    <interface type="network">
      <mac address="52:fd:fc:07:21:82"></mac>
      <source network="network"></source>
      <model type="virtio"></model>
    </interface>
    <console></console>
    <graphics type="vnc" autoport="yes" listen="127.0.0.1">
      <listen type="address" address="127.0.0.1"></listen>
    </graphics>
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
	})
	assert.NoError(t, err)
	assert.Regexp(t, `(?s)<devices>(.*?)<vsock model="virtio">\s*<cid auto="yes">\s*</cid>\s*</vsock>(.*?)</devices>`, xml)
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
	})
	assert.NoError(t, err)
	assert.Contains(t, xml, `<interface type="network">
      <mac address="52:fd:fc:07:21:82"></mac>
      <source network="crc"></source>
      <model type="virtio"></model>
    </interface>`)
}
