//+build darwin

// TODO later exclude cgo when cgo version is implemented

package edid

// TODO Apple M1

import (
	"bufio"
	"encoding/hex"
	"os"
	"os/exec"
	"strings"
)

func queryEDIDIOReg() ([]*EDID, error) {
	var (
		ioregPath = `/usr/sbin/ioreg`
		// TODO Apple M1 adjust for "IODPDevice"
		// https://developer.apple.com/forums/thread/667608
		args = []string{
			`-l`,      // list properties of all objects
			`-w`, `0`, // clip line length (0 is unlimited)
			`-r`,                     // show subtrees rooted by the given criteria
			`-c`, `IODisplayConnect`, // list properties of objects with the given class
			`-d`, `2`, // limit tree depth
		}
	)

	cmd := exec.Command(ioregPath, args...)

	return queryEDIDIORegCmd(cmd)
}
