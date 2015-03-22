package mtu

import "testing"

func TestPayloadGeneratorLargerThen11(t *testing.T) {
	length := 100
	payload100 := generatePayload(length)
	if len(payload100) != length {
		t.Error("Expected", length, "bytes of payload, instead got", len(payload100), "bytes.")
	}
}

func TestPayloadGeneratorSmallerThen11(t *testing.T) {
	length := 5
	payload100 := generatePayload(length)
	if len(payload100) != length {
		t.Error("Expected", length, "bytes of payload, instead got", len(payload100), "bytes.")
	}
}

func TestPayloadGeneratorZero(t *testing.T) {
	length := 0
	payload100 := generatePayload(length)
	if len(payload100) != length {
		t.Error("Expected", length, "bytes of payload, instead got", len(payload100), "bytes.")
	}
}
