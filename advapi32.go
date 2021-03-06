// +build windows
// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"errors"
	"fmt"
	"os"
	"runtime"
	"syscall"
	"unsafe"
)

var (
	modadvapi32 = syscall.NewLazyDLL("advapi32.dll")

	//	procRegSetKeyValue     = modadvapi32.NewProc("RegSetKeyValueW")
	procCloseEventLog                = modadvapi32.NewProc("CloseEventLog")
	procCloseServiceHandle           = modadvapi32.NewProc("CloseServiceHandle")
	procControlService               = modadvapi32.NewProc("ControlService")
	procControlTrace                 = modadvapi32.NewProc("ControlTraceW")
	procInitializeSecurityDescriptor = modadvapi32.NewProc("InitializeSecurityDescriptor")
	procLogonUserW                   = modadvapi32.NewProc("LogonUserW")
	procOpenEventLog                 = modadvapi32.NewProc("OpenEventLogW")
	procOpenSCManager                = modadvapi32.NewProc("OpenSCManagerW")
	procOpenService                  = modadvapi32.NewProc("OpenServiceW")
	procReadEventLog                 = modadvapi32.NewProc("ReadEventLogW")
	procRegCloseKey                  = modadvapi32.NewProc("RegCloseKey")
	procRegCreateKeyEx               = modadvapi32.NewProc("RegCreateKeyExW")
	procRegDeleteKeyValue            = modadvapi32.NewProc("RegDeleteKeyValueW")
	procRegDeleteTree                = modadvapi32.NewProc("RegDeleteTreeW")
	procRegDeleteValue               = modadvapi32.NewProc("RegDeleteValueW")
	procRegEnumKeyEx                 = modadvapi32.NewProc("RegEnumKeyExW")
	procRegGetValue                  = modadvapi32.NewProc("RegGetValueW")
	procRegOpenKeyEx                 = modadvapi32.NewProc("RegOpenKeyExW")
	procRegSetValueEx                = modadvapi32.NewProc("RegSetValueExW")
	procSetSecurityDescriptorDacl    = modadvapi32.NewProc("SetSecurityDescriptorDacl")
	procStartService                 = modadvapi32.NewProc("StartServiceW")
	procStartTrace                   = modadvapi32.NewProc("StartTraceW")
)

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa379583(v=vs.85).aspx
func SetSecurityDescriptorDacl(pSecurityDescriptor *SECURITY_DESCRIPTOR, pDacl *ACL) (e error) {

	if pSecurityDescriptor == nil {
		return errors.New("null descriptor")
	}

	var ret uintptr
	if pDacl == nil {
		ret, _, _ = procSetSecurityDescriptorDacl.Call(
			uintptr(unsafe.Pointer(pSecurityDescriptor)),
			uintptr(1), // DaclPresent
			uintptr(0), // pDacl
			uintptr(0), // DaclDefaulted
		)
	} else {
		ret, _, _ = procSetSecurityDescriptorDacl.Call(
			uintptr(unsafe.Pointer(pSecurityDescriptor)),
			uintptr(1), // DaclPresent
			uintptr(unsafe.Pointer(pDacl)),
			uintptr(0), //DaclDefaulted
		)
	}

	if ret != 0 {
		return
	}
	e = syscall.GetLastError()
	return
}

func ControlTrace(hTrace TRACEHANDLE, lpSessionName string, props *EVENT_TRACE_PROPERTIES, dwControl uint32) (success bool, e error) {

	ret, _, _ := procControlTrace.Call(
		uintptr(unsafe.Pointer(hTrace)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpSessionName))),
		uintptr(unsafe.Pointer(props)),
		uintptr(dwControl))

	if ret == ERROR_SUCCESS {
		return true, nil
	}
	e = errors.New(fmt.Sprintf("error: 0x%x", ret))
	return
}

func StartTrace(lpSessionName string, props *EVENT_TRACE_PROPERTIES) (hTrace TRACEHANDLE, e error) {

	ret, _, _ := procStartTrace.Call(
		uintptr(unsafe.Pointer(&hTrace)),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpSessionName))),
		uintptr(unsafe.Pointer(props)))

	if ret == ERROR_SUCCESS {
		return
	}
	e = errors.New(fmt.Sprintf("error: 0x%x", ret))
	return
}

// http://msdn.microsoft.com/en-us/library/windows/desktop/aa378863(v=vs.85).aspx
func InitializeSecurityDescriptor(rev uint16) (pSecurityDescriptor *SECURITY_DESCRIPTOR, e error) {

	pSecurityDescriptor = &SECURITY_DESCRIPTOR{}

	ret, _, _ := procInitializeSecurityDescriptor.Call(
		uintptr(unsafe.Pointer(pSecurityDescriptor)),
		uintptr(rev),
	)

	if ret != 0 {
		return
	}
	e = syscall.GetLastError()
	return
}

