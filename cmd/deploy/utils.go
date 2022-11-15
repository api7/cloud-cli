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
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/google/uuid"
	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
	"github.com/api7/cloud-cli/internal/types"
	"github.com/api7/cloud-cli/internal/utils"
)

type deployContext struct {
	cloudLuaModuleDir string
	tlsDir            string
	apisixEtcdCertDir string
	essentialConfig   []byte
	apisixIDFile      string
	apisixID          string
	// ControlPlane is the current control plane.
	ControlPlane   *types.ControlPlane
	KubernetesOpts *options.KubernetesDeployOptions
}

type config struct {
	CloudModuleDir string
	TLSDir         string
}

type helmConfig struct {
	ImageRepository string
	ImageTag        string
	ReplicaCount    uint
}

func getEssentialConfigTpl(ctx *deployContext, configType cloud.StartupConfigType) (*template.Template, error) {
	config, err := cloud.DefaultClient.GetStartupConfig(ctx.ControlPlane.ID, configType)
	if err != nil {
		return nil, fmt.Errorf("failed to get startup config: %s", err.Error())
	}

	configTemplate, err := template.New("essential config").Parse(config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse startup config template: %s", err.Error())
	}

	return configTemplate, nil
}

func deployPreRunForDocker(ctx *deployContext) error {
	err := deployPreRun(ctx)
	if err != nil {
		return err
	}

	essentialConfigTpl, err := getEssentialConfigTpl(ctx, cloud.APISIX)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err := essentialConfigTpl.Execute(buf, &config{
		CloudModuleDir: "/cloud_lua_module",
		TLSDir:         "/cloud/tls",
	}); err != nil {
		return fmt.Errorf("Failed to execute essential config template: %s", err)
	}

	ctx.essentialConfig = buf.Bytes()

	// We generate the APISIX instance ID and mount to /usr/local/apisix/conf/apisix.uid
	ctx.apisixID = options.Global.Deploy.APISIXInstanceID
	if ctx.apisixID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return errors.Wrap(err, "failed to generate APISIX instance ID")
		}
		ctx.apisixID = id.String()
	}

	ctx.apisixIDFile = filepath.Join(persistence.HomeDir, "apisix.uid")
	if err := os.WriteFile(ctx.apisixIDFile, []byte(ctx.apisixID), 0644); err != nil {
		return errors.Wrap(err, "failed to save APISIX instance ID")
	}

	return nil
}

func deployPreRunForBare(ctx *deployContext) error {
	err := deployPreRun(ctx)
	if err != nil {
		return err
	}

	essentialConfigTemplate, err := getEssentialConfigTpl(ctx, cloud.APISIX)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err := essentialConfigTemplate.Execute(buf, &config{
		CloudModuleDir: ctx.cloudLuaModuleDir,
		TLSDir:         "/usr/local/apisix/conf/ssl",
	}); err != nil {
		return fmt.Errorf("Failed to execute essential config template: %s", err)
	}

	ctx.essentialConfig = buf.Bytes()
	ctx.apisixEtcdCertDir = "/usr/local/apisix/conf/ssl"
	return nil
}

func deployPreRun(ctx *deployContext) error {
	cp, err := cloud.DefaultClient.GetDefaultControlPlane()
	if err != nil {
		return fmt.Errorf("Failed to get default control plane: %s", err.Error())
	}
	if err := persistence.PrepareCertificate(cp.ID); err != nil {
		return fmt.Errorf("Failed to prepare certificate: %s", err.Error())
	}
	ctx.tlsDir = filepath.Join(persistence.TLSDir, cp.ID)

	cloudLuaModuleDir, err := persistence.SaveCloudLuaModule()
	if err != nil {
		return fmt.Errorf("Failed to save cloud lua module: %s", err)
	}
	output.Verbosef("Saved cloud lua module to: %s", cloudLuaModuleDir)

	ctx.cloudLuaModuleDir = cloudLuaModuleDir
	ctx.ControlPlane = cp
	return nil
}

