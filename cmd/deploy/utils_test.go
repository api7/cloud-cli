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

package deploy

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/persistence"
	"github.com/api7/cloud-cli/internal/types"
)

func TestDeployPreRunForDocker(t *testing.T) {
	testCases := []struct {
		name          string
		errorReason   string
		mockFn        func(t *testing.T)
		filledContext deployContext
	}{
		{
			name:        "failed to get default control plane",
			errorReason: "Failed to get default control plane: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultControlPlane().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "failed to prepare cert",
			errorReason: "Failed to prepare certificate: download tls bundle: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "3",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "get cloud lua module failed",
			errorReason: "Failed to save cloud lua module: failed to get cloud lua module: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "3",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)
				mockClient.EXPECT().GetCloudLuaModule().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name: "success",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "3",
					},
					Domain: "foo.com",
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)
				buffer := bytes.NewBuffer(nil)
				gzipWriter, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
				assert.NoError(t, err, "create gzip writer")
				tarWriter := tar.NewWriter(gzipWriter)
				body := "hello world"
				hdr := &tar.Header{
					Name: "foo.txt",
					Size: int64(len(body)),
				}
				err = tarWriter.WriteHeader(hdr)
				assert.NoError(t, err, "write tar header")
				_, err = tarWriter.Write([]byte(body))
				assert.NoError(t, err, "write tar body")
				err = tarWriter.Flush()
				assert.NoError(t, err, "flush tar body")
				err = tarWriter.Close()
				assert.NoError(t, err, "close tar writer")
				err = gzipWriter.Close()
				assert.NoError(t, err, "close gzip writer")

				mockClient.EXPECT().GetCloudLuaModule().Return(buffer.Bytes(), nil)

				cloud.DefaultClient = mockClient
			},
			filledContext: deployContext{
				cloudLuaModuleDir: filepath.Join(os.TempDir(), ".api7cloud"),
				essentialConfig: []byte(`apisix:
  enable_admin: false
  ssl:
    ssl_trusted_certificate: /cloud/tls/ca.crt
  lua_module_hook: cloud
  extra_lua_path: /cloud_lua_module/?.ljbc;
nginx_config:
  http:
    custom_lua_shared_dict:
      cloud: 1m
etcd:
  host:
    - "https://foo.com:443"
  tls:
    cert: /cloud/tls/tls.crt
    key: /cloud/tls/tls.key
    sni: foo.com
    verify: true
`),
			},
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			persistence.HomeDir = filepath.Join(os.TempDir(), ".api7cloud")
			if err := persistence.Init(); err != nil {
				panic(err)
			}

			defer func() {
				os.Remove(filepath.Join(persistence.HomeDir, "tls", "tls.crt"))
				os.Remove(filepath.Join(persistence.HomeDir, "tls", "tls.key"))
				os.Remove(filepath.Join(persistence.HomeDir, "tls", "ca.crt"))
			}()
			ctx := &deployContext{}
			tc.mockFn(t)

			err := deployPreRunForDocker(ctx)
			if tc.errorReason != "" {
				assert.Equal(t, tc.errorReason, err.Error(), "check error")
			} else {
				assert.Equal(t, tc.filledContext.cloudLuaModuleDir, ctx.cloudLuaModuleDir, "check cloud lua module dir")
				assert.Equal(t, tc.filledContext.essentialConfig, ctx.essentialConfig, "check essential config")
			}
		})
	}
}

