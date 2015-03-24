package mtu

import (
	"code.google.com/p/gopacket/examples/util"
	"log"
	"time"
	"strconv"
	"math/rand"
)

//Package setup
var appID int
var srcIP string
var destIP string
var destPort int
var incStep int

//Package internal temp. variables
var currentMTU int

//Setup configures the MTU-daemon with the necessary information to
//determine the MTU between two nodes. At some point it will likely get
//it's information automatically from a central config reader within the
//application. But it's still useful if you want to use the MTU-detection
//directly in your application.
//If you set the application ID 0, a random one will be automatically generated.
func Setup(applicationID int, sourceIP string, destinationIP string, destinationPort int, incrementationStep int) {
	if(applicationID == 0){
		rand.Seed(time.Now().UnixNano()) //Seed is required otherwise we always get the same number
		appID = rand.Intn(100000)
	} else {
		appID = applicationID
	}
	srcIP = sourceIP
	destIP = destinationIP
	destPort = destinationPort
	incStep = incrementationStep
	currentMTU = 500 //Starting MTU
}

//Analyze determines the ideal MTU (Maximum Transmission Unit) between two nodes
//by sending increasingly big packets between them. Analyze determine the MTU
//exactly once and return the value of the ideal MTU. To continuously determine
//the MTU you should run [not implemented yet].
func Analyze() {
	defer util.Run()()
	setDefaultValues()
	log.Println("Analyzing MTU..")

	//Capture all traffic via goroutine in separate thread
	go startCapture("tcp port " + strconv.Itoa(destPort))

	//Fire first packet to determine MTU. Later this should be done at
	//certain times or via outside input in form of a cronjob.
	time.Sleep(1000 * time.Millisecond)
	go sendPacket(srcIP, destIP, destPort, currentMTU, "MTU?")

	//TODO:
	//-Record packet
	//-Loop several times to find ideal MTU
}

//setDefaultValues is run when the user doesn't configure the MTU package via Setup().
func setDefaultValues() {
	if srcIP == "" {
		Setup(0, "127.0.0.1","127.0.0.1",22, 100)
		log.Println("Setting default values, because Analyze() was called before or without Setup()")
	}
}
