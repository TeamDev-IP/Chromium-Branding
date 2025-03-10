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

package core

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
)

// Signs all the required Chromium binaries on macOS and Windows.
func SignAppBinaries(outDir string, params common.BrandingParams) (bool, error) {
	if runtime.GOOS == "linux" {
		return false, nil
	}

	filesToSign, err := getFilesToSign(outDir, params)
	if err != nil {
		return false, err
	}

	return SignBinaries(params, filesToSign, "application")
}

// Signs the provided `binaries` if there is an available sign tool.
//
// If the sign tool is configured incorrectly, skips signing.
// For example, on macOS it can be caused by the invalid `codesign` parameters.
// On Windows, by the empty sign command or the command without the placeholder to substitute a binary file path.
//
// This function also reports signing status updates based on the provided `binariesGroupName` in the following format:
// [STATUS] Signing `binariesGroupName` <current status>
//
// If signing has succeeded, returns `true`.
// If signing has been skipped, returns `false`.
// If signing has failed, returns `false` and error.
func SignBinaries(params common.BrandingParams, binaries []string, binariesGroupName string) (bool, error) {
	signTool, err := GetSignTool(params)
	if err != nil {
		return false, nil
	}

	for _, binaryPath := range binaries {
		if err := signTool.SignBinary(binaryPath); err != nil {
			return false, unableToSign(err)
		}
	}

	return true, nil
}

func unableToSign(internalError error) error {
	return errors.New("unable to sign binaries: " + internalError.Error())
}

func getFilesToSign(outBinDir string, params common.BrandingParams) ([]string, error) {
	switch runtime.GOOS {
	case "windows":
		return getFilesToSignWin(outBinDir)
	case "darwin":
		return getFilesToSignMac(outBinDir, *params.Mac.Bundle.Name)
	default:
		return []string{}, errors.New("cannot sign binaries on the platform: " + runtime.GOOS)
	}
}

func getFilesToSignWin(outBinDir string) ([]string, error) {
	return getFilesFromDirectoryRoot(outBinDir, []string{"exe", "dll"})
}

func getFilesToSignMac(outBinDir string, bundleName string) ([]string, error) {
	bundlePath := filepath.Join(outBinDir, bundleName+".app")
	currentVersionDir := filepath.Join(bundlePath, "Contents",
		"Frameworks", "Chromium Framework.framework", "Versions", "Current")
	libraries := filepath.Join(currentVersionDir, "Libraries")
	helpers := filepath.Join(currentVersionDir, "Helpers")
	fileExtensionsToSeek := map[string][]string{
		libraries: {"dylib"},
		helpers:   {"", "app"},
	}
	filesToSign := []string{}
	for directory, binaryFilesExtensions := range fileExtensionsToSeek {
		if files, err := getFilesFromDirectoryRoot(directory, binaryFilesExtensions); err == nil {
			filesToSign = append(filesToSign, files...)
		} else {
			return filesToSign, err
		}
	}

	filesToSign = append(filesToSign, currentVersionDir)
	filesToSign = append(filesToSign, bundlePath)

	return filesToSign, nil
}

func getFilesFromDirectoryRoot(directory string, extensions []string) ([]string, error) {
	files := []string{}
	entries, err := os.ReadDir(directory)
	if err != nil {
		return files, err
	}
	for _, entry := range entries {
		if !entry.IsDir() || getFileExtension(entry.Name()) != "" {
			if mathesOneOfExtensions(entry.Name(), extensions) {
				files = append(files, filepath.Join(directory, entry.Name()))
			}
		}
	}
	return files, nil
}

func mathesOneOfExtensions(filename string, extensions []string) bool {
	for _, extension := range extensions {
		if getFileExtension(filename) == extension {
			return true
		}
	}
	return false
}

func getFileExtension(filename string) string {
	splitParts := strings.Split(filename, ".")
	if len(splitParts) == 1 {
		return ""
	} else {
		return splitParts[len(splitParts)-1]
	}
}
