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
	"errors"
	"fmt"
	"strings"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
)

const originalChromiumAppBundleName = "Chromium"

// CrBundleType identifies different types of Chromium app bundles on macOS.
// These constants define whether the bundle is the main app or a helper.
type CrBundleType int

const (
	// CrBundleMain is the main Chromium .app bundle.
	CrBundleMain CrBundleType = iota - 2

	// CrBundleHelper is the default helper .app bundle.
	CrBundleHelper

	// CrBundleHelperAlerts is a helper specifically for alert notifications.
	CrBundleHelperAlerts

	// CrBundleHelperGPU is a helper for GPU processes.
	CrBundleHelperGPU

	// CrBundleHelperPlugin is a helper for plugin processes.
	CrBundleHelperPlugin

	// CrBundleHelperRenderer is a helper for renderer processes.
	CrBundleHelperRenderer
)

// AppInfo holds metadata about a Chromium macOS bundle, including
// version, executable name, and bundle identifier.
type AppInfo struct {
	// Version indicates the application version (e.g., "1.0.0").
	Version *string

	// ExecutableName is the name of the main executable within the bundle.
	ExecutableName *string

	// BundleId is the unique identifier for this bundle (e.g., "com.example.chromium").
	BundleId *string
}

// ChromiumAppBundle describes a generic Chromium macOS .app bundle,
// whether main or helper. It includes methods to retrieve the path,
// icon path, and Info.plist location, as well as the type of bundle.
type ChromiumAppBundle interface {
	// IconPath returns the path to the .icns icon file inside this bundle.
	IconPath() base.AbsPath

	// PlistFilePath returns the path to this bundle's Info.plist file.
	PlistFilePath() base.AbsPath

	// Path returns the absolute path to the .app bundle directory.
	Path() base.AbsPath

	// GetType returns the CrBundleType indicating what kind of bundle this is
	// (main app, helper, renderer, etc.).
	GetType() CrBundleType
}

// ChromiumMainAppBundle extends ChromiumAppBundle for the main .app
// and includes methods for renaming and enumerating helper bundles.
type ChromiumMainAppBundle interface {
	ChromiumAppBundle() ChromiumAppBundle

	// Rename updates the main bundle directory name (and internal binary name)
	// to the specified newName. It also renames any discovered helper bundles.
	Rename(newName string) error

	// Helpers returns a slice of ChromiumAppBundle representing
	// the various helper bundles within this main app.
	Helpers() []ChromiumAppBundle
}

// GetChromiumAppBundle locates the main Chromium .app bundle inside
// the provided binariesDir by using the brandedAppName. It checks that
// the expected .app directory exists and returns a ChromiumMainAppBundle
// interface to manipulate it.
//
// Returns an error if the directory does not exist or cannot be validated
// as a directory.
func GetChromiumAppBundle(binariesDir base.Directory, brandedAppName string) (ChromiumMainAppBundle, error) {
	bundleDirRelpath := base.RelPathFromEntries(getBrandedCrBundleName(CrBundleMain, brandedAppName))
	_, err := binariesDir.AbsPath().Join(bundleDirRelpath).AsDirectory()
	if err != nil {
		return nil, fmt.Errorf("failed to locate Chromium app bundle in %s: %w", binariesDir.AbsPath().String(), err)
	}
	return &ChromiumBundle{location: binariesDir, brandedName: brandedAppName}, nil
}

// ChromiumBundle implements ChromiumMainAppBundle for the main .app bundle.
// It knows how to locate helper bundles, rename itself, and provide paths
// to the main .app Info.plist, icon, etc.
type ChromiumBundle struct {
	location    base.Directory
	brandedName string
}

// ChromiumHelperBundle implements ChromiumAppBundle for one of the Chromium helper .app
// bundles. It references its parent ChromiumBundle and knows its own helper type.
type ChromiumHelperBundle struct {
	parent     *ChromiumBundle
	helperType CrBundleType
}

func (bundle *ChromiumBundle) Rename(newName string) error {
	rootPath := bundle.Path()
	rootDir, err := rootPath.AsDirectory()
	if err != nil {
		return err
	}
	if err := rootDir.Rename(getBrandedCrBundleName(CrBundleMain, newName)); err != nil {
		return err
	}
	bundle.brandedName = newName

	initialName := strings.ReplaceAll(rootPath.Base(), ".app", "")
	exeFile, err := bundle.Path().Join(base.RelPathFromEntries("Contents", "MacOS", initialName)).AsFile()
	if err != nil {
		return err
	}

	if err := exeFile.Rename(newName); err != nil {
		return err
	}

	for _, helperType := range crHelperTypes {
		helperBundle, err := bundle.FindHelper(helperType, initialName)
		if err != nil {
			return err
		}
		updatedHelperPath := helperBundle.Path()
		currentHelerRelPath := base.RelPathFromEntries(getBrandedCrBundleName(helperType, initialName))
		currentHelperDir, err := updatedHelperPath.Parent().Join(currentHelerRelPath).AsDirectory()
		if err != nil {
			return err
		}

		currentHelperExeName := strings.ReplaceAll(getBrandedCrBundleName(helperType, initialName), ".app", "")
		newHelperExeName := strings.ReplaceAll(getBrandedCrBundleName(helperType, newName), ".app", "")
		helperExeFile, err := currentHelperDir.AbsPath().Join(base.RelPathFromEntries("Contents", "MacOS", currentHelperExeName)).AsFile()
		if err != nil {
			return err
		}

		if err := helperExeFile.Rename(newHelperExeName); err != nil {
			return err
		}

		if err := currentHelperDir.Rename(getBrandedCrBundleName(helperType, newName)); err != nil {
			return err
		}
	}

	return nil
}

