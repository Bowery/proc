package proc

import (
	"os"
)

// Proc describes a process and its children.
type Proc struct {
	Pid      int     `json:"pid"`
	Ppid     int     `json:"ppid"`
	Children []*Proc `json:"children"`
}

// Kill kills a process and its children.
func (proc *Proc) Kill() error {
	p, err := os.FindProcess(proc.Pid)
	if err != nil {
		return err
	}

	err = p.Kill()
	if err != nil {
		return err
	}

	if proc.Children != nil {
		for _, p := range proc.Children {
			err = p.Kill()
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// GetPidTree gets the process tree for a pid. The returned process structure
// is nil if the given pid is not a running process.
func GetPidTree(cpid int) (*Proc, error) {
	var root *Proc

	procs, err := listProcs()
	if err != nil {
		return nil, err
	}
	var children []*Proc

	for _, proc := range procs {
		// We've found the root process.
		if proc.Pid == cpid {
			root = proc
			continue
		}

		// Found a child process.
		if proc.Ppid == cpid {
			p, err := GetPidTree(proc.Pid)
			if err != nil {
				return nil, err
			}

			if p != nil {
				children = append(children, p)
			}
		}
	}

	if root != nil {
		root.Children = children
	}
	return root, nil
}
