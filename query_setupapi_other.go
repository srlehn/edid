//+build !windows

package edid

import "errors"

func queryEDIDSetupAPI() ([]*EDID, error) {
	return nil, errors.New(`SetupAPI not available on this platform`)
}
