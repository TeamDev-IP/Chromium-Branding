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

package win

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
)

const (
	rceditUrl = "https://storage.googleapis.com/molybden-tools/rcedit/2.0.0/rcedit.zip"

	setIconFlag           = "--set-icon"
	setFileVersionFlag    = "--set-file-version"
	setProductVersionFlag = "--set-product-version"
	setVersionStringFlag  = "--set-version-string"

	fileDescriptionVersionString = "FileDescription"
	authorVersionString          = "CompanyName"
	productNameVersionString     = "ProductName"
	copyrightVersionString       = "LegalCopyright"
)

// Rcedit wraps a path to an rcedit.exe tool and provides methods
// to modify version information, icons, and other metadata of a
// Windows executable file.s
type Rcedit struct {
	toolPath string
}

// FetchRcedit downloads the rcedit tool from a predetermined URL,
// extracts it into the OS temp directory, and removes the downloaded
// zip file. It returns a pointer to an Rcedit instance pointing to
// the extracted rcedit.exe, or an error if any step fails.
func FetchRcedit() (*Rcedit, error) {
	tempDir := os.TempDir()
	rceditZipPath := filepath.Join(tempDir, "rcedit.zip")
	rceditPath := filepath.Join(tempDir, "rcedit.exe")

	if err := base.DownloadFile(rceditUrl, rceditZipPath); err != nil {
		return nil, err
	}

	if err := base.ExtractZip(rceditZipPath, tempDir); err != nil {
		return nil, err
	}

	if err := os.Remove(rceditZipPath); err != nil {
		return nil, err
	}

	return &Rcedit{toolPath: rceditPath}, nil
}

// SetIcon uses rcedit to replace the icon resource of the
// specified chromiumBinaryPath with the file at iconPath.
func (rcedit *Rcedit) SetIcon(chromiumBinary base.File, icon base.File) error {
	fmt.Println("Setting icon for " + chromiumBinary.AbsPath().String())
	return base.ExecCommand(rcedit.toolPath, []string{chromiumBinary.AbsPath().String(), setIconFlag, icon.AbsPath().String()})
}

// SetVersion uses rcedit to set both the file version and
// product version of the specified chromiumBinaryPath to version.
// It first sets the file version, then the product version.
func (rcedit *Rcedit) SetVersion(chromiumBinary base.File, version string) error {
	err := base.ExecCommand(rcedit.toolPath, []string{
		chromiumBinary.AbsPath().String(),
		setFileVersionFlag,
		version})
	if err != nil {
		return err
	}

	return base.ExecCommand(rcedit.toolPath, []string{
		chromiumBinary.AbsPath().String(),
		setProductVersionFlag,
		version})
}

// SetVersionString uses rcedit to set an arbitrary version string
// field (e.g., CompanyName, ProductName, etc.) in the specified
// chromiumBinaryPath to the provided versionStringValue.
func (rcedit *Rcedit) SetVersionString(chromiumBinary base.File, versionStringKey, versionStringValue string) error {
	return base.ExecCommand(rcedit.toolPath, []string{
		chromiumBinary.AbsPath().String(),
		setVersionStringFlag,
		versionStringKey,
		versionStringValue})
}

// SetProcessDescription sets the FileDescription version string
// for the given chromiumBinaryPath to the provided description.
func (rcedit *Rcedit) SetProcessDescription(chromiumBinary base.File, description string) error {
	fmt.Println("Setting description for " + chromiumBinary.AbsPath().String() + " : " + description)
	return rcedit.SetVersionString(chromiumBinary, fileDescriptionVersionString, description)
}

// SetAuthor sets the CompanyName version string for the given
// chromiumBinaryPath to the provided author name.
func (rcedit *Rcedit) SetAuthor(chromiumBinary base.File, author string) error {
	fmt.Println("Setting author for " + chromiumBinary.AbsPath().String() + " : " + author)
	return rcedit.SetVersionString(chromiumBinary, authorVersionString, author)
}

// SetProductName sets the ProductName version string for the
// given chromiumBinaryPath to the provided product name.
func (rcedit *Rcedit) SetProductName(chromiumBinary base.File, productName string) error {
	fmt.Println("Setting product name for " + chromiumBinary.AbsPath().String() + " : " + productName)
	return rcedit.SetVersionString(chromiumBinary, productNameVersionString, productName)
}

// SetCopyright sets the LegalCopyright version string
// for the given chromiumBinaryPath to the provided text.
func (rcedit *Rcedit) SetCopyright(chromiumBinary base.File, copyright string) error {
	fmt.Println("Setting copyright for " + chromiumBinary.AbsPath().String() + " : " + copyright)
	return rcedit.SetVersionString(chromiumBinary, copyrightVersionString, copyright)
}
