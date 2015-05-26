package proc

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// listProcs returns a list of the running processes.
func listProcs() ([]*Proc, error) {
	var (
		comm  string
		state byte
		ppid  int
	)

	procfs, err := os.Open("/proc")
	if err != nil {
		return nil, err
	}
	defer procfs.Close()

	names, err := procfs.Readdirnames(0)
	if err != nil {
		return nil, err
	}

	var procs []*Proc
	for _, name := range names {
		// Skip non pid paths.
		pid, err := strconv.Atoi(name)
		if err != nil {
			continue
		}

		stat, err := os.Open(filepath.Join("/proc", strconv.Itoa(pid), "stat"))
		if err != nil {
			// If it can't be opened the process has exited.
			if os.IsNotExist(err) {
				continue
			}

			return nil, err
		}
		defer stat.Close()

		// Store the stat info needed to retreive the ppid.
		_, err = fmt.Fscanf(stat, "%d %s %c %d", &pid, &comm, &state, &ppid)
		if err != nil {
			return nil, err
		}

		procs = append(procs, &Proc{Pid: pid, Ppid: ppid})
	}

	return procs, nil
}
