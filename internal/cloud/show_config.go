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
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/pkg/errors"
)

var (
	_validResources = map[string]struct{}{
		"application": {},
		"api":         {},
		"consumer":    {},
		"certificate": {},
	}
)

func (a *api) DebugShowConfig(cpID, resource, id string) (string, error) {
	if _, ok := _validResources[resource]; !ok {
		return "", fmt.Errorf("invalid resource type: %s", resource)
	}

	var rawData json.RawMessage
	if err := a.makeGetRequest(&url.URL{
		Path: fmt.Sprintf("/api/v1/debug/config/controlplanes/%s/%s/%s", cpID, resource, id),
	}, &rawData); err != nil {
		return "", err
	}

	return formatJSONData(rawData)
}

func formatJSONData(raw []byte) (string, error) {
	var data map[string]interface{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return "", errors.Wrap(err, "unmarshal data")
	}
	for key, arr := range data {
		// arr is a specific apisix resource array.
		for _, rawObj := range arr.([]interface{}) {
			var structualValue map[string]interface{}
			obj := rawObj.(map[string]interface{})
			// value is a JSON string, and we want to show it structurally, so
			// here we unmarshal and reset it.
			if err := json.Unmarshal([]byte(obj["value"].(string)), &structualValue); err != nil {
				return "", errors.Wrap(err, fmt.Sprintf("unmarshal %s", key))
			}
			obj["value"] = structualValue
		}
	}
	newData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(newData), nil
}
