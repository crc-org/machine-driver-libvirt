package libvirt

import (
	"fmt"
	"os"

	"github.com/libvirt/libvirt-go"
	libvirtxml "github.com/libvirt/libvirt-go-xml"

	"github.com/code-ready/machine/libmachine/log"
)

func (d *Driver) activateStoragePool(pool *libvirt.StoragePool) error {
	log.Debugf("Activating pool '%s'", d.getStoragePoolName())

	if err := os.MkdirAll(d.ResolveStorePath("."), 0755); err != nil {
		return err
	}

	if err := pool.Create(libvirt.STORAGE_POOL_CREATE_NORMAL); err != nil {
		log.Warnf("Failed to start storage pool: %s", err)
		return err
	}

	return nil
}

// Create, or verify the private storage pool is properly configured
// storage pool must be preexisting, which breaks upgrades
func (d *Driver) validateStoragePool() error {
	log.Debug("Validating storage pool")
	pool, err := d.getPool()
	if err != nil {
		/* FIXME: not the right place to talk about 'crc setup' */
		return fmt.Errorf("Use 'crc setup' to define the machine driver storage pool, %+v", err)
	}
	defer pool.Free() // nolint:errcheck

	return nil
}

func (d *Driver) getStoragePoolName() string {
	if d.StoragePool != "" {
		return d.StoragePool
	}
	if d.MachineName != "" {
		return d.MachineName
	}
	return DefaultPool
}

func (d *Driver) refreshStoragePool() error {
	pool, err := d.getPool()
	if err != nil {
		return err
	}
	return pool.Refresh(0)
}

func (d *Driver) createStoragePool() (*libvirt.StoragePool, error) {
	log.Debug("Creating storage pool")

	conn, err := d.getConn()
	if err != nil {
		return nil, err
	}
	poolName := d.getStoragePoolName()
	poolConfig := libvirtxml.StoragePool{
		Name: poolName,
		Type: "dir",
		Target: &libvirtxml.StoragePoolTarget{
			Path: d.ResolveStorePath("."),
		},
	}
	poolXML, err := poolConfig.Marshal()
	if err != nil {
		return nil, err
	}
	log.Infof("Creating storage pool with XML %s", poolXML)
	pool, err := conn.StoragePoolDefineXML(poolXML, uint32(libvirt.STORAGE_POOL_CREATE_NORMAL))
	if err != nil {
		log.Debugf("Could not create storage pool %s", d.StoragePool)
		return nil, fmt.Errorf("Use 'crc setup' to define the storage pool, %+v", err)
	}
	err = d.activateStoragePool(pool)
	if err != nil {
		return nil, err
	}
	return pool, nil
}

func (d *Driver) getPool() (*libvirt.StoragePool, error) {
	conn, err := d.getConn()
	if err != nil {
		return nil, err
	}
	pool, err := conn.LookupStoragePoolByName(d.getStoragePoolName())
	if err != nil {
		log.Debugf("Could not find storage pool '%s', trying to create it", d.getStoragePoolName())
		return d.createStoragePool()
	}

	// Corner case, but might happen...
	if active, _ := pool.IsActive(); !active {
		err = d.activateStoragePool(pool)
		if err != nil {
			return nil, err
		}
	}

	return pool, nil
}