func TestDeployPreRunForKubernetes(t *testing.T) {
	testCases := []struct {
		name          string
		errorReason   string
		mockFn        func(t *testing.T)
		filledContext deployContext
		globalOptions options.Options
	}{
		{
			name:        "failed to get default control plane",
			errorReason: "Failed to get default control plane: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultControlPlane().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "failed to prepare cert",
			errorReason: "Failed to prepare certificate: download tls bundle: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "3",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "get cloud lua module failed",
			errorReason: "Failed to save cloud lua module: failed to get cloud lua module: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "3",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)
				mockClient.EXPECT().GetCloudLuaModule().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name: "create namespace, secret or configMap on kubernetes failed",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "3",
					},
					Domain: "foo.com",
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				buffer := bytes.NewBuffer(nil)
				gzipWriter, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
				assert.NoError(t, err, "create gzip writer")
				tarWriter := tar.NewWriter(gzipWriter)
				body := "hello world"
				hdr := &tar.Header{
					Name: "foo.txt",
					Size: int64(len(body)),
				}
				err = tarWriter.WriteHeader(hdr)
				assert.NoError(t, err, "write tar header")
				_, err = tarWriter.Write([]byte(body))
				assert.NoError(t, err, "write tar body")
				err = tarWriter.Flush()
				assert.NoError(t, err, "flush tar body")
				err = tarWriter.Close()
				assert.NoError(t, err, "close tar writer")
				err = gzipWriter.Close()
				assert.NoError(t, err, "close gzip writer")
				mockClient.EXPECT().GetCloudLuaModule().Return(buffer.Bytes(), nil)

				cloud.DefaultClient = mockClient
			},
			globalOptions: options.Options{
				Verbose: true,
				Deploy: options.DeployOptions{
					Kubernetes: options.KubernetesDeployOptions{
						NameSpace:      "apisix",
						APISIXImage:    "apache/apisix:2.11.0-centos",
						ReplicaCount:   1,
						KubectlCLIPath: "/tmp/kubectl",
					},
				},
			},
			errorReason: "/tmp/kubectl",
		},
		{
			name: "deploy on kubernetes pre run was succeed",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultControlPlane().Return(&types.ControlPlane{
					TypeMeta: types.TypeMeta{
						ID: "3",
					},
					Domain: "foo.com",
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&types.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				buffer := bytes.NewBuffer(nil)
				gzipWriter, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
				assert.NoError(t, err, "create gzip writer")
				tarWriter := tar.NewWriter(gzipWriter)
				body := "hello world"
				hdr := &tar.Header{
					Name: "foo.txt",
					Size: int64(len(body)),
				}
				err = tarWriter.WriteHeader(hdr)
				assert.NoError(t, err, "write tar header")
				_, err = tarWriter.Write([]byte(body))
				assert.NoError(t, err, "write tar body")
				err = tarWriter.Flush()
				assert.NoError(t, err, "flush tar body")
				err = tarWriter.Close()
				assert.NoError(t, err, "close tar writer")
				err = gzipWriter.Close()
				assert.NoError(t, err, "close gzip writer")
				mockClient.EXPECT().GetCloudLuaModule().Return(buffer.Bytes(), nil)

				cloud.DefaultClient = mockClient
			},
			filledContext: deployContext{
				cloudLuaModuleDir: filepath.Join(os.TempDir(), ".api7cloud"),
				essentialConfig: []byte(`apisix:
  image:
    repository: apache/apisix
    tag: 2.11.0-centos
  replicaCount: 1
  setIDFromPodUID: true
  luaModuleHook:
    enabled: true
    luaPath: "/lua-module-hook/?.ljbc"
    hookPoint: cloud
    configMapRef:
      name: "cloud-module"
      mounts:
        - key: cloud.ljbc
          path: /lua-module-hook/cloud.ljbc
        - key: cloud-agent.ljbc
          path: /lua-module-hook/cloud/agent.ljbc
        - key: cloud-metrics.ljbc
          path: /lua-module-hook/cloud/metrics.ljbc
        - key: cloud-utils.ljbc
          path: /lua-module-hook/cloud/utils.ljbc
  customLuaSharedDicts:
    - name: cloud
      size: 1m
gateway:
  tls:
    enabled: true
    existingCASecret: cloud-ssl
    certCAFilename: "ca.crt"
admin:
  enabled: false
etcd:
  enabled: false
  host:
    - "https://foo.com:443"
  timeout: 30
  auth:
    tls:
      enabled: true
      sni: foo.com
      existingSecret: cloud-ssl
      certFilename: tls.crt
      certKeyFilename: tls.key
`),
				KubernetesOpts: &options.KubernetesDeployOptions{
					NameSpace:    "apisix",
					APISIXImage:  "apache/apisix:2.11.0-centos",
					ReplicaCount: 1,
				},
			},
			globalOptions: options.Options{
				DryRun:  true,
				Verbose: true,
				Deploy: options.DeployOptions{
					Kubernetes: options.KubernetesDeployOptions{
						NameSpace:    "apisix",
						APISIXImage:  "apache/apisix:2.11.0-centos",
						ReplicaCount: 1,
					},
				},
			},
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			options.Global = tc.globalOptions
			persistence.HomeDir = filepath.Join(os.TempDir(), ".api7cloud")
			if err := persistence.Init(); err != nil {
				panic(err)
			}

			defer func() {
				os.Remove(filepath.Join(persistence.HomeDir, "tls", "tls.crt"))
				os.Remove(filepath.Join(persistence.HomeDir, "tls", "tls.key"))
				os.Remove(filepath.Join(persistence.HomeDir, "tls", "ca.crt"))
			}()
			ctx := &deployContext{}
			tc.mockFn(t)

			err := deployPreRunForKubernetes(ctx)
			if tc.errorReason != "" {
				assert.Equal(t, tc.errorReason, err.Error(), "check error")
			} else {
				assert.Equal(t, tc.filledContext.cloudLuaModuleDir, ctx.cloudLuaModuleDir, "check cloud lua module dir")
				assert.Equal(t, tc.filledContext.essentialConfig, ctx.essentialConfig, "check essential config")
			}
		})
	}
}
