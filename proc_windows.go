package proc

import (
	"syscall"
	"unsafe"
)

var (
	kernel         = syscall.NewLazyDLL("kernel32.dll")
	processFirst   = kernel.NewProc("Process32First")
	processNext    = kernel.NewProc("Process32Next")
	createSnapshot = kernel.NewProc("CreateToolhelp32Snapshot")
)

// processEntry describes a process snapshot.
type processEntry struct {
	dwSize            uint32 // REQUIRED: FILL THIS OUT WITH unsafe.Sizeof(processEntry{})
	cntUsage          uint32
	pid               uint32
	th32DefaultHeapID uintptr
	th32ModuleID      uint32
	cntThreads        uint32
	ppid              uint32
	pcPriClassBase    int32
	dwFlags           uint32
	szExeFile         [260]byte // MAX_PATH is 260, only use byte if using ascii ver procs.
}

// listProcs returns a list of the running processes.
func listProcs() ([]*Proc, error) {
	handle, _, err := createSnapshot.Call(2, 0)
	if syscall.Handle(handle) == syscall.InvalidHandle {
		return nil, err
	}
	defer syscall.CloseHandle(syscall.Handle(handle))

	procs := make([]*Proc, 0)
	procEntry := new(processEntry)
	procEntry.dwSize = uint32(unsafe.Sizeof(*procEntry))

	ret, _, err := processFirst.Call(handle, uintptr(unsafe.Pointer(procEntry)))
	if ret == 0 {
		if err == syscall.ERROR_NO_MORE_FILES {
			return procs, nil
		}

		return nil, err
	}
	procs = append(procs, &Proc{Pid: int(procEntry.pid), Ppid: int(procEntry.ppid)})

	for {
		ret, _, err := processNext.Call(handle, uintptr(unsafe.Pointer(procEntry)))
		if ret == 0 {
			if err == syscall.ERROR_NO_MORE_FILES {
				break
			}

			return nil, err
		}

		procs = append(procs, &Proc{Pid: int(procEntry.pid), Ppid: int(procEntry.ppid)})
	}

	return procs, nil
}
