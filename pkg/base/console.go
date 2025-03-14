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
	"fmt"
)

var Verbose = false

// Flush clears the current line in the console.
func Flush() {
	fmt.Print("\033[0K")
}

// Print prints the message to the console.
func Print(message string) {
	fmt.Print(message)
	Flush()
}

// Println prints the message to the console and adds a newline.
func Println(message string) {
	fmt.Println(message)
	Flush()
}

// Printf prints the formatted message to the console.
func Printf(format string, a ...any) {
	fmt.Printf(format, a...)
	Flush()
}

// Log prints the message to the console if the verbose mode is enabled.
func Log(message string) {
	if Verbose && message != "" {
		fmt.Println(message)
	}
}

// Logf prints the formatted message to the console if the verbose mode is enabled.
func Logf(format string, a ...any) {
	if Verbose {
		fmt.Printf(format, a...)
		Print("\n")
		Flush()
	}
}
