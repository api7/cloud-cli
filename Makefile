# Licensed to the Apache Software Foundation (ASF) under one or more
# contributor license agreements.  See the NOTICE file distributed with // this work for additional information regarding copyright ownership.
# The ASF licenses this file to You under the Apache License, Version 2.0
# (the "License"); you may not use this file except in compliance with
# the License.  You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

BUILD_DATE ?= "$(shell date +"%Y-%m-%dT%H:%M")"
GITSHA=$(shell git rev-parse --short=7 HEAD)

MAJORSYM="$(shell go list -m)/internal/pkg/version._major"
MINORSYM="$(shell go list -m)/internal/pkg/version._minor"
BUILDDATESYM="$(shell go list -m)/internal/pkg/version._buildDate"
GITCOMMITSYM="$(shell go list -m)/internal/pkg/version._gitCommit"
VERSION_MAJOR=0
VERSION_MINOR=1
BINDIR=bin

GO_LDFLAGS ?= "-X=$(MAJORSYM)=$(VERSION_MAJOR) -X=$(MINORSYM)=$(VERSION_MINOR) -X=$(BUILDDATESYM)=$(BUILD_DATE) -X=$(GITCOMMITSYM)=$(GITSHA)"

default: help

.PHONY: help
help:  ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z0-9_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: build
build: create-bin-dir ## Build the binary
	go build -ldflags $(GO_LDFLAGS) -o bin/cloud-cli github.com/api7/cloud-cli

create-bin-dir:
	@mkdir -p $(BINDIR)
.PHONY: create-bin-dir

gofmt: ## Format the source code
	@find . -type f -name "*.go" | xargs gofmt -w -s
.PHONY: gofmt

lint: ## Apply go lint check
	@golangci-lint run --timeout 10m ./...
.PHONY: lint

test: ## Run the unit tests
	# go test run cases in different package parallel by default, but cloud cli config file is referenced by multi test cases, so we need to run them in sequence with -p=1
	@go test -p 1 -v ./...

.PHONY: install-tools
install-tools: ## Install necessary tools
	@bash -c 'go install github.com/golang/mock/mockgen@v1.6.0'

.PHONY: codegen
codegen: install-tools ## Run code generation
	./scripts/mockgen.sh