func RegCreateKey(hKey HKEY, subKey string) HKEY {
	var result HKEY
	ret, _, _ := procRegCreateKeyEx.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(0),
		uintptr(0),
		uintptr(0),
		uintptr(KEY_ALL_ACCESS),
		uintptr(0),
		uintptr(unsafe.Pointer(&result)),
		uintptr(0))
	_ = ret
	return result
}

func RegOpenKeyEx(hKey HKEY, subKey string, samDesired uint32) HKEY {
	var result HKEY
	ret, _, _ := procRegOpenKeyEx.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(0),
		uintptr(samDesired),
		uintptr(unsafe.Pointer(&result)))

	if ret != ERROR_SUCCESS {
		panic(fmt.Sprintf("RegOpenKeyEx(%d, %s, %d) failed", hKey, subKey, samDesired))
	}
	return result
}

func RegCloseKey(hKey HKEY) error {
	var err error
	ret, _, _ := procRegCloseKey.Call(
		uintptr(hKey))

	if ret != ERROR_SUCCESS {
		err = errors.New("RegCloseKey failed")
	}
	return err
}

func RegGetRaw(hKey HKEY, subKey string, value string) []byte {
	var bufLen uint32
	var valptr unsafe.Pointer
	if len(value) > 0 {
		valptr = unsafe.Pointer(syscall.StringToUTF16Ptr(value))
	}
	procRegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(valptr),
		uintptr(RRF_RT_ANY),
		0,
		0,
		uintptr(unsafe.Pointer(&bufLen)))

	if bufLen == 0 {
		return nil
	}

	buf := make([]byte, bufLen)
	ret, _, _ := procRegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(valptr),
		uintptr(RRF_RT_ANY),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)))

	if ret != ERROR_SUCCESS {
		return nil
	}

	return buf
}

func RegSetBinary(hKey HKEY, subKey string, value []byte) (errno int) {
	var lptr, vptr unsafe.Pointer
	if len(subKey) > 0 {
		lptr = unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))
	}
	if len(value) > 0 {
		vptr = unsafe.Pointer(&value[0])
	}
	ret, _, _ := procRegSetValueEx.Call(
		uintptr(hKey),
		uintptr(lptr),
		uintptr(0),
		uintptr(REG_BINARY),
		uintptr(vptr),
		uintptr(len(value)))

	return int(ret)
}

func RegSetString(hKey HKEY, subKey string, value string) (errno int) {
	var lptr, vptr unsafe.Pointer
	if len(subKey) > 0 {
		lptr = unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))
	}
	var buf []uint16
	if len(value) > 0 {
		buf, err := syscall.UTF16FromString(value)
		if err != nil {
			return ERROR_BAD_FORMAT
		}
		vptr = unsafe.Pointer(&buf[0])
	}
	ret, _, _ := procRegSetValueEx.Call(
		uintptr(hKey),
		uintptr(lptr),
		uintptr(0),
		uintptr(REG_SZ),
		uintptr(vptr),
		uintptr(unsafe.Sizeof(buf)+2)) // 2 is the size of the terminating null character

	return int(ret)
}

func RegSetUint32(hKey HKEY, subKey string, value uint32) (errno int) {
	var lptr unsafe.Pointer
	if len(subKey) > 0 {
		lptr = unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))
	}
	vptr := unsafe.Pointer(&value)
	ret, _, _ := procRegSetValueEx.Call(
		uintptr(hKey),
		uintptr(lptr),
		uintptr(0),
		uintptr(REG_DWORD),
		uintptr(vptr),
		uintptr(unsafe.Sizeof(value)))

	return int(ret)
}

func RegGetString(hKey HKEY, subKey string, value string) string {
	var bufLen uint32
	procRegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
		uintptr(RRF_RT_REG_SZ),
		0,
		0,
		uintptr(unsafe.Pointer(&bufLen)))

	if bufLen == 0 {
		return ""
	}

	buf := make([]uint16, bufLen)
	ret, _, _ := procRegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
		uintptr(RRF_RT_REG_SZ),
		0,
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)))

	if ret != ERROR_SUCCESS {
		return ""
	}

	return syscall.UTF16ToString(buf)
}

func RegGetUint32(hKey HKEY, subKey string, value string) (data uint32, errno int) {
	var dataLen uint32 = uint32(unsafe.Sizeof(data))
	ret, _, _ := procRegGetValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(value))),
		uintptr(RRF_RT_REG_DWORD),
		0,
		uintptr(unsafe.Pointer(&data)),
		uintptr(unsafe.Pointer(&dataLen)))
	errno = int(ret)
	return
}

