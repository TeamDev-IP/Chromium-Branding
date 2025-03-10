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

package common

import (
	"encoding/json"
	"fmt"
	"os"
)

// Win holds Windows-specific branding parameters that can be
// embedded into executables or used for process naming.
type Win struct {
	// IcoPath is a path to a .ico file used as the icon
	// for the Windows executable.
	IcoPath *string

	// ExecutableName is the name of the executable file.
	ExecutableName *string

	// ProcessDisplayName is the display name shown in Task Manager
	// or when hovering over the executable in Windows.
	ProcessDisplayName *string

	// LegalCopyright holds a copyright notice that can be
	// embedded into the Windows executable file properties.
	LegalCopyright *string

	// Author represents the name of the organization or person
	// who produced the software, often displayed in file properties.
	Author *string

	// ProductName is a user-friendly name for the product.
	ProductName *string

	SignCommand string
}

// Bundle holds macOS-specific metadata about application bundles.
type Bundle struct {
	// Name is the user-friendly name of the application bundle.
	Name *string

	// Id is the unique bundle identifier (e.g., com.example.app).
	Id *string
}

// Mac holds macOS-specific branding parameters.
type Mac struct {
	// IcnsPath is a path to an .icns file used as the application icon.
	IcnsPath *string

	// Bundle contains metadata related to the macOS application bundle.
	Bundle *Bundle

	TeamId               string
	CodesignIdentity     string
	CodesignEntitlements string
	AppleId              string
	Password             string
}

// Linux holds Linux-specific branding parameters.
type Linux struct {
	// ProcessName is the name of the running process (e.g., the
	// display in system monitors or process listings).
	ProcessName *string
}

// BrandingParams holds versioning and platform-specific branding
// details used to customize executables and app bundles across
// different operating systems.
type BrandingParams struct {
	// Version specifies the version string (e.g., "1.0.0").
	Version *string

	Win   Win
	Mac   Mac
	Linux Linux
}

// GetBrandingParams reads a JSON file from paramsFilePath and
// unmarshals its contents into a BrandingParams struct. If the
// file cannot be read or unmarshaled, it returns an error.
func GetBrandingParams(paramsFilePath string) (*BrandingParams, error) {
	var params BrandingParams
	jsonText, err := os.ReadFile(paramsFilePath)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(jsonText, &params)
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	return &params, nil
}
