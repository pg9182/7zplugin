//go:build windows

package z7

import "github.com/lxn/win"

// CPP/7zip/IDecl.h

func Z7_DECL_IFACE_7ZIP___IID(groupID, subID byte) win.IID {
	return win.IID{0x23170F69, 0x40C1, 0x278A, [8]byte{0, 0, 0, groupID, 0, subID, 0, 0}}
}
