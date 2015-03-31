package mtu

import (
	"testing"
	"net"
)

func TestConfirmMTUNoResponse(t *testing.T) {
	var result = confirmMTU(net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"), 22, 200, 1)
	if(result){
		t.Error("confirmMTU should return false. Check if there's some server still running.")
	}
}
