// +build darwin freebsd netbsd openbsd

package proc

import (
	"encoding/binary"
	"syscall"
	"unsafe"
)

var (
	allProcs = []int32{1, 14, 0} // CTL_KERN, KERN_PROC, KERN_PROC_ALL
	// The following are sizes and offsets so we can easily get the fields out
	// of the output without having to map out the entirety of kinfo_proc since
	// it is quite large.
	kinfoProcSize    = 648
	kinfoProcPidOff  = 40
	kinfoProcPpidOff = 560
)

// sysctl implements access to the sysctl calling interface.
func sysctl(mib []int32, old *byte, oldlen *uintptr, new *byte, newlen uintptr) error {
	var err error
	mibptr := unsafe.Pointer(&mib[0])

	_, _, e1 := syscall.Syscall6(syscall.SYS___SYSCTL, uintptr(mibptr), uintptr(len(mib)),
		uintptr(unsafe.Pointer(old)), uintptr(unsafe.Pointer(oldlen)),
		uintptr(unsafe.Pointer(new)), uintptr(newlen))
	if e1 != 0 {
		err = e1
	}

	return err
}

// listProcs returns a list of the running processes.
func listProcs() ([]*Proc, error) {
	var size uintptr
	var procs []*Proc

	// Get the byte size of all the process structures.
	err := sysctl(allProcs, nil, &size, nil, 0)
	if err != nil {
		return nil, err
	}
	data := make([]byte, size)

	// Fill in data with the processes structures.
	err = sysctl(allProcs, &data[0], &size, nil, 0)
	if err != nil {
		return nil, err
	}

	// Get the process information for the processes.
	num := int(size) / kinfoProcSize
	for i := 0; i < num; i++ {
		structure := data[i*kinfoProcSize:]
		pid := binary.LittleEndian.Uint32(structure[kinfoProcPidOff : kinfoProcPidOff+4])
		ppid := binary.LittleEndian.Uint32(structure[kinfoProcPpidOff : kinfoProcPpidOff+4])

		procs = append(procs, &Proc{Pid: int(pid), Ppid: int(ppid)})
	}

	return procs, nil
}