func (bundle *ChromiumBundle) GetType() CrBundleType {
	return CrBundleMain
}

func (bundle *ChromiumBundle) ChromiumAppBundle() ChromiumAppBundle {
	return bundle
}

func (bundle *ChromiumBundle) IconPath() base.AbsPath {
	return bundle.Path().Join(bundleIconRelPath)
}

func (bundle *ChromiumBundle) Helpers() []ChromiumAppBundle {
	helpersList := []ChromiumAppBundle{}
	for _, helperType := range crHelperTypes {
		helper, err := bundle.FindHelper(helperType, bundle.brandedName)
		if err != nil {
			continue
		}
		helpersList = append(helpersList, helper)
	}

	return helpersList
}

func (bundle *ChromiumBundle) Path() base.AbsPath {
	return bundle.location.AbsPath().Join(base.RelPathFromEntries(getBrandedCrBundleName(CrBundleMain, bundle.brandedName)))
}

func (bundle *ChromiumHelperBundle) GetType() CrBundleType {
	return bundle.helperType
}

func (bundle *ChromiumHelperBundle) IconPath() base.AbsPath {
	return bundle.Path().Join(bundleIconRelPath)
}

func (helper *ChromiumHelperBundle) Path() base.AbsPath {
	bundleHelpersPath, _ := getAppBundleHelpersPath(helper.parent)

	bundleHelperRelpath := base.RelPathFromEntries(getBrandedCrBundleName(helper.helperType, helper.parent.brandedName))
	return bundleHelpersPath.AbsPath().Join(bundleHelperRelpath)
}

func (helper *ChromiumHelperBundle) PlistFilePath() base.AbsPath {
	return helper.Path().Join(base.RelPathFromEntries("Contents", "Info.plist"))
}

func (bundle *ChromiumBundle) PlistFilePath() base.AbsPath {
	return bundle.Path().Join(base.RelPathFromEntries("Contents", "Info.plist"))
}

func (bundle *ChromiumBundle) FindHelper(helperType CrBundleType, brandedAppName string) (*ChromiumHelperBundle, error) {
	bundleHelpersDir, err := bundle.getHelpersDir()
	if err != nil {
		return nil, err
	}

	helperDirRelpath := base.RelPathFromEntries(getBrandedCrBundleName(helperType, brandedAppName))
	_, err = bundleHelpersDir.AbsPath().Join(helperDirRelpath).AsDirectory()
	if err != nil {
		return nil, err
	}

	return &ChromiumHelperBundle{parent: bundle, helperType: helperType}, nil
}

func (bundle *ChromiumBundle) getHelpersDir() (base.Directory, error) {
	return getAppBundleHelpersPath(bundle)
}

var crHelperTypes = []CrBundleType{
	CrBundleHelper,
	CrBundleHelperAlerts,
	CrBundleHelperGPU,
	CrBundleHelperPlugin,
	CrBundleHelperRenderer,
}

var crHelperTypeNames = []string{
	"Alerts",
	"GPU",
	"Plugin",
	"Renderer",
}

var bundleIconRelPath = base.RelPathFromEntries("Contents", "Resources", "app.icns")

func getAppBundleHelpersPath(bundle ChromiumMainAppBundle) (base.Directory, error) {
	bundleAppHelpersVersionsRelpath := base.RelPathFromEntries("Contents", "Frameworks", "Chromium Framework.framework", "Versions")
	versionsDir, err := bundle.ChromiumAppBundle().Path().Join(bundleAppHelpersVersionsRelpath).AsDirectory()
	if err != nil {
		return base.Directory{}, err
	}

	versions := versionsDir.ChildDirs()
	if len(versions) != 1 {
		return base.Directory{}, errors.New("no versions Chromium Framework found in the Chromium app bundle")
	}

	return versions[0].AbsPath().Join(base.RelPathFromEntries("Helpers")).AsDirectory()
}

func getBrandedCrBundleName(bundleType CrBundleType, brandedAppName string) string {
	return getBrandedCrBundleExeName(bundleType, brandedAppName) + ".app"
}

func getBrandedCrBundleExeName(bundleType CrBundleType, brandedAppName string) string {
	if bundleType == CrBundleMain {
		return brandedAppName
	} else if bundleType == CrBundleHelper {
		return fmt.Sprintf("%s Helper", brandedAppName)
	} else {
		return fmt.Sprintf("%s Helper (%s)", brandedAppName, crHelperTypeNames[bundleType])
	}
}

func getBrandedCrBundleId(bundleType CrBundleType, brandedBundleId string) string {
	if bundleType == CrBundleMain {
		return brandedBundleId
	} else if bundleType == CrBundleHelperAlerts {
		return brandedBundleId + ".framework.AlertNotificationService"
	} else if bundleType == CrBundleHelperPlugin {
		return brandedBundleId + ".helper.plugin"
	} else if bundleType == CrBundleHelperRenderer {
		return brandedBundleId + ".helper.renderer"
	} else {
		return brandedBundleId + ".helper"
	}
}
