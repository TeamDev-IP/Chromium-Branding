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
