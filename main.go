package main

import (
	//GO default packages
	"fmt"
	"os"

	//Our packages
	"code.google.com/p/gopacket"
	"github.com/ipsecdiagtool/ipsecdiagtool/capture"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"github.com/ipsecdiagtool/ipsecdiagtool/mtu"
	"github.com/ipsecdiagtool/ipsecdiagtool/packetloss"
)

var configuration config.Config
var capQuit chan bool

func main() {
	configuration = config.LoadConfig()

	if configuration.Debug {
		fmt.Println("Debug-Mode:")
		//Everything we need for testing belongs in here. E.g. if we're testing a new function
		//we can add it here and set the debug flag in the config to "true". Then we don't
		//need to mess with the flow of the real application.

		icmpPackets := make(chan gopacket.Packet, 100)
		capQuit = capture.Start(configuration, icmpPackets)
		go mtu.FindAll(configuration, icmpPackets)

		/*
			packetloss.InitLoger(configuration.SyslogServer, configuration.AlertCounter, configuration.AlertTime)
			go packetloss.Detect(configuration)
			packetloss.InfoLog("Dies ist eine kurze Info")
			packetloss.AlertLog("Dies ist ein Alert")
		*/
	} else {
		handleArgs()
	}

	//Keep main open forever
	//http://stackoverflow.com/questions/9543835/how-best-do-i-keep-a-long-running-go-program-running
	//might be the better solution, but for now scanln is enough.
	fmt.Println("Press any key to exit IPSecDiagTool")
	fmt.Scanln()
	capQuit <- true
}

//Handle commandline arguments. Arg0 = path where program is running,
//Arg1+ raw arguments.
func handleArgs() {
	if len(os.Args) > 1 {
		if os.Args[1] == "about" {
			fmt.Println("IPSecDiagTool is being developed at HSR (Hoschschule für Technik Rapperswil)\n" +
				"as a semester/bachelor thesis. For more information please visit our repository on\n" +
				"Github: https://github.com/IPSecDiagTool/IPSecDiagTool\n")
		} else if os.Args[1] == "help" {
			fmt.Println("IPSecDiagTool Help")
			fmt.Println("==================")
			fmt.Println("\n  Commands:")
			fmt.Println("   + mtu: Discover the ideal MTU between two nodes.")
			fmt.Println("   + packetloss: Passivly listen to incomming traffic and detect packet loss.")
			fmt.Println("   + about: Learn more about IPSecDiagTool")
		} else if os.Args[1] == "mtu" {
			icmpPackets := make(chan gopacket.Packet, 100)
			capQuit = capture.Start(configuration, icmpPackets)
			go mtu.FindAll(configuration, icmpPackets)
		} else if os.Args[1] == "mtu-listen" {
			icmpPackets := make(chan gopacket.Packet, 100)
			capQuit = capture.Start(configuration, icmpPackets)
			//TODO: doesn't reply.. --> make it reply
		} else if os.Args[1] == "packetloss" {
			go packetloss.Detect(configuration)
		}
	} else if len(os.Args) == 1 {
		fmt.Println("Run ipsecdiagtool help to learn how to use this application.")
	}
}
