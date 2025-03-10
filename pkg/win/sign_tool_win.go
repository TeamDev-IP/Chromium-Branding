// Copyright (c) 2000-2024 TeamDev. All rights reserved.
// TeamDev PROPRIETARY and CONFIDENTIAL.
// Use is subject to license terms.

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
