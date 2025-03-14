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

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
)

// ChromiumBinaries represents a set of Chromium binaries on Windows,
// including the main executable (which may be renamed) and related DLLs.
type ChromiumBinaries struct {
	binariesDir     base.Directory
	chromiumExeName string
}

// getChromiumBinaries initializes a ChromiumBinaries instance for
// the given binariesDir path with the default "chrome.exe" executable name.
func getChromiumBinaries(binariesDir base.Directory) (*ChromiumBinaries, error) {
	binaries := &ChromiumBinaries{binariesDir: binariesDir, chromiumExeName: originalChromiumExeName}
	for _, binary := range binaries.List() {
		if _, err := binary.AsFile(); err != nil {
			return nil, fmt.Errorf("could not locate %s: %w", binary.Base(), err)
		}
	}

	return binaries, nil
}

// List returns a list of file paths that constitute the main chromium executable
// and its associated "chrome.dll".
func (binaries *ChromiumBinaries) List() []base.AbsPath {
	return []base.AbsPath{binaries.ChromiumExePath(), binaries.ChromeDllPath()}
}

// ChromiumExePath returns the absolute path to the primary Chromium executable,
// which may be renamed if branding is applied.
func (binaries *ChromiumBinaries) ChromiumExePath() base.AbsPath {
	return binaries.binariesDir.AbsPath().Join(base.RelPathFromEntries(binaries.chromiumExeName + ".exe"))
}

// ChromeDllPath returns the absolute path to the "chrome.dll" file in
// the same directory as the Chromium executable.
func (binaries *ChromiumBinaries) ChromeDllPath() base.AbsPath {
	return binaries.binariesDir.AbsPath().Join(base.RelPathFromEntries("chrome.dll"))
}
