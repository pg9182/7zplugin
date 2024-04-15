//go:build windows

package z7

import "github.com/lxn/win"

// CPP/7zip/ICoder.h

func Z7_IFACE_CONSTR_CODER___IID(n byte) win.IID {
	return Z7_DECL_IFACE_7ZIP___IID(4, n)
}

var (
	IID_ICompressCoder  = Z7_IFACE_CONSTR_CODER___IID(0x5)
	IID_ICompressCoder2 = Z7_IFACE_CONSTR_CODER___IID(0x18)
	IID_ICompressFilter = Z7_IFACE_CONSTR_CODER___IID(0x40)
	IID_IHasher         = Z7_IFACE_CONSTR_CODER___IID(0xC0)
)

type NModulePropID = uint32

const (
	NModulePropID_kInterfaceType NModulePropID = iota // VT_UI4
	NModulePropID_kVersion                            // VT_UI4
)

const (
	/*
	  virtual destructor in IUnknown:
	  - no  : 7-Zip (Windows)
	  - no  : 7-Zip (Linux) (v23) in default mode
	  - yes : p7zip
	  - yes : 7-Zip (Linux) before v23
	  - yes : 7-Zip (Linux) (v23), if Z7_USE_VIRTUAL_DESTRUCTOR_IN_IUNKNOWN is defined
	*/
	NModuleInterfaceType_k_IUnknown_VirtDestructor_Yes        uint32 = 1
	NModuleInterfaceType_k_IUnknown_VirtDestructor_No         uint32 = 0
	NModuleInterfaceType_k_IUnknown_VirtDestructor_ThisModule uint32 = NModuleInterfaceType_k_IUnknown_VirtDestructor_No
)
