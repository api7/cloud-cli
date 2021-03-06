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
	"encoding/json"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	ver := Version{
		Major:     "0",
		Minor:     "1",
		GitCommit: "2ad4hz",
		BuildDate: time.Now().String(),
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}
	s := ver.String()
	var (
		ver2 Version
	)
	err := json.Unmarshal([]byte(s), &ver2)
	assert.Nil(t, err, "unmarshalling version info")
	assert.Equal(t, ver, ver2, "checking version")
}
