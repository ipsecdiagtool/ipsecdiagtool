package mtu

import (
	"testing"
)

func TestSetupRandomAppID(t *testing.T) {
	Setup(0, "127.0.0.1", "127.0.0.1", 22, 10)
	var firstAppID = appID

	Setup(0, "127.0.0.1", "127.0.0.1", 22, 10)
	var secondAppID = appID

	if firstAppID == secondAppID {
		t.Error("Expected random AppID got two times the same AppID instead.")
	}
}

func TestSetupSettingEverything(t *testing.T) {
	var id, ip1, ip2, port, step = 737, "127.0.0.1", "127.0.0.2", 22, 10
	Setup(id, ip1, ip2, port, step)
	if appID != id || srcIP.String() != ip1 || destIP.String() != ip2 || destPort != port || incStep != step {
		t.Error("Expected",id,ip1,ip2,port,step,"got",appID,srcIP.String(),destIP.String(),destPort,incStep,"instead.")
	}
}
