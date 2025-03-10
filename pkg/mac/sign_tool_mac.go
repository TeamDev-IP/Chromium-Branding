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
