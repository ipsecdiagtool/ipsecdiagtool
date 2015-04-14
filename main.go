package main

import (
	//GO default packages
	"fmt"
	"os"

	//Our packages
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"github.com/ipsecdiagtool/ipsecdiagtool/mtu"
	"github.com/ipsecdiagtool/ipsecdiagtool/packetloss"
)

var configuration config.Config

func main() {
	configuration = config.LoadConfig()

	if configuration.Debug {
		//Everything we need for testing belongs in here. E.g. if we're testing a new function
		//we can add it here and set the debug flag in the config to "true". Then we don't
		//need to mess with the flow of the real application.
		//go packetloss.Detect(configuration.WindowSize)
		fmt.Println("Debug-Mode:")
		go mtu.Analyze(configuration)
	} else {
		handleArgs()
	}

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
			go mtu.Analyze(configuration)
		} else if os.Args[1] == "mtu-listen" {
			go mtu.Listen(configuration, 3000) //TODO: maybe increase
		} else if os.Args[1] == "packetloss" {
			go packetloss.Detect(configuration)
		}
	} else if len(os.Args) == 1 {
		fmt.Println("Run ipsecdiagtool help to learn how to use this application.")
	}
}
