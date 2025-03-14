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

// The sign tool for macOS.
//
// To create one, use `GetSignToolMac`.
type SignToolMac struct {
	params common.BrandingParams
}

// Creates a new sign tool for macOS.
func GetSignToolMac(params common.BrandingParams) (*SignToolMac, error) {
	return &SignToolMac{params}, nil
}

// Signs the Chromium binaries located at `binaryPath`.
func (tool *SignToolMac) SignBinary(binaryPath string) error {
	if err := tool.sign(binaryPath); err != nil {
		return err
	} else if err := tool.verify(binaryPath); err != nil {
		return err
	}
	return nil
}

func (tool *SignToolMac) sign(binaryPath string) error {
	return base.ExecCommand("codesign",
		[]string{
			"--force",
			"--options", "runtime",
			"--timestamp",
			"--entitlements", tool.params.Mac.CodesignEntitlements,
			"--verbose",
			"--sign",
			base.GetValue(tool.params.Mac.CodesignIdentity),
			binaryPath})
}

func (tool *SignToolMac) verify(binaryPath string) error {
	return base.ExecCommand("codesign",
		[]string{
			"-vvv",
			"--deep",
			"--strict",
			binaryPath})
}