func deployPreRunForKubernetes(ctx *deployContext, kubectl commands.Cmd) error {
	opts := &options.Global.Deploy.Kubernetes
	image := strings.Split(opts.APISIXImage, ":")
	opts.APISIXImageRepo = image[0]
	if len(image) > 1 {
		opts.APISIXImageTag = image[1]
	} else {
		opts.APISIXImageTag = "latest"
	}
	ctx.KubernetesOpts = opts

	cp, err := cloud.DefaultClient.GetDefaultControlPlane()
	if err != nil {
		return fmt.Errorf("Failed to get default control plane: %v", err.Error())
	}
	if err = persistence.PrepareCertificate(cp.ID); err != nil {
		return fmt.Errorf("Failed to prepare certificate: %s", err.Error())
	}
	ctx.tlsDir = filepath.Join(persistence.TLSDir, cp.ID)

	cloudLuaModuleDir, err := persistence.SaveCloudLuaModule()
	if err != nil {
		return fmt.Errorf("Failed to save cloud lua module: %s", err.Error())
	}
	output.Verbosef("Saved cloud lua module to: %s", cloudLuaModuleDir)
	ctx.cloudLuaModuleDir = cloudLuaModuleDir

	ctx.ControlPlane = cp

	helmEssentialConfigTemplate, err := getEssentialConfigTpl(ctx, cloud.HELM)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	if err = helmEssentialConfigTemplate.Execute(buf, helmConfig{
		ImageRepository: ctx.KubernetesOpts.APISIXImageRepo,
		ImageTag:        ctx.KubernetesOpts.APISIXImageTag,
		ReplicaCount:    ctx.KubernetesOpts.ReplicaCount,
	}); err != nil {
		return fmt.Errorf("Failed to execute helm essential config template: %s", err.Error())
	}
	ctx.essentialConfig = buf.Bytes()

	if err = createOnKubernetes(ctx, types.Namespace, kubectl); err != nil {
		return fmt.Errorf("Failed to create namespace on kubernetes: %s", err.Error())
	}
	if err = createOnKubernetes(ctx, types.Secret, kubectl); err != nil {
		return fmt.Errorf("Failed to create secret on kubernetes: %s", err.Error())
	}
	if err = createOnKubernetes(ctx, types.ConfigMap, kubectl); err != nil {
		return fmt.Errorf("Failed to create configmap on kubernetes: %s", err.Error())
	}

	return nil
}

