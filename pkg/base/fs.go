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
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	cp "github.com/otiai10/copy"
)

// AbsPath represents an absolute path on the filesystem.
type AbsPath struct {
	absPath string
}

// RelPath represents a relative path on the filesystem.
type RelPath struct {
	pathEntries []string
}

// File represents a path that points to a single file on the filesystem.
type File struct {
	path AbsPath
}

// Directory represents a path that points to a directory on the filesystem.
type Directory struct {
	path AbsPath
}

// AbsPathFromPathString returns an AbsPath corresponding to the given (possibly relative) path string.
// An error is returned if the path cannot be resolved to an absolute path.
func AbsPathFromPathString(path string) (AbsPath, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return AbsPath{}, err
	}
	return AbsPath{absPath}, nil
}

// FileFromPathString tries to convert the provided path string to AbsPath and then to File.
// If an error occurs on any stage, the function returns it.
func FileFromPathString(path string) (File, error) {
	fileAbsPath, err := AbsPathFromPathString(path)
	if err != nil {
		return File{}, err
	}
	return fileAbsPath.AsFile()
}

// FileFromPathString tries to convert the provided path string to AbsPath and then to Directory.
// If an error occurs on any stage, the function returns it.
func DirectoryFromPathString(path string) (Directory, error) {
	fileAbsPath, err := AbsPathFromPathString(path)
	if err != nil {
		return Directory{}, err
	}
	return fileAbsPath.AsDirectory()
}

// RelPathFromEntries constructs a RelPath from the given path entries.
func RelPathFromEntries(entries ...string) RelPath {
	return RelPath{pathEntries: entries}
}

// Base returns the last element of the absolute path (akin to filepath.Base).
func (path AbsPath) Base() string {
	return filepath.Base(path.absPath)
}

// Parent returns the absolute path of the parent directory of the given path.
func (path AbsPath) Parent() AbsPath {
	return AbsPath{filepath.Dir(path.String())}
}

// String returns the underlying string representation of the AbsPath.
func (path AbsPath) String() string {
	return path.absPath
}

// AsFile treats the AbsPath as a file path. It returns a File and validates
// that the path refers to a file.
// An error is returned if validation fails.
func (path AbsPath) AsFile() (File, error) {
	file := File{path: path}
	if err := file.validate(); err != nil {
		return File{}, err
	}
	return file, nil
}

// AsDirectory treats the AbsPath as a directory path. It returns a Directory
// and validates that the path refers to a directory.
// An error is returned if validation fails.
func (path AbsPath) AsDirectory() (Directory, error) {
	directory := Directory{path: path}
	if err := directory.validate(); err != nil {
		return Directory{}, err
	}
	return directory, nil
}

// Join appends a relative path (RelPath) to the AbsPath and returns the resulting absolute path.
func (path AbsPath) Join(relPath RelPath) AbsPath {
	absPath := []string{path.absPath}
	entries := append(absPath, relPath.pathEntries...)
	return AbsPath{filepath.Join(entries...)}
}

// AbsPath returns the absolute path of the directory.
func (directory Directory) AbsPath() AbsPath {
	return directory.path
}

// Copy recursively copies the directory to the specified destination path.
// If the destination already exists, it merges or overwrites contents as determined by cp.Copy.
func (directory Directory) Copy(destination AbsPath) error {
	return cp.Copy(directory.AbsPath().String(), destination.String())
}

// Rename renames the directory to the given new name (changing only the last path element).
// If the new name matches the existing name, it is a no-op.
func (directory *Directory) Rename(newName string) error {
	newPath, err := renameLastPathEntry(directory.AbsPath(), newName)
	if err != nil {
		return err
	}
	directory.path = newPath
	return nil
}

// ChildDirs walks the *immediate* subdirectories of this directory and returns them.
// It does not descend recursively into deeper levels, but it does create a Directory entry
// for each subdirectory found directly under this directory.
func (directory Directory) ChildDirs() []Directory {
	entries := []Directory{}
	root := directory.AbsPath().String()
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if filepath.Dir(path) == root && info.IsDir() {
			abs, err := AbsPathFromPathString(path)
			if err != nil {
				return err
			}
			dir, err := abs.AsDirectory()
			if err != nil {
				return err
			}
			entries = append(entries, dir)
		}
		return nil
	})

	return entries
}

