package edid

import (
	"bytes"
	"errors"
	"time"
)

var (
	errNilReceiver = errors.New(`nil receiver`)
	errNotEDID     = errors.New(`data seems to be something else than EDID`)
)

func (e *EDID) parse() error {
	if e == nil {
		return errNilReceiver
	}
	if !isEDID(e.data) {
		return errNotEDID
	}

	var err error
	e.header()
	e.getBasicDisplayParameter()
	e.getResolutionPreferred()
	_, err = e.getManufacturer()

	return err
}

var edidPrefix = []byte{0x00, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0xFF, 0x00}

func isEDID(b []byte) bool {
	return len(b) >= 128 && bytes.Equal(b[:len(edidPrefix)], edidPrefix)
}

func (e *EDID) header() {
	if e == nil {
		return
	}
	e.Model = 0
	for i, x := range e.data[10:12] {
		e.Model |= uint32(x) << uint32(i*8)
	}
	e.Serial = 0
	for i, x := range e.data[12:16] {
		e.Serial |= uint64(x) << uint64(i*8)
	}
	e.ProductionWeekFF = e.data[17] == 0xFF
	e.ProductionWeek = uint(e.data[17]) // +-1 week ???
	e.ProductionYear = 1990 + uint(e.data[17])
	e.VersionMajor = uint(e.data[18])
	e.VersionMinor = uint(e.data[19])
}

func (e *EDID) getResolutionPreferred() {
	if e == nil {
		return
	}
	// https://wiki.osdev.org/EDID
	e.ResolutionPreferredX = int(e.data[0x38]) | (int(e.data[0x3A]&0xF0) << 4)
	e.ResolutionPreferredY = int(e.data[0x3B]) | (int(e.data[0x3D]&0xF0) << 4)
}

type PNPEntry struct {
	ID             string
	DateOfApproval string
	Company        string
}

func (e *EDID) getManufacturer() (*PNPEntry, error) {
	if e == nil {
		return nil, errNilReceiver
	}
	var name [3]byte
	// manufacturer_name() in parse-base-block.cpp (edid-decode)
	name[0] = ((e.data[8] & 0x7c) >> 2) + '@'
	name[1] = ((e.data[8] & 0x03) << 3) + ((e.data[9] & 0xe0) >> 5) + '@'
	name[2] = (e.data[9] & 0x1f) + '@'
	e.PNPID = name

	isUpper := func(r byte) bool { return r <= 'Z' && r >= 'A' }
	for _, r := range name {
		if !isUpper(r) {
			return nil, errors.New(`erroneous manufacturer field`)
		}
	}

	var ret *PNPEntry
	entry, ok := pnpIDEntries[string(name[:])]
	if ok {
		e.Manufacturer = entry.Company
		tm, err := time.Parse(`01/02/2006`, entry.DateOfApproval)
		if err != nil {
			return nil, err
		}
		e.ManufacturerPNPDateOfApproval = tm
		ret = &entry
	} else {
		ret = &PNPEntry{}
	}
	ret.ID = string(name[:])

	return ret, nil
}

func (e *EDID) getBasicDisplayParameter() {
	if e == nil {
		return
	}
	// byte 20: video input parameters bitmap
	if e.data[20]&1<<7 == 1<<7 {
		e.Input = `digital`
		switch vi := e.data[20] & 0b1111; vi {
		case 0b0000:
			e.VideoInterface = "undefined"
		case 0b0010:
			e.VideoInterface = "HDMIa"
		case 0b0011:
			e.VideoInterface = "HDMIb"
		case 0b0100:
			e.VideoInterface = "MDDI"
		case 0b0101:
			e.VideoInterface = "DisplayPort"
		}
		switch bd := e.data[20] >> 4 & 0b111; bd {
		case 0b000:
			e.BitDepth = "undefined"
		case 0b001:
			e.BitDepth = "6"
		case 0b010:
			e.BitDepth = "8"
		case 0b011:
			e.BitDepth = "10"
		case 0b100:
			e.BitDepth = "12"
		case 0b101:
			e.BitDepth = "14"
		case 0b110:
			e.BitDepth = "16"
		case 0b111:
			e.BitDepth = "reserved"
		}
	} else {
		e.Input = `analog`
		// VSync pulse must be serrated when composite or sync-on-green is used
		e.VSyncPulseMustBeSerratedWhenCompositeOrSyncOnGreenIsUsed = e.data[20]&1 == 1
		// Sync on green supported
		e.SyncOnGreenSupport = e.data[20]>>1&1 == 1
		// Composite sync (on HSync) supported
		e.CompositeSyncOnHSyncSupport = e.data[20]>>2&1 == 1
		// Separate sync supported
		e.SeparateSyncSupport = e.data[20]>>3&1 == 1
		// Blank-to-black setup (pedestal) expected
		e.BlankToBlackSetupPedestalExpected = e.data[20]>>4&1 == 1
		// video white and sync levels, relative to blank:
		switch bd := e.data[20] >> 5 & 0b11; bd {
		case 0b00:
			e.VideoWhiteAndSyncLevelsRelativeToBlank = "+0.7/-0.3 V"
		case 0b01:
			e.VideoWhiteAndSyncLevelsRelativeToBlank = "+0.714/-0.286 V"
		case 0b10:
			e.VideoWhiteAndSyncLevelsRelativeToBlank = "+1.0/-0.4 V"
		case 0b11:
			e.VideoWhiteAndSyncLevelsRelativeToBlank = "+0.7/0 V"
		}
	}

	// byte 21
	// Horizontal screen size, in centimetres (range 1–255).
	// If vertical screen size is 0, landscape aspect ratio (range 1.00–3.54), datavalue = (AR×100) − 99 (example: 16:9, 79; 4:3, 34.)
	// byte 22
	// Vertical screen size, in centimetres.
	// If horizontal screen size is 0, portrait aspect ratio (range 0.28–0.99), datavalue = (100/AR) − 99 (example: 9:16, 79; 3:4, 34.)
	// If either byte is 0, screen size and aspect ratio are undefined (e.g. projector)

	// byte 23
	// Display gamma, factory default (range 1.00–3.54), datavalue = (gamma×100) − 100 = (gamma − 1)×100. If 255, gamma is defined by DI-EXT block.
}
