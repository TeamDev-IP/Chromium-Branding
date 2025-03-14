// Copyright (c) 2025 TeamDev
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

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
