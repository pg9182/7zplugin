package z7plugin

import (
	"github.com/lxn/win"
	"github.com/pg9182/7zplugin/winext"
	"github.com/pg9182/7zplugin/z7"
	"github.com/pg9182/7zplugin/z7plugin/internal"
)

// CPP/7zip/Archive/DllExports.cpp
// CPP/7zip/Archive/DllExports2.cpp

var CaseSensitive bool

func init() {
	internal.Archive2.CreateObject = _CreateObject
	internal.Archive2.SetCodecs = _SetCodecs
	internal.Archive2.SetLargePageMode = _SetLargePageMode
	internal.Archive2.SetCaseSensitive = _SetCaseSensitive
}

func _CreateObject(clsid win.CLSID, iid win.IID, outObject *uintptr) uint32 {
	switch *outObject = 0; iid {
	case z7.IID_ICompressCoder, z7.IID_ICompressCoder2, z7.IID_ICompressFilter:
		return _CreateCoder(clsid, iid, outObject)
	case z7.IID_IHasher:
		return _CreateHasher(clsid, iid, outObject)
	default:
		return _CreateArchiver(clsid, iid, outObject)
	}
}

func _SetCodecs(codecs uintptr) winext.HRESULT {
	_ = codecs
	return win.S_OK
}

func _SetLargePageMode() winext.HRESULT {
	// we don't do anything with this, so don't bother getting it: windows.GetLargePageMinimum()
	return win.S_OK
}

func _SetCaseSensitive(caseSensitive int32) winext.HRESULT {
	CaseSensitive = caseSensitive != 0
	return win.S_OK
}
