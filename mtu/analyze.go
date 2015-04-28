package mtu

import (
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
	"sort"
	"strconv"
	"time"
	"code.google.com/p/gopacket"
)

//FindAll accepts a configuration and a mtuOK channel. It finds the MTU for each connection specified in the
//configuration. Use Find() if you're only looking for a specific MTU.
func FindAll(c config.Config, icmpPackets chan gopacket.Packet) {

	//Setup a mtuOK channel for each config
	var mtuOkChannels = make(map[int]chan int)
	for conf := range c.MTUConfList {
		mtuOkChannels[conf] = make(chan int, 100)
	}

	go handlePackets(icmpPackets, c.ApplicationID, mtuOkChannels)

	for conf := range c.MTUConfList {
		log.Println("------------------------- MTU Conf",conf," -------------------------")
		go Find(
			c.MTUConfList[conf].SourceIP,
			c.MTUConfList[conf].DestinationIP,
			c.MTUConfList[conf].Timeout, c.ApplicationID,
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
func Find(srcIP string, destIP string, timeoutInSeconds time.Duration, appID int, chanID int, mtuOK chan int) int {

	var rangeStart = 0
	var rangeEnd = 2000
	var itStep = ((rangeEnd - rangeStart) / 20)
	var roughMTU = 0
	var mtuDetected = false
	var retries = 0

	for !mtuDetected {
		if itStep == 0 {
			itStep = 1
			mtuDetected = true
		}
		roughMTU = sendBatch(srcIP, destIP, rangeStart, rangeEnd, itStep, timeoutInSeconds, appID, chanID, mtuOK)

		if roughMTU == rangeEnd {
			rangeStart = rangeEnd
			rangeEnd = 2 * rangeEnd
		} else if roughMTU == 0 {
			//Retry
			if retries < 1 {
				retries++
				log.Println("Reported 0.. trying again.")
				roughMTU = sendBatch(srcIP, destIP, rangeStart, rangeEnd, itStep, timeoutInSeconds, appID, chanID, mtuOK)
			} else {
				log.Println("ERROR: Reported MTU 0.. ")
				mtuDetected = true //TODO: better name for mtuDetected needed?
			}
		} else {
			rangeStart = roughMTU
			rangeEnd = roughMTU + itStep
			itStep = ((rangeEnd - rangeStart) / 20)
		}
	}
	log.Println("Exact MTU found", roughMTU)
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
