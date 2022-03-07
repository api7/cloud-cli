//  Licensed to the Apache Software Foundation (ASF) under one or more
//  contributor license agreements.  See the NOTICE file distributed with
//  this work for additional information regarding copyright ownership.
//  The ASF licenses this file to You under the Apache License, Version 2.0
//  (the "License"); you may not use this file except in compliance with
//  the License.  You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS,
//  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//  See the License for the specific language governing permissions and
//  limitations under the License.

package output

import (
	"fmt"
	"os"

	"github.com/fatih/color"
)

var (
	// VerboseMode controls whether the output should be verbose.
	VerboseMode bool
)

// Errorf prints the error message to the console and quit the program.
func Errorf(format string, args ...interface{}) {
	color.Red("ERROR: " + fmt.Sprintf(format, args...))
	os.Exit(-1)
}

// Verbosef prints the verbose message to the console.
func Verbosef(format string, args ...interface{}) {
	if VerboseMode {
		color.Green(fmt.Sprintf(format, args...))
	}
}
