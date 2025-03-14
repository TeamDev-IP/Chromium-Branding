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
