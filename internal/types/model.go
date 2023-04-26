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

import (
	"time"
)

// Status represents an error type, it contains the error code and its
// description.
type Status struct {
	// Code is the error code number.
	// example: 0
	Code int `json:"code"`
	// Message describes the error code.
	// example: OK
	Message string `json:"message"`
}

// ResponseWrapper wraps the response with its original payload,
// and sets the Status field to codes.OK if everything is OK, but when
// the response is invalid, ErrorReason could be filled to show the error
// details and in such a case, Status is not codes.OK but a specific error
// code to show the kind.
type ResponseWrapper struct {
	// Payload carries the original data.
	// discriminator: true
	Payload interface{} `json:"payload,omitempty"`
	// Status shows the operation status for current request.
	Status Status `json:"status"`
	// ErrorReason is the error details, it's exclusive with Payload.
	ErrorReason string `json:"error,omitempty"`
	// Warning attaches a warning message to the response.
	Warning string `json:"warning,omitempty"`
}

// TypeMeta contains some common and basic items, like id, name.
type TypeMeta struct {
	// ID is the unique identify to mark an object.
	ID string `json:"id,inline" yaml:"id"`
	// Name is the object name.
	Name string `json:"name" yaml:"name"`
	// CreatedAt is the object creation time.
	CreatedAt time.Time `json:"created_at" yaml:"created_at"`
	// UpdatedAt is the last modified time of this object.
	UpdatedAt time.Time `json:"updated_at" yaml:"updated_at"`
}

// Organization is the specification of organization.
type Organization struct {
	TypeMeta `json:",inline" yaml:",inline"`
	// PlanID indicates which plan is used by this organization.
	// PlanID should refer to a valid Plan object.
	PlanID string `json:"plan_id" yaml:"plan_id"`
	// PlanExpireTime indicates the binding plan expire time for this organization.
	PlanExpireTime time.Time `json:"plan_expire_time" yaml:"plan_expire_time"`
	// SubscriptionStartedAt is the time when the organization subscribed to the plan.
	SubscriptionStartedAt *time.Time `json:"subscription_started_at" yaml:"subscription_started_at"`
	// OwnerID indicates who create the organization.
	OwnerID string `json:"owner_id" yaml:"owner_id"`
}

// Cluster is the specification of cluster.
type Cluster struct {
	TypeMeta `json:",inline" yaml:",inline"`
	// OrganizationID refers to an Organization object, which
	// indicates the belonged organization for this cluster.
	OrganizationID string `json:"org_id" yaml:"org_id"`
	// RegionID refers to a Region object, which indicates the
	// region that the Cloud Plane resides.
	RegionID string `json:"region_id" yaml:"region_id"`
	// Status indicates the cluster status, candidate values are:
	// * ClusterBuildInProgress: the cluster is being created.
	// * ClusterCreating means a cluster is being created.
	// * ClusterNormal: the cluster is built, and can be used normally.
	// * ClusterCreateFailed means a cluster was not created successfully.
	// * ClusterDeleting means a cluster is being deleted.
	// * ClusterDeleted means a cluster was deleted.
	// enum: ClusterBuildInProgress:1,ClusterNormal:2,ClusterCreateFailed:3,ClusterDeleting:4,ClusterDeleted:5
	Status int `json:"status" yaml:"status"`
	// Domain is the domain assigned by APISEVEN Cloud and has correct
	// records so that DP instances can access APISEVEN Cloud by it.
	Domain string `json:"domain" yaml:"domain"`
	// ConfigPayload is the customize gateway config for specific cluster
	ConfigPayload string `json:"config_payload" yaml:"config_payload"`
}

// GetOrganizationClusterResponsePayload contains list clusters request
type GetOrganizationClusterResponsePayload struct {
	// Count is total count of clusters
	Count uint64 `json:"count" uri:"count"`
	// List is array of clusters
	List []*ClusterSummary `json:"list"`
}

// ClusterStartupConfigResponsePayload contains APISIX startup config.
type ClusterStartupConfigResponsePayload struct {
	// Configuration is the startup config
	Configuration string `json:"configuration"`
}

// ClusterSummary is cluster with region and org summary
type ClusterSummary struct {
	Cluster `json:",inline" yaml:",inline"`
	// OrgName is the org name of the cluster
	OrgName string `json:"org_name" yaml:"org_name"`
}

// TLSBundle contains a pair of certificate, private key,
// and the issuing certificate.
type TLSBundle struct {
	Certificate   string `json:"certificate" yaml:"certificate"`
	PrivateKey    string `json:"private_key" yaml:"private_key"`
	CACertificate string `json:"ca_certificate" yaml:"ca_certificate"`
}
