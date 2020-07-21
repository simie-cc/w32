// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"syscall"
	"unsafe"
)

var (
	modpowr = syscall.NewLazyDLL("PowrProf.dll")

	procCallNtPowerInformation = modpowr.NewProc("CallNtPowerInformation")
)

func CallNtPowerInformation(informationLevel uintptr,
	inputBuffer PVOID, inputBufferLength uintptr,
	outputBuffer unsafe.Pointer, outputBufferLength uintptr) uint32 {
	ret, _, _ := procCallNtPowerInformation.Call(
		uintptr(informationLevel),
		uintptr(inputBuffer), uintptr(inputBufferLength),
		uintptr(outputBuffer), uintptr(outputBufferLength))

	return uint32(ret)
}
