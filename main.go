package main

import (
	//GO default packages
	"fmt"
	"os"

	//Our packages
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"github.com/ipsecdiagtool/ipsecdiagtool/mtu"
	//"github.com/ipsecdiagtool/ipsecdiagtool/packetloss"
)

var Configuration config.Config

func main() {
	fmt.Printf("Hello, IPSec.\n")

	handleArgs()

	//capture.LiveCapture("")
	//capture.ReadPcapFile("/home/parallels/Desktop/capture.pcap")

	//go packetloss.Detect(512)

	Configuration = config.LoadConfig()

	go mtu.Analyze(Configuration)

	//Keep main open forever
	//http://stackoverflow.com/questions/9543835/how-best-do-i-keep-a-long-running-go-program-running
	//might be the better solution, but for now scanln is enough.
	fmt.Println("Press any key to exit IPSecDiagTool")
	fmt.Scanln()
}

//Handle commandline arguments. Arg0 = path where program is running,
//Arg1+ raw arguments.
func handleArgs() {
	if len(os.Args) > 1 {
		if os.Args[1] == "about" {
			fmt.Println("IPSecDiagTool is being developed at HSR (Hoschschule für Technik Rapperswil)" +
				"\n as a semester/bachelor thesis. For more information please visit our repository on" +
				"\n Github: https://github.com/IPSecDiagTool/IPSecDiagTool")
		} else if os.Args[1] == "help" {
			//TODO: help infos
			fmt.Println("TODO: Help.")
		}
	}
}
