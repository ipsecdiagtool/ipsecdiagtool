package mtu

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"log"
	"testing"
)

//TODO: test doesn't work yet, I'll probably need to split the send function.
func TestHandlePacket(t *testing.T) {
	setDefaultValues()
	var data = sendPacket(srcIP, destIP, destPort, 200, "Testing.")
	var packet = gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)
	log.Println(packet.Dump())
	//handlePacket(packet)

	/*
		if firstAppID == secondAppID {
			t.Error("Expected random AppID got two times the same AppID instead.")
		}*/
}
