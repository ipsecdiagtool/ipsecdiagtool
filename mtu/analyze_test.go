package mtu

import (
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"testing"
	"github.com/ipsecdiagtool/ipsecdiagtool/capture"
	"code.google.com/p/gopacket"
	"github.com/ipsecdiagtool/ipsecdiagtool/logging"
)

func testFind(simulatedMTU int, rangeStart int, rangeEnd int) int {
	//Test Setup
	mtu := config.MTUConfig{"127.0.0.1", "127.0.0.1", 2, rangeStart, rangeEnd}
	mtuList := []config.MTUConfig{mtu, mtu}

	conf := config.Config{0, false, "localhost:514", int32(simulatedMTU+16), mtuList, 32, "any", 60, 10, "", 0}
	logging.InitLoger(conf.SyslogServer, conf.AlertCounter, conf.AlertTime)

	icmpPackets := make(chan gopacket.Packet, 100)
	ipsecPackets := make(chan gopacket.Packet, 100)
	Init(conf, icmpPackets)
	var capQuit chan bool
	capQuit = capture.Start(conf, icmpPackets, ipsecPackets)

	var mtuOkChannels = make(map[int]chan int)
	for conf := range conf.MTUConfList {
		mtuOkChannels[conf] = make(chan int, 100)
	}

	go distributeMtuOkPackets(icmpPacketsStage2, mtuOkChannels)

	//TEST
	result := Find(mtu.SourceIP, mtu.DestinationIP, mtu.Timeout, conf.ApplicationID, 0, mtuOkChannels[0])
	capQuit <- true
	return result
}

//Start with a range of 0-2000 and detect the simulated MTU which is 500.
func TestDetectMTU500(t *testing.T) {
	tMTU := 500
	result := testFind(tMTU, 0, 2000)

	if result != (tMTU) {
		t.Error("Expected", (tMTU), "got", result, "instead.")
	}
}

//Start with a range of 0-2000 and detect the simulated MTU which is 1500.
func TestDetectMTU1500(t *testing.T) {
	tMTU := 1500
	result := testFind(tMTU, 0, 2000)

	if result != (tMTU) {
		t.Error("Expected", (tMTU), "got", result, "instead.")
	}
}

//Start with a range of 0-500 and detect the simulated MTU which is 1500.
func TestDetectMTU1500withSmallRange(t *testing.T) {
	tMTU := 1500
	result := testFind(tMTU, 0, 500)

	if result != (tMTU) {
		t.Error("Expected", (tMTU), "got", result, "instead.")
	}
}
