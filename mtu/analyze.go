package mtu

import (
	"code.google.com/p/gopacket/examples/util"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
	"time"
)

//Package internal temp. variables
var mtuOKchan = make(chan int, 100)

//Analyze determines the ideal MTU (Maximum Transmission Unit) between two nodes
//by sending increasingly big packets between them. Analyze determine the MTU
//by running FindMTU multiple times. However it is not a daemon. Once it has found
//the ideal MTU it reports it and then closes shop.
func Analyze(c config.Config) int {
	defer util.Run()()
	log.Println("Analyzing MTU..")
	log.Println(c)

	//Setup a channel for communication with capture
	mtuOKchan = make(chan int) // Allocate a channel.

	go startCapture("icmp", 1600, c.ApplicationID)

	//TODO: currently required to give capture enough time to boot..
	time.Sleep(1000 * time.Millisecond)

	return FastMTU(
	c.MTUConfList[0].SourceIP, //TODO: use additional configs as well, not just first. --> Iterate
	c.MTUConfList[0].DestinationIP,
	10, c.ApplicationID); //TODO: use value from config
}

//TODO: reduce duplicate code
//Listen only listens to MTU requests and replies with OK-Packets.
func Listen(c config.Config, snaplen int32){
	defer util.Run()()

	//Setup a channel for communication with capture
	mtuOKchan = make(chan int) // Allocate a channel

	log.Println("Listener", c)

	go startCapture("icmp", snaplen, c.ApplicationID)
}

//FindMTU discovers the MTU between two nodes and returns it as an int value. FindMTU currently
//increases the packet size until it runs into a timeout. Once it runs into the timeout it returns
//the last known as good MTU.
func FindMTU(srcIP string, destIP string, startMTU int, increment int, appID int) int {
	var goodMTU, nextMTU = 0, startMTU

	//1. Initiate MTU discovery by sending first packet.
	sendPacket(srcIP, destIP, nextMTU, "MTU?", appID)

	//2. Either we get a message from our mtu channel or the timeout channel will message us after 10s.
	for {
		//2.1 Setting up the timeout channel
		//http://blog.golang.org/go-concurrency-patterns-timing-out-and
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(config.Timeout * time.Second)
			timeout <- true
		}()

		select {
		case <-mtuOKchan:
			goodMTU = nextMTU
			nextMTU += increment
			time.Sleep(1000 * time.Millisecond)
			sendPacket(srcIP, destIP,nextMTU, "MTU?", appID)
		case <-timeout:
			log.Println("Timeout has occured. We've steped over the MTU!")
			return goodMTU
		}
	}
}

//Detects the exact MTU asap.
func FastMTU(srcIP string, destIP string, timeoutInSeconds time.Duration, appID int) int{

	var rangeStart = 0
	var rangeEnd = 2000
	var itStep = ((rangeEnd-rangeStart)/20)
	var roughMTU = 0
	var mtuDetected = false

	for !mtuDetected {
		if itStep == 0 {
			itStep = 1
			mtuDetected = true
		}
		roughMTU = sendBatch(srcIP, destIP, rangeStart, rangeEnd, itStep, timeoutInSeconds, appID)

		if(roughMTU == rangeEnd){
			rangeStart = rangeEnd
			rangeEnd = 2*rangeEnd
		} else if (roughMTU == 0){
			log.Println("ERROR: Reported MTU 0.. ")
			mtuDetected = true //TODO: better name for mtuDetected needed?
		} else {
			rangeStart = roughMTU
			rangeEnd = roughMTU+itStep
			itStep = ((rangeEnd-rangeStart)/20)
		}
	}
	log.Println("Exact MTU found", roughMTU)
	return roughMTU
}

func sendBatch(srcIP string, destIP string, rangeStart int, rangeEnd int, itStep int, timeoutInSeconds time.Duration, appID int) int {
	//1. Send a batch of packets
	var results = make(map[int]bool)
	for i := rangeStart; i < (rangeEnd+itStep); i+=itStep {
		sendPacket(srcIP, destIP, i, "MTU?", appID)
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
		case goodPacket := <-mtuOKchan:
			if(goodPacket > largestSuccessfulPacket){
				//Check if the packet we received was one that we sent. (Based on size)
				if _, ok := results[goodPacket]; ok {
					largestSuccessfulPacket = goodPacket
					results[goodPacket] = true
				}
			}
		case <- timeout:
			gatherPackets = false
		}
	}

	log.Println("---------------------------------------------------")
	log.Println("Range:",rangeStart,"-",rangeEnd,"  itStep:",itStep, "  Timeout:",timeoutInSeconds)
	log.Println("Largest successful packet:", largestSuccessfulPacket)
	log.Println(results)

	return largestSuccessfulPacket
}
