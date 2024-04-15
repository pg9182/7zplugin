package z7plugin

import (
	"fmt"

	"github.com/lxn/win"
	"github.com/pg9182/7zplugin/winext"
	"github.com/pg9182/7zplugin/z7"
)

type ArchiveHandler struct {
	Name         string
	GUID         win.CLSID
	Extension    string
	AddExtension string
	Flags        z7.NArchive_NArcInfoFlags
	Signature    string
	EnableUpdate bool
	CreateObject func() win.IUnknown // TODO: replace with wrapped interface
}

var handlers []ArchiveHandler

func RegisterArchiveHandler(x ArchiveHandler) {
	handlers = append(handlers, x)
}

func _CreateObject(clsid win.CLSID, iid win.IID, outObject *uintptr) uint32 {
	for _, handler := range handlers {
		if handler.GUID == clsid && iid == z7.IID_IInArchive { // TODO: check iid for
			//obj := handler.CreateObject()
			win.MessageBox(0, win.SysAllocString("test"), win.SysAllocString(fmt.Sprint(iid)), win.MB_OK)
			break
		}
	}
	return win.S_OK
}

func _GetHandlerProperty(propID uint32, value *win.VARIANT) uint32 {
	return _GetHandlerProperty2(0, propID, value)
}

func _GetNumberOfFormats(numFormats *uint32) uint32 {
	*numFormats = uint32(len(handlers))
	return win.S_OK
}

func _GetHandlerProperty2(formatIndex uint32, propID uint32, value *win.VARIANT) uint32 {
	if int(formatIndex) < len(handlers) {
		switch handler := &handlers[formatIndex]; propID {
		case z7.NArchive_NHandlerPropID_kName:
			value.SetBSTR(win.SysAllocString(handler.Name))
		case z7.NArchive_NHandlerPropID_kClassID:
			value.SetBSTR(winext.SysAllocStringByteLenGUID(handler.GUID))
		case z7.NArchive_NHandlerPropID_kExtension:
			value.SetBSTR(win.SysAllocString(handler.Extension))
		case z7.NArchive_NHandlerPropID_kAddExtension:
			value.SetBSTR(win.SysAllocString(handler.AddExtension))
		case z7.NArchive_NHandlerPropID_kFlags:
			value.SetULong(uint32(handler.Flags))
		case z7.NArchive_NHandlerPropID_kUpdate:
			if handler.EnableUpdate {
				value.SetBool(1)
			} else {
				value.SetBool(0)
			}
		case z7.NArchive_NHandlerPropID_kTimeFlags:
		case z7.NArchive_NHandlerPropID_kSignature:
			value.SetBSTR(winext.SysAllocStringByteLenStr(handler.Signature))
		case z7.NArchive_NHandlerPropID_kMultiSignature:
		case z7.NArchive_NHandlerPropID_kSignatureOffset:
			value.SetULong(0)
		default:
			return uint32(win.E_FAIL)
		}
		return win.S_OK
	}
	return win.E_INVALIDARG
}
