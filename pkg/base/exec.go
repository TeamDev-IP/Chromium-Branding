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
