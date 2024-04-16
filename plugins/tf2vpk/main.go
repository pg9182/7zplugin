// Package tf2vpk is a 7-Zip archive format plugin for Respawn VPKs as used in
// Titanfall 2.
package tf2vpk

import (
	"github.com/lxn/win"
	"github.com/pg9182/7zplugin/winext"
	"github.com/pg9182/7zplugin/z7"
	"github.com/pg9182/7zplugin/z7plugin"
)

func init() {
	z7plugin.RegisterArc(&z7plugin.CArcInfo{
		Name:            "VPK0203",
		CLSID:           winext.MustGUID[win.CLSID]("{3a128a09-88fe-45db-8727-565dff106ebe}"),
		Ext:             "vpk",
		AddExt:          "",
		Flags:           z7.NArchive_NArcInfoFlags_kPureStartOpen,
		Signature:       "\x55\xaa\x12\x34\x02\x00\x03\x00",
		CreateInArchive: func() {}, // TODO
	})
}
