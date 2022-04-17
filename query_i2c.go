package edid

import (
	"bytes"
	"errors"
	"os"
	"strings"

	"periph.io/x/conn/v3/i2c"
	"periph.io/x/conn/v3/i2c/i2creg"
	"periph.io/x/host/v3"
)

func queryEDIDI2C() ([]*EDID, error) {
	const (
		addr = 0x50 // EDID
		size = 256
	)
	if _, err := host.Init(); err != nil {
		return nil, err
	}

	refs := i2creg.All()
	if len(refs) == 0 {
		return nil, errors.New(`i2creg: no bus found`)
	}
	var edids []*EDID
	var missingRights bool
	for _, ref := range refs {
		if ref == nil {
			continue
		}
		edid, err := queryEDIDI2CPeriph(ref, addr, size)
		if err == nil {
			// fmt.Printf("periph %s 0x%x\n", ref.Name, edid.Data[:16])
			// fmt.Printf("periph %s 0x%x\n", ref.Name, edid.Data)
			edid.Querier = `i2c-periph`
			edids = append(edids, edid)
			continue
		} else {
			// underlying error is not wrapped
			// https://github.com/periph/host/blob/17c4f529ff591729c88f949d97c83d3a093d9e7e/sysfs/i2c.go#L172
			// TODO: use os.IsPermission() when fixed.
			if s := err.Error(); strings.HasPrefix(s, `sysfs-i2c: are you member of group 'plugdev'? open `) &&
				strings.HasSuffix(s, `: permission denied`) {
				missingRights = true
			}
		}
		edid, err = queryEDIDI2CFallback(ref.Name, addr, size)
		if err == nil {
			// fmt.Printf("fallback %s 0x%x\n", ref.Name, edid.Data[:16])
			// fmt.Printf("fallback %s 0x%x\n", ref.Name, edid.Data)
			edid.Querier = `i2c-fallback`
			edids = append(edids, edid)
		} else if os.IsPermission(err) {
			missingRights = true
		}
	}
	if len(edids) == 0 && missingRights {
		return nil, errors.New(`missing access rights for i2c devices`)
	}
	return edids, nil
}

/*
TODO: move from ioctlSlave to ioctlRdwr
// periph.io/x/host/v3/sysfs/i2c.go - (*I2C).Tx()

// original
/not working
if err := i.f.Ioctl(ioctlRdwr, pp); err != nil {
	return fmt.Errorf("sysfs-i2c: %v", err)
}

// fix
// working
if err := i.f.Ioctl(ioctlSlave, uintptr(addr)); err != nil {
	return fmt.Errorf("sysfs-i2c: %v", err)
}
if _, err := io.ReadFull(i.f.(io.Reader), r); err != nil {
	return err
}
*/

func queryEDIDI2CPeriph(ref *i2creg.Ref, addr uint16, size uint) (*EDID, error) {
	b, err := ref.Open()
	if err != nil {
		return nil, err
	}
	defer b.Close()
	d := &i2c.Dev{Addr: addr, Bus: b}
	edid := make([]byte, size)
	err = d.Tx(nil, edid)
	if err != nil {
		return nil, err
	}
	return findEDID(edid)
}

func findEDID(b []byte) (*EDID, error) {
	// fmt.Printf("0x%x 0x%x\n", b[:8], b[128:128+8])
	switch len(b) {
	case 128:
		if bytes.Equal(b[:8], edidPrefix) {
			return &EDID{data: b}, nil
		}
	case 256:
		if bytes.Equal(b[:8], edidPrefix) {
			return &EDID{data: b[:128]}, nil
		}
		if bytes.Equal(b[128:128+8], edidPrefix) {
			return &EDID{data: b[128:]}, nil
		}
	default:
		return nil, errors.New(`received byte slice of wrong length`)
	}
	return nil, errors.New(`no edid found in byte slice`)
}
