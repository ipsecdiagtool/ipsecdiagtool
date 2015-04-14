package mtu

import (
	"testing"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"time"
)

func TestDetectMTU500(t *testing.T) {
	//Setup Listener
	mtuSample := config.MTUConfig{"127.0.0.1", "127.0.0.1", 4, 0, 2000}
	mtuList := []config.MTUConfig{mtuSample,mtuSample}
	confListener := config.Config{1, true, mtuList, 32, "any", 0}
	Listen(confListener, 1500)
	time.Sleep(1000 * time.Millisecond)

	//Setup Analyzer
	var confAnalyzer = config.Config{2, true, mtuList, 32, "any", 0}

	//Run test
	Analyze(confAnalyzer)


}
