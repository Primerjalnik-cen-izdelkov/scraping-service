package common

import (
	"os/exec"
	"runtime"
)

func MultiOSCommand(name string, args ...string) *exec.Cmd {
	if runtime.GOOS == "windows" {
		// prepend name of the command to args
		args = append(args, "")
		copy(args[1:], args)
		args[0] = name

		// prepend /C to the args
		args = append(args, "")
		copy(args[1:], args)
		args[0] = "/C"

		return exec.Command("cmd", args...)
	} else {
		return exec.Command(name, args...)
	}
}
