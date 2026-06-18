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
	"runtime"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/mac"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/win"
)

// The utility to sign the application platform binaries.
type SignTool interface {
	// Signs the platform binary located at `binaryPath`.
	SignBinary(binaryPath string) error
}

// MacSignTool extends SignTool with entitlements-aware signing, used to sign
// helper bundles and dylibs with a filtered entitlements file that omits
// keychain-access-groups (required on macOS 26+).
type MacSignTool interface {
	SignTool
	SignBinaryWithEntitlements(binaryPath, entitlements string) error
}

// Tries to obtain the sign tool for the current platform.
func GetSignTool(params common.BrandingParams) (SignTool, error) {
	switch runtime.GOOS {
	case "windows":
		return win.GetSignToolWin(params)
	case "darwin":
		return mac.GetSignToolMac(params)
	default:
		return nil, errors.New("signing app binaries on " + runtime.GOOS + " is not supported")
	}
}
