package utils

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/consts"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
)

var (
	getAPISIXIDRetry = 3
)

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

func GetAPISIXID(kubectl commands.Cmd, podName string) (string, error) {
	var (
		stdout     string
		err        error
		retry      int
		deployOpts = options.Global.Deploy
	)

	for {
		retry++

		kubectl.AppendArgs("exec", podName, "-n", deployOpts.Kubernetes.NameSpace)
		kubectl.AppendArgs("--", "cat", "/usr/local/apisix/conf/apisix.uid")
		if stdout, err = runKubectl(kubectl); err == nil {
			break
		} else if strings.Contains(err.Error(), "container not found") && retry < getAPISIXIDRetry {
			output.Warnf("After %d seconds will be auto retry!", retry)
			time.Sleep(time.Second * time.Duration(retry))
		} else {
			return "", err
		}
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
		return "", errors.Wrap(err, stderr)
	}

	return stdout, nil
}
