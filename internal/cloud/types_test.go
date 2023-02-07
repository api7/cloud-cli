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

package cloud

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name        string
		apiServer   string
		wantErr     bool
		errorReason string
		want        *api
	}{
		{
			name:        "malformed api server url",
			apiServer:   "--",
			wantErr:     true,
			errorReason: "invalid API7 Cloud server URL",
		},

		{
			name:      "custom api server url",
			apiServer: "http://abc.example.com",
			wantErr:   false,
			want: &api{
				host:        "abc.example.com",
				scheme:      "http",
				accessToken: "access-token",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a, err := newClient(tt.apiServer, "access-token")
			if tt.wantErr {
				assert.Error(t, err, "checking error")
				assert.Equal(t, tt.errorReason, err.Error(), "checking error reason")
			} else {
				assert.NoError(t, err, "checking error")
				a := a.(*api)
				assert.Equal(t, tt.want.host, a.host, "checking host")
				assert.Equal(t, tt.want.scheme, a.scheme, "checking scheme")
				assert.Equal(t, tt.want.accessToken, a.accessToken, "checking access token")
			}
		})
	}
}
