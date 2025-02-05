package main

import (
	"log"
	"syscall"
	"unsafe"

	"golang.org/x/sys/windows"
)

var (
	// Library
	libgdi32 *windows.LazyDLL

	// Functions
	getFontResourceInfoW *windows.LazyProc
)

func init() {
	// Library
	libgdi32 = windows.NewLazySystemDLL("gdi32.dll")

	// Functions
	getFontResourceInfoW = libgdi32.NewProc("GetFontResourceInfoW")
}

func GetFontResourceInfoW(str string, size *uint32, buffer uintptr, aType int) bool { // DWORD
	strStr, err := windows.UTF16FromString(str)
	if err != nil {
		log.Println(err)
	}

	ret1, _, _ := syscall.SyscallN(getFontResourceInfoW.Addr(),
		uintptr(unsafe.Pointer(&strStr[0])),
		uintptr(unsafe.Pointer(size)),
		buffer,
		uintptr(aType))
	return ret1 != 0
}

const QFR_DESCRIPTION = 1

func GetFontResourceInfo(font string) string {
	var size uint32 = 0
	if !GetFontResourceInfoW(font, &size, 0, QFR_DESCRIPTION) {
		return ""
	}

	buff := make([]uint16, size/2)
	if GetFontResourceInfoW(font, &size, uintptr(unsafe.Pointer(&buff[0])), QFR_DESCRIPTION) {
		return syscall.UTF16ToString(buff)
	}

	return ""
}
