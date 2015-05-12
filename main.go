package main

import (
	"bufio"
	"code.google.com/p/gopacket"
	"flag"
	"fmt"
	"github.com/ipsecdiagtool/ipsecdiagtool/capture"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"github.com/ipsecdiagtool/ipsecdiagtool/logging"
	"github.com/ipsecdiagtool/ipsecdiagtool/mtu"
	"github.com/ipsecdiagtool/ipsecdiagtool/packetloss"
	"github.com/kardianos/service"
	"log"
	"os"
	"strings"
)

var configuration config.Config
var capQuit chan bool
var icmpPackets = make(chan gopacket.Packet, 100)
var ipsecPackets = make(chan gopacket.Packet, 100)
var logger service.Logger

// Program structures.
//  Define Start and Stop methods.
type program struct {
	exit chan struct{}
}

func (p *program) Start(s service.Service) error {
	p.exit = make(chan struct{})
	configuration = config.LoadConfig(os.Args[0])
	logging.InitLoger(configuration.SyslogServer, configuration.AlertCounter, configuration.AlertTime)

	if service.Interactive() {
		logger.Info("Running in terminal.")

		if configuration.Debug {
			//Code tested directly in the IDE belongs in here
			mtu.Init(configuration, icmpPackets)
			capQuit = capture.Start(configuration, icmpPackets, ipsecPackets)
			go mtu.FindAll()
		} else {
			if len(os.Args) > 1 {
				command := os.Args[1]
				switch command {
				case "install":
					installService(s)
				case "uninstall", "remove":
					uninstallService(s)
				case "interactive", "demo":
					log.Println("Interactive testing")
					go p.run()
					interactiveMode()
				case "mtu-discovery", "mtu":
					mtu.RequestDaemonMTU(configuration.ApplicationID)
				case "about":
					printAbout()
				case "debug":
					printDebug(configuration)
				case "help":
					printHelp()
				default:
					fmt.Println("Argument not reconized. Run 'ipsecdiagtool help' to learn how to use this application.")
				}
			} else {
				fmt.Println("Run 'ipsecdiagtool help' to learn how to use this application.")
			}
			os.Exit(0)
		}
	} else {
		logger.Info("Running under service manager.")
		go p.run()
	}
	return nil
}

func (p *program) run() error {
	logger.Infof("I'm running %v.", service.Platform())
	go packetloss.Detectnew(configuration, ipsecPackets)
	mtu.Init(configuration, icmpPackets)
	capQuit = capture.Start(configuration, icmpPackets, ipsecPackets)

	<-p.exit
	return nil
}

func installService(s service.Service) {
	err := s.Install()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("IPSecDiagTool Daemon successfully installed.")
	}
}

func uninstallService(s service.Service) {
	err := s.Uninstall()
	if err != nil {
		log.Println(err)
	} else {
		log.Println("IPSecDiagTool Daemon successfully uninstalled.")
	}
}

func printAbout() {
	fmt.Print("IPSecDiagTool is being developed at HSR (Hoschschule für Technik Rapperswil)\n" +
			    "as a semester/bachelor thesis. For more information please visit our repository on\n" +
			    "Github: https://github.com/IPSecDiagTool/IPSecDiagTool\n")
}

func printDebug(conf config.Config) {
	fmt.Println("IPSecDiagTool Debug Information")
	fmt.Println(conf.ToString())
}

func printHelp() {
	fmt.Println("IPSecDiagTool Help")
	fmt.Println("==================")
	fmt.Println("\n  Commands:")
	fmt.Println("   + mtu: Discover the ideal MTU between two nodes.")
	fmt.Println("   + packetloss: Passivly listen to incomming traffic and detect packet loss.")
	fmt.Println("   + install: Install this application as a service/daemon.")
	fmt.Println("   + uninstall: Uninstall this application's service/daemon.")
	fmt.Println("   + about: Learn more about IPSecDiagTool")
}

//TODO: make better
func interactiveMode() {
	reader := bufio.NewReader(os.Stdin)
	for {
		printHelp()
		fmt.Print("Enter a command: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimRight(input, "\n")
		//TODO proper error handling
		switch input {
		case "mtu":
			mtu.FindAll()
		case "packetloss":
			//TODO:
		case "about":
			printAbout()
		default:
			log.Println("Command", input, "not recognized")
		}
	}
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	logger.Info("I'm Stopping!")
	close(p.exit)
	return nil
}

// Service setup.
//   Define service config.
//   Create the service.
//   Setup the logger.
//   Handle service controls (optional).
//   Run the service.
func main() {
	svcFlag := flag.String("service", "", "Control the system service.")
	flag.Parse()

	svcConfig := &service.Config{
		Name:        "IPSecDiagTool",
		DisplayName: "A service for IPSecDiagTool",
		Description: "Detects packet loss & periodically reports the MTU for all configured tunnels.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	logger, err = s.Logger(errs)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		for {
			err := <-errs
			if err != nil {
				log.Print(err)
			}
		}
	}()

	if len(*svcFlag) != 0 {
		err := service.Control(s, *svcFlag)
		if err != nil {
			log.Printf("Valid actions: %q\n", service.ControlAction)
			log.Fatal(err)
		}
		return
	}

	err = s.Run()
	if err != nil {
		logger.Error(err)
	}
}
