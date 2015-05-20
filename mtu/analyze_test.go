package mtu

import (
	"code.google.com/p/gopacket"
	"github.com/ipsecdiagtool/ipsecdiagtool/capture"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"github.com/ipsecdiagtool/ipsecdiagtool/logging"
	"testing"
	"sync"
)

/*
 * These tests will only work if there's a IPSecDiagTool instance running that answers
 * the "MTU?" requests. You can start a local instance of IPSecDiagTool via 'ipsecdiagtool install' -->
 * 'service ipsecdiagtool start'.
 * It's also possible to use a remote instance of IPSecDiagTool by changing the src & dstIP below.
 */

const srcIP string = "127.0.0.1"
const dstIP string = "127.0.0.1"

func testFind(simulatedMTU int, rangeStart int, rangeEnd int) int {
	mtu := config.MTUConfig{srcIP, dstIP, 5, rangeStart, rangeEnd, 20}
	mtuList := []config.MTUConfig{mtu, mtu}

	//Experimental: AppID=1337 is allowed to answer it's own packets.
	conf := config.Config{0, false, "localhost:514", int32(simulatedMTU + 16), mtuList, 32, "any", 60, 10, "", 0}
	logging.InitLoger(conf.SyslogServer, conf.AlertCounter, conf.AlertTime)

	icmpPackets := make(chan gopacket.Packet, 500)
	ipsecPackets := make(chan gopacket.Packet, 500)
	Init(conf, icmpPackets)
	var capQuit chan bool
	capQuit = capture.Start(conf, icmpPackets, ipsecPackets)

	var mtuOkChannels = make(map[int]chan int)
	for conf := range conf.MTUConfList {
		mtuOkChannels[conf] = make(chan int, 100)
	}

	var quitDistribute = make(chan bool)
	go distributeMtuOkPackets(icmpPacketsStage2, mtuOkChannels, quitDistribute)

	//TEST
	var wg sync.WaitGroup
	wg.Add(1)
	result := Find(mtu, conf.ApplicationID, 0, mtuOkChannels[0], &wg)
	wg.Wait()
	quitDistribute <- true
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
