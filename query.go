package edid

import (
	"errors"
	"runtime"
)

// QueryEDID - valid methods: "x11","i2c","setupapi","ioreg"
func QueryEDID(method string) ([]*EDID, error) {
	var edids []*EDID
	var errRet error
	switch method {
	case ``:
		var queriers []func() ([]*EDID, error)
		switch runtime.GOOS {
		case `windows`:
			queriers = append(queriers, queryEDIDSetupAPI, queryEDIDI2C, queryEDIDX11)
		case `darwin`:
			queriers = append(queriers, queryEDIDIOReg, queryEDIDI2C, queryEDIDX11, queryEDIDSetupAPI)
		case `android`:
			queriers = append(queriers, queryEDIDI2C, queryEDIDX11, queryEDIDSetupAPI)
		default:
			queriers = append(queriers, queryEDIDX11, queryEDIDI2C, queryEDIDSetupAPI)
		}
		var atLeastOneSucceededMethod bool
		for _, querier := range queriers {
			es, err := querier()
			if err != nil {
				continue
			}
			atLeastOneSucceededMethod = true
			if len(es) == 0 {
				continue
			}
			edids = append(edids, es...)
		}
		if len(edids) == 0 && !atLeastOneSucceededMethod {
			edids, errRet = nil, errors.New(`all methods failed`)
		}
	case `windows`:
		fallthrough
	case `setupapi`:
		edids, errRet = queryEDIDSetupAPI()
	case `randr`: // TODO: x11 + wayland
		fallthrough
	case `x11`:
		edids, errRet = queryEDIDX11()
	case `i2c`:
		edids, errRet = queryEDIDI2C()
	case `darwin`:
		fallthrough
	case `ioreg`:
		edids, errRet = queryEDIDIOReg()
	default:
		return nil, errors.New(`unknown method`)
	}
	// TODO tmp
	var edidsRet []*EDID
	for _, ed := range edids {
		_ = ed.parse()
		edidsRet = append(edidsRet, ed)
	}
	return edidsRet, errRet
}

func filterActive(edids []*EDID) ([]*EDID, error) {
	var ret []*EDID
	for _, ed := range edids {
		if !ed.IsActive {
			continue
		}
		ret = append(ret, ed)
	}
	return ret, nil
}
