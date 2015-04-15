package mtu

import (
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"testing"
)

func TestDetectMTU1584(t *testing.T) {
	//Test Setup
	mtuSample := config.MTUConfig{"127.0.0.1", "127.0.0.1", 1, 0, 2000}
	mtuList := []config.MTUConfig{mtuSample, mtuSample}
	config := config.Config{1337, true, mtuList, 32, "any", 0}

	//Run test
	var exactMTU = Analyze(config, 1600)
	if exactMTU != 1584 {
		t.Error("Expected 1584 got", exactMTU, "instead.")
	}
}
