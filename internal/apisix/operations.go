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

package apisix

import (
	"context"
	"github.com/api7/cloud-cli/internal/commands"
	"github.com/api7/cloud-cli/internal/options"
)

// Reload reloads APISIX.
// This function only supports for APISIX running on bare metal.
func Reload(ctx context.Context) error {
	bin := options.Global.Deploy.Bare.APISIXBinPath
	dryrun := options.Global.DryRun
	reload := commands.New(bin, dryrun)
	reload.AppendArgs("reload")

	return reload.Execute(ctx)
}
