package main

import (
	//GO default packages
	"fmt"

	//Our packages
	"hsr/ipsecdiagtool/capture"
)

func main() {
	fmt.Printf("Hello, IPSec.\n")

	//capture.LiveCapture()
	capture.ReadPcapFile("/home/parallels/Desktop/capture.pcap")
}
