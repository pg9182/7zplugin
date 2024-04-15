package z7plugin

import (
	"github.com/lxn/win"
	"github.com/pg9182/7zplugin/winext"
	"github.com/pg9182/7zplugin/z7"
	"github.com/pg9182/7zplugin/z7plugin/internal"
	"golang.org/x/sys/windows"
)

// CPP/7zip/Archive/ArchiveExports.cpp

func init() {
	internal.Archive2.CreateObject = _CreateObject
	internal.Archive2.GetHandlerProperty = _GetHandlerProperty
	internal.Archive2.GetNumberOfFormats = _GetNumberOfFormats
	internal.Archive2.GetHandlerProperty2 = _GetHandlerProperty2
	internal.Archive2.GetIsArc = _GetIsArc
}

type CArcInfo struct {
	Flags            z7.NArchive_NArcInfoFlags
	CLSID            win.CLSID
	SignatureOffset  uint16
	Signature        string
	Name             string
	Ext              string
	AddExt           string
	TimeFlags        uint32
	CreateInArchive  func() // TODO: accept interface, internally allocate cgo mem like { *InArchiveVtbl { syscall.NewCallback(inArchiveVtblThunkXXX)... } cgo.Handle(iface) }
	CreateOutArchive func() // TODO: accept interface, internally allocate cgo mem like { *OutArchiveVtbl { syscall.NewCallback(inArchiveVtblThunkXXX)... } cgo.Handle(iface) }
	IsArc            func(b []byte) z7.NArchive_k_IsArc_Res
}

func (arcInfo CArcInfo) IsMultiSignature() bool {
	return arcInfo.Flags&z7.NArchive_NArcInfoFlags_kMultiSignature != 0
}

var _Arcs []*CArcInfo
var isArcCache []internal.Func_IsArc

func RegisterArc(arcInfo *CArcInfo) {
	_Arcs = append(_Arcs, arcInfo)
	isArcCache = append(isArcCache, internal.Func_IsArc_Wrap(arcInfo.IsArc))
}

func _CreateArchiver(clsid win.CLSID, iid win.IID, outObject *uintptr) uint32 {
	var (
		needIn  = iid == z7.IID_IInArchive
		needOut = iid == z7.IID_IOutArchive
	)
	if !needIn && !needOut {
		return win.E_NOINTERFACE
	}
	for _, arc := range _Arcs {
		if arc.CLSID == clsid {
			if needIn && arc.CreateInArchive != nil {
				win.MessageBox(0, win.SysAllocString(winext.GUIDToString(clsid)), win.SysAllocString("TODO"), win.MB_OK) // TODO: actually create it
				// *outObject = ...
				// AddRef()
				return win.S_OK
			}
			if needOut {
				win.MessageBox(0, win.SysAllocString(winext.GUIDToString(clsid)), win.SysAllocString("TODO"), win.MB_OK) // TODO: actually create it
				// *outObject = ...
				// AddRef()
				return win.S_OK
			}
		}
	}
	return winext.HRESULT(windows.CLASS_E_CLASSNOTAVAILABLE)
}

func _GetHandlerProperty(propID uint32, value *win.VARIANT) uint32 {
	return _GetHandlerProperty2(0, propID, value)
}

func _GetNumberOfFormats(numFormats *uint32) uint32 {
	*numFormats = uint32(len(_Arcs))
	return win.S_OK
}

func _GetHandlerProperty2(formatIndex uint32, propID uint32, value *win.VARIANT) uint32 {
	value.Vt = win.VT_EMPTY
	if int(formatIndex) >= len(_Arcs) {
		return win.E_INVALIDARG
	}
	switch arc := _Arcs[formatIndex]; propID {
	case z7.NArchive_NHandlerPropID_kName:
		value.SetBSTR(win.SysAllocString(arc.Name))
	case z7.NArchive_NHandlerPropID_kClassID:
		value.SetBSTR(winext.SysAllocStringByteLenGUID(arc.CLSID))
	case z7.NArchive_NHandlerPropID_kExtension:
		value.SetBSTR(win.SysAllocString(arc.Ext))
	case z7.NArchive_NHandlerPropID_kAddExtension:
		value.SetBSTR(win.SysAllocString(arc.AddExt))
	case z7.NArchive_NHandlerPropID_kUpdate:
		if arc.CreateOutArchive != nil {
			value.SetBool(win.VARIANT_TRUE)
		} else {
			value.SetBool(win.VARIANT_FALSE)
		}
	case z7.NArchive_NHandlerPropID_kKeepName:
		if arc.Flags&z7.NArchive_NArcInfoFlags_kKeepName != 0 {
			value.SetBool(win.VARIANT_TRUE)
		} else {
			value.SetBool(win.VARIANT_FALSE)
		}
	case z7.NArchive_NHandlerPropID_kAltStreams:
		if arc.Flags&z7.NArchive_NArcInfoFlags_kAltStreams != 0 {
			value.SetBool(win.VARIANT_TRUE)
		} else {
			value.SetBool(win.VARIANT_FALSE)
		}
	case z7.NArchive_NHandlerPropID_kNtSecure:
		if arc.Flags&z7.NArchive_NArcInfoFlags_kNtSecure != 0 {
			value.SetBool(win.VARIANT_TRUE)
		} else {
			value.SetBool(win.VARIANT_FALSE)
		}
	case z7.NArchive_NHandlerPropID_kFlags:
		value.SetULong(uint32(arc.Flags))
	case z7.NArchive_NHandlerPropID_kTimeFlags:
		value.SetULong(uint32(arc.TimeFlags))
	case z7.NArchive_NHandlerPropID_kSignatureOffset:
		value.SetULong(uint32(arc.SignatureOffset))
	case z7.NArchive_NHandlerPropID_kSignature:
		if len(arc.Signature) != 0 && !arc.IsMultiSignature() {
			value.SetBSTR(winext.SysAllocStringByteLenStr(arc.Signature))
		}
	case z7.NArchive_NHandlerPropID_kMultiSignature:
		if len(arc.Signature) != 0 && arc.IsMultiSignature() {
			value.SetBSTR(winext.SysAllocStringByteLenStr(arc.Signature))
		}
	}
	return win.S_OK
}

func _GetIsArc(formatIndex uint32, isArc *internal.Func_IsArc) winext.HRESULT {
	*isArc = 0
	if int(formatIndex) >= len(_Arcs) {
		return win.E_INVALIDARG
	}
	*isArc = internal.Func_IsArc(isArcCache[formatIndex]) // note: can be null
	return win.S_OK
}
