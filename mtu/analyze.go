package mtu

import (
	"code.google.com/p/gopacket/examples/util"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
	"time"
)

//Package internal temp. variables
var mtuOKchan = make(chan int, 100)
var conf config.Config

//Analyze determines the ideal MTU (Maximum Transmission Unit) between two nodes
//by sending increasingly big packets between them. Analyze determine the MTU
//by running FindMTU multiple times. However it is not a daemon. Once it has found
//the ideal MTU it reports it and then closes shop.
func Analyze(c config.Config) {
	defer util.Run()()
	log.Println("Analyzing MTU..")

	//Setup a channel for communication with capture
	mtuOKchan = make(chan int) // Allocate a channel.
	conf = c

	//Capture all traffic via goroutine in separate thread
	go startCapture("icmp")

	//TODO: currently required to give capture enough time to boot..
	time.Sleep(1000 * time.Millisecond)

	FastMTU(
	conf.SourceIP,
	conf.DestinationIP, 10);
}

//Listen only listens to MTU requests and replies with OK-Packets.
func Listen(c config.Config){
	defer util.Run()()

	//Setup a channel for communication with capture
	mtuOKchan = make(chan int) // Allocate a channel.
	conf = c

	go startCapture("icmp")
}

//FindMTU discovers the MTU between two nodes and returns it as an int value. FindMTU currently
//increases the packet size until it runs into a timeout. Once it runs into the timeout it returns
//the last known as good MTU.
func FindMTU(srcIP string, destIP string, startMTU int, increment int) int {
	var goodMTU, nextMTU = 0, startMTU

	//1. Initiate MTU discovery by sending first packet.
	sendPacket(srcIP, destIP, nextMTU, "MTU?")

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
			sendPacket(srcIP, destIP,nextMTU, "MTU?")
		case <-timeout:
			log.Println("Timeout has occured. We've steped over the MTU!")
			return goodMTU
		}
	}
}

//Detects the exact MTU asap.
func FastMTU(srcIP string, destIP string, timeoutInSeconds time.Duration){
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
		roughMTU = sendBatch(srcIP, destIP, rangeStart, rangeEnd, itStep, timeoutInSeconds)

		if(roughMTU == rangeEnd){
			rangeStart = rangeEnd
			rangeEnd = 2*rangeEnd
		} else if (roughMTU == 0){
			log.Println("ERROR: Reported MTU 0.. ")
		} else {
			rangeStart = roughMTU
			rangeEnd = roughMTU+itStep
			itStep = ((rangeEnd-rangeStart)/20)
		}
	}
	log.Println("Exact MTU found", roughMTU)
}

func sendBatch(srcIP string, destIP string, rangeStart int, rangeEnd int, itStep int, timeoutInSeconds time.Duration) int {
	//1. Send a batch of packets
	var results = make(map[int]bool)
	for i := rangeStart; i < (rangeEnd+itStep); i+=itStep {
		sendPacket(srcIP, destIP, i, "MTU?")
		log.Println(i)
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
				if _, ok := results[goodPacket]; ok {
					largestSuccessfulPacket = goodPacket
					results[goodPacket] = true
				} else {
					log.Println("Received a packet of a size that wasn't sent. Truncation!")
				}
			}
		case <- timeout:
			log.Println("Time's up")
			gatherPackets = false
		}
	}

	if(conf.Debug){
		log.Println("Done...")
		log.Println("Largest successful packet", largestSuccessfulPacket)
		log.Println(results)
	}

	return largestSuccessfulPacket
}

func confirmMTU(srcIP string, destIP string, mtu int, timeoutInSeconds time.Duration) bool {
	sendPacket(srcIP, destIP, mtu, "MTU?")

	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(timeoutInSeconds * time.Second)
		timeout <- true
	}()

	select {
	case <-mtuOKchan:
		return true
	case <-timeout:
		return false
	}
}
