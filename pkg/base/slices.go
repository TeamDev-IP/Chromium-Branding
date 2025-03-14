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
	"errors"
	"fmt"
)

// AnySlice takes a slice of any type T and returns a new slice of `[]any`
// containing the same elements. This allows you to treat all elements
// as the empty interface type (any).
//
// Example:
//
//	ints := []int{1, 2, 3}
//	converted := AnySlice(ints) // []any{1, 2, 3}
func AnySlice[T any](slice []T) []any {
	anySlice := []any{}
	for _, elem := range slice {
		anySlice = append(anySlice, elem)
	}
	return anySlice
}

// IndexOf searches for the given value in the provided slice and returns
// the zero-based index where the value is found, or -1 along with an error
// if the value does not exist in the slice. The type parameter T must be
// comparable to allow for equality checks.
//
// Example:
//
//	names := []string{"alice", "bob", "charlie"}
//	idx, err := IndexOf(names, "bob") // idx=1, err=nil
//	idx, err := IndexOf(names, "david") // idx=-1, err=...
func IndexOf[T comparable](slice []T, value T) (int, error) {
	for index, val := range slice {
		if val == value {
			return index, nil
		}
	}

	return -1, errors.New(fmt.Sprint("value ", value, "is not in slice ", slice))
}

// Contains indicates if the given slice contains the given value.
func Contains[T comparable](slice []T, value T) bool {
	_, err := IndexOf(slice, value)
	return err == nil
}
