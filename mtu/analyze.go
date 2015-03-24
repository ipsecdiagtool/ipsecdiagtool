package mtu

import (
	"code.google.com/p/gopacket/examples/util"
	"log"
	"time"
	"strconv"
	"math/rand"
)

var appID int
var srcIP string
var destIP string
var destPort int

//Setup configures the MTU-daemon with the necessary information to
//determine the MTU between two nodes. At some point it will likely get
//it's information automatically from a central config reader within the
//application. But it's still useful if you want to use the MTU-detection
//directly in your application.
//If you set the application ID 0, a random one will be automatically generated.
func Setup(applicationID int, sourceIP string, destinationIP string, destinationPort int) {
	if(applicationID == 0){
		rand.Seed(time.Now().UnixNano()) //Seed is required otherwise we always get the same number
		appID = rand.Intn(100000)
	} else {
		appID = applicationID
	}
	srcIP = sourceIP
	destIP = destinationIP
	destPort = destinationPort
}

//Analyze computes the ideal MTU for a connection between two computers.
func Analyze() {
	defer util.Run()()
	setDefaultValues()
	log.Println("Analyzing MTU..")

	//Capture all traffic via goroutine in separate thread
	go StartCapture("tcp port " + strconv.Itoa(destPort))

	//TODO: remove later
	//Temporary delay to wait until the filter is properly setup.
	time.Sleep(1000 * time.Millisecond)

	//Send packet via goroutine in separate thread
	go sendPacket(srcIP, destIP, destPort, 120, "MTU?")


	//TODO:
	//-Record packet
	//-Loop several times to find ideal MTU
}

func setDefaultValues() {
	if srcIP == "" {
		Setup(0, "127.0.0.1","127.0.0.1",22)
		log.Println("Setting default values, because Analyze() was called before or without Setup()")
	}
}
