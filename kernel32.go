//go:build windows
// +build windows

// Copyright 2010-2012 The W32 Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package w32

import (
	"encoding/binary"
	"errors"
	"syscall"
	"unsafe"
)

var (
	modkernel32                    = syscall.NewLazyDLL("kernel32.dll")
	procCloseHandle                = modkernel32.NewProc("CloseHandle")
	procCopyMemory                 = modkernel32.NewProc("RtlCopyMemory")
	procCreateProcessA             = modkernel32.NewProc("CreateProcessA")
	procCreateProcessW             = modkernel32.NewProc("CreateProcessW")
	procCreateRemoteThread         = modkernel32.NewProc("CreateRemoteThread")
	procCreateThread               = modkernel32.NewProc("CreateThread")
	procCreateToolhelp32Snapshot   = modkernel32.NewProc("CreateToolhelp32Snapshot")
	procFindResource               = modkernel32.NewProc("FindResourceW")
	procGetConsoleScreenBufferInfo = modkernel32.NewProc("GetConsoleScreenBufferInfo")
	procGetConsoleWindow           = modkernel32.NewProc("GetConsoleWindow")
	procGetCurrentThread           = modkernel32.NewProc("GetCurrentThread")
	procGetCurrentThreadId         = modkernel32.NewProc("GetCurrentThreadId")
	procGetDiskFreeSpaceEx         = modkernel32.NewProc("GetDiskFreeSpaceExW")
	procGetExitCodeProcess         = modkernel32.NewProc("GetExitCodeProcess")
	procGetLastError               = modkernel32.NewProc("GetLastError")
	procGetLogicalDrives           = modkernel32.NewProc("GetLogicalDrives")
	procGetModuleHandle            = modkernel32.NewProc("GetModuleHandleW")
	procGetPerformanceInfo         = modkernel32.NewProc("K32GetPerformanceInfo")
	procGetProcAddress             = modkernel32.NewProc("GetProcAddress")
	procGetProcessTimes            = modkernel32.NewProc("GetProcessTimes")
	procGetSystemPowerStatus       = modkernel32.NewProc("GetSystemPowerStatus")
	procGetSystemTime              = modkernel32.NewProc("GetSystemTime")
	procGetSystemTimes             = modkernel32.NewProc("GetSystemTimes")
	procGetUserDefaultLCID         = modkernel32.NewProc("GetUserDefaultLCID")
	procGlobalAlloc                = modkernel32.NewProc("GlobalAlloc")
	procGlobalFree                 = modkernel32.NewProc("GlobalFree")
	procGlobalLock                 = modkernel32.NewProc("GlobalLock")
	procGlobalMemoryStatusEx       = modkernel32.NewProc("GlobalMemoryStatusEx")
	procGlobalUnlock               = modkernel32.NewProc("GlobalUnlock")
	procLoadLibraryA               = modkernel32.NewProc("LoadLibraryA")
	procLoadResource               = modkernel32.NewProc("LoadResource")
	procLockResource               = modkernel32.NewProc("LockResource")
	procLstrcpy                    = modkernel32.NewProc("lstrcpyW")
	procLstrlen                    = modkernel32.NewProc("lstrlenW")
	procModule32First              = modkernel32.NewProc("Module32FirstW")
	procModule32Next               = modkernel32.NewProc("Module32NextW")
	procMoveMemory                 = modkernel32.NewProc("RtlMoveMemory")
	procMulDiv                     = modkernel32.NewProc("MulDiv")
	procOpenProcess                = modkernel32.NewProc("OpenProcess")
	procQueryFullProcessImageName  = modkernel32.NewProc("QueryFullProcessImageNameW")
	procQueryPerformanceCounter    = modkernel32.NewProc("QueryPerformanceCounter")
	procQueryPerformanceFrequency  = modkernel32.NewProc("QueryPerformanceFrequency")
	procReadProcessMemory          = modkernel32.NewProc("ReadProcessMemory")
	procSetConsoleCtrlHandler      = modkernel32.NewProc("SetConsoleCtrlHandler")
	procSetConsoleTextAttribute    = modkernel32.NewProc("SetConsoleTextAttribute")
	procSetSystemTime              = modkernel32.NewProc("SetSystemTime")
	procSizeofResource             = modkernel32.NewProc("SizeofResource")
	procTerminateProcess           = modkernel32.NewProc("TerminateProcess")
	procVirtualAlloc               = modkernel32.NewProc("VirtualAlloc")
	procVirtualAllocEx             = modkernel32.NewProc("VirtualAllocEx")
	procVirtualFreeEx              = modkernel32.NewProc("VirtualFreeEx")
	procVirtualProtect             = modkernel32.NewProc("VirtualProtect")
	procVirtualQuery               = modkernel32.NewProc("VirtualQuery")
	procVirtualQueryEx             = modkernel32.NewProc("VirtualQueryEx")
	procWaitForSingleObject        = modkernel32.NewProc("WaitForSingleObject")
	procWriteProcessMemory         = modkernel32.NewProc("WriteProcessMemory")
)

