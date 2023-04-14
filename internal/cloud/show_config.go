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

package cloud

import (
	"context"
	"fmt"
	"github.com/api7/cloud-go-sdk"
)

var (
	_validResources = map[string]struct{}{
		"application":      {},
		"api":              {},
		"consumer":         {},
		"certificate":      {},
		"cluster_settings": {},
	}
)

func (a *api) DebugShowConfig(clusterID cloud.ID, resource string, id cloud.ID) (string, error) {
	if _, ok := _validResources[resource]; !ok {
		return "", fmt.Errorf("invalid resource type: %s", resource)
	}

	var (
		data string
		err  error
	)

	switch resource {
	case "application":
		data, err = a.sdk.DebugApplicationResources(context.TODO(), id, &cloud.ResourceGetOptions{
			Cluster: &cloud.Cluster{
				ID: clusterID,
			},
		})
	case "api":
		data, err = a.sdk.DebugAPIResources(context.TODO(), id, &cloud.ResourceGetOptions{
			Cluster: &cloud.Cluster{
				ID: clusterID,
			},
		})
	case "consumer":
		data, err = a.sdk.DebugConsumerResources(context.TODO(), id, &cloud.ResourceGetOptions{
			Cluster: &cloud.Cluster{
				ID: clusterID,
			},
		})
	case "certificate":
		data, err = a.sdk.DebugCertificateResources(context.TODO(), id, &cloud.ResourceGetOptions{
			Cluster: &cloud.Cluster{
				ID: clusterID,
			},
		})
	case "cluster_settings":
		data, err = a.sdk.DebugClusterSettings(context.TODO(), &cloud.ResourceGetOptions{
			Cluster: &cloud.Cluster{
				ID: clusterID,
			},
		})
	}

	if err != nil {
		return "", err
	}
	return data, nil
}
