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
