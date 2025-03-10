// Copyright 2025, TeamDev. All rights reserved.
//
// Redistribution and use in source and/or binary forms, with or without
// modification, must retain the above copyright notice and the following
// disclaimer.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
// "AS IS" AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT
// LIMITED TO, THE IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR
// A PARTICULAR PURPOSE ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT
// OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL,
// SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT
// LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES; LOSS OF USE,
// DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND ON ANY
// THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
// (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

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
