// Copyright (c) 2000-2024 TeamDev. All rights reserved.
// TeamDev PROPRIETARY and CONFIDENTIAL.
// Use is subject to license terms.

//go:build windows

package base

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

func createCommand(command string, args []string, envVariables ...string) *exec.Cmd {
	var cmd *exec.Cmd
	if len(args) == 0 {
		// On Windows, to run the raw terminal command in Go, we need 2 workarouds:
		// - executing `cmd` directly providing the actual command with `/C` option;
		// - using `SysProcAttr` for providing the raw command line to execute.
		// See https://github.com/golang/go/issues/29841.
		// See https://pkg.go.dev/os/exec#Command.
		//
		// Otherwise, the `argv` of the executable can be parsed incorrectly with respect to
		// the arguments quoting, producing parts of the quoted string as the separate `argv` entries.
		// That's why using `strings.Fields` is also wrong way to split the raw command line.
		cmd = exec.Command("cmd")
		cmd.SysProcAttr = &syscall.SysProcAttr{CmdLine: fmt.Sprintf(`/c "%s"`, command)}
	} else {
		cmd = exec.Command(command, args...)
		return cmd
	}
	for _, variable := range envVariables {
		cmd.Env = append(os.Environ(), variable)
	}
	return cmd
}