// ListFiles walks the directory to find all immediate child files (non-directories)
// and returns them as a slice of File.
// Any error encountered during the walk is returned from the callback, aborting
// the walk, but is not otherwise exposed by this function.
func (directory Directory) ListFiles() []File {
	entries := []File{}
	root := directory.AbsPath().String()
	filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if filepath.Dir(path) == root && !info.IsDir() {
			abs, err := AbsPathFromPathString(path)
			if err != nil {
				return err
			}
			file, err := abs.AsFile()
			if err != nil {
				return err
			}
			entries = append(entries, file)
		}
		return nil
	})

	return entries
}

// Returns the absolute path of the `file`.
func (file File) AbsPath() AbsPath {
	return file.path
}

// Open validates the file path and, if valid, opens the file from the filesystem.
// The caller is responsible for closing the returned *os.File handle.
func (file File) Open() (*os.File, error) {
	if err := file.validate(); err != nil {
		return nil, err
	}
	return os.Open(file.path.absPath)
}

// Rewrite truncates the file and calls the provided writer function to write new content.
// If the file path is invalid or cannot be opened for writing, an error is returned.
func (file File) Rewrite(writer func(io.Writer) error) error {
	out, err := os.OpenFile(file.AbsPath().String(), os.O_TRUNC|os.O_APPEND|os.O_RDWR, os.ModePerm)
	if err != nil {
		return err
	}
	defer out.Close()
	return writer(out)
}

// Read validates the file path and returns the file's content as a byte slice.
// If reading fails or the path does not represent a file, an error is returned.
func (file File) Read() ([]byte, error) {
	if err := file.validate(); err != nil {
		return nil, err
	}
	text, err := os.ReadFile(file.path.absPath)
	return text, err
}

// Remove deletes the file from the filesystem. Returns an error if removal fails.
func (file File) Remove() error {
	return os.Remove(file.AbsPath().absPath)
}

// Copy reads the entire content of the file and writes it into a new file at dstPath.
// Returns an error if reading from the original or writing to the destination fails.
func (file File) Copy(dstPath AbsPath) error {
	srcBytes, err := file.Read()
	if err != nil {
		return err
	}
	return os.WriteFile(dstPath.absPath, srcBytes, os.ModePerm)
}

// Replace replaces the content of the current file with the content of the provided other File.
// Returns an error if reading from 'other' fails or writing to the current file fails.
func (file File) Replace(other File) error {
	contents, err := other.Read()
	if err != nil {
		return err
	}
	return os.WriteFile(file.path.absPath, contents, os.ModePerm)
}

// Rename renames the file to the given new name (changing only the last path element).
// If the new name matches the existing name, it is a no-op.
func (file *File) Rename(newName string) error {
	newPath, err := renameLastPathEntry(file.AbsPath(), newName)
	if err != nil {
		return err
	}
	file.path = newPath
	return nil
}

// CopyFile copies a file from sourcePath to destPath using io.Copy.
// It creates the destination file if it does not exist. Returns an error on failure.
func CopyFile(sourcePath, destPath string) error {
	sourceFile, err := os.Open(sourcePath)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destinationFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destinationFile.Close()

	if _, err = io.Copy(destinationFile, sourceFile); err != nil {
		return err
	}
	return nil
}

func renameLastPathEntry(entryPath AbsPath, newName string) (AbsPath, error) {
	if newName == entryPath.Base() {
		return entryPath, nil
	}
	newPath := entryPath.Parent().Join(RelPathFromEntries(newName))
	return newPath, os.Rename(entryPath.String(), newPath.String())
}

func (directory Directory) validate() error {
	stat, err := os.Stat(directory.path.absPath)
	if err != nil {
		return err
	}
	if !stat.IsDir() {
		return fmt.Errorf("the path %s does not represent a directory", directory.path.absPath)
	}

	return nil
}

func (file File) validate() error {
	stat, err := os.Stat(file.path.absPath)
	if err != nil {
		return err
	}
	if stat.IsDir() {
		return fmt.Errorf("the path %s represents a directory, not a file", file.path.absPath)
	}

	return nil
}

func PathExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