func GetExitCodeProcess(hProcess HANDLE) (code uintptr, e error) {
	ret, _, lastErr := procGetExitCodeProcess.Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(&code)),
	)

	if ret == 0 {
		e = lastErr
	}

	return
}

// DWORD WINAPI WaitForSingleObject(
//   _In_ HANDLE hHandle,
//   _In_ DWORD  dwMilliseconds
// );

func WaitForSingleObject(hHandle HANDLE, msecs uint32) (ok bool, e error) {

	ret, _, lastErr := procWaitForSingleObject.Call(
		uintptr(hHandle),
		uintptr(msecs),
	)

	if ret == WAIT_OBJECT_0 {
		ok = true
		return
	}

	// don't set e for timeouts, or it will be ERROR_SUCCESS which is
	// confusing
	if ret != WAIT_TIMEOUT {
		e = lastErr
	}
	return

}

func CreateProcessW(
	lpApplicationName, lpCommandLine string,
	lpProcessAttributes, lpThreadAttributes *SECURITY_ATTRIBUTES,
	bInheritHandles BOOL,
	dwCreationFlags uint32,
	lpEnvironment unsafe.Pointer,
	lpCurrentDirectory string,
	lpStartupInfo *STARTUPINFOW,
	lpProcessInformation *PROCESS_INFORMATION,
) (e error) {

	var lpAN, lpCL, lpCD *uint16
	if len(lpApplicationName) > 0 {
		lpAN, e = syscall.UTF16PtrFromString(lpApplicationName)
		if e != nil {
			return
		}
	}
	if len(lpCommandLine) > 0 {
		lpCL, e = syscall.UTF16PtrFromString(lpCommandLine)
		if e != nil {
			return
		}
	}
	if len(lpCurrentDirectory) > 0 {
		lpCD, e = syscall.UTF16PtrFromString(lpCurrentDirectory)
		if e != nil {
			return
		}
	}

	ret, _, lastErr := procCreateProcessW.Call(
		uintptr(unsafe.Pointer(lpAN)),
		uintptr(unsafe.Pointer(lpCL)),
		uintptr(unsafe.Pointer(lpProcessAttributes)),
		uintptr(unsafe.Pointer(lpProcessInformation)),
		uintptr(bInheritHandles),
		uintptr(dwCreationFlags),
		uintptr(lpEnvironment),
		uintptr(unsafe.Pointer(lpCD)),
		uintptr(unsafe.Pointer(lpStartupInfo)),
		uintptr(unsafe.Pointer(lpProcessInformation)),
	)

	if ret == 0 {
		e = lastErr
	}

	return
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366902(v=vs.85).aspx
func VirtualQuery(lpAddress uintptr, lpBuffer *MEMORY_BASIC_INFORMATION, dwLength int) int {
	ret, _, _ := procVirtualQuery.Call(
		lpAddress,
		uintptr(unsafe.Pointer(lpBuffer)),
		uintptr(dwLength))
	return int(ret) // TODO check for errors
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366907(v=vs.85).aspx
func VirtualQueryEx(hProcess HANDLE, lpAddress uintptr, lpBuffer *MEMORY_BASIC_INFORMATION, dwLength int) int {
	ret, _, _ := procVirtualQueryEx.Call(
		uintptr(hProcess), // The handle to a process.
		lpAddress,
		uintptr(unsafe.Pointer(lpBuffer)),
		uintptr(dwLength))
	return int(ret) // TODO check for errors
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366898(v=vs.85).aspx
func VirtualProtect(lpAddress uintptr, dwSize int, flNewProtect int, lpflOldProtect *DWORD) bool {
	ret, _, _ := procVirtualProtect.Call(
		lpAddress,
		uintptr(dwSize),
		uintptr(flNewProtect),
		uintptr(unsafe.Pointer(lpflOldProtect)))
	return ret != 0
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms682425(v=vs.85).aspx
func CreateProcessA(lpApplicationName *string,
	lpCommandLine string,
	lpProcessAttributes *syscall.SecurityAttributes,
	lpThreadAttributes *syscall.SecurityAttributes,
	bInheritHandles bool,
	dwCreationFlags uint32,
	lpEnvironment *string,
	lpCurrentDirectory *uint16,
	lpStartupInfo *syscall.StartupInfo,
	lpProcessInformation *syscall.ProcessInformation) {

	inherit := 0
	if bInheritHandles {
		inherit = 1
	}

	procCreateProcessA.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(*lpApplicationName))),
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(lpCommandLine))),
		uintptr(unsafe.Pointer(lpProcessAttributes)),
		uintptr(unsafe.Pointer(lpThreadAttributes)),
		uintptr(inherit),
		uintptr(dwCreationFlags),
		uintptr(unsafe.Pointer(lpEnvironment)),
		uintptr(unsafe.Pointer(lpCurrentDirectory)),
		uintptr(unsafe.Pointer(lpStartupInfo)),
		uintptr(unsafe.Pointer(lpProcessInformation)))
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366890(v=vs.85).aspx
func VirtualAllocEx(hProcess HANDLE, lpAddress int, dwSize int, flAllocationType int, flProtect int) (addr uintptr, err error) {
	ret, _, err := procVirtualAllocEx.Call(
		uintptr(hProcess),  // The handle to a process.
		uintptr(lpAddress), // The pointer that specifies a desired starting address for the region of pages that you want to allocate.
		uintptr(dwSize),    // The size of the region of memory to allocate, in bytes.
		uintptr(flAllocationType),
		uintptr(flProtect))
	if int(ret) == 0 {
		return ret, err
	}
	return ret, nil
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/aa366887(v=vs.85).aspx
func VirtualAlloc(lpAddress int, dwSize int, flAllocationType int, flProtect int) (addr uintptr, err error) {
	ret, _, err := procVirtualAlloc.Call(
		uintptr(lpAddress), // The starting address of the region to allocate
		uintptr(dwSize),    // The size of the region of memory to allocate, in bytes.
		uintptr(flAllocationType),
		uintptr(flProtect))
	if int(ret) == 0 {
		return ret, err
	}
	return ret, nil
}

// https://github.com/AllenDang/w32/pull/62/commits/08a52ff508c6b2b9b9bf5ee476109b903dcf219d
func VirtualFreeEx(hProcess HANDLE, lpAddress, dwSize uintptr, dwFreeType uint32) bool {
	ret, _, _ := procVirtualFreeEx.Call(
		uintptr(hProcess),
		lpAddress,
		dwSize,
		uintptr(dwFreeType),
	)

	return ret != 0
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms683212(v=vs.85).aspx
func GetProcAddress(hProcess HANDLE, procname string) (addr uintptr, err error) {
	var pn uintptr

	if procname == "" {
		pn = 0
	} else {
		pn = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(procname)))
	}

	ret, _, err := procGetProcAddress.Call(uintptr(hProcess), pn)
	if int(ret) == 0 {
		return ret, err
	}
	return ret, nil
}

