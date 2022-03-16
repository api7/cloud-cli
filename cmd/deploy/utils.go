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
	"bytes"
	"context"
	_ "embed"
	"fmt"
	"html/template"
	"path/filepath"
	"strings"
	"time"

	"github.com/api7/cloud-cli/internal/cloud"
	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-cli/internal/persistence"
	"github.com/api7/cloud-cli/internal/types"
	"github.com/api7/cloud-cli/internal/utils"
)

var (
	//go:embed apisix.yaml
	essentialConfig string
	//go:embed apisix_chart_values.yaml
	helmEssentialConfig string

	essentialConfigTemplate     = template.Must(template.New("essential config").Parse(essentialConfig))
	helmEssentialConfigTemplate = template.Must(template.New("helm essential config").Parse(helmEssentialConfig))
)

type kind uint8

const (
	configMap kind = iota
	secret
	namespace
)

type deployContext struct {
	cloudLuaModuleDir string
	tlsDir            string
	essentialConfig   []byte
	// ControlPlane is the current control plane.
	ControlPlane   *types.ControlPlane
	KubernetesOpts *options.KubernetesDeployOptions
}

func deployPreRunForDocker(ctx *deployContext) error {
	cp, err := cloud.DefaultClient.GetDefaultControlPlane()
	if err != nil {
		return fmt.Errorf("Failed to get default control plane: %s", err.Error())
	}
	if err := persistence.PrepareCertificate(cp.ID); err != nil {
		return fmt.Errorf("Failed to prepare certificate: %s", err.Error())
	}
	ctx.tlsDir = persistence.TLSDir

	cloudLuaModuleDir, err := persistence.SaveCloudLuaModule()
	if err != nil {
		return fmt.Errorf("Failed to save cloud lua module: %s", err)
	}
	output.Verbosef("Saved cloud lua module to: %s", cloudLuaModuleDir)

	ctx.cloudLuaModuleDir = cloudLuaModuleDir
	ctx.ControlPlane = cp

	buf := bytes.NewBuffer(nil)
	if err := essentialConfigTemplate.Execute(buf, ctx); err != nil {
		return fmt.Errorf("Failed to execute essential config template: %s", err)
	}

	ctx.essentialConfig = buf.Bytes()
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
	ctx.tlsDir = persistence.TLSDir

	cloudLuaModuleDir, err := persistence.SaveCloudLuaModule()
	if err != nil {
		return fmt.Errorf("Failed to save cloud lua module: %s", err.Error())
	}
	output.Verbosef("Saved cloud lua module to: %s", cloudLuaModuleDir)
	ctx.cloudLuaModuleDir = cloudLuaModuleDir

	ctx.ControlPlane = cp
	buf := bytes.NewBuffer(nil)
	if err = helmEssentialConfigTemplate.Execute(buf, ctx); err != nil {
		return fmt.Errorf("Failed to execute helm essential config template: %s", err.Error())
	}
	ctx.essentialConfig = buf.Bytes()

	if err = createOnKubernetes(ctx, namespace, kubectl); err != nil {
		return fmt.Errorf("Failed to create namespace on kubernetes: %s", err.Error())
	}
	if err = createOnKubernetes(ctx, secret, kubectl); err != nil {
		return fmt.Errorf("Failed to create secret on kubernetes: %s", err.Error())
	}
	if err = createOnKubernetes(ctx, configMap, kubectl); err != nil {
		return fmt.Errorf("Failed to create configmap on kubernetes: %s", err.Error())
	}

	return nil
}

// createOnKubernetes create namespace, secret or configmap on Kubernetes
func createOnKubernetes(ctx *deployContext, k kind, kubectl commands.Cmd) error {
	var (
		err  error
		opts = ctx.KubernetesOpts
	)

	switch k {
	case secret:
		kubectl.AppendArgs("create", "secret", "generic", "cloud-ssl")
		kubectl.AppendArgs("--from-file", fmt.Sprintf("tls.crt=%s", filepath.Join(ctx.tlsDir, "tls.crt")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("tls.key=%s", filepath.Join(ctx.tlsDir, "tls.key")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("ca.crt=%s", filepath.Join(ctx.tlsDir, "ca.crt")))
		kubectl.AppendArgs("--namespace", opts.NameSpace)
	case configMap:
		kubectl.AppendArgs("create", "configmap", "cloud-module")
		kubectl.AppendArgs("--from-file", fmt.Sprintf("cloud.ljbc=%s", filepath.Join(ctx.cloudLuaModuleDir, "cloud.ljbc")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("cloud-agent.ljbc=%s", filepath.Join(ctx.cloudLuaModuleDir, "cloud-agent.ljbc")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("cloud-metrics.ljbc=%s", filepath.Join(ctx.cloudLuaModuleDir, "cloud-metrics.ljbc")))
		kubectl.AppendArgs("--from-file", fmt.Sprintf("cloud-utils.ljbc=%s", filepath.Join(ctx.cloudLuaModuleDir, "cloud-utils.ljbc")))
		kubectl.AppendArgs("--namespace", opts.NameSpace)
	case namespace:
		kubectl.AppendArgs("create", "ns", opts.NameSpace)
	default:
		panic(fmt.Sprintf("invaild kind: %d", k))
	}

	if options.Global.DryRun {
		output.Infof("Running:\n%s\n", kubectl.String())
	} else {
		output.Verbosef("Running:\n%s\n", kubectl.String())
	}

	newCtx, cancel := context.WithTimeout(context.TODO(), time.Minute*1)
	defer cancel()
	go utils.WaitForSignal(func() {
		cancel()
	})

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
