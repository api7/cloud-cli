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

package deploy

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	sdk "github.com/api7/cloud-go-sdk"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/persistence"
)

var (
	_apisixStartupConfigTpl = `apisix:
  enable_admin: false
  ssl:
    ssl_trusted_certificate: {{ .TLSDir }}/ca.crt
  lua_module_hook: cloud
  extra_lua_path: {{ .CloudModuleDir }}/?.ljbc;
nginx_config:
  http:
    custom_lua_shared_dict:
      cloud: 1m
etcd:
  host:
    - "https://foo.com:443"
  tls:
    cert: {{ .TLSDir }}/tls.crt
    key: {{ .TLSDir }}/tls.key
    sni: foo.com
    verify: true
`

	_helmStartupConfigTpl = `apisix:
  image:
    repository: {{ .ImageRepository }}
    tag: {{ .ImageTag }}
  replicaCount: {{ .ReplicaCount }}
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
`
)

func mockCloudModule(t *testing.T) []byte {
	buffer := bytes.NewBuffer(nil)
	gzipWriter, err := gzip.NewWriterLevel(buffer, gzip.BestCompression)
	assert.NoError(t, err, "create gzip writer")
	files := []struct {
		Name, Body string
	}{
		{
			Name: "cloud_lua_module_beta",
		},
		{
			Name: "foo.txt",
			Body: "hello world",
		},
		{
			Name: "cloud_lua_module_beta/apisix",
		},
		{
			Name: "cloud_lua_module_beta/apisix/cli",
		},
		{
			Name: "cloud_lua_module_beta/apisix/cli/etcd.ljbc",
			Body: "this is apisix.cli.etcd",
		},
		{
			Name: "cloud_lua_module_beta/apisix/cli/local_storage.ljbc",
			Body: "this is apisix.cli.local_storage",
		},
	}
	tarWriter := tar.NewWriter(gzipWriter)
	for _, file := range files {
		hdr := &tar.Header{
			Name: file.Name,
			Size: int64(len(file.Body)),
		}
		if len(file.Body) == 0 {
			hdr.Typeflag = tar.TypeDir
		} else {
			hdr.Typeflag = tar.TypeReg
		}

		err = tarWriter.WriteHeader(hdr)
		assert.NoError(t, err, "write tar header")
		_, err = tarWriter.Write([]byte(file.Body))
		assert.NoError(t, err, "write tar body")
	}
	err = tarWriter.Flush()
	assert.NoError(t, err, "flush tar body")
	err = tarWriter.Close()
	assert.NoError(t, err, "close tar writer")
	err = gzipWriter.Close()
	assert.NoError(t, err, "close gzip writer")
	return buffer.Bytes()
}

func TestDeployPreRunForDocker(t *testing.T) {
	testCases := []struct {
		name              string
		errorReason       string
		mockFn            func(t *testing.T)
		specifiedAPISIXID string
		filledContext     deployContext
	}{
		{
			name:        "failed to get default cluster",
			errorReason: "Failed to get default cluster: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "failed to prepare cert",
			errorReason: "Failed to prepare certificate: download tls bundle: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
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
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)
				mockClient.EXPECT().GetCloudLuaModule().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "get startup config failed",
			errorReason: "failed to get startup config: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
					ClusterSpec: sdk.ClusterSpec{
						Domain: "foo.com",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				mockClient.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)

				mockClient.EXPECT().GetStartupConfig(sdk.ID(3), cloud.APISIX).Return("", errors.New("mock error"))

				cloud.DefaultClient = mockClient
			},
		},
		{
			name:              "success with specified apisix id",
			specifiedAPISIXID: "abcabc",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
					ClusterSpec: sdk.ClusterSpec{
						Domain: "foo.com",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				mockClient.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)

				mockClient.EXPECT().GetStartupConfig(sdk.ID(3), cloud.APISIX).Return(_apisixStartupConfigTpl, nil)
				cloud.DefaultClient = mockClient
			},
			filledContext: deployContext{
				cloudLuaModuleDir: filepath.Join(os.TempDir(), ".api7cloud", "cloud_lua_module_beta"),
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
				os.RemoveAll(filepath.Join(persistence.HomeDir, "tls"))
			}()
			ctx := &deployContext{}
			tc.mockFn(t)

			if tc.specifiedAPISIXID != "" {
				options.Global.Deploy.APISIXInstanceID = tc.specifiedAPISIXID
			}

			err := deployPreRunForDocker(ctx)
			if tc.errorReason != "" {
				assert.Equal(t, tc.errorReason, err.Error(), "check error")
			} else {
				assert.NoError(t, err, "check error")
				assert.Equal(t, tc.filledContext.cloudLuaModuleDir, ctx.cloudLuaModuleDir, "check cloud lua module dir")
				assert.Equal(t, string(tc.filledContext.essentialConfig), string(ctx.essentialConfig), "check essential config")

				id, err := os.ReadFile(filepath.Join(persistence.HomeDir, "apisix.uid"))
				assert.Nil(t, err, "read apisix.uid")
				// We cannot add an assertion if the ID was generated randomly.
				if tc.specifiedAPISIXID != "" {
					assert.Equal(t, tc.specifiedAPISIXID, string(id), "check apisix id")
				}
			}
		})
	}
}

