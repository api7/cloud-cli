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

package utils

import (
	"context"
	"sync"

	"github.com/api7/cloud-go-sdk"

	"github.com/api7/cloud-cli/internal/output"
)

type traceVerbose struct {
	ctx  context.Context
	Exit context.CancelFunc
	Wg   sync.WaitGroup
}

var TraceVerbose traceVerbose

func init() {
	TraceVerbose.Wg = sync.WaitGroup{}
	TraceVerbose.ctx, TraceVerbose.Exit = context.WithCancel(context.Background())
}

func VerboseGoroutine(traceChan <-chan *cloud.TraceSeries) {
	TraceVerbose.Wg.Add(1)
	defer TraceVerbose.Wg.Done()

	for {
		select {
		case data := <-traceChan:
			DumpTrace(data)
		case <-TraceVerbose.ctx.Done():
			return
		}
	}
}

func DumpTrace(data *cloud.TraceSeries) {
	req := data.Request

	output.Verbosef("[%v] Send a %v request to %v", data.ID, req.Method, req.URL.String())
	if len(data.RequestBody) != 0 {
		output.Verbosef("[%v] Send a request body: %s", data.ID, string(data.RequestBody))
	}

	output.Verbosef("[%v] Receive a response with status: %v", data.ID, data.Response.StatusCode)
	if len(data.ResponseBody) != 0 {
		output.Verbosef("[%v] Receive a response body: %s", data.ID, string(data.ResponseBody))
	}

	evts := data.Events
	output.Verbosef("[%v] Dump %d events:", data.ID, len(evts))
	for i, evt := range evts {
		output.Verbosef("[%v] Event#%d %s : %s", data.ID, i, evt.HappenedAt.Format("2006-01-02 15:04:05"), evt.Message)
	}
}
