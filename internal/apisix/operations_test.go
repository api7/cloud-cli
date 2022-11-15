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

package apisix

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/options"
)

func TestReload(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name               string
		apisixBinPath      string
		expectedErrMessage string
	}{
		{
			name:          "success",
			apisixBinPath: "echo",
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			options.Global.Deploy.Bare.APISIXBinPath = tc.apisixBinPath

			err := Reload(context.Background(), "/tmp")
			if tc.expectedErrMessage == "" {
				assert.Nil(t, err, "check reload error")
			} else {
				assert.NotNil(t, err, "check reload error")
				assert.Containsf(t, err.Error(), tc.expectedErrMessage, "check reload error message")
			}
		})
	}
}
