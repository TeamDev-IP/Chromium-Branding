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
	"fmt"
	"runtime"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/linux"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/mac"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/win"
)

// BrandBinaries applies platform-specific branding to the binaries located in
// the provided binariesDirPath and copies them to outputDirectoryPath.
//
// Parameters:
//   - params: The BrandingParams struct with platform-specific metadata.
//   - binariesDirPath: The path (absolute or relative) to the source binaries directory.
//   - outputDirectoryPath: The destination path where the branded binaries should reside.
//
// Returns an error if any file I/O or branding operations fail.
func BrandBinaries(params common.BrandingParams, binariesDirPath string, outputDirPath string) error {
	branding, err := GetBrandingForParams(params)
	if err != nil {
		return err
	}

	binariesDir, err := base.DirectoryFromPathString(binariesDirPath)
	if err != nil {
		return err
	}

	if err := branding.CheckBinariesExist(binariesDir); err != nil {
		return err
	}

	outputDirAbsPath, err := base.AbsPathFromPathString(outputDirPath)
	if err != nil {
		return err
	}

	outputDir, err := copyBinaries(binariesDir, outputDirAbsPath)
	if err != nil {
		return err
	}
	return brandBinariesInDirectory(branding, params, outputDir)
}

// GetBrandingForParams returns a new Branding instance populated
// with the provided BrandingParams and the appropriate PlatformBranding
// for the current runtime OS.
//
// Returns an error if the platform-specific branding cannot be determined.
func GetBrandingForParams(params common.BrandingParams) (*Branding, error) {
	branding := Branding{params: params}
	if platformBranding, err := GetPlatformBranding(); err != nil {
		return nil, err
	} else {
		branding.platform = platformBranding
	}

	return &branding, nil
}

// PlatformBranding defines the interface for applying platform-specific
// branding logic and retrieving the main executable name.
//
// Implementations of this interface exist for each supported runtime.GOOS:
// Windows, macOS, and Linux.
type PlatformBranding interface {
	// CheckBinaries ensures that the binariesDir contains the Chromium binaries.
	CheckBinariesExist(binariesDir base.Directory) error

	// Apply applies the branding to the binaries located in binariesDir
	// according to the provided BrandingParams.
	Apply(params *common.BrandingParams, binariesDir base.Directory) error

	// ExecutableName returns the name of the main executable file,
	// derived from the BrandingParams.
	ExecutableName(params *common.BrandingParams) string
}

// Branding wraps a set of BrandingParams and a PlatformBranding
// implementation to manage applying branding across different
// operating systems.
type Branding struct {
	params   common.BrandingParams
	platform PlatformBranding
}

// GetPlatformBranding inspects runtime.GOOS and returns the corresponding
// PlatformBranding instance. An error is returned if the platform is not supported.
func GetPlatformBranding() (PlatformBranding, error) {
	if runtime.GOOS == "windows" {
		return win.GetPlatformBranding()
	} else if runtime.GOOS == "darwin" {
		return mac.GetPlatformBranding()
	} else if runtime.GOOS == "linux" {
		return linux.GetPlatformBranding()
	} else {
		return nil, errors.New("Branding is not available for platform: " + runtime.GOOS)
	}
}

// Apply calls the underlying platform's Apply method, passing the stored
// BrandingParams and the provided binariesDir.
func (branding *Branding) CheckBinariesExist(binariesDir base.Directory) error {
	return branding.platform.CheckBinariesExist(binariesDir)
}

// Apply calls the underlying platform's Apply method, passing the stored
// BrandingParams and the provided binariesDir.
func (branding *Branding) Apply(binariesDir base.Directory) error {
	return branding.platform.Apply(&branding.params, binariesDir)
}

func copyBinaries(binariesDir base.Directory, outputDirPath base.AbsPath) (base.Directory, error) {
	if binariesDir.AbsPath().String() == outputDirPath.String() {
		return base.Directory{}, fmt.Errorf("chromium binaries directory %s must not be equal to the output directory", binariesDir.AbsPath().String())
	}
	if err := binariesDir.Copy(outputDirPath); err != nil {
		return base.Directory{}, err
	}
	return outputDirPath.AsDirectory()
}

func brandBinariesInDirectory(branding *Branding, params common.BrandingParams, outputDir base.Directory) error {
	if err := branding.Apply(outputDir); err != nil {
		return err
	}
	if err := common.WriteExecutableName(branding.platform.ExecutableName(&params), outputDir); err != nil {
		return err
	}

	return nil
}
