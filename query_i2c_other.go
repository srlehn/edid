//+build !linux

package edid

import (
	"errors"
)

func queryEDIDI2CFallback(devFileName string, addr uint16, size uint) (*EDID, error) {
	return nil, errors.New(`not implemented on this platform`)
}
