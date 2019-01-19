// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"syscall"
)

var (
	modwtsapi32 = syscall.NewLazyDLL("wtsapi32.dll")

	procWTSRegisterSessionNotification = modwtsapi32.NewProc("WTSRegisterSessionNotification")
)

// BOOL WTSRegisterSessionNotification(
// 	HWND  hWnd,
// 	DWORD dwFlags
//   );
func WTSRegisterSessionNotification(hwnd HWND, dwFlags DWORD) bool {
	ret, _, _ := procWTSRegisterSessionNotification.Call(
		uintptr(hwnd),
		uintptr(dwFlags),
	)

	return ret == 0
}
