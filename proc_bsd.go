// +build darwin freebsd netbsd openbsd

package proc

import (
	"bufio"
	"bytes"
	"os/exec"
	"strconv"
	"strings"
)

func ps(args ...string) (*bytes.Buffer, error) {
	var stdout bytes.Buffer
	cmd := exec.Command("ps", args...)
	cmd.Stdout = &stdout

	if err := cmd.Run(); err != nil {
		eErr, ok := err.(*exec.ExitError)
		if ok && !eErr.Success() {
			return &stdout, nil
		}

		return nil, err
	}

	return &stdout, nil
}

// listProcs returns a list of the running processes.
func listProcs() ([]*Proc, error) {
	buf, err := ps("-x", "-o", "pid= ppid=")
	if err != nil {
		return nil, err
	}

	var procs []*Proc
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())

		// If the field can't be converted to a number skip.
		// It's probably the name of the column or something similar.
		pid, err := strconv.Atoi(fields[0])
		if err != nil {
			continue
		}

		ppid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		procs = append(procs, &Proc{Pid: pid, Ppid: ppid})
	}

	return procs, scanner.Err()
}
