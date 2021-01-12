// +build windows
// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"syscall"
	"unsafe"
)

var (
	modpsapi = syscall.NewLazyDLL("psapi.dll")

	procEnumProcesses           = modpsapi.NewProc("EnumProcesses")
	procGetModuleFileNameEx     = modpsapi.NewProc("GetModuleFileNameExW")
	procGetProcessImageFileName = modpsapi.NewProc("GetProcessImageFileNameW")
)

func EnumProcesses(processIds []uint32, cb uint32, bytesReturned *uint32) bool {
	ret, _, _ := procEnumProcesses.Call(
		uintptr(unsafe.Pointer(&processIds[0])),
		uintptr(cb),
		uintptr(unsafe.Pointer(bytesReturned)))

	return ret != 0
}

// Retrieves the fully qualified path for the file containing the specified module.
// https://docs.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-getmodulefilenameexw
func GetModuleFileNameEx(hProcess HANDLE, hModule HMODULE) (string, DWORD) {

	bufLen := 1024
	buf := make([]uint16, bufLen)

	ret, _, _ := procGetModuleFileNameEx.Call(
		uintptr(hProcess),
		uintptr(hModule),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(bufLen))

	if ret <= 0 {
		return "", DWORD(ret)

	}

	return syscall.UTF16ToString(buf), DWORD(ret)
}

// https://docs.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-getprocessimagefilenamew
func GetProcessImageFileName(hProcess HANDLE) (bool, string) {
	bufLen := 1024
	buf := make([]uint16, bufLen)

	ret, _, _ := procGetProcessImageFileName.Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(bufLen))

	if ret != 0 {
		return true, syscall.UTF16ToString(buf)
	} else {
		return false, ""
	}
}
