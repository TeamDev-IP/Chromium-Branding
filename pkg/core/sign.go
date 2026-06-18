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

package core

import (
	"bytes"
	"encoding/xml"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
)

// Signs all the required Chromium binaries on macOS and Windows.
func SignAppBinaries(outDir string, params common.BrandingParams) (bool, error) {
	if runtime.GOOS == "linux" {
		return false, nil
	}
	if runtime.GOOS == "darwin" {
		return signMacAppBinaries(outDir, params)
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

func signMacAppBinaries(outDir string, params common.BrandingParams) (bool, error) {
	signTool, err := GetSignTool(params)
	if err != nil {
		return false, nil
	}

	mst, ok := signTool.(MacSignTool)
	if !ok {
		return false, errors.New("macOS sign tool does not support entitlements-aware signing")
	}

	if base.GetValue(params.Mac.CodesignIdentity) == "" {
		return false, nil
	}

	helperEntitlements, tempFile, err := prepareForSigning(outDir, params)
	if err != nil {
		return false, err
	}
	if tempFile != "" {
		defer os.Remove(tempFile)
	}

	bundleName := *params.Mac.Bundle.Name
	bundlePath := filepath.Join(outDir, bundleName+".app")

	filesToSign, err := getFilesToSignMac(outDir, bundleName)
	if err != nil {
		return false, err
	}

	for _, path := range filesToSign {
		if path == bundlePath {
			continue
		}
		if err := mst.SignBinaryWithEntitlements(path, helperEntitlements); err != nil {
			return false, unableToSign(err)
		}
	}

	if err := signTool.SignBinary(bundlePath); err != nil {
		return false, unableToSign(err)
	}

	return true, nil
}

// prepareForSigning copies the provisioning profile and builds a filtered
// helper entitlements file when keychain-access-groups is present.
// Returns the helper entitlements path and the temp file path (empty if no
// temp file was created). The caller must remove the temp file when done.
func prepareForSigning(outDir string, params common.BrandingParams) (helperEntitlements, tempFile string, err error) {
	entitlementsPath := params.Mac.CodesignEntitlements

	hasKAG, err := hasKeychainAccessGroups(entitlementsPath)
	if err != nil {
		return "", "", err
	}

	if !hasKAG {
		return entitlementsPath, "", nil
	}

	profilePath := params.Mac.ProvisioningProfile
	if profilePath == "" {
		return "", "", errors.New("entitlements contain keychain-access-groups but no provisioning profile is configured; set mac.provisioningProfile in params.json")
	}
	if _, statErr := os.Stat(profilePath); os.IsNotExist(statErr) {
		return "", "", fmt.Errorf("provisioning profile not found: %s", profilePath)
	}

	bundleName := *params.Mac.Bundle.Name
	destProfile := filepath.Join(outDir, bundleName+".app", "Contents", "embedded.provisionprofile")
	if err := base.CopyFile(profilePath, destProfile); err != nil {
		return "", "", fmt.Errorf("copying provisioning profile: %w", err)
	}

	tempPath, err := writeHelperEntitlements(entitlementsPath)
	if err != nil {
		return "", "", err
	}

	return tempPath, tempPath, nil
}

// hasKeychainAccessGroups reports whether the plist at entitlementsPath
// contains a keychain-access-groups key.
func hasKeychainAccessGroups(entitlementsPath string) (bool, error) {
	data, err := os.ReadFile(entitlementsPath)
	if err != nil {
		return false, fmt.Errorf("reading entitlements %s: %w", entitlementsPath, err)
	}
	decoder := xml.NewDecoder(bytes.NewReader(data))
	inKey := false
	for {
		tok, err := decoder.Token()
		if err == io.EOF {
			break
		}
		if err != nil {
			return false, fmt.Errorf("parsing entitlements: %w", err)
		}
		switch t := tok.(type) {
		case xml.StartElement:
			if t.Name.Local == "key" {
				inKey = true
			}
		case xml.EndElement:
			inKey = false
		case xml.CharData:
			if inKey && strings.TrimSpace(string(t)) == "keychain-access-groups" {
				return true, nil
			}
		}
	}
	return false, nil
}

// writeHelperEntitlements writes a copy of the plist at entitlementsPath
// with the keychain-access-groups key/value pair removed to a temp file.
// Returns the temp file path; the caller is responsible for removing it.
func writeHelperEntitlements(entitlementsPath string) (string, error) {
	data, err := os.ReadFile(entitlementsPath)
	if err != nil {
		return "", fmt.Errorf("reading entitlements: %w", err)
	}

	filtered, err := plistWithoutKeychainAccessGroups(data)
	if err != nil {
		return "", err
	}

	tmpFile, err := os.CreateTemp("", "helper-entitlements-*.plist")
	if err != nil {
		return "", fmt.Errorf("creating temp entitlements: %w", err)
	}
	defer tmpFile.Close()

	if _, err := tmpFile.Write(filtered); err != nil {
		os.Remove(tmpFile.Name())
		return "", fmt.Errorf("writing temp entitlements: %w", err)
	}
	return tmpFile.Name(), nil
}

// plistWithoutKeychainAccessGroups returns a copy of the plist bytes with the
// keychain-access-groups key/value pair removed.
func plistWithoutKeychainAccessGroups(data []byte) ([]byte, error) {
	decoder := xml.NewDecoder(bytes.NewReader(data))

	var (
		removeFrom      int64 = -1
		removeTo        int64 = -1
		prevWSStart     int64 = -1
		candidateFrom   int64 = -1
		inKey                 = false
		foundKAGKey           = false
		waitingForValue       = false
		skipDepth             = 0
	)

	for removeTo < 0 {
		tokStart := decoder.InputOffset()
		tok, err := decoder.Token()
		tokEnd := decoder.InputOffset()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("parsing entitlements plist: %w", err)
		}

		if skipDepth > 0 {
			switch tok.(type) {
			case xml.StartElement:
				skipDepth++
			case xml.EndElement:
				skipDepth--
				if skipDepth == 0 {
					removeTo = tokEnd
				}
			}
			continue
		}

		switch t := tok.(type) {
		case xml.CharData:
			if strings.TrimSpace(string(t)) == "" {
				if !waitingForValue && !inKey {
					prevWSStart = tokStart
				}
			} else if inKey && strings.TrimSpace(string(t)) == "keychain-access-groups" {
				foundKAGKey = true
			}
		case xml.StartElement:
			if t.Name.Local == "key" {
				inKey = true
				if prevWSStart >= 0 {
					candidateFrom = prevWSStart
				} else {
					candidateFrom = tokStart
				}
			} else if waitingForValue {
				removeFrom = candidateFrom
				skipDepth = 1
				waitingForValue = false
			} else {
				prevWSStart = -1
			}
		case xml.EndElement:
			if t.Name.Local == "key" && inKey {
				inKey = false
				if foundKAGKey {
					waitingForValue = true
					foundKAGKey = false
				} else {
					candidateFrom = -1
				}
			}
		}
	}

	if removeFrom < 0 || removeTo < 0 {
		return data, nil
	}

	result := make([]byte, 0, int64(len(data))-(removeTo-removeFrom))
	result = append(result, data[:removeFrom]...)
	result = append(result, data[removeTo:]...)
	return result, nil
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
