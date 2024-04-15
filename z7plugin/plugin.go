package z7plugin

import (
	"github.com/lxn/win"
)

func _CreateObject(clsid win.CLSID, iid win.IID, outObject *uintptr) uint32 {
	return win.E_NOTIMPL
}

func _GetHandlerProperty(propID uint32, value *win.VARIANT) uint32 {
	return win.E_NOTIMPL
}

func _GetNumberOfFormats(numFormats *uint32) uint32 {
	return win.E_NOTIMPL
}

func _GetHandlerProperty2(formatIndex uint32, propID uint32, value *win.VARIANT) uint32 {
	return win.E_NOTIMPL
}
