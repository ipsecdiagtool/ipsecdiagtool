package mtu

import (
	"testing"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"time"
)

func TestDetectMTU500(t *testing.T) {
	//Setup Listener
	mtuSample := config.MTUConfig{"127.0.0.1", "127.0.0.1", 1, 0, 2000}
	mtuList := []config.MTUConfig{mtuSample,mtuSample}
	config := config.Config{1337, true, mtuList, 32, "any", 0}
	Listen(config, 1600)
	time.Sleep(1000 * time.Millisecond)

	//Run test
	var exactMTU = Analyze(config)
	if(exactMTU != 1584){
		t.Error("Expected 1584 got",exactMTU,"instead.")
	}
}
