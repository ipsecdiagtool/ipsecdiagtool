package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
	"testing"
)

//Check that there's no deadlock when setting up and tearing down the capture-routine.
func TestStartStopCapture(t *testing.T) {
	mtuSample := config.MTUConfig{"127.0.0.1", "127.0.0.1", 10, 0, 2000, 20}
	mtuList := []config.MTUConfig{mtuSample, mtuSample}
	conf := config.Config{0, false, "localhost:514", 3000, mtuList, 32, "any", 60, 10, "", 0}

	icmpPackets := make(chan gopacket.Packet, 100)
	ipsecPackets := make(chan gopacket.Packet, 100)

	quitChannel := Start(conf, icmpPackets, ipsecPackets)
	quitChannel <- true

	if quitChannel == nil {
		t.Error("Quit Channel did not get initialized properly.")
	}
}

//Check that startCapture tries to read a pcap file when it is specified and that it
//returns a error if the file doesn't exist.
func TestReadFromPcap(t *testing.T) {
	icmpPackets := make(chan gopacket.Packet, 100)
	ipsecPackets := make(chan gopacket.Packet, 100)
	initChannels(icmpPackets, ipsecPackets)

	quit := make(chan bool)
	captureReady := make(chan bool)

	err := capture(500, quit, captureReady, "/test.pcap")
	if err.Error() != "/test.pcap: No such file or directory" {
		t.Error("Tried reading a pcap file that doesn't exist. Didn't get the correct error. Got", err, "instead.")
	}
	log.Println(err)
}

//Check that there is an error if the packet channel is full
func TestChannelFull(t *testing.T) {
	icmpPackets := make(chan gopacket.Packet, 100)
	var data []byte
	data = nil
	packet := gopacket.NewPacket(data, layers.LayerTypeEthernet, gopacket.Default)

	var err error
	err = nil
	i := 0
	for err == nil {
		i++
		err = putChannel(packet, icmpPackets)
		if i > 200 {
			t.Error("Channel should be full and there should be an error but there isn't.")
		}
	}
}
