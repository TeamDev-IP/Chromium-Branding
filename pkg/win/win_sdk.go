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
	"fmt"
	"sort"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
)

const defaultWinSdksInstallPath = "C:\\Program Files (x86)\\Windows Kits"
const supportedArch = "x64"

var availableSdk *WinSdk

// WinSdk encapsulates a Windows SDK installation by holding a reference
// to its binaries directory.
type WinSdk struct {
	binDir base.Directory
}

// Path returns the absolute path to the Windows SDK bin directory.
func (sdk *WinSdk) Path() base.AbsPath {
	return sdk.binDir.AbsPath()
}

// SigntoolPath returns the absolute path to the signtool.exe located within
// the Windows SDK bin directory.
func (sdk *WinSdk) SigntoolPath() base.AbsPath {
	return sdk.Path().Join(base.RelPathFromEntries("signtool.exe"))
}

// FindWinSdk attempts to locate a Windows SDK installation by first checking
// the system PATH for signtool.exe and, if not found, falling back to the
// default Windows Kits installation location.
// It returns a WinSdk instance if found, or an error if no suitable SDK can be located.
func FindWinSdk() (*WinSdk, error) {
	if availableSdk != nil {
		return availableSdk, nil
	}

	winSdk, err := WinSdkFromPathEnv()
	if err == nil {
		availableSdk = winSdk
		return winSdk, nil
	}
	fmt.Println("Searching for Windows SDK in the default install location...")
	defaultWinSdk, err := DefaultWinSdk()
	if err == nil {
		availableSdk = defaultWinSdk
		fmt.Printf("Found Windows SDK at %s.\n", defaultWinSdk.Path().String())
	}
	return defaultWinSdk, err
}

// WinSdkFromPathEnv tries to locate signtool.exe by invoking the "where" command.
// If found, it converts the result to a base.Directory and returns a WinSdk instance.
func WinSdkFromPathEnv() (*WinSdk, error) {
	where, err := base.ExecCommandAndGetOutput("where", []string{"signtool"})
	if err != nil {
		return nil, err
	}
	sdkBinDir, err := base.DirectoryFromPathString(where)
	if err != nil {
		return nil, err
	}
	return &WinSdk{sdkBinDir}, err
}

// DefaultWinSdk searches the default Windows SDK install path for a Windows 10/11 SDK.
// It looks for a directory structure matching "bin/10/<version>/<supportedArch>" and returns
// the latest valid WinSdk found. If none is found, an error with installation instructions is returned.
func DefaultWinSdk() (*WinSdk, error) {
	defaultWinSdksInstallDir, err := base.DirectoryFromPathString(defaultWinSdksInstallPath)
	if err != nil {
		return nil, err
	}

	// Seek only for Win 10/11 SDKs
	sdkVersionsDir, err := defaultWinSdksInstallDir.AbsPath().Join(base.RelPathFromEntries("10", "bin")).AsDirectory()
	if err == nil {
		versions := sdkVersionsDir.ChildDirs()
		sort.SliceStable(versions, func(i, j int) bool { return versions[i].AbsPath().Base() > versions[j].AbsPath().Base() })
		for _, version := range versions {
			if sdkBinPath, err := version.AbsPath().Join(base.RelPathFromEntries(supportedArch)).AsDirectory(); err == nil {
				return &WinSdk{sdkBinPath}, nil
			} else {
				fmt.Printf("WARNING: Windows SDK of version %s has been found, but contains no binaries for the appropriate architecure.\n", version.AbsPath().Base())
			}
		}
	}

	return nil, fmt.Errorf(`
	cannot find any Windows SDKs in the default install location: %s.
	Please, install up-to date windows SDK https://developer.microsoft.com/en-us/windows/downloads/windows-sdk/
	or add the appropriate bin directory (e. g. %s\10\bin\10.0.26100.0\x64) to PATH
	if you already have Windows SDK installed to the custom location.
	`, defaultWinSdksInstallPath, defaultWinSdksInstallPath)
}
