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

package mac

import (
	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
)

// MacBranding implements branding logic specific to macOS.
// It overrides Info.plist properties, configures icons, and renames
// the top-level .app directory to match the branding parameters.
type MacBranding struct{}

// GetPlatformBranding returns a new MacBranding instance for macOS platforms.
func GetPlatformBranding() (*MacBranding, error) {
	return &MacBranding{}, nil
}

// ApplyToBundle applies macOS-specific branding changes to the provided ChromiumAppBundle.
// This includes setting Info.plist keys (bundle identifier, bundle name, etc.)
// and optionally replacing the .icns file if params.Mac.IcnsPath
// is provided.
//
// Parameters:
//   - params: The BrandingParams containing user-specified overrides.
//   - appBundle: A ChromiumAppBundle that points to the .app directory to brand.
//
// Returns an error if any of the file or plist operations fail.
func (branding *MacBranding) ApplyToBundle(params *common.BrandingParams, appBundle ChromiumAppBundle) error {
	appInfo := AppInfo{
		Version:        params.Version,
		ExecutableName: params.Mac.Bundle.Name,
		BundleId:       params.Mac.Bundle.Id,
	}

	if err := cofigureBundlePlist(appBundle, appInfo); err != nil {
		return err
	}

	if iconExpectedFor(appBundle) && params.Mac.IcnsPath != nil {
		iconPath, err := base.AbsPathFromPathString(*params.Mac.IcnsPath)
		if err != nil {
			return err
		}
		if err := configureBundleIcon(appBundle, iconPath); err != nil {
			return err
		}
	}

	return nil
}

func (branding *MacBranding) CheckBinariesExist(binariesDir base.Directory) error {
	if _, err := GetChromiumAppBundle(binariesDir, originalChromiumAppBundleName); err != nil {
		return err
	}
	return nil
}

func (branding *MacBranding) Apply(params *common.BrandingParams, binariesDir base.Directory) error {
	rootBundle, err := GetChromiumAppBundle(binariesDir, originalChromiumAppBundleName)
	if err != nil {
		return err
	}

	if params.Mac.Bundle.Name != nil {
		if err := rootBundle.Rename(*params.Mac.Bundle.Name); err != nil {
			return err
		}
	}

	allBundles := append([]ChromiumAppBundle{rootBundle.ChromiumAppBundle()}, rootBundle.Helpers()...)

	for _, bundle := range allBundles {
		if err := branding.ApplyToBundle(params, bundle); err != nil {
			return err
		}
	}

	return nil
}

func (branding *MacBranding) ExecutableName(params *common.BrandingParams) string {
	if params.Mac.Bundle.Name != nil {
		return *params.Mac.Bundle.Name
	}
	return "Chromium"
}

var bundleIdProperties = []string{
	"CFBundleIdentifier",
}

var bundleNameProperties = []string{
	"CFBundleName",
	"CFBundleDisplayName",
	"CFBundleExecutable",
}

var bundleVersionProperties = []string{
	"CFBundleShortVersionString",
}

func iconExpectedFor(appBundle ChromiumAppBundle) bool {
	return appBundle.GetType() == CrBundleMain || appBundle.GetType() == CrBundleHelperAlerts
}

func cofigureBundlePlist(bundle ChromiumAppBundle, appInfo AppInfo) error {
	return base.AnyErrorFrom(
		overrideBundleName(bundle, appInfo.ExecutableName),
		overrideBundleId(bundle, appInfo.BundleId),
		overrideBundleVersion(bundle, appInfo.Version),
	)
}

func configureBundleIcon(appBundle ChromiumAppBundle, iconPath base.AbsPath) error {
	bundleIcon, err := appBundle.IconPath().AsFile()
	if err != nil {
		return err
	}
	customIconFile, err := iconPath.AsFile()
	if err != nil {
		return err
	}

	if err := bundleIcon.Replace(customIconFile); err != nil {
		return err
	}

	return nil
}

func overrideBundleProperties(bundle ChromiumAppBundle, properties []string, value string) error {
	plist, err := bundle.PlistFilePath().AsFile()
	if err != nil {
		return err
	}

	for _, prop := range properties {
		if err := setPlistProperty(plist, prop, value); err != nil {
			return err
		}
	}

	return nil
}

func setPlistProperty(plist base.File, key, value string) error {
	return base.ExecCommand("defaults", []string{"write", plist.AbsPath().String(), key, "\"" + value + "\""})
}

func overrideBundleName(bundle ChromiumAppBundle, name *string) error {
	if name == nil {
		return nil
	}
	return overrideBundleProperties(bundle, bundleNameProperties, getBrandedCrBundleExeName(bundle.GetType(), *name))
}

func overrideBundleId(bundle ChromiumAppBundle, id *string) error {
	if id == nil {
		return nil
	}
	return overrideBundleProperties(bundle, bundleIdProperties, getBrandedCrBundleId(bundle.GetType(), *id))
}

func overrideBundleVersion(bundle ChromiumAppBundle, version *string) error {
	if version == nil {
		return nil
	}
	return overrideBundleProperties(bundle, bundleVersionProperties, *version)
}
