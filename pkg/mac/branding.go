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
