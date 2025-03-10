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

package win

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
)

const binaryFilePathPlaceholder = "@@BINARY_PATH@@"

// The sign tool on Windows is defined by the custom user commands.
//
// To create one, use `GetSignToolWin`.
type SignToolWin struct {
	signCommandTemplate string
}

func GetSignToolWin(params common.BrandingParams) (*SignToolWin, error) {
	if params.Win.SignCommand == "" {
		return nil, errors.New("the sign command is empty")
	}
	if !canSubstituteBinaryPath(params.Win.SignCommand) {
		return nil, fmt.Errorf("the sign command is invalid: requires %s placeholder", binaryFilePathPlaceholder)
	}

	return &SignToolWin{params.Win.SignCommand}, nil
}

func (tool *SignToolWin) SignBinary(binaryPath string) error {
	return tool.execCommand(substituteBinaryPath(tool.signCommandTemplate, binaryPath), binaryPath)
}

// Executes the given `command` with the `binaryPath` substituted.
//
// If the command is `nil`, this is just no-op.
func (tool *SignToolWin) execCommand(command string, binaryPath string) error {
	if _, err := os.Stat(binaryPath); err != nil {
		return err
	}

	if err := base.ExecCommand(command, []string{}); err != nil {
		return err
	}
	return nil
}

// Indicates if the given `commandTemplate` contains the `binaryFilePathPlaceholder`
// to substitute a value instead.
func canSubstituteBinaryPath(commandTemplate string) bool {
	return strings.Contains(commandTemplate, binaryFilePathPlaceholder)
}

// Substitutes the given `binaryPath` instead of the `binaryFilePathPlaceholder` into the `commandTemplate`.
func substituteBinaryPath(commandTemplate string, binaryPath string) string {
	return strings.Replace(commandTemplate, binaryFilePathPlaceholder, binaryPath, -1)
}
