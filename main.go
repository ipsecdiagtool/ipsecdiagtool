package main

import (
	//GO default packages
	"fmt"

	//Our packages
	//"hsr/ipsecdiagtool/capture"
	//"hsr/ipsecdiagtool/send"
	"github.com/ipsecdiagtool/ipsecdiagtool/capture"
)

func main() {
	fmt.Printf("Hello, IPSec.\n")

	//capture.LiveCapture()
	capture.ReadPcapFile("/home/parallels/Desktop/capture.pcap")
	//send.Run()
	fmt.Println("End")
}

/*
	##Temporary Notes:##
	+ local godocs can be compiled and accessed via: godoc -http=:6060
	+ Functions starting with big letters are public, small letters private
	+ Ints can be converted to string via. +trconv.Itoa

 */
