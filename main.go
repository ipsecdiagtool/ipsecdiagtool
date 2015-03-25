package main

import (
	//GO default packages
	"fmt"

	//Our packages
	"github.com/ipsecdiagtool/ipsecdiagtool/mtu"
	"github.com/ipsecdiagtool/ipsecdiagtool/packetloss"
)

func main() {
	fmt.Printf("Hello, IPSec.\n")

	//capture.LiveCapture("")
	//capture.ReadPcapFile("/home/parallels/Desktop/capture.pcap")

	go packetloss.Detect()

	go mtu.Analyze()

	//Keep main open forever
	//http://stackoverflow.com/questions/9543835/how-best-do-i-keep-a-long-running-go-program-running
	//might be the better solution, but for now scanln is enough.
	fmt.Println("Press any key to exit IPSecDiagTool")
	fmt.Scanln()
}
