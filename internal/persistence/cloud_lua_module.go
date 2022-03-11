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

package persistence

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/cloud"
)

var (
	api7CloudHomeDir string
)

func init() {
	api7CloudHomeDir = filepath.Join(os.Getenv("HOME"), ".api7cloud")
	if err := os.MkdirAll(api7CloudHomeDir, 0755); err != nil {
		panic(err)
	}
}

// SaveCloudLuaModule downloads the cloud lua module and unzip and untar it,
// finally it'll be saved to the filesystem and the directory will be returned.
func SaveCloudLuaModule() (string, error) {
	data, err := cloud.DefaultClient.GetCloudLuaModule()
	if err != nil {
		return "", errors.Wrap(err, "failed to get cloud lua module")
	}
	tempReader, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return "", errors.Wrap(err, "failed to create gzip reader")
	}

	entryDir := api7CloudHomeDir
	reader := tar.NewReader(tempReader)
	for {
		hdr, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				return entryDir, nil
			}
			return "", errors.Wrap(err, "failed to read tar")
		}
		if hdr.FileInfo().IsDir() {
			dir := filepath.Join(api7CloudHomeDir, hdr.Name)
			if entryDir == "" {
				entryDir = dir
			}
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", errors.Wrap(err, "failed to create dir")
			}
		} else {
			filename := filepath.Join(api7CloudHomeDir, hdr.Name)
			buffer := make([]byte, 0, hdr.Size)
			if _, err := reader.Read(buffer); err != nil {
				return "", errors.Wrap(err, "failed to read tar")
			}
			if err := ioutil.WriteFile(filename, buffer, 0600); err != nil {
				return "", errors.Wrap(err, "failed to save file")
			}
		}
	}
}
