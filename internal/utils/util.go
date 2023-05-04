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

package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	sdk "github.com/api7/cloud-go-sdk"
	"gopkg.in/yaml.v3"
)

func readFile(filename string, v interface{}) error {
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %s", err)
	}

	switch filepath.Ext(filename) {
	case ".json":
		err = json.Unmarshal(data, v)
	case ".yml", ".yaml":
		err = yaml.Unmarshal(data, v)
	}
	if err != nil {
		return fmt.Errorf("failed to unmarshal: %s", err)
	}

	return nil
}
func ReadServiceFromFile(filename string) (*sdk.Application, error) {
	var app *sdk.Application
	if err := readFile(filename, &app); err != nil {
		return nil, fmt.Errorf("failed to read service from file: %s", err)
	}
	return app, nil
}

func ReadConsumerFromFile(filename string) (*sdk.Consumer, error) {
	var consumer *sdk.Consumer
	if err := readFile(filename, &consumer); err != nil {
		return nil, fmt.Errorf("failed to read consumer from file: %s", err)
	}
	return consumer, nil
}

func ReadRouterFromFile(filename string) (*sdk.API, error) {
	var router *sdk.API
	if err := readFile(filename, &router); err != nil {
		return nil, fmt.Errorf("failed to read router from file: %s", err)
	}
	return router, nil
}
