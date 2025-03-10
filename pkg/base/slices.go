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
