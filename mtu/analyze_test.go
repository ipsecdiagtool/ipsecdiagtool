package mtu

import (
	"testing"
)

func TestConfirmMTUNoResponse(t *testing.T) {
	var result = confirmMTU("127.0.0.1", "127.0.0.1", 200, 1)
	if result {
		t.Error("confirmMTU should return false. Check if there's some server still running.")
	}
}
