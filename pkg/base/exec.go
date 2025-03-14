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

package base

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime"
	"strings"
)

// Exec executes the given command in the current working directory.
func Exec(command string) error {
	return ExecInWorkingDir(command, "")
}

// ExecInWorkingDir executes the command in the specified working directory.
func ExecInWorkingDir(command string, workingDir string) error {
	// Windows requires a different command execution approach described in
	// https://github.com/TeamDev-IP/Molybden/pull/677
	if runtime.GOOS == "windows" {
		_, err := ExecCommandInWorkingDir(command, []string{}, workingDir)
		return err
	} else {
		args := strings.Fields(command)
		_, err := ExecCommandInWorkingDir(args[0], args[1:], workingDir)
		return err
	}
}

// ExecCommand executes the given command with
// the specified arguments in the current working directory.
func ExecCommand(command string, args []string) error {
	_, err := ExecCommandInWorkingDir(command, args, "")
	return err
}

// ExecCommandAndGetOutput executes the given command with
// the specified arguments in the current working directory
func ExecCommandAndGetOutput(command string, args []string) (string, error) {
	out, err := ExecCommandInWorkingDir(command, args, "")
	return string(out), err
}

// ExecCommandInWorkingDir executes the command with the given arguments
// in the specified working directory and environment variables.
func ExecCommandInWorkingDir(command string, args []string, workingDir string, envVariables ...string) ([]byte, error) {
	if Verbose {
		fmt.Println(command + " " + strings.Join(args, " "))
	}
	cmd := createCommand(command, args, envVariables...)
	if workingDir != "" {
		cmd.Dir = workingDir
	}

	var stdBuffer bytes.Buffer
	var mw io.Writer
	if Verbose {
		mw = io.MultiWriter(os.Stdout, &stdBuffer)
	} else {
		mw = io.MultiWriter(&stdBuffer)
	}

	cmd.Stdout = mw
	cmd.Stderr = mw

	if err := cmd.Start(); err != nil {
		return stdBuffer.Bytes(), err
	}

	if err := cmd.Wait(); err != nil {
		return stdBuffer.Bytes(), err
	}
	return stdBuffer.Bytes(), nil
}
