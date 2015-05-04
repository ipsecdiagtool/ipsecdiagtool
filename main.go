package main

import (
	//GO default packages
	"fmt"
	"os"

	//Our packages
	"code.google.com/p/gopacket"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"github.com/ipsecdiagtool/ipsecdiagtool/capture"
	"github.com/ipsecdiagtool/ipsecdiagtool/logging"
	"github.com/ipsecdiagtool/ipsecdiagtool/mtu"
	"github.com/ipsecdiagtool/ipsecdiagtool/packetloss"
	//"github.com/ipsecdiagtool/ipsecdiagtool/service"
	//"flag"
	//"time"
	"github.com/kardianos/service"
	"log"
	"time"
)

var configuration config.Config
var capQuit chan bool
var icmpPackets = make(chan gopacket.Packet, 100)
var ipsecPackets = make(chan gopacket.Packet, 100)


var doStuff = true
var logger service.Logger
type program struct{}

func (p *program) Start(s service.Service) error {
	// Start should not block. Do the actual work async.
	go p.run()
	return nil
}
func (p *program) run() {
	// Do work here
	for(doStuff){
		log.Println("hi")
		time.Sleep(10*time.Second)
	}
}
func (p *program) Stop(s service.Service) error {
	// Stop should not block. Return with a few seconds.
	doStuff = false
	log.Println("done done")
	return nil
}

func main() {
	svcConfig := &service.Config{
		Name:        "ipsecdiagtool",
		DisplayName: "Go Service Test",
		Description: "This is a test Go service.",
	}
	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	logger, err = s.Logger(nil)
	if err != nil {
		log.Fatal(err)
	}
	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
/*
func main() {
	daemon, err := service.New("ipsecdiagtool", "IPSecDiag Tool Daemon")
	if err != nil {
		fmt.Println("Error: ", err)
		os.Exit(1)
	}
	flag.Parse()

	status, err := daemon.Manage()
	if err != nil {
		fmt.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	// Wait for logger output
	time.Sleep(100 * time.Millisecond)
	fmt.Println(status)
}*/


/*
func main() {
	service, err := daemon.New("IPSecDiagTool", "Detects IPSec packet loss and discovers the MTU periodically.")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	status, err := service.Install()
	if err != nil {
		log.Fatal(status, "\nError: ", err)
	}
	fmt.Println(status)

	configuration = config.LoadConfig(os.Args[0])

	if configuration.Debug {
		//Everything we need for testing belongs in here. E.g. if we're testing a new function
		//we can add it here and set the debug flag in the config to "true". Then we don't
		//need to mess with the flow of the real application.

		fmt.Println("Debug-Mode:")
		logging.InitLoger(configuration.SyslogServer, configuration.AlertCounter, configuration.AlertTime)
		go packetloss.Detectnew(configuration, ipsecPackets)
		capQuit = capture.Start(configuration, icmpPackets, ipsecPackets)
		//go mtu.FindAll(configuration, icmpPackets)

	} else {
		handleArgs()
	}

	//Keep main open forever
	//http://stackoverflow.com/questions/9543835/how-best-do-i-keep-a-long-running-go-program-running
	//might be the better solution, but for now scanln is enough.
	fmt.Println("Press any key to exit IPSecDiagTool")
	fmt.Scanln()

	if capQuit != nil {
		capQuit <- true
	}
}
*/
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
			capQuit = capture.Start(configuration, icmpPackets, ipsecPackets)
			go mtu.FindAll(configuration, icmpPackets)
		} else if os.Args[1] == "mtu-listen" {
			capQuit = capture.Start(configuration, icmpPackets, ipsecPackets)
			//TODO: doesn't reply.. --> make it reply
		} else if os.Args[1] == "packetloss" {
			logging.InitLoger(configuration.SyslogServer, configuration.AlertCounter, configuration.AlertTime)
			go packetloss.Detectnew(configuration, ipsecPackets)
			capQuit = capture.Start(configuration, icmpPackets, ipsecPackets)
		}
	} else if len(os.Args) == 1 {
		fmt.Println("Run ipsecdiagtool help to learn how to use this application.")
	}
}

