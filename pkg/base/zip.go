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

package base

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

// ExtractZip extracts the contents of the .zip file located at zipPath
// into the directory targetDir. It preserves the directory structure
// and creates directories as needed with os.ModePerm permissions.
//
// It returns an error if zipPath cannot be opened, if directories cannot
// be created, or if files cannot be written.
func ExtractZip(zipPath, targetDir string) (err error) {
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer func(reader *zip.ReadCloser) {
		err = reader.Close()
	}(reader)

	for _, file := range reader.File {
		targetPath := filepath.Join(targetDir, file.Name)

		if file.FileInfo().IsDir() {
			err := os.MkdirAll(targetPath, os.ModePerm)
			if err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(targetPath), os.ModePerm); err != nil {
			return err
		}

		err = func() error {
			source, err := file.Open()
			if err != nil {
				return err
			}
			defer func(source io.ReadCloser) {
				err = source.Close()
			}(source)

			destination, err := os.Create(targetPath)
			if err != nil {
				return err
			}
			defer func(destination *os.File) {
				err = destination.Close()
			}(destination)

			_, err = io.Copy(destination, source)
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}
	return nil
}
