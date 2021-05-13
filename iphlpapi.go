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
	modiphlpapi = syscall.NewLazyDLL("iphlpapi.dll")

	procNotifyAddrChange = modiphlpapi.NewProc("NotifyAddrChange")
)

// https://docs.microsoft.com/en-us/windows/win32/api/iphlpapi/nf-iphlpapi-notifyaddrchange
func NotifyAddrChange(handle *HANDLE, overlapped *OVERLAPPED) DWORD {

	ret, _, _ := procNotifyAddrChange.Call(
		uintptr(unsafe.Pointer(handle)),
		uintptr(unsafe.Pointer(overlapped)),
	)

	return DWORD(ret)
}
