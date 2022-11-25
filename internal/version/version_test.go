// Copyright 2022 API7.ai, Inc
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

package version

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	ver := Version{
		Major:     "0",
		Minor:     "1",
		GitCommit: "2ad4hz",
		BuildDate: "2022-11-25 10:15:14.230457 +0800 CST m=+0.001464835",
		GoVersion: "go1.19.1",
		Compiler:  runtime.Compiler,
		Platform:  "darwin/arm64",
	}
	s := ver.String()
	res := "version 0.1, git_commit 2ad4hz, build_date 2022-11-25 10:15:14.230457 +0800 CST m=+0.001464835, go_version go1.19.1, compiler gc, platform darwin/arm64"
	assert.Equal(t, res, s, "checking version")

}