// https://msdn.microsoft.com/en-us/library/windows/desktop/ms682437(v=vs.85).aspx
// Credit: https://github.com/contester/runlib/blob/master/win32/win32_windows.go#L577
func CreateRemoteThread(hprocess HANDLE, sa *syscall.SecurityAttributes,
	stackSize uint32, startAddress uint32, parameter uintptr, creationFlags uint32) (syscall.Handle, uint32, error) {
	var threadId uint32
	r1, _, e1 := procCreateRemoteThread.Call(
		uintptr(hprocess),
		uintptr(unsafe.Pointer(sa)),
		uintptr(stackSize),
		uintptr(startAddress),
		uintptr(parameter),
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(&threadId)))

	if int(r1) == 0 {
		return syscall.InvalidHandle, 0, e1
	}
	return syscall.Handle(r1), threadId, nil
}

// https://docs.microsoft.com/en-us/windows/win32/api/processthreadsapi/nf-processthreadsapi-createthread
func CreateThread(sa *syscall.SecurityAttributes,
	stackSize uint32, threadFn ThreadFunc, parameter uintptr, creationFlags uint32) (syscall.Handle, uint32, error) {

	var threadId uint32
	r1, _, e1 := procCreateThread.Call(
		uintptr(unsafe.Pointer(sa)),
		uintptr(stackSize),
		uintptr(syscall.NewCallback(threadFn)),
		uintptr(parameter),
		uintptr(creationFlags),
		uintptr(unsafe.Pointer(&threadId)))

	if int(r1) == 0 {
		return syscall.InvalidHandle, 0, e1
	}

	return syscall.Handle(r1), threadId, nil
}

