// Command tf2vpk is a 7-Zip archive format plugin for Respawn VPKs as used in
// Titanfall 2.
package main

import (
	"github.com/lxn/win"
	"github.com/pg9182/7zplugin/winext"
	"github.com/pg9182/7zplugin/z7"
	"github.com/pg9182/7zplugin/z7plugin"
)

// GOOS=windows GOARCH=amd64 CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc go build -buildmode=c-shared -ldflags '-extldflags=-static' -o tf2vpk.dll ./plugins/tf2vpk
// put in Program Files/7-Zip/Formats/ (note: arch must match)

func init() {
	z7plugin.RegisterArchiveHandler(z7plugin.ArchiveHandler{
		Name:         "VPK0203",
		GUID:         winext.MustGUID[win.CLSID]("{3a128a09-88fe-45db-8727-565dff106ebe}"),
		Extension:    "vpk",
		AddExtension: "",
		Flags:        z7.NArchive_NArcInfoFlags_kPureStartOpen,
		Signature:    "\x55\xaa\x12\x34\x02\x00\x03\x00",
		EnableUpdate: false,
		CreateObject: nil,
	})
}

func main() {}
