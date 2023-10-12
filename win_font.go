package main

import (
	"syscall"
	"unsafe"

	"github.com/tailscale/win"
	"golang.org/x/sys/windows"
)

type LPLOGFONT *win.LOGFONT

type CHOOSEFONT struct {
	LStructSize    uint32
	HwndOwner      win.HWND
	HDC            win.HDC
	LpLogFont      LPLOGFONT
	IPointSize     int32
	Flags          uint32
	RgbColors      win.COLORREF
	LCustData      uintptr
	LpfnHook       uintptr
	LpTemplateName *uint16
	HInstance      win.HWND
	LpszStyle      *uint16
	NFontType      uint16
	NSizeMin       int32
	NSizeMax       int32
}

// CHOOSEFONT flags
const (
	CF_APPLY                = 0x00000200
	CF_ANSIONLY             = 0x00000400
	CF_BOTH                 = 0x00000003
	CF_EFFECTS              = 0x00000100
	CF_ENABLEHOOK           = 0x00000008
	CF_ENABLETEMPLATE       = 0x00000010
	CF_ENABLETEMPLATEHANDLE = 0x00000020
	CF_FIXEDPITCHONLY       = 0x00004000
	CF_FORCEFONTEXIST       = 0x00010000
	CF_INACTIVEFONTS        = 0x02000000
	CF_INITTOLOGFONTSTRUCT  = 0x00000040
	CF_LIMITSIZE            = 0x00002000
	CF_NOOEMFONTS           = 0x00000800
	CF_NOFACESEL            = 0x00080000
	CF_NOSCRIPTSEL          = 0x00800000
	CF_NOSIMULATIONS        = 0x00001000
	CF_NOSIZESEL            = 0x00200000
	CF_NOSTYLESEL           = 0x00100000
	CF_NOVECTORFONTS        = 0x00000800
	CF_NOVERTFONTS          = 0x01000000
	CF_PRINTERFONTS         = 0x00000002
	CF_SCALABLEONLY         = 0x00020000
	CF_SCREENFONTS          = 0x00000001
	CF_SCRIPTSONLY          = 0x00000400
	CF_SELECTSCRIPT         = 0x00400000
	CF_SHOWHELP             = 0x00000004
	CF_TTONLY               = 0x00040000
	CF_USESTYLE             = 0x00000080
	CF_WYSIWYG              = 0x00008000

	SCREEN_FONTTYPE = 0x2000
)

var (
	// Library
	libcomdlg32 *windows.LazyDLL

	// Functions
	chooseFontW *windows.LazyProc
)

func init() {
	// Library
	libcomdlg32 = windows.NewLazySystemDLL("comdlg32.dll")

	// Functions
	chooseFontW = libcomdlg32.NewProc("ChooseFontW")
}

func ChooseFontW(lpcc *CHOOSEFONT) bool {
	ret, _, _ := syscall.SyscallN(chooseFontW.Addr(), uintptr(unsafe.Pointer(lpcc)))
	return ret != 0
}

func (mw MyMainWindow) ChooseFont(targetFont string) string {
	name := [32]uint16{}
	copy(name[:], windows.StringToUTF16(targetFont))

	// https://learn.microsoft.com/en-us/windows/win32/api/commdlg/ns-commdlg-choosefonta
	cfont := CHOOSEFONT{
		HwndOwner: mw.Handle(),
		LpLogFont: &win.LOGFONT{ // https://learn.microsoft.com/en-us/windows/win32/api/wingdi/ns-wingdi-logfontw
			LfFaceName:       name,
			LfCharSet:        win.ANSI_CHARSET,
			LfWidth:          0,
			LfEscapement:     0,
			LfOrientation:    0,
			LfWeight:         win.FW_NORMAL,
			LfQuality:        win.DEFAULT_QUALITY,
			LfPitchAndFamily: win.DEFAULT_PITCH | win.FF_DONTCARE,
		},
		Flags:     CF_SCREENFONTS | CF_INITTOLOGFONTSTRUCT | CF_NOSCRIPTSEL | CF_FORCEFONTEXIST,
		NFontType: SCREEN_FONTTYPE,
	}
	cfont.LStructSize = uint32(unsafe.Sizeof(cfont))

	if ChooseFontW(&cfont) {
		return windows.UTF16ToString(cfont.LpLogFont.LfFaceName[:])
	}

	return ""
}

const LF_FULLFACESIZE = 64

const (
	MM_MAX_NUMAXES      = 16
	MM_MAX_AXES_NAMELEN = 16
)