func GetModuleHandle(modulename string) HINSTANCE {
	var mn uintptr
	if modulename == "" {
		mn = 0
	} else {
		mn = uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(modulename)))
	}
	ret, _, _ := procGetModuleHandle.Call(mn)
	return HINSTANCE(ret)
}

func MulDiv(number, numerator, denominator int) int {
	ret, _, _ := procMulDiv.Call(
		uintptr(number),
		uintptr(numerator),
		uintptr(denominator))

	return int(ret)
}

func GetConsoleWindow() HWND {
	ret, _, _ := procGetConsoleWindow.Call()

	return HWND(ret)
}

func GetCurrentThread() HANDLE {
	ret, _, _ := procGetCurrentThread.Call()

	return HANDLE(ret)
}

func GetCurrentThreadId() DWORD {
	ret, _, _ := procGetCurrentThreadId.Call()

	return DWORD(ret)
}

func GetLogicalDrives() uint32 {
	ret, _, _ := procGetLogicalDrives.Call()

	return uint32(ret)
}

func GetUserDefaultLCID() uint32 {
	ret, _, _ := procGetUserDefaultLCID.Call()

	return uint32(ret)
}

func Lstrlen(lpString *uint16) int {
	ret, _, _ := procLstrlen.Call(uintptr(unsafe.Pointer(lpString)))

	return int(ret)
}

func Lstrcpy(buf []uint16, lpString *uint16) {
	procLstrcpy.Call(
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(lpString)))
}

func GlobalAlloc(uFlags uint, dwBytes uint32) HGLOBAL {
	ret, _, _ := procGlobalAlloc.Call(
		uintptr(uFlags),
		uintptr(dwBytes))

	if ret == 0 {
		panic("GlobalAlloc failed")
	}

	return HGLOBAL(ret)
}

func GlobalFree(hMem HGLOBAL) {
	ret, _, _ := procGlobalFree.Call(uintptr(hMem))

	if ret != 0 {
		panic("GlobalFree failed")
	}
}

func GlobalLock(hMem HGLOBAL) unsafe.Pointer {
	ret, _, _ := procGlobalLock.Call(uintptr(hMem))

	if ret == 0 {
		panic("GlobalLock failed")
	}

	return unsafe.Pointer(ret)
}

func GlobalMemoryStatusEx(lpBuffer *MEMORYSTATUSEX) bool {
	ret, _, _ := procGlobalMemoryStatusEx.Call(
		uintptr(unsafe.Pointer(lpBuffer)))

	return ret != 0
}

func GlobalUnlock(hMem HGLOBAL) bool {
	ret, _, _ := procGlobalUnlock.Call(uintptr(hMem))

	return ret != 0
}

func MoveMemory(destination, source unsafe.Pointer, length uint32) {
	procMoveMemory.Call(
		uintptr(unsafe.Pointer(destination)),
		uintptr(source),
		uintptr(length))
}

func FindResource(hModule HMODULE, lpName, lpType *uint16) (HRSRC, error) {
	ret, _, _ := procFindResource.Call(
		uintptr(hModule),
		uintptr(unsafe.Pointer(lpName)),
		uintptr(unsafe.Pointer(lpType)))

	if ret == 0 {
		return 0, syscall.GetLastError()
	}

	return HRSRC(ret), nil
}

func SizeofResource(hModule HMODULE, hResInfo HRSRC) uint32 {
	ret, _, _ := procSizeofResource.Call(
		uintptr(hModule),
		uintptr(hResInfo))

	if ret == 0 {
		panic("SizeofResource failed")
	}

	return uint32(ret)
}

func LockResource(hResData HGLOBAL) unsafe.Pointer {
	ret, _, _ := procLockResource.Call(uintptr(hResData))

	if ret == 0 {
		panic("LockResource failed")
	}

	return unsafe.Pointer(ret)
}

func LoadResource(hModule HMODULE, hResInfo HRSRC) HGLOBAL {
	ret, _, _ := procLoadResource.Call(
		uintptr(hModule),
		uintptr(hResInfo))

	if ret == 0 {
		panic("LoadResource failed")
	}

	return HGLOBAL(ret)
}

func GetLastError() uint32 {
	ret, _, _ := procGetLastError.Call()
	return uint32(ret)
}

