package w32

import (
	//"encoding/binary"
	//"errors"
	"errors"

	"syscall"
	//"unicode/utf8"
	"unsafe"
)

var (
	modShcore = syscall.NewLazyDLL("Shcore.dll")

	procGetDpiForMonitor         = modShcore.NewProc("GetDpiForMonitor")
	procGetScaleFactorForMonitor = modShcore.NewProc("GetScaleFactorForMonitor")
)

// https://docs.microsoft.com/en-us/windows/win32/api/shellscalingapi/nf-shellscalingapi-getdpiformonitor
func GetDpiForMonitor(hmonitor HMONITOR, dpiType int) (dpiX, dpiY uint32, err error) {
	ret, _, _ := procGetDpiForMonitor.Call(
		uintptr(hmonitor),
		uintptr(dpiType),
		uintptr(unsafe.Pointer(&dpiX)),
		uintptr(unsafe.Pointer(&dpiY)),
	)
	if ret != 0 {
		err = errors.New("GetDpiForMonitor failed")
	}

	return
}

// https://docs.microsoft.com/en-us/windows/win32/api/shellscalingapi/nf-shellscalingapi-getscalefactorformonitor
func GetScaleFactorForMonitor(hmonitor HMONITOR) (scale uint32, err error) {
	ret, _, _ := procGetScaleFactorForMonitor.Call(
		uintptr(hmonitor),
		uintptr(unsafe.Pointer(&scale)),
	)
	if ret != 0 {
		err = errors.New("GetScaleFactorForMonitor failed")
	}

	return
}
