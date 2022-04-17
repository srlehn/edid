# NOTES
https://github.com/linuxhw/EDID

https://thinkwiki.de/Display-EDID_ver%C3%A4ndern

https://www.extron.com/article/uedid

https://bugs.chromium.org/p/chromium/issues/detail?id=240341

## VESA PNP
https://uefi.org/pnp_id_list

https://uefi.org/uefi-pnp-export (.xls with browser, .html with curl)

https://github.com/kontron/keapi/blob/master/include/eapi/PNPID.h

### EISA/ISA PNP ID (probably irrelevant)
https://web.archive.org/web/20040408140943/http://www.microsoft.com/whdc/hwdev/pnpid.mspx

https://web.archive.org/web/20040411214849/http://www.microsoft.com/whdc/hwdev/tech/pnp/default.mspx

## I2C
accessible via the I²C-bus at address A0

http://www.polypux.org/projects/read-edid/

sudo modprobe i2-dev

## VESA VBE
https://wiki.osdev.org/Getting_VBE_Mode_Info

https://cgit.freedesktop.org/~libv/vbe-edid

http://www.polypux.org/projects/read-edid/

https://wiki.osdev.org/Getting_VBE_Mode_Info

## X11 RandR
https://github.com/burntsushi/xgb/blob/deaf085860bc/examples/randr/main.go

https://chromium.googlesource.com/chromium/src/base/+/a3305756b9f14bb8a3d6961e79b490b8671c075d/x11/edid_parser_x11.cc

## x86 - V86 mode
https://wiki.osdev.org/Virtual_8086_Mode

https://git.linuxtv.org/edid-decode.git/

## Windows
https://social.msdn.microsoft.com/Forums/windowsdesktop/en-US/1a19a278-c296-4d34-ade7-83bf3315db96/how-to-read-edid-data-direct-from-monitor-not-registry#e184c327-15a1-4c88-a223-c96ac694912b-isAnswer

https://github.com/chromium/chromium/blob/d7da0240cae77824d1eda25745c4022757499131/ui/gfx/win/physical_size.cc

* SetupDiGetClassDevs(GUID_DEVINTERFACE_MONITOR)
* SetupDiEnumDeviceInterfaces()
* SetupDiOpenDevRegKey()

> This means that Windows has such a functionality but does not expose it.

Right Correct.
the method to read EDID is private to the display/monitor driver.
it does not expose a public API for you to get at the EDID directly.
why can't you use the data in the device instance (what you refer to as HKEY_LOCAL_MACHINE\SYSTEM\CurrentControlSet\Enum\DISPLAY\{Type}\{SWnummer}\Device Parameters)?
is it because you can't construct this path in the registry ahead of time? that is by design.
you need to use SetupDiGetClassDevices(GUID_DEVINTERFACE_MONITOR) / SetupDiEnumDeviceInterfaces / SetupDiOpenDevRegKey to get an HKEY to this path.
in other words, the path is an abstraction that you should not parse, the system provides APIs for you to get to the Device Parameters key without knowing the path.

https://www.winvistatips.com/threads/how-to-read-monitors-edid-information.181727/#post-916636

https://gist.github.com/texus/3212ebc1ed1502ecd265cc7cf1322b02

https://github.com/distatus/battery/blob/master/battery_windows.go

https://github.com/qemu/qemu/blob/master/qga/commands-win32.c

https://www.nirsoft.net/utils/monitor_info_view.html

## macos
https://developer.apple.com/documentation/kernel/ioframebuffer/1813183-getddcblock

https://opensource.apple.com/source/IOGraphics/IOGraphics-123/IOGraphicsFamily/IOKit/graphics/IOFramebuffer.h.auto.html

getDDCBlock()

https://developer.apple.com/forums/thread/667608

IODisplayConnect() (old)

IODPDevice()

https://developer.apple.com/forums/thread/666383

https://stackoverflow.com/a/60658692

`ioreg -lw0 -r -c "IODisplayConnect" -d 2 | grep IODisplayEDID`

http://support.touch-base.com/Documentation/50730/EDID-Structure

`ioreg -l | grep IODisplayEDID`

https://www.hackintosh-forum.de/forum/thread/36617-sammelthread-f%C3%BCr-ioreg-befehle-welche-kennt-ihr-denn-noch/
```
usage: ioreg [-abfilrtx] [-c class] [-d depth] [-k key] [-n name] [-p plane] [-w width]
where options are:
	-a archive output
	-b show object name in bold
	-c list properties of objects with the given class
	-d limit tree to the given depth
	-f enable smart formatting
	-i show object inheritance
	-k list properties of objects with the given key
	-l list properties of all objects
	-n list properties of objects with the given name
	-p traverse registry over the given plane (IOService is default)
	-r show subtrees rooted by the given criteria
	-t show location of each subtree
	-w clip output to the given line width (0 is unlimited)
	-x show data and numbers as hexadecimal
```

https://github.com/Akemi/macOS-edid-modification/blob/master/ioreg-short-info.txt

ioreg example output

https://stackoverflow.com/a/56637436

IORegistryEntryCreateCFProperty()

https://developer.apple.com/documentation/iokit/1514293-ioregistryentrycreatecfproperty?language=objc

## Android
https://source.android.com/devices/tech/display/multi_display/displays?hl=en

## Parse
https://github.com/chromium/chromium/blob/d7da0240cae77824d1eda25745c4022757499131/ui/display/util/edid_parser.cc

https://github.com/chromium/chromium/blob/d7da0240cae77824d1eda25745c4022757499131/ui/display/types/display_snapshot.cc
