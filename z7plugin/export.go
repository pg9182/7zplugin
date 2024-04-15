// Package z7plugin exposes 7-Zip's plugin interface.
//
// Must be built with `-buildmode=c-shared`. and the architecture must match
// 7-Zip. Put it in the `Formats“ directory in the install location. If the
// library doesn't load for some reason, try clicking on About in 7zFM, which
// will probably display an odd error message of some kind.
package z7plugin

import "C"

import (
	"unsafe"

	"github.com/lxn/win"
)

//export CreateObject
func CreateObject(clsid uintptr, iid uintptr, outObject uintptr) int32 {
	return int32(win.HRESULT(_CreateObject(
		*(*win.CLSID)(unsafe.Pointer(clsid)),
		*(*win.IID)(unsafe.Pointer(iid)),
		(*uintptr)(unsafe.Pointer(outObject)),
	)))
}

//export GetHandlerProperty
func GetHandlerProperty(propID uint32, value uintptr) int32 {
	return int32(win.HRESULT(_GetHandlerProperty(
		propID,
		(*win.VARIANT)(unsafe.Pointer(value)),
	)))
}

//export GetNumberOfFormats
func GetNumberOfFormats(numFormats uintptr) int32 {
	return int32(win.HRESULT(_GetNumberOfFormats(
		(*uint32)(unsafe.Pointer(numFormats)),
	)))
}

//export GetHandlerProperty2
func GetHandlerProperty2(formatIndex uint32, propID uint32, value uintptr) int32 {
	return int32(win.HRESULT(_GetHandlerProperty2(
		formatIndex,
		propID,
		(*win.VARIANT)(unsafe.Pointer(value)),
	)))
}
