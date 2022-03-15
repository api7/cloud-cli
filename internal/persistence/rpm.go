package persistence

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"

	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
	"github.com/api7/cloud-cli/internal/output"
)

var _OpenRestyRepoURL = "https://repos.apiseven.com/packages/centos/apache-apisix-repo-1.0-1.noarch.rpm"
var _APISIXRepoURL = "https://repos.apiseven.com/packages/centos/apache-apisix.repo"

// DownloadRPM installs related repositories and download RPM package of Apache APISIX with dependencies
// return empty string when apisix is installed
func DownloadRPM(ctx context.Context, version string) (file string, err error) {
	path := filepath.Join(os.Getenv("HOME"), ".api7cloud/rpm")
	if err := os.MkdirAll(path, 0700); err != nil {
		panic(err)
	}
	err = checkIfAPISIXAvailable(ctx, version)
	if err != nil {
		if os.IsExist(err) {
			return "", nil
		}
		return "", err
	}

	if existed := checkIfReposExisted(ctx); !existed {
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
	}
	if options.Global.DryRun {
		return filepath.Join(path, fmt.Sprintf("apisix-%s-0.el7.x86_64.rpm", version)), nil
	}
	// find out the apisix rpm package name
	dir, err := ioutil.ReadDir(path)
	if err != nil {
		return "", err
	}
	var name string
	for _, f := range dir {
		if strings.HasPrefix(f.Name(), "apisix-") {
			name = f.Name()
			break
		}
	}
	if name == "" {
		return "", errors.New("Failed to download Apache APISIX RPM package")
	}
	return filepath.Join(path, name), nil
}

func checkIfAPISIXAvailable(ctx context.Context, version string) error {
	if options.Global.DryRun {
		return nil
	}
	cmd := commands.New("apisix", false).AppendArgs("version")
	stdout, _, err := cmd.Run(ctx)
	if err != nil {
		return nil
	}
	if !strings.Contains(stdout, version) {
		return errors.New("other version of Apache APISIX already installed")
	} else {
		return os.ErrExist
	}
}

func checkIfReposExisted(ctx context.Context) bool {
	if options.Global.DryRun {
		return false
	}
	cmd := commands.New("yum", false).AppendArgs("list", "apisix")
	stdout, stderr, err := cmd.Run(ctx)
	if err != nil {
		output.Errorf(err.Error())
		return false
	}
	if stderr != "" {
		output.Verbosef(stderr)
		return false
	}
	return strings.Contains(stdout, "apisix.*")
}
