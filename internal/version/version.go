// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package version

import (
	"encoding/json"
	"fmt"
	"runtime"
)

// Version contains version information.
type Version struct {
	Major     string `json:"major"`
	Minor     string `json:"minor"`
	GitCommit string `json:"git_commit"`
	BuildDate string `json:"build_date"`
	GoVersion string `json:"go_version"`
	Compiler  string `json:"compiler"`
	Platform  string `json:"platform"`
}

var (
	_major     string
	_minor     string
	_buildDate string
	_gitCommit string

	// V contains the version information about the component.
	V Version
)

func init() {
	V = Version{
		Major:     _major,
		Minor:     _minor,
		BuildDate: _buildDate,
		GoVersion: runtime.Version(),
		GitCommit: _gitCommit,
		Compiler:  runtime.Compiler,
		Platform:  fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	}
}

// String shows the version info.
func (v Version) String() string {
	data, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return string(data)
}
