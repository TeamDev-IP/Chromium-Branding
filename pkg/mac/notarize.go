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

package mac

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
)

func archiveApp(bundleName string, outDir string) (string, error) {
	bundlePath := filepath.Join(outDir, bundleName)
	bundleZip := filepath.Join(outDir, bundleName+".zip")
	commandArgs := []string{
		"--recurse-paths",
		"--symlinks",
		"--quiet",
		bundleZip,
		bundleName,
	}
	base.Log("Compressing " + bundlePath + "...")
	_, err := base.ExecCommandInWorkingDir("zip", commandArgs, outDir)
	if err != nil {
		return "", err
	}
	return bundleZip, nil
}

func notarize(appBundlePath string, teamID string, appleID string, password string) error {
	base.Log("Notarizing " + appBundlePath + " (it may take a while)...")
	commandArgs := []string{
		"notarytool",
		"submit", appBundlePath,
		"--team-id", base.GetValue(teamID),
		"--apple-id", base.GetValue(appleID),
		"--password", base.GetValue(password),
		"--output-format", "plist",
		"--wait"}
	output, err := base.ExecCommandInWorkingDir("xcrun", commandArgs, "")
	if err != nil {
		return err
	} else if !strings.Contains(string(output), "<string>Accepted</string>") {
		return errors.New("failed to notarize the application. The status is not \"Accepted\"")
	}
	base.Log("The application has been notarized successfully")
	return nil
}

func verify(appBundlePath string) error {
	base.Log("Verifying notarization " + appBundlePath + " ...")
	commandArgs := []string{
		"-a",
		"-v",
		appBundlePath,
	}
	output, err := base.ExecCommandInWorkingDir("spctl", commandArgs, "")
	if err != nil {
		return err
	}
	if strings.Contains(string(output), ": accepted") {
		base.Log("Notarization is verified.")
		return nil
	}
	return errors.New("verification failed")
}

func stapleTicket(appPath string) error {
	base.Log("Stapling a ticket...")
	commandArgs := []string{
		"stapler",
		"staple",
		appPath,
	}
	if err := base.ExecCommand("xcrun", commandArgs); err != nil {
		return err
	}
	return nil
}

func validateStapling(appPath string) error {
	base.Log("Validating the ticket...")
	commandArgs := []string{
		"stapler",
		"validate",
		appPath,
	}
	if err := base.ExecCommand("xcrun", commandArgs); err != nil {
		return err
	}
	return nil
}

// ValidateNotarizationParams checks if the notarization parameters are set.
func ValidateNotarizationParams(params common.BrandingParams) error {
	p := map[string]string{
		"Team ID":  base.GetValue(params.Mac.TeamId),
		"Apple ID": base.GetValue(params.Mac.AppleId),
		"Password": base.GetValue(params.Mac.Password),
	}
	for paramName, paramValue := range p {
		if paramValue == "" {
			return errors.New(paramName + " is empty")
		}
	}
	return nil
}

// Notarize notarizes the application bundle with the provided parameters.
func Notarize(outDir string, params common.BrandingParams) (bool, error) {
	err := ValidateNotarizationParams(params)
	if err != nil {
		return false, nil
	}

	outDirPath, err := filepath.Abs(outDir)
	if err != nil {
		return false, err
	}
	appBundleName := *params.Mac.Bundle.Name + ".app"
	appBundlePath := filepath.Join(outDirPath, appBundleName)
	appBundleZip, err := archiveApp(appBundleName, outDirPath)
	if err != nil {
		return false, err
	}

	defer func(name string) {
		if base.PathExists(name) {
			_ = os.Remove(name)
		}
	}(appBundleZip)

	if _, err := os.Stat(appBundleZip); err != nil {
		return false, err
	}

	if err := notarize(
		appBundleZip,
		base.GetValue(params.Mac.TeamId),
		base.GetValue(params.Mac.AppleId),
		base.GetValue(params.Mac.Password),
	); err != nil {
		return false, err
	}

	if err := verify(appBundlePath); err != nil {
		return false, err
	}

	if err := stapleTicket(appBundlePath); err != nil {
		return false, err
	}

	if err := validateStapling(appBundlePath); err != nil {
		return false, err
	}

	return true, nil
}
