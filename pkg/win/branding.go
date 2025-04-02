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

func (branding *WinBranding) ExecutableNameFile(params *common.BrandingParams, binariesDir base.Directory) (common.ExecutableNameFile, error) {
	return common.ExecutableNameFile{
		Location: binariesDir,
		Content:  branding.ExecutableName(params),
	}, nil
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
func (branding *WinBranding) SetIcon(binaryFile UnsignedBinary, icon base.File) error {
	return branding.rceditTool.SetIcon(binaryFile, icon)
}

// SetFileDescription updates the FileDescription resource of the
// specified Windows executable or DLL to the provided description.
// This is often displayed in Task Manager or file properties.
func (branding *WinBranding) SetFileDescription(description string, binaryFile UnsignedBinary) {
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

	initialChromiumExecutable, err := binariesToBrand.ChromiumExePath().AsFile()
	if err != nil {
		return err
	}

	chromiumExecutable, err := RemoveSignature(initialChromiumExecutable)
	if err != nil {
		return err
	}

	if params.Win.ExecutableName != nil {
		newChromiumExeFilename := *params.Win.ExecutableName + ".exe"
		fmt.Println("Renaming ", chromiumExecutable.AbsPath().String(), "==>", newChromiumExeFilename)

		if err := chromiumExecutable.File().Rename(newChromiumExeFilename); err != nil {
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

		initialChromeDll, err := binariesToBrand.ChromeDllPath().AsFile()
		if err != nil {
			return err
		}

		chromeDll, err := RemoveSignature(initialChromeDll)
		if err != nil {
			return err
		}

		for _, binaryFile := range []UnsignedBinary{chromiumExecutable, chromeDll} {
			if err := branding.SetIcon(binaryFile, icon); err != nil {
				return err
			}
		}
	}

	return nil
}
