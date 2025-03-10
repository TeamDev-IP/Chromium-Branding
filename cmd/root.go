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

package cmd

import (
	"errors"
	"fmt"
	"os"

	"github.com/TeamDev-IP/Chromium-Branding/pkg/base"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/common"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/core"
	"github.com/TeamDev-IP/Chromium-Branding/pkg/mac"
	"github.com/spf13/cobra"
)

const (
	binariesDirFlag       = "binaries_dir"
	jsonPathFlag          = "params"
	outputBinariesDirFlag = "output_dir"
	verboseFlag           = "verbose"
)

var jsonPath string
var binariesDir string
var outputDirPath string
var verbose bool

var rootCmd = &cobra.Command{
	Use:   `chromium_branding`,
	Short: `chromium_branding is a command line tool for branding Chromium binaries`,
	Long:  `chromium_branding is a command line tool for branding JxBrowser's and DotNetBrowser's Chromium binaries`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if !cmd.Flags().Changed(jsonPathFlag) {
			return errors.New("missing flag: " + jsonPathFlag)
		}
		if !cmd.Flags().Changed(binariesDirFlag) {
			return errors.New("missing flag: " + binariesDirFlag)
		}
		if !cmd.Flags().Changed(outputBinariesDirFlag) {
			return errors.New("missing flag: " + outputBinariesDirFlag)
		}

		// Obtain the branding parameters from the JSON file.
		params, err := common.GetBrandingParams(jsonPath)
		if err != nil {
			return fmt.Errorf("could not obtain branding info: %w", err)
		}

		// Set the verbose flag for logging and command line output.
		base.Verbose = verbose

		if err := core.BrandBinaries(*params, binariesDir, outputDirPath); err != nil {
			return fmt.Errorf("failed to brand Chromium binaries: %w", err)
		}

		if _, err := core.SignAppBinaries(outputDirPath, *params); err != nil {
			return err
		}

		if _, err = mac.Notarize(outputDirPath, *params); err != nil {
			return err
		}

		return nil
	},
}

// Execute adds all child commands to the root command and sets flags
// appropriately. If an error occurs while executing the CLI command,
// the process exits with a non-zero status.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(-1)
	}
}

func init() {
	rootCmd.Flags().StringVarP(&binariesDir, binariesDirFlag, "b", "",
		`absolute path to the directory with Chromium binaries`)
	rootCmd.Flags().StringVarP(&jsonPath, jsonPathFlag, "p", "",
		`absolute path to the JSON file with the custom branding parameters`)
	rootCmd.Flags().StringVarP(&outputDirPath, outputBinariesDirFlag, "o", "",
		`absolute path to the directory where the branded Chromium binaries will be stored`)
	rootCmd.Flags().BoolVarP(&verbose, verboseFlag, "v", false,
		`enable verbose output`)
}
