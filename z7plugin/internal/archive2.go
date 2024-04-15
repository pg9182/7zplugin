//go:build windows

package internal

import "C"

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"github.com/pg9182/7zplugin/winext"
	"github.com/pg9182/7zplugin/z7"
)

// CPP/Common/MyWindows.h
// NOTE: calling convention for exported API functions is stdcall

// CPP/7zip/Archive/Archive2.def
var Archive2 struct {

	// CPP/7zip/Archive/DllExports2.cpp
	CreateObject func(clsid win.CLSID, iid win.IID, outObject *uintptr) winext.HRESULT

	// CPP/7zip/Archive/ArchiveExports.cpp
	GetHandlerProperty  func(propID winext.PROPID, value *winext.PROPVARIANT) winext.HRESULT
	GetNumberOfFormats  func(numFormats *uint32) winext.HRESULT
	GetHandlerProperty2 func(formatIndex uint32, propID winext.PROPID, value *winext.PROPVARIANT) winext.HRESULT
	GetIsArc            func(formatIndex uint32, isArc *Func_IsArc) winext.HRESULT

	// CPP/7zip/Compress/CodecExports.cpp
	GetNumberOfMethods func(numCodecs *uint32) winext.HRESULT
	GetMethodProperty  func(codecIndex uint32, propID winext.PROPID, value *winext.PROPVARIANT) winext.HRESULT
	CreateDecoder      func(index uint32, iid win.IID, outObject *uintptr) winext.HRESULT
	CreateEncoder      func(index uint32, iid win.IID, outObject *uintptr) winext.HRESULT

	// CPP/7zip/Compress/CodecExports.cpp
	GetHashers func(hashers *uintptr) winext.HRESULT

	// CPP/7zip/Archive/DllExports2.cpp
	// NOTE: we don't need to support this -- it's only if you need to import codecs from other plugins
	SetCodecs func(codecs uintptr) winext.HRESULT // codecs is a pointer to a ICompressCodecsInfo

	// CPP/7zip/Archive/DllExports.cpp
	SetLargePageMode func() winext.HRESULT
	SetCaseSensitive func(caseSensitive int32) winext.HRESULT

	// CPP/7zip/Compress/CodecExports.cpp
	GetModuleProp func(propID winext.PROPID, value *winext.PROPVARIANT) winext.HRESULT
}

type Func_IsArc uintptr // syscall.NewCallback: func(p *byte, size uintptr) HRESULT

func Func_IsArc_Wrap(fn func(b []byte) z7.NArchive_k_IsArc_Res) Func_IsArc {
	if fn == nil {
		return 0
	}
	return Func_IsArc(syscall.NewCallback(func(p *byte, size uintptr) z7.NArchive_k_IsArc_Res {
		return fn(unsafe.Slice(p, int(size)))
	}))
}

// STDAPI CreateObject(const GUID *clsid, const GUID *iid, void **outObject);
//
//export CreateObject
func CreateObject(clsid uintptr, iid uintptr, outObject uintptr) int32 {
	return int32(win.HRESULT(Archive2.CreateObject(
		*(*win.CLSID)(unsafe.Pointer(clsid)),
		*(*win.IID)(unsafe.Pointer(iid)),
		(*uintptr)(unsafe.Pointer(outObject)),
	)))
}

// STDAPI GetHandlerProperty(PROPID propID, PROPVARIANT *value);
//
//export GetHandlerProperty
func GetHandlerProperty(propID uint32, value uintptr) int32 {
	return int32(win.HRESULT(Archive2.GetHandlerProperty(
		propID,
		(*winext.PROPVARIANT)(unsafe.Pointer(value)),
	)))
}

// STDAPI GetNumberOfFormats(UINT32 *numFormats);
//
//export GetNumberOfFormats
func GetNumberOfFormats(numFormats uintptr) int32 {
	return int32(win.HRESULT(Archive2.GetNumberOfFormats(
		(*uint32)(unsafe.Pointer(numFormats)),
	)))
}

// STDAPI GetHandlerProperty2(UInt32 formatIndex, PROPID propID, PROPVARIANT *value);
//
//export GetHandlerProperty2
func GetHandlerProperty2(formatIndex uint32, propID uint32, value uintptr) int32 {
	return int32(win.HRESULT(Archive2.GetHandlerProperty2(
		formatIndex,
		propID,
		(*winext.PROPVARIANT)(unsafe.Pointer(value)),
	)))
}

// STDAPI GetIsArc(UInt32 formatIndex, Func_IsArc *isArc);
//
//export GetIsArc
func GetIsArc(formatIndex uint32, isArc uintptr) int32 {
	return int32(win.HRESULT(Archive2.GetIsArc(
		formatIndex,
		(*Func_IsArc)(unsafe.Pointer(isArc)),
	)))
}

// STDAPI GetNumberOfMethods(UInt32 *numCodecs);
//
//export GetNumberOfMethods
func GetNumberOfMethods(numCodecs uintptr) int32 {
	return int32(win.HRESULT(Archive2.GetNumberOfMethods(
		(*uint32)(unsafe.Pointer(numCodecs)),
	)))
}

// STDAPI GetMethodProperty(UInt32 codecIndex, PROPID propID, PROPVARIANT *value);
//
//export GetMethodProperty
func GetMethodProperty(codecIndex uint32, propID uint32, value uintptr) int32 {
	return int32(win.HRESULT(Archive2.GetMethodProperty(
		codecIndex,
		propID,
		(*winext.PROPVARIANT)(unsafe.Pointer(value)),
	)))
}

// STDAPI CreateDecoder(UInt32 index, const GUID *iid, void **outObject);
//
//export CreateDecoder
func CreateDecoder(index uint32, iid uintptr, outObject uintptr) int32 {
	return int32(win.HRESULT(Archive2.CreateDecoder(
		index,
		*(*win.IID)(unsafe.Pointer(iid)),
		(*uintptr)(unsafe.Pointer(outObject)),
	)))
}

// STDAPI CreateEncoder(UInt32 index, const GUID *iid, void **outObject);
//
//export CreateEncoder
func CreateEncoder(index uint32, iid uintptr, outObject uintptr) int32 {
	return int32(win.HRESULT(Archive2.CreateEncoder(
		index,
		*(*win.IID)(unsafe.Pointer(iid)),
		(*uintptr)(unsafe.Pointer(outObject)),
	)))
}

// STDAPI GetHashers(IHashers **hashers);
//
//export GetHashers
func GetHashers(hashers uintptr) int32 {
	return int32(win.HRESULT(Archive2.GetHashers(
		(*uintptr)(unsafe.Pointer(hashers)),
	)))
}

// STDAPI SetLargePageMode();
//
//export SetLargePageMode
func SetLargePageMode() int32 {
	return int32(win.HRESULT(Archive2.SetLargePageMode()))
}

// STDAPI SetCaseSensitive(Int32 caseSensitive);
//
//export SetCaseSensitive
func SetCaseSensitive(caseSensitive int32) int32 {
	return int32(win.HRESULT(Archive2.SetCaseSensitive(
		caseSensitive,
	)))
}

// STDAPI GetModuleProp(PROPID propID, PROPVARIANT *value);
//
//export GetModuleProp
func GetModuleProp(propID uint32, value uintptr) int32 {
	return int32(win.HRESULT(Archive2.GetModuleProp(
		propID,
		(*winext.PROPVARIANT)(unsafe.Pointer(value)),
	)))
}
