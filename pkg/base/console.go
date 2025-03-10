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
