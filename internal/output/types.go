// Copyright 2023 API7.ai, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"

	"github.com/api7/cloud-cli/internal/options"
)

// Errorf prints the error message to the console and quit the program.
func Errorf(format string, args ...interface{}) {
	color.Red("ERROR: " + fmt.Sprintf(format, args...))
	os.Exit(-1)
}

// Warnf prints the warning message to the console.
func Warnf(format string, args ...interface{}) {
	color.Yellow("WARNING: " + fmt.Sprintf(format, args...))
}

// Verbosef prints the verbose message to the console.
func Verbosef(format string, args ...interface{}) {
	if options.Global.Verbose {
		color.Cyan(fmt.Sprintf(format, args...))
	}
}

// Infof prints the info message to the console.
func Infof(format string, args ...interface{}) {
	color.Cyan(fmt.Sprintf(format, args...))
}
