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
	modnetapi32 = syscall.NewLazyDLL("netapi32.dll")

	procNetUserChangePassword = modnetapi32.NewProc("NetUserChangePassword")
)

const (
	NERR_Success          = 0 /* Success */
	NERR_BASE             = 2100
	NERR_UserNotFound     = (NERR_BASE + 121)
	NERR_NotPrimary       = (NERR_BASE + 126)
	NERR_PasswordTooShort = (NERR_BASE + 145)
	NERR_InvalidComputer  = (NERR_BASE + 251)
)

// https://docs.microsoft.com/en-us/windows/win32/api/lmaccess/nf-lmaccess-netuserchangepassword
// NET_API_STATUS NET_API_FUNCTION
// NetUserChangePassword (
//     _In_opt_ IN  LPCWSTR   domainname OPTIONAL,
//     _In_opt_ IN  LPCWSTR   username OPTIONAL,
//     _In_ IN  LPCWSTR   oldpassword,
//     _In_ IN  LPCWSTR   newpassword
//     );
func NetUserChangePassword(domainname, username, oldpassword, newpassword string) DWORD {
	var domainVal uintptr = uintptr(0)
	if len(domainname) > 0 {
		domainVal = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(domainname)))
	}
	ret, _, _ := procNetUserChangePassword.Call(
		domainVal,
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(username))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(oldpassword))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(newpassword))))

	return DWORD(ret)
}
