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

package common

import (
	"os"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
)

// ExecutableNameFile represents a file that holds the name of the main executable.
// This file is named "executable.name" and is stored in a specific resources directory.
type ExecutableNameFile struct {
	// Location is the directory where the executable name file is to be located.
	Location base.Directory
	// Content is the name of the main executable to be written into the file.
	Content string
}

// CreateOrUpdate creates or updates the executable name file in the specified location.
// It writes the Content field to a file named "executable.name" inside the Location directory.
// Returns an error if the file cannot be created or written to.
func (executableName *ExecutableNameFile) CreateOrUpdate() error {
	file, err := os.Create(executableName.Location.AbsPath().Join(base.RelPathFromEntries("executable.name")).String())
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(executableName.Content)
	return err
}
