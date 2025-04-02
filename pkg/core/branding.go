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

	// ExecutableNameFile returns common.ExecutableNameFile for the Chromium binaries from the given
	// binariesDir assuming they are branded with the given params.
	// Returns an error if the file destination is invalid or cannot be determined.
	ExecutableNameFile(params *common.BrandingParams, binariesDir base.Directory) (common.ExecutableNameFile, error)
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
	executableNameFile, err := branding.platform.ExecutableNameFile(&params, outputDir)
	if err != nil {
		return fmt.Errorf("the executable.name file destination is invalid: %w", err)
	}

	return executableNameFile.CreateOrUpdate()
}
