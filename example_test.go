package proc_test

import (
	"os"

	"github.com/Bowery/proc"
)

func ExampleGetPidTree() {
	tree, err := proc.GetPidTree(os.Getpid())
	if err != nil {
		panic(err)
	}

	err = tree.Kill()
	if err != nil {
		panic(err)
	}
}
