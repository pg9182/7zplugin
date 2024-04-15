package winext

import (
	"syscall"
	"unsafe"

	"github.com/lxn/win"
	"golang.org/x/sys/windows"
)

type GUID interface {
	syscall.GUID | windows.GUID | win.CLSID | win.IID
}

func SysAllocStringByteLenGUID[T GUID](g T) *uint16 /*BSTR*/ {
	ret, _, _ := syscall.Syscall(sysAllocStringByteLen.Addr(), 1,
		uintptr(unsafe.Pointer((*uint16)(unsafe.Pointer(&g)))),
		uintptr(16),
		0)

	return (*uint16) /*BSTR*/ (unsafe.Pointer(ret))
}

func GUIDToString[T GUID](guid T) string {
	return (windows.GUID)(guid).String()
}

func MustGUID[T GUID](s string) T {
	guid, err := windows.GUIDFromString(s)
	if err != nil {
		panic(err)
	}
	return (T)(guid)
}
