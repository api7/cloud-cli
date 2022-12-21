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

	"github.com/api7/cloud-go-sdk"
	"github.com/bitly/go-simplejson"
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

func (a *api) DebugShowConfig(cpID cloud.ID, resource, id string) (string, error) {
	if _, ok := _validResources[resource]; !ok {
		return "", fmt.Errorf("invalid resource type: %s", resource)
	}

	var rawData json.RawMessage
	if err := a.makeGetRequest(&url.URL{
		Path: fmt.Sprintf("/api/v1/debug/config/controlplanes/%s/%s/%s", cpID.String(), resource, id),
	}, &rawData); err != nil {
		return "", err
	}

	return formatJSONData(rawData)
}

func formatJSONData(raw []byte) (string, error) {
	js, err := simplejson.NewJson(raw)
	if err != nil {
		return "", errors.Wrap(err, "invalid json")
	}

	for _, resName := range []string{"routes", "services", "upstreams", "certificates", "consumers"} {
		res, ok := js.CheckGet(resName)
		if !ok {
			continue
		}

		for i := 0; i < len(res.MustArray()); i++ {
			var structualValue map[string]interface{}
			value := res.GetIndex(i).Get("value").MustString()

			// value is a JSON string, and we want to show it structurally, so
			// here we unmarshal and reset it.
			if err := json.Unmarshal([]byte(value), &structualValue); err != nil {
				return "", errors.Wrap(err, fmt.Sprintf("unmarshal %s", value))
			}

			res.GetIndex(i).Set("value", structualValue)
		}
	}
	newData, err := js.MarshalJSON()
	if err != nil {
		return "", err
	}
	return string(newData), nil
}
