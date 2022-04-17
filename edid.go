package edid

import (
	"time"
)

type EDID struct {
	Display                                                  string
	data                                                     []byte
	IsActive                                                 bool
	Querier                                                  string
	Model                                                    uint32
	Serial                                                   uint64
	ProductionYear                                           uint
	ProductionWeek                                           uint
	ProductionWeekFF                                         bool
	VersionMajor                                             uint
	VersionMinor                                             uint
	PNPID                                                    [3]byte
	Manufacturer                                             string
	ManufacturerPNPDateOfApproval                            time.Time
	Input                                                    string // "digital" / "analog"
	VideoInterface                                           string
	BitDepth                                                 string
	VSyncPulseMustBeSerratedWhenCompositeOrSyncOnGreenIsUsed bool
	SyncOnGreenSupport                                       bool
	CompositeSyncOnHSyncSupport                              bool
	SeparateSyncSupport                                      bool
	BlankToBlackSetupPedestalExpected                        bool
	VideoWhiteAndSyncLevelsRelativeToBlank                   string
	ResolutionPreferredX                                     int
	ResolutionPreferredY                                     int
}

func New(b []byte) (*EDID, error) {
	if !isEDID(b) {
		return nil, errNotEDID
	}
	e := &EDID{data: b}
	err := e.parse()
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (e *EDID) Data() []byte {
	if e == nil {
		return nil
	}
	return e.data
}
