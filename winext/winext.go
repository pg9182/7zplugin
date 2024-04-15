// Package winext is like github.com/lxn/win, but has some additional functions.
package winext

import (
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

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