func OpenProcess(desiredAccess uint32, inheritHandle bool, processId uint32) (handle HANDLE, err error) {
	inherit := 0
	if inheritHandle {
		inherit = 1
	}

	ret, _, err := procOpenProcess.Call(
		uintptr(desiredAccess),
		uintptr(inherit),
		uintptr(processId))
	if err != nil && err.Error() == "The operation completed successfully." {
		err = nil
	}
	handle = HANDLE(ret)
	return
}

func QueryFullProcessImageName(hProcess HANDLE, dwFlags DWORD) (exeName string, retB bool) {

	buf := make([]uint16, 1024)
	bufLen := uint32(len(buf))

	ret, _, _ := procQueryFullProcessImageName.Call(
		uintptr(hProcess),
		uintptr(dwFlags),
		uintptr(unsafe.Pointer(&buf[0])),
		uintptr(unsafe.Pointer(&bufLen)))

	retB = ret != 0
	if retB {
		exeName = syscall.UTF16ToString(buf)
	}

	return
}

func TerminateProcess(hProcess HANDLE, uExitCode uint) bool {
	ret, _, _ := procTerminateProcess.Call(
		uintptr(hProcess),
		uintptr(uExitCode))
	return ret != 0
}

func CloseHandle(object HANDLE) bool {
	ret, _, _ := procCloseHandle.Call(
		uintptr(object))
	return ret != 0
}

func CreateToolhelp32Snapshot(flags, processId uint32) HANDLE {
	ret, _, _ := procCreateToolhelp32Snapshot.Call(
		uintptr(flags),
		uintptr(processId))

	if ret <= 0 {
		return HANDLE(0)
	}

	return HANDLE(ret)
}

func Module32First(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32First.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)))

	return ret != 0
}

func Module32Next(snapshot HANDLE, me *MODULEENTRY32) bool {
	ret, _, _ := procModule32Next.Call(
		uintptr(snapshot),
		uintptr(unsafe.Pointer(me)))

	return ret != 0
}

func GetSystemTimes(lpIdleTime, lpKernelTime, lpUserTime *FILETIME) bool {
	ret, _, _ := procGetSystemTimes.Call(
		uintptr(unsafe.Pointer(lpIdleTime)),
		uintptr(unsafe.Pointer(lpKernelTime)),
		uintptr(unsafe.Pointer(lpUserTime)))

	return ret != 0
}

func GetProcessTimes(hProcess HANDLE, lpCreationTime, lpExitTime, lpKernelTime, lpUserTime *FILETIME) bool {
	ret, _, _ := procGetProcessTimes.Call(
		uintptr(hProcess),
		uintptr(unsafe.Pointer(lpCreationTime)),
		uintptr(unsafe.Pointer(lpExitTime)),
		uintptr(unsafe.Pointer(lpKernelTime)),
		uintptr(unsafe.Pointer(lpUserTime)))

	return ret != 0
}

func GetConsoleScreenBufferInfo(hConsoleOutput HANDLE) *CONSOLE_SCREEN_BUFFER_INFO {
	var csbi CONSOLE_SCREEN_BUFFER_INFO
	ret, _, _ := procGetConsoleScreenBufferInfo.Call(
		uintptr(hConsoleOutput),
		uintptr(unsafe.Pointer(&csbi)))
	if ret == 0 {
		return nil
	}
	return &csbi
}

func SetConsoleTextAttribute(hConsoleOutput HANDLE, wAttributes uint16) bool {
	ret, _, _ := procSetConsoleTextAttribute.Call(
		uintptr(hConsoleOutput),
		uintptr(wAttributes))
	return ret != 0
}

func GetDiskFreeSpaceEx(dirName string) (r bool,
	freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes uint64) {
	ret, _, _ := procGetDiskFreeSpaceEx.Call(
		uintptr(unsafe.Pointer(syscall.StringToUTF16Ptr(dirName))),
		uintptr(unsafe.Pointer(&freeBytesAvailable)),
		uintptr(unsafe.Pointer(&totalNumberOfBytes)),
		uintptr(unsafe.Pointer(&totalNumberOfFreeBytes)))
	return ret != 0,
		freeBytesAvailable, totalNumberOfBytes, totalNumberOfFreeBytes
}

func GetSystemTime() (time SYSTEMTIME, err error) {
	_, _, err = procGetSystemTime.Call(
		uintptr(unsafe.Pointer(&time)))
	if err.Error() != ErrSuccess {
		return
	}
	err = nil
	return
}

