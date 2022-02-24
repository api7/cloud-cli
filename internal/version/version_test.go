package version

import (
	"encoding/json"
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestVersion(t *testing.T) {
	ver := Version{
		Major:     "0",
		Minor:     "1",
		GitCommit: "2ad4hz",
		BuildDate: time.Now().String(),
		GoVersion: runtime.Version(),
		Compiler:  runtime.Compiler,
		Platform:  runtime.GOOS + "/" + runtime.GOARCH,
	}
	s := ver.String()
	var (
		ver2 Version
	)
	err := json.Unmarshal([]byte(s), &ver2)
	assert.Nil(t, err, "unmarshalling version info")
	assert.Equal(t, ver, ver2, "checking version")
}