// createOnKubernetes create namespace, secret or configmap on Kubernetes
func createOnKubernetes(ctx *deployContext, k types.K8sResourceKind, kubectl commands.Cmd) error {
	var (
		err  error
		opts = ctx.KubernetesOpts
	)

	switch k {
	case types.Secret:
		kubectl.AppendArgs("create", "secret", "generic", consts.DefaultSecretName)
		kubectl.AppendArgs("--from-file", fmt.Sprintf("tls.crt=%s", filepath.Join(ctx.tlsDir, "tls.crt")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("tls.key=%s", filepath.Join(ctx.tlsDir, "tls.key")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("ca.crt=%s", filepath.Join(ctx.tlsDir, "ca.crt")))
		kubectl.AppendArgs("--namespace", opts.NameSpace)
	case types.ConfigMap:
		kubectl.AppendArgs("create", "configmap", consts.DefaultConfigMapName)
		// TODO: dynamic list files in cloud lua module instead of hard code maybe better
		kubectl.AppendArgs("--from-file", fmt.Sprintf("cloud.ljbc=%s", filepath.Join(ctx.cloudLuaModuleDir, "cloud.ljbc")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("cloud-agent.ljbc=%s", filepath.Join(ctx.cloudLuaModuleDir, "cloud/agent.ljbc")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("cloud-metrics.ljbc=%s", filepath.Join(ctx.cloudLuaModuleDir, "cloud/metrics.ljbc")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("cloud-utils.ljbc=%s", filepath.Join(ctx.cloudLuaModuleDir, "cloud/utils.ljbc")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("cloud-file.ljbc=%s", filepath.Join(ctx.cloudLuaModuleDir, "cloud/file.ljbc")))
		kubectl.AppendArgs("--namespace", opts.NameSpace)
	case types.Namespace:
		kubectl.AppendArgs("create", "ns", opts.NameSpace)
	default:
		panic(fmt.Sprintf("invaild kind: %d", k))
	}

	if options.Global.DryRun {
		output.Infof("Running:\n%s\n", kubectl.String())
	} else {
		output.Verbosef("Running:\n%s\n", kubectl.String())
	}

	newCtx, cancel := context.WithTimeout(context.TODO(), consts.DefaultKubectlTimeout)
	defer cancel()

	stdout, stderr, err := kubectl.Run(newCtx)
	if stderr != "" {
		output.Warnf(stderr)
		// if stderr contains AlreadyExists flag, we don't think it's an error
		if strings.Contains(stderr, "AlreadyExists") {
			err = nil
		}
	}
	if stdout != "" {
		output.Verbosef(stdout)
	}
	if err != nil {
		return err
	}

	return nil
}

func printInstallDetailForKubernetes(kubectl commands.Cmd) {
	var (
		deploymentName string
		serviceName    string
		podsNames      []string
		APISIXID       string
		err            error
	)

	defer func() {
		if err != nil {
			output.Errorf("Failed to print APISIX installation details. Please view related resources of Kubernetes manually.")
		}
	}()

	output.Infof("\nCongratulations! Your APISIX cluster was deployed successfully on Kubernetes.\n")
	output.Infof("The Helm release name is: %s", options.Global.Deploy.Name)

	if deploymentName, err = utils.GetDeploymentName(kubectl); err != nil {
		output.Warnf("Failed to get Deployment: %s", err.Error())
		return
	}
	output.Infof("The APISIX Deployment name is: %s", deploymentName)

	if serviceName, err = utils.GetServiceName(kubectl); err != nil {
		output.Warnf("Failed to get Service: %s", err.Error())
		return
	}
	output.Infof("The APISIX Service name is: %s", serviceName)

	output.Infof("\nWorkloads:")
	if podsNames, err = utils.GetPodsNames(kubectl); err != nil {
		output.Warnf("Failed to get pods: %s", err.Error())
		return
	}

	for _, podName := range podsNames {
		if APISIXID, err = utils.GetAPISIXID(kubectl, podName); err != nil {
			output.Warnf("Failed to get APISIXID: %s", err.Error())
			return
		}
		output.Infof("Pod Name: %s APISIX ID: %s", podName, APISIXID)
	}
}

func getDockerContainerIDByName(ctx context.Context, docker commands.Cmd, name string) (string, error) {
	docker.AppendArgs("ps", "--filter", "name="+name, "--format", "{{.ID}}")
	stdout, stderr, err := docker.Run(ctx)
	if err != nil {
		return "", err
	}
	if stderr != "" {
		return "", fmt.Errorf("get container id: stderr: %s", stderr)
	}
	return strings.TrimRight(stdout, "\r\n"), nil
}

func deployOnBareMetal(ctx context.Context, deployCtx *deployContext, opts *options.BareDeployOptions, configFile string) {
	buf := bytes.NewBuffer(nil)
	err := _installer.Execute(buf, &installContext{
		APISIXRepoURL:     _apisixRepoURL,
		TLSDir:            deployCtx.tlsDir,
		APISIXEtcdCertDir: deployCtx.apisixEtcdCertDir,
		ConfigFile:        configFile,
		Version:           opts.APISIXVersion,
		InstanceID:        options.Global.Deploy.APISIXInstanceID,
	})
	if err != nil {
		output.Errorf(err.Error())
		return
	}
	installerPath := filepath.Join(persistence.HomeDir, "scripts")
	err = os.Mkdir(installerPath, 0755)
	if err != nil {
		if !os.IsExist(err) {
			output.Errorf(err.Error())
			return
		}
	}

	installerFile := filepath.Join(installerPath, "install.sh")
	err = os.WriteFile(installerFile, buf.Bytes(), 0755)
	if err != nil {
		output.Errorf(err.Error())
		return
	}

	bare := commands.New("/usr/bin/bash", options.Global.DryRun)
	bare.AppendArgs("-C")
	bare.AppendArgs(installerFile)

	if err = bare.Execute(ctx); err != nil {
		output.Errorf(err.Error())
		return
	}
}
