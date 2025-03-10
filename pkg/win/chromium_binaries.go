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
