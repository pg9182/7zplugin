// Package winext is like github.com/lxn/win, but has some additional functions.
package winext

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

type PROPID = uint32
type PROPVARIANT = win.VARIANT // TODO: is this the same as PROPVARIANT internally?
type HRESULT = uint32          // note: actually an int32, but the constants in the win pkg are untyped and overflow its HRESULT...

var (
	liboleaut32 = windows.NewLazySystemDLL("oleaut32.dll")

	sysAllocStringByteLen = liboleaut32.NewProc("SysAllocStringByteLen")
)

func SysAllocStringByteLen(b []byte) *uint16 /*BSTR*/ {
	ret, _, _ := syscall.Syscall(sysAllocStringByteLen.Addr(), 1,
		uintptr(unsafe.Pointer((*uint16)(unsafe.Pointer(unsafe.SliceData(b))))),
		uintptr(len(b)),
		0)

	return (*uint16) /*BSTR*/ (unsafe.Pointer(ret))
}

func SysAllocStringByteLenStr(b string) *uint16 /*BSTR*/ {
	ret, _, _ := syscall.Syscall(sysAllocStringByteLen.Addr(), 1,
		uintptr(unsafe.Pointer((*uint16)(unsafe.Pointer(unsafe.StringData(b))))),
		uintptr(len(b)),
		0)

	return (*uint16) /*BSTR*/ (unsafe.Pointer(ret))
}
