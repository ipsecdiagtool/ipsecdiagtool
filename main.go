package main

import (
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
	"time"
)

var configuration config.Config
var capQuit chan bool
var icmpPackets = make(chan gopacket.Packet, 100)
var ipsecPackets = make(chan gopacket.Packet, 100)

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
		if configuration.Debug {
			log.Println("Running in terminal.")
			//Code tested directly in the IDE belongs in here
			mtu.Init(configuration, icmpPackets)
			capQuit = capture.Start(configuration, icmpPackets, ipsecPackets)
			log.Println("Waiting 2seconds for partner. Can be disabled if partner is already running.")
			time.Sleep(2 * time.Second)
			go mtu.FindAll()

		} else {
			if len(os.Args) > 1 {
				command := os.Args[1]
				switch command {
				case "install":
					installService(s)
				case "uninstall", "remove":
					uninstallService(s)
				case "interactive", "i":
					log.Println("Interactive testing")
					go p.run()
					if len(os.Args) > 2 {
						handleInteractiveArg(os.Args[2])
					} else {
						fmt.Println("Please specify an additional argument when using 'ipsecdiagtool " + command + "'")
						fmt.Println("Use 'ipsecdiagtool help' to get additional information.")
					}
					os.Exit(0)
				case "mtu-discovery", "mtu":
					mtu.RequestDaemonMTU(configuration.ApplicationID)
					os.Exit(0) //TODO: remove os.Exit find better solution.
				case "about", "a":
					printAbout()
					os.Exit(0)
				case "debug", "d":
					printDebug(configuration)
					os.Exit(0)
				case "help", "--help", "h":
					printHelp()
					os.Exit(0)
				default:
					fmt.Println("Argument not reconized. Run 'ipsecdiagtool help' to learn how to use this application.")
					os.Exit(0)
				}
			} else {
				fmt.Println("Run 'ipsecdiagtool help' to learn how to use this application.")
				os.Exit(0)
			}
		}
	} else {
		if config.Debug {
			log.Println("Running under service manager.")
		}
		go packetloss.Detectnew(configuration, ipsecPackets, false)
		go p.run()
	}
	return nil
}

func handleInteractiveArg(arg string) {
	switch arg {
	case "mtu", "m":
		go mtu.FindAll()
	case "packetloss", "p", "pl":
		if len(os.Args) > 3 {
			pcapPath := os.Args[3]
			configuration.PcapFile = pcapPath
			log.Println("Reading packetloss test data from file.")
			go packetloss.Detectnew(configuration, ipsecPackets, true)
		}
		if configuration.PcapFile == "" {
			log.Println("Detecting packetloss from ethernet")
		} else {
			log.Println("Detecting packetloss from configured file:", configuration.PcapFile)
		}
		go packetloss.Detectnew(configuration, ipsecPackets, true)
	default:
		fmt.Println("Command", arg, "not recognized")
		fmt.Println("See 'ipsecdiagtool help' for additional information.")
	}
}

func (p *program) run() error {
	if configuration.Debug {
		log.Println("Running Daemon via", service.Platform())
	}
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
		"as a semester/bachelor thesis. For more information please visit our repository on\n \n" +
		"Authors: Jan Balmer, Theo Winter \n" +
		"Github: https://github.com/IPSecDiagTool/IPSecDiagTool\n")
}

func printDebug(conf config.Config) {
	fmt.Println("IPSecDiagTool Debug Information")
	fmt.Println(conf.ToString())
}

func printHelp() {
	var spac = "\n   "
	var help = "Usage: ipsecdiagtool [OPTION ...] \n" +
		"IPSecDiagTool detects packet loss for all connected IPSec VPN tunnels. It can also discover the MTU for all configured connections. " +
		"IPSecDiagTool is intended to be run as a daemon on both sides of a VPN tunnel." + spac +
		"\n" +
		"Daemon operation mode:" + spac +
		"ipsecdiagtool install    #Install the daemon/service on your system." + spac +
		"ipsecdiagtool uninstall  #Uninstall the daemon/service from system." + spac +
		"ipsecdiagtool mtu        #Tell locally running daemon to start discoverying the MTU." + spac +
		"\n" +
		"Interactive opertation mode:" + spac +
		"ipsecdiagtool i mtu        #Run the mtu discovery locally, without a daemon." + spac +
		"ipsecdiagtool i pl         #Run the packetloss detection locally, without a daemon." + spac +
		"ipsecdiagtool i pl [path]  #Run the packetloss detection locally, reading pcap data from a file." + spac +
		"\n" +
		"Information commands:" + spac +
		"ipsecdiagtool debug        #Show debug information." + spac +
		"ipsecdiagtool help         #Display this help menu." + spac +
		"ipsecdiagtool about        #Who created this application."
	fmt.Println(help)
}

func (p *program) Stop(s service.Service) error {
	// Any work in Stop should be quick, usually a few seconds at most.
	log.Println("Stopping IPSecDiagTool")
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
		Name:        "ipsecdiagtool",
		DisplayName: "A service for IPSecDiagTool",
		Description: "Detects packet loss & periodically reports the MTU for all configured tunnels.",
	}

	prg := &program{}
	s, err := service.New(prg, svcConfig)
	if err != nil {
		log.Fatal(err)
	}
	errs := make(chan error, 5)
	//logger, err = s.Logger(errs)
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
		//TODO: check if we need a system logger rather then the panics we've been using so far.
		//logger.Error(err)
		panic(err)
	}
}
