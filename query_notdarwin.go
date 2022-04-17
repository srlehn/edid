//+build !darwin

package edid

import "errors"

func queryEDIDIOReg() ([]*EDID, error) {
	return nil, errors.New(`not implemented on this platform`)
}
