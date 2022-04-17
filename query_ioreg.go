package edid

import (
	"bufio"
	"encoding/hex"
	"errors"
	"os/exec"
	"strings"
)

func queryEDIDIORegCmd(cmd *exec.Cmd) ([]*EDID, error) {
	if cmd == nil {
		return nil, errors.New(`received nil parameter`)
	}

	// cmd.Env = os.Environ()
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(stdout)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	var edids []*EDID
	for scanner.Scan() {
		line := strings.TrimLeft(scanner.Text(), `| `)
		if !strings.HasPrefix(line, `"IODisplayEDID" = `) {
			continue
		}
		sp := strings.SplitN(line, ` = <`, 2)
		if len(sp) < 2 {
			continue
		}
		line = strings.TrimRight(strings.TrimSpace(sp[1]), `>`)
		b, err := hex.DecodeString(line)
		if err != nil {
			continue
		}
		edids = append(edids, &EDID{Querier: `ioreg-command`, data: b})
	}
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	return edids, nil
}

/*
https://stackoverflow.com/a/60658692
ioreg -lw0 -r -c "IODisplayConnect" -d 2 | grep IODisplayEDID

https://github.com/Akemi/macOS-edid-modification/blob/618b259c0fdfad00ae664cf5ea9bf54df07065eb/ioreg-short-info.txt
"    | |   |           | |           "IODisplayEDID" = <00ffffffffffff0010ac65d04c534c3027...00000000000000fe>"

https://developer.apple.com/forums/thread/667608
Apple M1: iterate over IODPDevice instead of IODisplayConnect
https://developer.apple.com/forums/thread/666383
*/

/*
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
*/