/*
func RegSetKeyValue(hKey HKEY, subKey string, valueName string, dwType uint32, data uintptr, cbData uint16) (errno int) {
	ret, _, _ := procRegSetKeyValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(valueName))),
		uintptr(dwType),
		data,
		uintptr(cbData))

	return int(ret)
}
*/

func RegDeleteKeyValue(hKey HKEY, subKey string, valueName string) (errno int) {
	ret, _, _ := procRegDeleteKeyValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(valueName))))

	return int(ret)
}

func RegDeleteValue(hKey HKEY, valueName string) (errno int) {
	ret, _, _ := procRegDeleteValue.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(valueName))))

	return int(ret)
}

func RegDeleteTree(hKey HKEY, subKey string) (errno int) {
	ret, _, _ := procRegDeleteTree.Call(
		uintptr(hKey),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(subKey))))

	return int(ret)
}

func RegEnumKeyEx(hKey HKEY, index uint32) string {
	var bufLen uint32 = 255
	buf := make([]uint16, bufLen)
	procRegEnumKeyEx.Call(
		uintptr(hKey),
		uintptr(index),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)),
		0,
		0,
		0,
		0)
	return syscall.UTF16ToString(buf)
}

func OpenEventLog(servername string, sourcename string) HANDLE {
	ret, _, _ := procOpenEventLog.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(servername))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(sourcename))))

	return HANDLE(ret)
}

func ReadEventLog(eventlog HANDLE, readflags, recordoffset uint32, buffer []byte, numberofbytestoread uint32, bytesread, minnumberofbytesneeded *uint32) bool {
	ret, _, _ := procReadEventLog.Call(
		uintptr(eventlog),
		uintptr(readflags),
		uintptr(recordoffset),
		uintptr(unsafe.Pointer(&buffer[0])),
		uintptr(numberofbytestoread),
		uintptr(unsafe.Pointer(bytesread)),
		uintptr(unsafe.Pointer(minnumberofbytesneeded)))

	return ret != 0
}

func CloseEventLog(eventlog HANDLE) bool {
	ret, _, _ := procCloseEventLog.Call(
		uintptr(eventlog))

	return ret != 0
}

func OpenSCManager(lpMachineName, lpDatabaseName string, dwDesiredAccess uint32) (HANDLE, error) {
	var p1, p2 uintptr
	if len(lpMachineName) > 0 {
		p1 = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpMachineName)))
	}
	if len(lpDatabaseName) > 0 {
		p2 = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpDatabaseName)))
	}
	ret, _, _ := procOpenSCManager.Call(
		p1,
		p2,
		uintptr(dwDesiredAccess))

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return HANDLE(ret), nil
}

func CloseServiceHandle(hSCObject HANDLE) error {
	ret, _, _ := procCloseServiceHandle.Call(uintptr(hSCObject))
	if ret == 0 {
		return syscall.GetLastError()
	}
	return nil
}

func OpenService(hSCManager HANDLE, lpServiceName string, dwDesiredAccess uint32) (HANDLE, error) {
	ret, _, _ := procOpenService.Call(
		uintptr(hSCManager),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpServiceName))),
		uintptr(dwDesiredAccess))

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return HANDLE(ret), nil
}

func StartService(hService HANDLE, lpServiceArgVectors []string) error {
	l := len(lpServiceArgVectors)
	var ret uintptr
	if l == 0 {
		ret, _, _ = procStartService.Call(
			uintptr(hService),
			0,
			0)
	} else {
		lpArgs := make([]uintptr, l)
		for i := 0; i < l; i++ {
			lpArgs[i] = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpServiceArgVectors[i])))
		}

		ret, _, _ = procStartService.Call(
			uintptr(hService),
			uintptr(l),
			uintptr(unsafe.Pointer(&lpArgs[0])))
	}

	if ret == 0 {
		return syscall.GetLastError()
	}

	return nil
}

func ControlService(hService HANDLE, dwControl uint32, lpServiceStatus *SERVICE_STATUS) bool {
	if lpServiceStatus == nil {
		panic("ControlService:lpServiceStatus cannot be nil")
	}

	ret, _, _ := procControlService.Call(
		uintptr(hService),
		uintptr(dwControl),
		uintptr(unsafe.Pointer(lpServiceStatus)))

	return ret != 0
}

func LogonUser(username, domain, password string, logonType uint32, logonProvider uint32) (token syscall.Handle, err error) {
	r1, _, e1 := procLogonUserW.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(username))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(domain))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(password))),
		uintptr(logonType),
		uintptr(logonProvider),
		uintptr(unsafe.Pointer(&token)))
	runtime.KeepAlive(username)
	runtime.KeepAlive(domain)
	runtime.KeepAlive(password)
	if int(r1) == 0 {
		return syscall.InvalidHandle, os.NewSyscallError("LogonUser", e1)
	}
	return token, nil
}
