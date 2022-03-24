package utils

import (
	"context"
	"fmt"
	"strings"

	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
)

// GetDeploymentName get the deploy name for APISIX instance
func GetDeploymentName(kubectl commands.Cmd) (string, error) {
	deployOpts := options.Global.Deploy
	kubectl.AppendArgs("get", "deployment", "-n", deployOpts.Kubernetes.NameSpace)
	kubectl.AppendArgs("-l", fmt.Sprintf("app.kubernetes.io/instance=%s", deployOpts.Name))
	kubectl.AppendArgs("-o", "jsonpath=\"{.items[0].metadata.name}\"")

	stdout, err := runKubectl(kubectl)
	if err != nil {
		return "", err
	}

	return stdout, nil
}

// GetPodsNames get the pod names for APISIX instance
func GetPodsNames(kubectl commands.Cmd) ([]string, error) {
	deployOpts := options.Global.Deploy
	kubectl.AppendArgs("get", "pods", "-n", deployOpts.Kubernetes.NameSpace)
	kubectl.AppendArgs("-l", fmt.Sprintf("app.kubernetes.io/instance=%s", deployOpts.Name))
	kubectl.AppendArgs("-o", "jsonpath=\"{.items[*].metadata.name}\"")
	stdout, err := runKubectl(kubectl)
	if err != nil {
		return nil, err
	}

	podsNames := strings.Split(strings.Replace(stdout, "\"", "", -1), " ")

	return podsNames, nil
}

// GetAPISIXID get the id for APISIX instance
func GetAPISIXID(kubectl commands.Cmd, podName string) (string, error) {
	deployOpts := options.Global.Deploy

	kubectl.AppendArgs("wait", "--for", "condition=Ready", "--timeout", "60s")
	kubectl.AppendArgs(fmt.Sprintf("pod/%s", podName), "-n", deployOpts.Kubernetes.NameSpace)
	if _, err := runKubectl(kubectl); err != nil {
		return "", err
	}

	kubectl.AppendArgs("exec", podName, "-n", deployOpts.Kubernetes.NameSpace)
	kubectl.AppendArgs("--", "cat", "/usr/local/apisix/conf/apisix.uid")

	stdout, err := runKubectl(kubectl)
	if err != nil {
		return "", err
	}

	return stdout, nil
}

// GetServiceName get the service name for APISIX instance
func GetServiceName(kubectl commands.Cmd) (string, error) {
	deployOpts := options.Global.Deploy
	kubectl.AppendArgs("get", "service", "-n", deployOpts.Kubernetes.NameSpace)
	kubectl.AppendArgs("-l", fmt.Sprintf("app.kubernetes.io/instance=%s", deployOpts.Name))
	kubectl.AppendArgs("-o", "jsonpath=\"{.items[0].metadata.name}\"")

	stdout, err := runKubectl(kubectl)
	if err != nil {
		return "", err
	}

	return stdout, nil
}

func runKubectl(kubectl commands.Cmd) (string, error) {
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
	}
	if stdout != "" {
		output.Verbosef(stdout)
	}
	if err != nil {
		return "", err
	}

	return stdout, nil
}