func TestDeployPreRunForBare(t *testing.T) {
	testCases := []struct {
		name          string
		errorReason   string
		mockFn        func(t *testing.T)
		filledContext deployContext
	}{
		{
			name:        "failed to get default cluster",
			errorReason: "Failed to get default cluster: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "failed to prepare cert",
			errorReason: "Failed to prepare certificate: download tls bundle: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
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
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)
				mockClient.EXPECT().GetCloudLuaModule().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "get startup config failed",
			errorReason: "failed to get startup config: mock error",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
					ClusterSpec: sdk.ClusterSpec{
						Domain: "foo.com",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				mockClient.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)

				mockClient.EXPECT().GetStartupConfig(sdk.ID(3), cloud.APISIX).Return("", errors.New("mock error"))

				cloud.DefaultClient = mockClient
			},
		},
		{
			name: "success",
			mockFn: func(t *testing.T) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
					ClusterSpec: sdk.ClusterSpec{
						Domain: "foo.com",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				mockClient.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)

				mockClient.EXPECT().GetStartupConfig(sdk.ID(3), cloud.APISIX).Return(_apisixStartupConfigTpl, nil)

				{
					file, err := os.CreateTemp(os.TempDir(), "apisix-cli-etcd-*.lua")
					assert.Nil(t, err, "create temp file path for /usr/local/apisix/apisix/cli/etcd.lua")
					_targetApisixCliEtcdPath = file.Name()
				}
				{
					file, err := os.CreateTemp(os.TempDir(), "apisix-cli-local-storage-*.lua")
					assert.Nil(t, err, "create temp file path for /usr/local/apisix/apisix/cli/local_storage.lua")
					_targetApisixCliLocalStoragePath = file.Name()
				}

				cloud.DefaultClient = mockClient
			},
			filledContext: deployContext{
				cloudLuaModuleDir: filepath.Join(os.TempDir(), ".api7cloud", "cloud_lua_module_beta"),
				essentialConfig: []byte(`apisix:
  enable_admin: false
  ssl:
    ssl_trusted_certificate: /usr/local/apisix/conf/ssl/ca\.crt
  lua_module_hook: cloud
  extra_lua_path: .*/.api7cloud/cloud_lua_module_beta/\?\.ljbc;
nginx_config:
  http:
    custom_lua_shared_dict:
      cloud: 1m
etcd:
  host:
    - "https://foo.com:443"
  tls:
    cert: /usr/local/apisix/conf/ssl/tls.crt
    key: /usr/local/apisix/conf/ssl/tls.key
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
				os.RemoveAll(filepath.Join(persistence.HomeDir, "tls"))
			}()
			ctx := &deployContext{}
			tc.mockFn(t)

			err := deployPreRunForBare(ctx)
			if tc.errorReason != "" {
				assert.Equal(t, tc.errorReason, err.Error(), "check error")
			} else {
				assert.NoError(t, err, "check error")
				assert.Equal(t, tc.filledContext.cloudLuaModuleDir, ctx.cloudLuaModuleDir, "check cloud lua module dir")
				assert.Regexp(t, string(tc.filledContext.essentialConfig), string(ctx.essentialConfig), "check essential config")
			}
		})
	}
}

func TestDeployPreRunForKubernetes(t *testing.T) {
	essentialConfig := []byte(`apisix:
  image:
    repository: apache/apisix
    tag: 2.15.0-centos
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
`)

	type testCase struct {
		name          string
		errorReason   string
		mockFn        func(t *testing.T, test *testCase)
		filledContext deployContext
		globalOptions options.Options
		kubectl       commands.Cmd
	}

	testCases := []testCase{
		{
			name:        "failed to get default cluster",
			errorReason: "Failed to get default cluster: mock error",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "failed to prepare cert",
			errorReason: "Failed to prepare certificate: download tls bundle: mock error",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "get cloud lua module failed",
			errorReason: "Failed to save cloud lua module: failed to get cloud lua module: mock error",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)
				mockClient.EXPECT().GetCloudLuaModule().Return(nil, errors.New("mock error"))
				cloud.DefaultClient = mockClient
			},
		},
		{
			name:        "get startup config failed",
			errorReason: "failed to get startup config: mock error",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
					ClusterSpec: sdk.ClusterSpec{
						Domain: "foo.com",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				mockClient.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)

				mockClient.EXPECT().GetStartupConfig(sdk.ID(3), cloud.HELM).Return("", errors.New("mock error"))

				cloud.DefaultClient = mockClient
			},
		},
		{
			name: "create namespace, secret or configMap on kubernetes failed",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
					ClusterSpec: sdk.ClusterSpec{
						Domain: "foo.com",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				mockClient.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)

				mockClient.EXPECT().GetStartupConfig(sdk.ID(3), cloud.HELM).Return(_helmStartupConfigTpl, nil)
				cloud.DefaultClient = mockClient

				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "", errors.New("mock error")).AnyTimes()
				test.kubectl = mockCmd
			},
			globalOptions: options.Options{
				Verbose: true,
				Deploy: options.DeployOptions{
					Kubernetes: options.KubernetesDeployOptions{
						Namespace:    "apisix",
						APISIXImage:  "apache/apisix:2.15.0-centos",
						ReplicaCount: 1,
					},
				},
			},
			errorReason: "mock error",
		},
		{
			name: "when namespace, secret or configMap already exists, create should succeed",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
					ClusterSpec: sdk.ClusterSpec{
						Domain: "foo.com",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				mockClient.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)

				mockClient.EXPECT().GetStartupConfig(sdk.ID(3), cloud.HELM).Return(_helmStartupConfigTpl, nil)
				cloud.DefaultClient = mockClient

				mockCmd := commands.NewMockCmd(ctrl)
				mockCmd.EXPECT().String().AnyTimes()
				mockCmd.EXPECT().AppendArgs(gomock.Any()).AnyTimes()
				mockCmd.EXPECT().Run(gomock.Any()).Return("", "AlreadyExists", errors.New("mock error")).AnyTimes()
				test.kubectl = mockCmd
			},
			filledContext: deployContext{
				cloudLuaModuleDir: filepath.Join(os.TempDir(), ".api7cloud", "cloud_lua_module_beta"),
				essentialConfig:   essentialConfig,
				KubernetesOpts: &options.KubernetesDeployOptions{
					Namespace:    "apisix",
					APISIXImage:  "apache/apisix:2.15.0-centos",
					ReplicaCount: 1,
				},
			},
			globalOptions: options.Options{
				DryRun:  true,
				Verbose: true,
				Deploy: options.DeployOptions{
					Kubernetes: options.KubernetesDeployOptions{
						Namespace:    "apisix",
						APISIXImage:  "apache/apisix:2.15.0-centos",
						ReplicaCount: 1,
					},
				},
			},
		},
		{
			name: "deploy on kubernetes pre run was succeed",
			mockFn: func(t *testing.T, test *testCase) {
				ctrl := gomock.NewController(t)
				mockClient := cloud.NewMockAPI(ctrl)
				mockClient.EXPECT().GetDefaultCluster().Return(&sdk.Cluster{
					ID: 3,
					ClusterSpec: sdk.ClusterSpec{
						Domain: "foo.com",
					},
				}, nil)
				mockClient.EXPECT().GetTLSBundle(gomock.Any()).Return(&sdk.TLSBundle{
					Certificate:   "1",
					PrivateKey:    "1",
					CACertificate: "1",
				}, nil)

				mockClient.EXPECT().GetCloudLuaModule().Return(mockCloudModule(t), nil)

				mockClient.EXPECT().GetStartupConfig(sdk.ID(3), cloud.HELM).Return(_helmStartupConfigTpl, nil)
				cloud.DefaultClient = mockClient

				test.kubectl = commands.New("kubectl", test.globalOptions.DryRun)
			},
			filledContext: deployContext{
				cloudLuaModuleDir: filepath.Join(os.TempDir(), ".api7cloud", "cloud_lua_module_beta"),
				essentialConfig:   essentialConfig,
				KubernetesOpts: &options.KubernetesDeployOptions{
					Namespace:    "apisix",
					APISIXImage:  "apache/apisix:2.15.0-centos",
					ReplicaCount: 1,
				},
			},
			globalOptions: options.Options{
				DryRun:  true,
				Verbose: true,
				Deploy: options.DeployOptions{
					Kubernetes: options.KubernetesDeployOptions{
						Namespace:    "apisix",
						APISIXImage:  "apache/apisix:2.15.0-centos",
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
				os.RemoveAll(filepath.Join(persistence.HomeDir, "tls"))
			}()

			ctx := &deployContext{}
			tc.mockFn(t, &tc)

			err := deployPreRunForKubernetes(ctx, tc.kubectl)
			if tc.errorReason != "" {
				assert.Contains(t, err.Error(), tc.errorReason, "check error")
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.filledContext.cloudLuaModuleDir, ctx.cloudLuaModuleDir, "check cloud lua module dir")
				assert.Equal(t, tc.filledContext.essentialConfig, ctx.essentialConfig, "check essential config")
			}
		})
	}
}

func TestGetDockerContainerIDbyName(t *testing.T) {
	testCases := []struct {
		name        string
		mockFn      func(t *testing.T) commands.Cmd
		output      string
		errorReason string
	}{
		{
			name: "mock error",
			mockFn: func(t *testing.T) commands.Cmd {
				ctrl := gomock.NewController(t)
				cmd := commands.NewMockCmd(ctrl)
				cmd.EXPECT().AppendArgs("ps", "--filter", "name=apisix", "--format", "{{.ID}}")
				cmd.EXPECT().Run(gomock.Any()).Return("", "", errors.New("mock error"))
				return cmd
			},
			errorReason: "mock error",
		},
		{
			name: "stderr is not empty",
			mockFn: func(t *testing.T) commands.Cmd {
				ctrl := gomock.NewController(t)
				cmd := commands.NewMockCmd(ctrl)
				cmd.EXPECT().AppendArgs("ps", "--filter", "name=apisix", "--format", "{{.ID}}")
				cmd.EXPECT().Run(gomock.Any()).Return("", "stderr", nil)
				return cmd
			},
			errorReason: "get container id: stderr: stderr",
		},
		{
			name: "success",
			mockFn: func(t *testing.T) commands.Cmd {
				ctrl := gomock.NewController(t)
				cmd := commands.NewMockCmd(ctrl)
				cmd.EXPECT().AppendArgs("ps", "--filter", "name=apisix", "--format", "{{.ID}}")
				cmd.EXPECT().Run(gomock.Any()).Return("2b68d1dcfe34", "", nil)
				return cmd
			},
			output: "2b68d1dcfe34",
		},
	}
	for _, tc := range testCases {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			cmd := tc.mockFn(t)
			out, err := getDockerContainerIDByName(context.TODO(), cmd, "apisix")
			if tc.errorReason != "" {
				assert.Empty(t, out, "check if output is empty")
				assert.Equal(t, tc.errorReason, err.Error(), "check the error reason")
			} else {
				assert.Nil(t, err, "check if error is nil")
				assert.Equal(t, tc.output, out, "check the output")
			}
		})
	}
}
