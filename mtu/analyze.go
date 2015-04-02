package mtu

import (
	"code.google.com/p/gopacket/examples/util"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
	"strconv"
	"time"
)

//Package internal temp. variables
var mtuOKchan (chan int)
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
	go startCapture("tcp port " + strconv.Itoa(conf.Port))

	//TODO: currently required to give capture enough time to boot..
	time.Sleep(1000 * time.Millisecond)

	//MTU detection algorithm
	var mtu = c.StartingMTU
	var incStep = conf.IncrementationStep
	var passed = false
	for i := 0; i < c.MTUIterations; i++ {
		result := FindMTU(
			conf.SourceIP,
			conf.DestinationIP,
			conf.Port,
			mtu,
			incStep)

		if result != 0 {
			passed = true
			mtu = result
		} else {
			log.Println("Unable to find MTU on Try:", i)
		}

		incStep = incStep / 2
	}
	if passed {
		if confirmMTU(conf.SourceIP, conf.DestinationIP, conf.Port, mtu, config.Timeout) {
			log.Println("MTU sucessfully detected:", mtu)
		} else {
			log.Println("ERROR: MTU detected as", mtu, "but unable to confirm it.")
		}
	} else {
		log.Println("Unable to detect MTU within", c.MTUIterations, "tries.")
	}
}

//FindMTU discovers the MTU between two nodes and returns it as an int value. FindMTU currently
//increases the packet size until it runs into a timeout. Once it runs into the timeout it returns
//the last known as good MTU.
func FindMTU(srcIP string, destIP string, destPort int, startMTU int, increment int) int {
	var goodMTU, nextMTU = 0, startMTU

	//1. Initiate MTU discovery by sending first packet.
	sendPacket(srcIP, destIP, destPort, nextMTU, "MTU?")

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
			sendPacket(srcIP, destIP, destPort, nextMTU, "MTU?")
		case <-timeout:
			log.Println("Timeout has occured. We've steped over the MTU!")
			return goodMTU
		}
	}
}

func confirmMTU(srcIP string, destIP string, destPort int, mtu int, timeoutInSeconds time.Duration) bool {
	sendPacket(srcIP, destIP, destPort, mtu, "MTU?")

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
