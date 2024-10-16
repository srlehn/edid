//go:build generate

//go:generate go run pnp_gen.go

package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"text/template"
	"time"
)

func main() {
	// https://uefi.org/PNP_ID_List?search=&order=field_pnp_id&sort=asc&page=1
	resp, err := http.Get(`https://uefi.org/uefi-pnp-export`)
	if err != nil {
		log.Fatal(err)
	}
	rdr := csv.NewReader(resp.Body)
	rdr.FieldsPerRecord = 3
	rdr.LazyQuotes = true
	pnpEntries := make(map[PNPID]*PNPEntry)
	isFirst := true
	for {
		rec, err := rdr.Read()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		if isFirst {
			isFirst = false
			continue
		}
		if len(rec) != 3 {
			log.Fatal(errors.New(`column count is not 3`))
		}
		id, err := NewPNPID(rec[1])
		if err != nil {
			log.Fatal(err)
		}
		entry := &PNPEntry{
			Company:        rec[0],
			DateOfApproval: rec[2],
		}
		pnpEntries[id] = entry
	}

	tmpl := "// Code generated by \"pnp_gen.go\" DO NOT EDIT.\n\n" +
		"// Entries were copied from `https://uefi.org/uefi-pnp-export on " + time.Now().Format(`2006.01.02`) + ".\n\n" +
		"package edid\n\n" +
		"var pnpIDEntries = map[string]PNPEntry{\n{{range $id, $val := .}}\t`{{$id}}`: {" + /*"ID: `{{$id}}`, "+*/ "DateOfApproval: `{{$val.DateOfApproval}}`, Company: `{{$val.Company}}`},\n{{end}}}\n"

	t := template.Must(template.New("tmpl").Parse(tmpl))

	f, err := os.OpenFile(`pnpids.go`, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	t.Execute(f, pnpEntries)

	/*
		for id, entry := range pnpEntries {
			fmt.Println(id, entry.DateOfApproval, entry.Company)
		}
	*/
}

type PNPEntry struct {
	ID             string
	DateOfApproval string
	Company        string
}

type PNPID [3]byte

func NewPNPID(s string) (PNPID, error) {
	var ret PNPID
	// ARMSTEL, Inc.	AMS 	02/25/2011
	// AMS ends with '\u00a0' ('NO-BREAK SPACE')
	s = strings.TrimSpace(s)
	if len(s) != 3 {
		fmt.Printf("%q\n", s)
		return ret, errors.New(`id char length is not 3`)
	}

	copy(ret[:], s[:3])

	/*
		non-upper case ids:
		Vision Quest	VQ@
		Inovatec S.p.A.	inu
	*/

	return ret, nil
}

func (p PNPID) String() string {
	return string(p[:])
}
