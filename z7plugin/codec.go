//go:build windows

package z7plugin

import (
	"github.com/lxn/win"
	"github.com/pg9182/7zplugin/winext"
	"github.com/pg9182/7zplugin/z7"
	"github.com/pg9182/7zplugin/z7plugin/internal"
	"golang.org/x/sys/windows"
)

// CPP/7zip/Compress/CodecExports.cpp

func init() {
	internal.Archive2.GetNumberOfMethods = _GetNumberOfMethods
	internal.Archive2.GetMethodProperty = _GetMethodProperty
	internal.Archive2.CreateDecoder = _CreateDecoder
	internal.Archive2.CreateEncoder = _CreateEncoder
	internal.Archive2.GetHashers = _GetHashers
	internal.Archive2.GetModuleProp = _GetModuleProp
}

// TODO: these are stubs

func _CreateCoder(clsid win.CLSID, iid win.IID, outObject *uintptr) winext.HRESULT {
	return winext.HRESULT(windows.CLASS_E_CLASSNOTAVAILABLE)
}

func _GetNumberOfMethods(numCodecs *uint32) winext.HRESULT {
	*numCodecs = 0
	return win.S_OK
}

func _GetMethodProperty(codecIndex uint32, propID winext.PROPID, value *winext.PROPVARIANT) winext.HRESULT {
	value.Vt = win.VT_EMPTY
	return win.S_OK
}

func _CreateCoder2(encode bool, index uint32, iid win.IID, outObject *uintptr) winext.HRESULT {
	return winext.HRESULT(windows.CLASS_E_CLASSNOTAVAILABLE)
}

func _CreateDecoder(index uint32, iid win.IID, outObject *uintptr) winext.HRESULT {
	return _CreateCoder2(false, index, iid, outObject)
}

func _CreateEncoder(index uint32, iid win.IID, outObject *uintptr) winext.HRESULT {
	return _CreateCoder2(true, index, iid, outObject)
}

func _GetHashers(hashers *uintptr) winext.HRESULT {
	// note: this is safe: CPP/7zip/UI/Common/LoadCodecs.cpp
	//
	// MY_GET_FUNC_LOC (getHashers, Func_GetHashers, lib.Lib, "GetHashers")
	// if (getHashers)
	//	{
	//	  RINOK(getHashers(&lib.ComHashers))
	//	  if (lib.ComHashers)
	//
	*hashers = 0
	return win.S_OK
}

func _GetModuleProp(propID winext.PROPID, value *winext.PROPVARIANT) winext.HRESULT {
	value.Vt = win.VT_EMPTY
	switch propID {
	case z7.NModulePropID_kInterfaceType:
		value.SetULong(z7.NModuleInterfaceType_k_IUnknown_VirtDestructor_ThisModule)
	case z7.NModulePropID_kVersion:
		value.SetULong((z7.MY_VER_MAJOR << 16) | z7.MY_VER_MINOR)
	}
	return win.S_OK
}

func _CreateHasher(clsid win.CLSID, iid win.IID, outObject *uintptr) winext.HRESULT {
	return winext.HRESULT(windows.CLASS_E_CLASSNOTAVAILABLE)
}
