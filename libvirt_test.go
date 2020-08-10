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
	assert.Equal(t, `<domain type='kvm'>
  <name>domain</name>
  <memory unit='MB'>4096</memory>
  <vcpu placement='static'>4</vcpu>
  <features><acpi/><apic/><pae/></features>
  <cpu mode='host-passthrough'>
    <feature policy="disable" name="rdrand"/>
  </cpu>
  <os>
    <type arch='x86_64'>hvm</type>
    <boot dev='hd'/>
    <bootmenu enable='no'/>
  </os>
  <features>
    <acpi/>
    <apic/>
    <pae/>
  </features>
  <clock offset='utc'/>
  <on_poweroff>destroy</on_poweroff>
  <on_reboot>restart</on_reboot>
  <on_crash>destroy</on_crash>
  <devices>
    <disk type='file' device='disk'>
      <driver name='qemu' type='qcow2' cache='default' io='threads' />
      <source file='machines/domain/domain.test'/>
      <target dev='vda' bus='virtio'/>
    </disk>
    <graphics type='vnc' autoport='yes' listen='127.0.0.1'>
      <listen type='address' address='127.0.0.1'/>
    </graphics>
    <console type='pty'></console>
    <channel type='pty'>
      <target type='virtio' name='org.qemu.guest_agent.0'/>
    </channel>
    <rng model='virtio'>
      <backend model='random'>/dev/urandom</backend>
    </rng>
    <interface type='network'>
      <mac address='52:fd:fc:07:21:82'/>
      <source network='network'/>
      <model type='virtio'/>
    </interface>
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
	assert.Regexp(t, `(?s)<devices>(.*?)<vsock model='virtio'><cid auto='yes'/></vsock>(.*?)</devices>`, xml)
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
	assert.Contains(t, xml, `<interface type='network'>
      <mac address='52:fd:fc:07:21:82'/>
      <source network='crc'/>
      <model type='virtio'/>
    </interface>`)
}