func SetSystemTime(time *SYSTEMTIME) (err error) {
	_, _, err = procSetSystemTime.Call(
		uintptr(unsafe.Pointer(time)))
	if err.Error() != ErrSuccess {
		return
	}
	err = nil
	return
}

// Writes data to an area of memory in a specified process. The entire area to be written to must be accessible or the operation fails.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms681674(v=vs.85).aspx
func WriteProcessMemory(hProcess HANDLE, lpBaseAddress uint32, data []byte, size uint) (err error) {
	var numBytesRead uintptr

	_, _, err = procWriteProcessMemory.Call(uintptr(hProcess),
		uintptr(lpBaseAddress),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&numBytesRead)))
	if err.Error() != ErrSuccess {
		return
	}
	err = nil
	return
}

// Write process memory with a source of uint32
func WriteProcessMemoryAsUint32(hProcess HANDLE, lpBaseAddress uint32, data uint32) (err error) {

	bData := make([]byte, 4)
	binary.LittleEndian.PutUint32(bData, data)
	err = WriteProcessMemory(hProcess, lpBaseAddress, bData, 4)
	if err != nil {
		return
	}
	return
}

// Reads data from an area of memory in a specified process. The entire area to be read must be accessible or the operation fails.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms680553(v=vs.85).aspx
func ReadProcessMemory(hProcess HANDLE, lpBaseAddress uint32, size uint) (data []byte, err error) {
	var numBytesRead uintptr
	data = make([]byte, size)

	_, _, err = procReadProcessMemory.Call(uintptr(hProcess),
		uintptr(lpBaseAddress),
		uintptr(unsafe.Pointer(&data[0])),
		uintptr(size),
		uintptr(unsafe.Pointer(&numBytesRead)))
	if err.Error() != ErrSuccess {
		return
	}
	err = nil
	return
}

// Read process memory and convert the returned data to uint32
func ReadProcessMemoryAsUint32(hProcess HANDLE, lpBaseAddress uint32) (buffer uint32, err error) {
	data, err := ReadProcessMemory(hProcess, lpBaseAddress, 4)
	if err != nil {
		return
	}
	buffer = binary.LittleEndian.Uint32(data)
	return
}

// Adds or removes an application-defined HandlerRoutine function from the list of handler functions for the calling process.
// https://msdn.microsoft.com/en-us/library/windows/desktop/ms686016(v=vs.85).aspx
func SetConsoleCtrlHandler(handlerRoutine func(DWORD) int32, add uint) (err error) {
	_, _, err = procSetConsoleCtrlHandler.Call(uintptr(unsafe.Pointer(&handlerRoutine)),
		uintptr(add))
	if err.Error() != ErrSuccess {
		return
	}
	err = nil
	return
}

func QueryPerformanceCounter() uint64 {
	result := uint64(0)
	procQueryPerformanceCounter.Call(
		uintptr(unsafe.Pointer(&result)),
	)

	return result
}

func QueryPerformanceFrequency() uint64 {
	result := uint64(0)
	procQueryPerformanceFrequency.Call(
		uintptr(unsafe.Pointer(&result)),
	)

	return result
}

func GetSystemPowerStatus(lpSystemPowerStatus *SYSTEM_POWER_STATUS) bool {

	ret, _, _ := procGetSystemPowerStatus.Call(
		uintptr(unsafe.Pointer(lpSystemPowerStatus)),
	)

	return ret != 0
}

type PERFORMANCE_INFORMATION struct {
	cb                uintptr
	CommitTotal       SIZE_T
	CommitLimit       SIZE_T
	CommitPeak        SIZE_T
	PhysicalTotal     SIZE_T
	PhysicalAvailable SIZE_T
	SystemCache       SIZE_T
	KernelTotal       SIZE_T
	KernelPaged       SIZE_T
	KernelNonpaged    SIZE_T
	PageSize          SIZE_T
	HandleCount       DWORD
	ProcessCount      DWORD
	ThreadCount       DWORD
}

// https://learn.microsoft.com/en-us/windows/win32/api/psapi/nf-psapi-getperformanceinfo
func GetPerformanceInfo() (PERFORMANCE_INFORMATION, error) {
	var pi PERFORMANCE_INFORMATION
	pi.cb = unsafe.Sizeof(pi)
	ret, _, _ := procGetPerformanceInfo.Call(
		uintptr(unsafe.Pointer(&pi)),
		uintptr(pi.cb),
	)

	if ret != 0 {
		return pi, nil
	} else {
		return PERFORMANCE_INFORMATION{}, errors.New("Call failed")
	}
}
