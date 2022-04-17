package edid

import (
	"errors"
	"log"

	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/randr"
	"github.com/BurntSushi/xgb/xproto"
)

func queryEDIDX11() ([]*EDID, error) {
	X, err := xgb.NewConn()
	if err != nil {
		return nil, err
	}
	defer X.Close()
	return QueryEDIDX11(X)
}

func QueryEDIDX11(X *xgb.Conn) ([]*EDID, error) {
	if X == nil {
		return nil, errors.New(`received nil parameter`)
	}

	resps, err := queryEDIDX11RandRConn(X)
	if err != nil {
		return nil, err
	}
	var edids []*EDID
	for _, resp := range resps {
		edids = append(
			edids,
			&EDID{
				Display:  string(resp.info.Name),
				data:     resp.prop.Data,
				IsActive: resp.isActive,
				Querier:  `x11-randr`,
			},
		)
	}

	return edids, nil
}

type edidRandRResponse struct {
	resources *randr.GetScreenResourcesReply
	output    *randr.Output
	isActive  bool
	info      *randr.GetOutputInfoReply
	props     *randr.ListOutputPropertiesReply
	prop      *randr.GetOutputPropertyReply
}

// based on:
// https://github.com/burntsushi/xgb/blob/deaf085860bc/examples/randr/main.go
// https://chromium.googlesource.com/chromium/src/base/+/a3305756b9f14bb8a3d6961e79b490b8671c075d/x11/edid_parser_x11.cc

func queryEDIDX11RandRConn(X *xgb.Conn) ([]*edidRandRResponse, error) {
	if X == nil {
		return nil, errors.New(`nil parameter`)
	}

	// Every extension must be initialized before it can be used.
	err := randr.Init(X)
	if err != nil {
		return nil, err
	}

	// Get the root window on the default screen.
	root := xproto.Setup(X).DefaultScreen(X).Root

	// Gets the current screen resources. Screen resources contains a list
	// of names, crtcs, outputs and modes, among other things.
	resources, err := randr.GetScreenResources(X, root).Reply()
	if err != nil {
		return nil, err
	}

	var ret []*edidRandRResponse
	// Iterate through all of the outputs and show some of their info.
	for _, output := range resources.Outputs {
		info, err := randr.GetOutputInfo(X, output, 0).Reply()
		if err != nil {
			return nil, err
		}

		propList, err := randr.ListOutputProperties(X, output).Reply()
		if err != nil {
			log.Println(err)
			continue
		}
		var hasEDID bool
		var atomEDID xproto.Atom
		for _, atom := range propList.Atoms {
			name, err := xproto.GetAtomName(X, atom).Reply()
			if err != nil {
				continue
			}
			// /usr/include/X11/extensions/randr.h
			// #define RR_PROPERTY_RANDR_EDID "EDID"
			// "EDID", before 2008: "RANDR_EDID"
			if name.Name == `EDID` || name.Name == `RANDR_EDID` {
				hasEDID = true
				atomEDID = atom
			}
		}
		if !hasEDID {
			continue
		}

		prop, err := randr.GetOutputProperty(
			X,
			output,
			atomEDID,
			xproto.AtomAny,
			0, // offset
			// 128,// length
			256,   // length
			false, // delete
			false, // pending
		).Reply()
		if err != nil {
			log.Println(err)
			continue
		}

		ret = append(
			ret,
			&edidRandRResponse{
				resources: resources,
				output:    &output,
				isActive:  info.Connection == 0,
				info:      info,
				props:     propList,
				prop:      prop,
			},
		)
	}

	return ret, nil
}
