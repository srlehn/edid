//+build windows

package edid

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
	"golang.org/x/sys/windows/registry"
)

var (
	setupapi = windows.NewLazySystemDLL("setupapi.dll")

	procSetupDiGetClassDevsW         = setupapi.NewProc("SetupDiGetClassDevsW")
	procSetupDiDestroyDeviceInfoList = setupapi.NewProc("SetupDiDestroyDeviceInfoList")
	procSetupDiEnumDeviceInfo        = setupapi.NewProc("SetupDiEnumDeviceInfo")
	procSetupDiOpenDevRegKey         = setupapi.NewProc("SetupDiOpenDevRegKey")

	// GUID_DEVINTERFACE_MONITOR
	guidDevInterfaceMonitor = &windows.GUID{0xE6F07B5F, 0xEE97, 0x4A90, [8]byte{0xB0, 0x76, 0x33, 0xF5, 0x7B, 0xF4, 0xEA, 0xA7}}
)

// SP_DEVINFO_DATA
// https://docs.microsoft.com/en-us/windows/win32/api/setupapi/ns-setupapi-sp_devinfo_data
type devInfoData struct {
	cbSize    uint32
	classGUID windows.GUID
	devInst   uint32
	reserved  uint
}

func queryEDIDSetupAPI() ([]*EDID, error) {
	const (
		DIGCF_PRESENT         uintptr = 0x00000002
		DIGCF_DEVICEINTERFACE uintptr = 0x00000010
		KEY_READ              uintptr = 0x00020019
		DICS_FLAG_GLOBAL      uintptr = 0x00000001
		DIREG_DEV             uintptr = 0x00000001
	)

	// get device information set
	// https://docs.microsoft.com/en-us/windows-hardware/drivers/install/device-information-sets
	dis, _, err := syscall.Syscall6(
		procSetupDiGetClassDevsW.Addr(),
		4,
		uintptr(unsafe.Pointer(guidDevInterfaceMonitor)),
		0,
		0,
		DIGCF_DEVICEINTERFACE|DIGCF_PRESENT,
		0, 0,
	)
	if syscall.Handle(dis) == syscall.InvalidHandle {
		if err == 0 {
			err = syscall.EINVAL
		}
		return nil, err
	}
	defer syscall.Syscall(
		procSetupDiDestroyDeviceInfoList.Addr(),
		1,
		uintptr(dis),
		0, 0,
	)

	var ret []*EDID
	for i := 0; ; i++ {
		var did devInfoData
		did.cbSize = uint32(unsafe.Sizeof(did))
		_, _, errNo := syscall.Syscall(
			procSetupDiEnumDeviceInfo.Addr(),
			3,
			uintptr(dis),
			uintptr(i),
			uintptr(unsafe.Pointer(&did)),
		)
		if errNo != 0 {
			break
		}
		r1, _, _ := syscall.Syscall6(
			procSetupDiOpenDevRegKey.Addr(),
			6,
			dis,
			uintptr(unsafe.Pointer(&did)),
			DICS_FLAG_GLOBAL, // scope: global configuration information
			0,                // profile: current hardware profile
			DIREG_DEV,        // key type: hardware key
			KEY_READ,         // sam desired
		)
		if syscall.Handle(r1) == syscall.InvalidHandle {
			continue
		}
		key := registry.Key(r1)
		edid, _, err := key.GetBinaryValue(`EDID`)
		if err == nil {
			ret = append(
				ret,
				&EDID{
					Data:    edid,
					Querier: `setupapi`,
				},
			)
		}
	}

	return ret, nil
}
