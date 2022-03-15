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
	defer tempReader.Close()

	var entryDir string

	reader := tar.NewReader(tempReader)
	HomeDir := filepath.Join(HomeDir, "cloud_lua_module")
	_ = os.MkdirAll(HomeDir, 0755)
	for {
		hdr, err := reader.Next()
		if err != nil {
			if err == io.EOF {
				if entryDir == "" {
					entryDir = HomeDir
				}
				return entryDir, nil
			}
			return "", errors.Wrap(err, "failed to read tar")
		}
		if hdr.Typeflag == tar.TypeDir {
			dir := filepath.Join(HomeDir, hdr.Name)
			if entryDir == "" {
				entryDir = dir
			}
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", errors.Wrap(err, "failed to create dir")
			}
		} else if hdr.Typeflag == tar.TypeReg {
			filename := filepath.Join(HomeDir, hdr.Name)
			buffer, err := ioutil.ReadAll(reader)
			if err != nil {
				return "", errors.Wrap(err, "failed to read tar")
			}
			if err := ioutil.WriteFile(filename, buffer, 0644); err != nil {
				return "", errors.Wrap(err, "failed to save file")
			}
		}
	}
}
