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

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
)

const originalChromiumExeName = "chromium"

// WinBranding implements platform-specific branding logic for Windows.
type WinBranding struct {
	rceditTool *Rcedit
}

// GetPlatformBranding initializes a new WinBranding instance by
// fetching the rcedit tool from a known URL, extracting it, and
// storing its location in a WinBranding struct.
//
// Returns an error if the tool cannot be downloaded or extracted.
func GetPlatformBranding() (*WinBranding, error) {
	if rceditTool, err := FetchRcedit(); err != nil {
		return nil, err
	} else {
		return &WinBranding{rceditTool: rceditTool}, nil
	}
}

// ExecutableName returns the user-specified Windows executable name
// if set in BrandingParams, or falls back to original chromium executable name.
func (branding *WinBranding) ExecutableName(params *common.BrandingParams) string {
	if params.Win.ExecutableName != nil {
		return *params.Win.ExecutableName
	}
	return originalChromiumExeName
}

// SetIcon calls the underlying rcedit tool to replace the icon
// resource in the specified Windows executable or DLL file.
func (branding *WinBranding) SetIcon(binaryFile base.File, icon base.File) error {
	return branding.rceditTool.SetIcon(binaryFile, icon)
}

// SetFileDescription updates the FileDescription resource of the
// specified Windows executable or DLL to the provided description.
// This is often displayed in Task Manager or file properties.
func (branding *WinBranding) SetFileDescription(description string, binaryFile base.File) {
	branding.rceditTool.SetProcessDescription(binaryFile, description)
}

func (branding *WinBranding) CheckBinariesExist(binariesDir base.Directory) error {
	if _, err := getChromiumBinaries(binariesDir); err != nil {
		return err
	}
	return nil
}

func (branding *WinBranding) Apply(params *common.BrandingParams, binariesDir base.Directory) error {
	binariesToBrand, err := getChromiumBinaries(binariesDir)
	if err != nil {
		return err
	}

	chromiumExecutable, err := binariesToBrand.ChromiumExePath().AsFile()
	if err != nil {
		return err
	}

	if params.Win.ExecutableName != nil {
		newChromiumExeFilename := *params.Win.ExecutableName + ".exe"
		fmt.Println("Renaming ", chromiumExecutable.AbsPath().String(), "==>", newChromiumExeFilename)

		if err := chromiumExecutable.Rename(newChromiumExeFilename); err != nil {
			return errors.Join(err, errors.New("failed to rename "+chromiumExecutable.AbsPath().String()))
		}
	}

	if params.Win.Author != nil {
		if err := branding.rceditTool.SetAuthor(chromiumExecutable, *params.Win.Author); err != nil {
			return err
		}
	}

	if params.Win.ProductName != nil {
		if err := branding.rceditTool.SetProductName(chromiumExecutable, *params.Win.ProductName); err != nil {
			return err
		}
	}

	if params.Version != nil {
		if err := branding.rceditTool.SetVersion(chromiumExecutable, *params.Version); err != nil {
			return err
		}
	}

	if params.Win.ProcessDisplayName != nil {
		if err := branding.rceditTool.SetProcessDescription(chromiumExecutable, *params.Win.ProcessDisplayName); err != nil {
			return err
		}
	}

	if params.Win.LegalCopyright != nil {
		if err := branding.rceditTool.SetCopyright(chromiumExecutable, *params.Win.LegalCopyright); err != nil {
			return err
		}
	}

	if params.Win.IcoPath != nil {
		icon, err := base.FileFromPathString(*params.Win.IcoPath)
		if err != nil {
			return err
		}

		chromeDll, err := binariesToBrand.ChromeDllPath().AsFile()
		if err != nil {
			return err
		}

		for _, binaryFile := range []base.File{chromiumExecutable, chromeDll} {
			if err := branding.SetIcon(binaryFile, icon); err != nil {
				return err
			}
		}
	}

	return nil
}
