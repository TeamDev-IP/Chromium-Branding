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

package linux

import (
	"fmt"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
)

const originalChromiumExeName = "chromium"

// GetPlatformBranding creates and returns a new LinuxBranding instance.
func GetPlatformBranding() (*LinuxBranding, error) {
	return &LinuxBranding{}, nil
}

// LinuxBranding implements branding logic specific to Linux platforms.
type LinuxBranding struct{}

func (branding *LinuxBranding) ExecutableNameFile(params *common.BrandingParams, binariesDir base.Directory) (common.ExecutableNameFile, error) {
	return common.ExecutableNameFile{
		Location: binariesDir,
		Content:  branding.ExecutableName(params),
	}, nil
}

// ExecutableName returns the Linux process name from BrandingParams if set,
// or defaults to "chromium" otherwise.
func (branding *LinuxBranding) ExecutableName(params *common.BrandingParams) string {
	if params.Linux.ProcessName != nil {
		return *params.Linux.ProcessName
	}

	return originalChromiumExeName
}

func (branding *LinuxBranding) CheckBinariesExist(binariesDir base.Directory) error {
	fileNames := []string{}
	for _, file := range binariesDir.ListFiles() {
		fileNames = append(fileNames, file.AbsPath().Base())
	}
	if !base.Contains(fileNames, originalChromiumExeName) {
		return fmt.Errorf("chromium executable has not been found in directory %s", binariesDir.AbsPath().String())
	}

	return nil
}

func (branding *LinuxBranding) Apply(params *common.BrandingParams, binariesDir base.Directory) error {
	chromiumExe, err := binariesDir.AbsPath().Join(base.RelPathFromEntries(originalChromiumExeName)).AsFile()
	if err != nil {
		return nil
	}

	if params.Linux.ProcessName != nil {
		if err := chromiumExe.Rename(*params.Linux.ProcessName); err != nil {
			return err
		}
	}

	return nil
}
