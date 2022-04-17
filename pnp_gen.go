//+build generate

//go:generate go run pnp_gen.go

package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/antchfx/htmlquery"
)

func main() {
	doc, err := htmlquery.LoadURL(`https://uefi.org/uefi-pnp-export`)
	// https://uefi.org/pnp_id_list?search=&order=field_pnp_id&sort=asc&page=1
	if err != nil {
		log.Fatal(err)
	}
	rows, err := htmlquery.QueryAll(doc, `/html/body/table/tbody/tr`)
	if err != nil {
		log.Fatal(err)
	}
	pnpEntries := make(map[PNPID]*PNPEntry)
	for _, row := range rows {
		cells, err := htmlquery.QueryAll(row, `/td`)
		if err != nil {
			log.Fatal(err)
		}
		if len(cells) != 3 {
			log.Fatal(errors.New(`column count is not 3`))
		}
		id, err := NewPNPID(htmlquery.InnerText(cells[1]))
		if err != nil {
			log.Fatal(err)
		}
		entry := &PNPEntry{
			Company:        htmlquery.InnerText(cells[0]),
			DateOfApproval: htmlquery.InnerText(cells[2]),
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
