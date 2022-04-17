//+build linux

package edid

import (
	"io"
	"os"
	"syscall"
)

func queryEDIDI2CFallback(devFileName string, addr uint16, size uint) (*EDID, error) {
	f, err := os.OpenFile(devFileName, os.O_RDWR, os.ModeDevice)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	const i2cSlave = 0x0703
	if _, _, errno := syscall.Syscall(syscall.SYS_IOCTL, f.Fd(), uintptr(i2cSlave), uintptr(addr)); errno != 0 {
		return nil, err
	}
	edid := make([]byte, size)
	if _, err := io.ReadFull(f, edid); err != nil {
		return nil, err
	}
	return findEDID(edid)
}
