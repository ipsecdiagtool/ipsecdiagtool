package mtu

import (
	"net"
	"testing"
)

func TestPayloadGeneratorLargerThen11(t *testing.T) {
	length := 100
	payload := generatePayload(length, "Hello IPSec")
	if len(payload) != length {
		t.Error("Expected", length, "bytes of payload, instead got", len(payload), "bytes.")
	}
}

func TestPayloadGeneratorSmallerThen11(t *testing.T) {
	length := 5
	payload := generatePayload(length, "Hello IPSec")
	if len(payload) != length {
		t.Error("Expected", length, "bytes of payload, instead got", len(payload), "bytes.")
	}
}

func TestPayloadGeneratorZero(t *testing.T) {
	length := 0
	payload := generatePayload(length, "Hello IPSec")
	if len(payload) != length {
		t.Error("Expected", length, "bytes of payload, instead got", len(payload), "bytes.")
	}
}

func TestSendPacketLarge(t *testing.T) {
	length := 1000
	payload := sendPacket(net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"), 22, length, "Hello IPSec")
	if len(payload) != length {
		t.Error("Expected", length, "bytes of payload, instead got", len(payload), "bytes.")
	}
}

func TestSendPacketSmall(t *testing.T) {
	length := 10
	expected := 40
	payload := sendPacket(net.ParseIP("127.0.0.1"), net.ParseIP("127.0.0.1"), 22, length, "Hello IPSec")
	if len(payload) != expected {
		t.Error("Expected", expected, "bytes of payload, instead got", len(payload), "bytes.")
	}
}
