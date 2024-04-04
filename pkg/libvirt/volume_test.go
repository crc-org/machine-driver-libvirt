package libvirt

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVolumeXML(t *testing.T) {
	xml, err := volumeXML("crc-second-disk.qcow2", 5000)
	assert.NoError(t, err, "unexpected error")
	assert.Equal(t, `<volume type="file">
  <name>crc-second-disk.qcow2</name>
  <capacity unit="bytes">5000</capacity>
  <target>
    <format type="qcow2"></format>
    <features>
      <lazy_refcounts></lazy_refcounts>
    </features>
  </target>
</volume>`, xml, "unexpected volume xml")
}
