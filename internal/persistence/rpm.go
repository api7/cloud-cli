package persistence

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
)

var _OpenRestyRepoURL = "https://repos.apiseven.com/packages/centos/apache-apisix-repo-1.0-1.noarch.rpm"
var _APISIXRepoURL = "https://repos.apiseven.com/packages/centos/apache-apisix.repo"

// DownloadRPM installs related repositories and download RPM package of Apache APISIX with dependencies
func DownloadRPM(ctx context.Context, version string) (string, error) {
	path := filepath.Join(os.Getenv("HOME"), ".api7cloud/rpm")
	if err := os.MkdirAll(path, 0700); err != nil {
		panic(err)
	}
	err := checkIfAPISIXAvailable(ctx, version)
	if err != nil {
		if os.IsExist(err) {
			return "", nil
		}
		return "", err
	}

	cmd := commands.New("yum", options.Global.DryRun)
	cmd.AppendArgs("install", "-y", _OpenRestyRepoURL)
	if err := cmd.Execute(ctx); err != nil {
		return "", errors.Wrap(err, "failed to install repositories of OpenResty")
	}

	cmd = commands.New("yum-config-manager", options.Global.DryRun)
	cmd.AppendArgs("--add-repo", _APISIXRepoURL)
	if err := cmd.Execute(ctx); err != nil {
		return "", errors.Wrap(err, "failed to install repositories of Apache APISIX")
	}

	cmd = commands.New("yum", options.Global.DryRun)
	cmd.AppendArgs("install", "-y", "--downloadonly")
	cmd.AppendArgs("--downloaddir=" + path)
	cmd.AppendArgs("apisix-" + version)
	if err := cmd.Execute(ctx); err != nil {
		return "", errors.Wrap(err, "failed to download Apache APISIX")
	}
	return path, nil
}

func checkIfAPISIXAvailable(ctx context.Context, version string) error {
	cmd := commands.New("apisix", false).
		AppendArgs("version")
	_, stdout, err := cmd.Run(ctx)
	if err != nil {
		return nil
	}
	if !strings.Contains(stdout, version) {
		return errors.New(fmt.Sprintf("other vesrion of Apache APISIX already installed"))
	} else {
		return os.ErrExist
	}
	return nil
}
