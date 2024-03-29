// Copyright 2023 API7.ai, Inc
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

package types

// K8sResourceKind is the resource kind of kubernetes
type K8sResourceKind uint8

const (
	// ConfigMap is the kubernetes Resource that kind is configmap
	ConfigMap K8sResourceKind = iota
	// Secret is the kubernetes Resource that kind is secret
	Secret
	// Namespace is the namespace of kubernetes
	Namespace
)
