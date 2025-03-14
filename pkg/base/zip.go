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
