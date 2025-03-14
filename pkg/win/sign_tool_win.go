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
