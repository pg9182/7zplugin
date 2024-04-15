package z7

import "github.com/lxn/win"

var IID_IInArchive = win.IID{0x23170F69, 0x40C1, 0x278A, [8]byte{0, 0, 0, 6, 0, 0x60, 0, 0}}

type NArchive_NHandlerPropID = uint32

const (
	NArchive_NHandlerPropID_kName            NArchive_NHandlerPropID = iota // VT_BSTR
	NArchive_NHandlerPropID_kClassID                                        // binary GUID in VT_BSTR
	NArchive_NHandlerPropID_kExtension                                      // VT_BSTR
	NArchive_NHandlerPropID_kAddExtension                                   // VT_BSTR
	NArchive_NHandlerPropID_kUpdate                                         // VT_BOOL
	NArchive_NHandlerPropID_kKeepName                                       // VT_BOOL
	NArchive_NHandlerPropID_kSignature                                      // binary in VT_BSTR
	NArchive_NHandlerPropID_kMultiSignature                                 // binary in VT_BSTR
	NArchive_NHandlerPropID_kSignatureOffset                                // VT_UI4
	NArchive_NHandlerPropID_kAltStreams                                     // VT_BOOL
	NArchive_NHandlerPropID_kNtSecure                                       // VT_BOOL
	NArchive_NHandlerPropID_kFlags                                          // VT_UI4
	NArchive_NHandlerPropID_kTimeFlags                                      // VT_UI4
)

type NArchive_NArcInfoFlags uint32

const (
	NArchive_NArcInfoFlags_kKeepName        NArchive_NArcInfoFlags = 1 << 0  // keep name of file in archive name
	NArchive_NArcInfoFlags_kAltStreams      NArchive_NArcInfoFlags = 1 << 1  // the handler supports alt streams
	NArchive_NArcInfoFlags_kNtSecure        NArchive_NArcInfoFlags = 1 << 2  // the handler supports NT security
	NArchive_NArcInfoFlags_kFindSignature   NArchive_NArcInfoFlags = 1 << 3  // the handler can find start of archive
	NArchive_NArcInfoFlags_kMultiSignature  NArchive_NArcInfoFlags = 1 << 4  // there are several signatures
	NArchive_NArcInfoFlags_kUseGlobalOffset NArchive_NArcInfoFlags = 1 << 5  // the seek position of stream must be set as global offset
	NArchive_NArcInfoFlags_kStartOpen       NArchive_NArcInfoFlags = 1 << 6  // call handler for each start position
	NArchive_NArcInfoFlags_kPureStartOpen   NArchive_NArcInfoFlags = 1 << 7  // call handler only for start of file
	NArchive_NArcInfoFlags_kBackwardOpen    NArchive_NArcInfoFlags = 1 << 8  // archive can be open backward
	NArchive_NArcInfoFlags_kPreArc          NArchive_NArcInfoFlags = 1 << 9  // such archive can be stored before real archive (like SFX stub)
	NArchive_NArcInfoFlags_kSymLinks        NArchive_NArcInfoFlags = 1 << 10 // the handler supports symbolic links
	NArchive_NArcInfoFlags_kHardLinks       NArchive_NArcInfoFlags = 1 << 11 // the handler supports hard links
	NArchive_NArcInfoFlags_kByExtOnlyOpen   NArchive_NArcInfoFlags = 1 << 12 // call handler only if file extension matches
	NArchive_NArcInfoFlags_kHashHandler     NArchive_NArcInfoFlags = 1 << 13 // the handler contains the hashes (checksums)
	NArchive_NArcInfoFlags_kCTime           NArchive_NArcInfoFlags = 1 << 14
	NArchive_NArcInfoFlags_kCTime_Default   NArchive_NArcInfoFlags = 1 << 15
	NArchive_NArcInfoFlags_kATime           NArchive_NArcInfoFlags = 1 << 16
	NArchive_NArcInfoFlags_kATime_Default   NArchive_NArcInfoFlags = 1 << 17
	NArchive_NArcInfoFlags_kMTime           NArchive_NArcInfoFlags = 1 << 18
	NArchive_NArcInfoFlags_kMTime_Default   NArchive_NArcInfoFlags = 1 << 19
)
