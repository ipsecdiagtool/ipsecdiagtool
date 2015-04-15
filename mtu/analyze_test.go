package mtu

import (
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"testing"
	"time"
)

//Test Settings
var tOverhead = 16
var tTimeout = 5

//Start with a range of 0-2000 and detect the simulated MTU which is 500.
func TestDetectMTU500(t *testing.T) {
	//Test Settings
	tMTU := 500
	tRangeStart := 0
	tRangeEnd := 2000

	//Test Setup
	mtuConfig := config.MTUConfig{"127.0.0.1", "127.0.0.1", time.Duration(tTimeout), tRangeStart, tRangeEnd}
	mtuList := []config.MTUConfig{mtuConfig, mtuConfig}
	appConfig := config.Config{1337, true, mtuList, 0, "_", 0}

	//Run test & validate result
	var detectedMTU = Analyze(appConfig, int32(tMTU))
	if detectedMTU != (tMTU-tOverhead) {
		t.Error("Expected", (tMTU-tOverhead), "got", detectedMTU, "instead.")
	}
}

//Start with a range of 0-2000 and detect the simulated MTU which is 1600.
func TestDetectMTU1600(t *testing.T) {
	//Test Settings
	tMTU := 1600
	tRangeStart := 0
	tRangeEnd := 2000

	//Test Setup
	mtuConfig := config.MTUConfig{"127.0.0.1", "127.0.0.1", time.Duration(tTimeout), tRangeStart, tRangeEnd}
	mtuList := []config.MTUConfig{mtuConfig, mtuConfig}
	appConfig := config.Config{1337, true, mtuList, 0, "_", 0}

	//Run test & validate result
	var detectedMTU = Analyze(appConfig, int32(tMTU))
	if detectedMTU != (tMTU-tOverhead) {
		t.Error("Expected", (tMTU-tOverhead), "got", detectedMTU, "instead.")
	}
}

//Start with a range of 0-2000 and detect the simulated MTU which is 3000.
func TestDetectMTU3000(t *testing.T) {
	//Test Settings
	tMTU := 3000
	tRangeStart := 0
	tRangeEnd := 2000

	//Test Setup
	mtuConfig := config.MTUConfig{"127.0.0.1", "127.0.0.1", time.Duration(tTimeout), tRangeStart, tRangeEnd}
	mtuList := []config.MTUConfig{mtuConfig, mtuConfig}
	appConfig := config.Config{1337, true, mtuList, 0, "_", 0}

	//Run test & validate result
	var detectedMTU = Analyze(appConfig, int32(tMTU))
	if detectedMTU != (tMTU-tOverhead) {
		t.Error("Expected", (tMTU-tOverhead), "got", detectedMTU, "instead.")
	}
}
