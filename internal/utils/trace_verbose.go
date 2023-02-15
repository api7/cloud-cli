package utils

import (
	"github.com/api7/cloud-cli/internal/output"
	"github.com/api7/cloud-go-sdk"
	"sync"
	"time"
)

var WantExit bool
var VerboseWg sync.WaitGroup

func VerboseGoroutine(traceChan <-chan *cloud.TraceSeries) {
	VerboseWg.Add(1)
	defer VerboseWg.Done()

	for {
		select {
		case data := <-traceChan:
			DumpTrace(data)
		default:
			if WantExit {
				return
			}
			time.Sleep(500 * time.Millisecond)
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
