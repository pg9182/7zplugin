// Command tf2vpk is a 7-Zip archive format plugin for Respawn VPKs as used in
// Titanfall 2.
package main

import (
	"github.com/lxn/win"
	"github.com/pg9182/7zplugin/winext"
	"github.com/pg9182/7zplugin/z7"
	"github.com/pg9182/7zplugin/z7plugin"
)

//go:generate go run github.com/pg9182/7zplugin/z7plugin/genver -owner pg9182 -company pg9182 -copyright "© 2024 pg9182 (github.com/pg9182)" -description "Respawn VPK Plugin for 7-Zip"

// go generate ./plugins/tf2vpk
// CGO_ENABLED=1 GOOS=windows GOARCH=amd64 CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -ldflags '-s -w -extldflags=-static' -trimpath -v -x -o tf2vpk64.dll ./plugins/tf2vpk
// CGO_ENABLED=1 GOOS=windows GOARCH=386 CC=i686-w64-mingw32-gcc go build -buildmode=c-shared -ldflags '-s -w -extldflags=-static' -trimpath -v -x -o tf2vpk32.dll ./plugins/tf2vpk
// put the correct one in Program Files/7-Zip/Formats/ (note: arch must match)

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

func main() {}
