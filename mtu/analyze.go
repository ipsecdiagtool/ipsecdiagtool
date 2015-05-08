package mtu

import (
	"code.google.com/p/gopacket"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"github.com/ipsecdiagtool/ipsecdiagtool/logging"
	"log"
	"sort"
	"strconv"
	"time"
)

var initalized = false
var conf config.Config
var icmpPacketsStage1 chan gopacket.Packet
var icmpPacketsStage2 chan gopacket.Packet

//Init the MTU package so that you can call FindAll()
func Init(config config.Config, icmpPackets chan gopacket.Packet) {
	conf = config
	icmpPacketsStage1 = icmpPackets
	icmpPacketsStage2 = make(chan gopacket.Packet, 100)
	initalized = true
	go handlePackets(icmpPacketsStage1, icmpPacketsStage2, conf.ApplicationID)
}

//FindAll finds the MTU for each connection specified in the
//configuration. Use Find() if you're only looking for a specific MTU.
func FindAll() {
	if !initalized {
		log.Println("Please make sure that the MTU package was configured with mtu.Init(.., ..)")
		return
	}
	c := conf
	//Setup a mtuOK channel for each config
	var mtuOkChannels = make(map[int]chan int)
	for conf := range c.MTUConfList {
		mtuOkChannels[conf] = make(chan int, 100)
	}

	//go handlePackets(icmpPacketsStage2, c.ApplicationID, mtuOkChannels)
	go distributeMtuOkPackets(icmpPacketsStage2, mtuOkChannels)

	for conf := range c.MTUConfList {
		log.Println("------------------------- MTU Conf", conf, " -------------------------")
		go Find(
			c.MTUConfList[conf],
			c.ApplicationID,
			conf,
			mtuOkChannels[conf])
	}
}

//Find finds the ideal MTU between two nodes by sending batches of packets with varying sizes
//to a remote node. The remote nodes is requires to respond to those packets if it received them.
//so it can determine the largest packet that was received on the remote node and the smallest packet that
//went missing. In a next step FastMTU sends again a batch of packets with sizes between the largest successful
//and smallest unsuccessful packet. This behaviour is continued until the size-difference between individual
//packets is no larger then 1Byte. Once that happens the largest successful packet is reported as MTU.
func Find(mtuConf config.MTUConfig, appID int, chanID int, mtuOK chan int) int {
	if !initalized {
		log.Println("Please make sure that the MTU package was configured with mtu.Init(.., ..)")
		return 0
	}
	var rangeStart = mtuConf.MTURangeStart
	var rangeEnd = mtuConf.MTURangeEnd
	var itStep = ((rangeEnd - rangeStart) / 20)
	var roughMTU = 0
	var mtuDetected = false
	var retries = 0

	for !mtuDetected {
		if itStep == 0 {
			itStep = 1
			mtuDetected = true
		}
		roughMTU = sendBatch(mtuConf.SourceIP, mtuConf.DestinationIP, rangeStart, rangeEnd, itStep, mtuConf.Timeout, appID, chanID, mtuOK)

		if roughMTU == rangeEnd {
			rangeStart = rangeEnd
			rangeEnd = 2 * rangeEnd
		} else if roughMTU == 0 {
			//Retry
			if retries < 1 {
				retries++
				log.Println("ERROR: Reported 0.. trying again.")
				roughMTU = sendBatch(mtuConf.SourceIP, mtuConf.DestinationIP, rangeStart, rangeEnd, itStep, mtuConf.Timeout, appID, chanID, mtuOK)
			} else {
				log.Println("ERROR: Reported MTU 0.. ")
				mtuDetected = true
			}
		} else {
			rangeStart = roughMTU
			rangeEnd = roughMTU + itStep
			itStep = ((rangeEnd - rangeStart) / 20)
		}
	}
	log.Println("Exact MTU found", roughMTU)
	logging.InfoLog("Exact MTU found " + strconv.Itoa(roughMTU))
	return roughMTU
}

func sendBatch(srcIP string, destIP string, rangeStart int, rangeEnd int, itStep int, timeoutInSeconds time.Duration, appID int, chanID int, mtuOK chan int) int {
	//1. Send a batch of packets
	var results = make(map[int]bool)
	for i := rangeStart; i < (rangeEnd + itStep); i += itStep {
		sendPacket(srcIP, destIP, i, "MTU?", appID, chanID)
		results[i] = false
	}

	//2. Wait until time's up then gather results
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(timeoutInSeconds * time.Second)
		timeout <- true
	}()

	var largestSuccessfulPacket = 0

	var gatherPackets = true
	for gatherPackets {
		select {
		case goodPacket := <-mtuOK:
			if goodPacket > largestSuccessfulPacket {
				//Check if the packet we received was one that we sent. (Based on size)
				if _, ok := results[goodPacket]; ok {
					largestSuccessfulPacket = goodPacket
					results[goodPacket] = true
				}
			}
		case <-timeout:
			gatherPackets = false
		}
	}

	log.Println("---------------------------------------------------")
	log.Println("Range:", rangeStart, "-", rangeEnd, "  itStep:", itStep, "  Timeout:", timeoutInSeconds)
	log.Println("Largest successful packet:", largestSuccessfulPacket)
	printResultMap(results)
	log.Println()

	return largestSuccessfulPacket
}

func printResultMap(input map[int]bool) {
	// To store the keys in slice in sorted order
	var keys []int
	for k := range input {
		keys = append(keys, k)
	}
	sort.Ints(keys)

	// To perform the opertion you want
	var received = "Received: "
	var missing = "Missing: "
	for _, k := range keys {
		if input[k] {
			received += (strconv.Itoa(k) + " ")
		} else {
			missing += (strconv.Itoa(k) + " ")
		}
	}
	log.Println(received)
	log.Println(missing)
}
