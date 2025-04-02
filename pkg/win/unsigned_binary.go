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

package win

import (
	"fmt"
	"strings"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
)

// UnsignedBinary represents a file from which the digital signature
// has been removed, or which was never signed to begin with.
type UnsignedBinary struct {
	file base.File
}

// File returns a pointer to the underlying File of the UnsignedBinary.
func (unsignedBinary *UnsignedBinary) File() *base.File {
	return &unsignedBinary.file
}

// AbsPath returns the absolute path of the underlying file.
func (unsignedBinary UnsignedBinary) AbsPath() base.AbsPath {
	return unsignedBinary.File().AbsPath()
}

// RemoveSignature uses the Windows SDK signtool.exe to remove a digital signature
// from the given binary. It returns an UnsignedBinary if the operation succeeds,
// or an error if it fails.
//
//   - binary: The file from which the signature should be removed.
//   - returns: An UnsignedBinary containing the same file with the signature removed,
//     or an error if removal fails (e.g., if the file isn't signed or can't be accessed).
func RemoveSignature(binary base.File) (UnsignedBinary, error) {
	winSdk, err := FindWinSdk()
	if err != nil {
		return UnsignedBinary{}, err
	}

	if _, err := winSdk.SigntoolPath().AsFile(); err != nil {
		return UnsignedBinary{}, fmt.Errorf("no signtool found in Windows SDK: %w", err)
	}

	if output, err := base.ExecCommandAndGetOutput(winSdk.SigntoolPath().String(), []string{"verify", "/pa", "/v", binary.AbsPath().String()}); err != nil {
		if strings.Contains(output, "No signature found.") {
			return UnsignedBinary{binary}, nil
		} else {
			return UnsignedBinary{}, fmt.Errorf("the file is corrupted: %w", err)
		}
	}

	if err := base.ExecCommand(winSdk.SigntoolPath().String(), []string{"remove", "/s", binary.AbsPath().String()}); err != nil {
		return UnsignedBinary{}, fmt.Errorf("cannot remove signature from %s: %w", binary.AbsPath().String(), err)
	}
	return UnsignedBinary{binary}, nil
}
