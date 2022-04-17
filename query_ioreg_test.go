//go:build ignore

package edid

import (
	"os/exec"
	"testing"
)

func TestQueryEDIDIORegCmd(t *testing.T) {
	cmd := exec.Command(`cat`, `./testdata/ioreg-short-info.txt`)

	edids, err := queryEDIDIORegCmd(cmd)

	if err != nil {
		t.Fatal(err)
	}

	_ = edids
}
