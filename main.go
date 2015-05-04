package main

import (
	//GO default packages
	"fmt"
	"os"

	//Our packages
	"code.google.com/p/gopacket"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	/*"github.com/ipsecdiagtool/ipsecdiagtool/capture"
	"github.com/ipsecdiagtool/ipsecdiagtool/logging"
	"github.com/ipsecdiagtool/ipsecdiagtool/mtu"
	"github.com/ipsecdiagtool/ipsecdiagtool/packetloss"*/
	"github.com/takama/daemon"
	"log"
	"os/signal"
	"syscall"
	"net"
)

var configuration config.Config
var capQuit chan bool
var icmpPackets = make(chan gopacket.Packet, 100)
var ipsecPackets = make(chan gopacket.Packet, 100)

const (
	// name of the service, match with executable file name
	name        = "IPSecDiagTool"
	description = "Detects packet loss and periodically reports the MTU for all configured tunnels."

	// port which daemon should be listen
	port = ":9978"
)

var stdlog, errlog *log.Logger

// Service has embedded daemon
type Service struct {
	daemon.Daemon
}

// Manage by daemon commands or run the daemon
func (service *Service) Manage() (string, error) {

	usage := "Usage: myservice install | remove | start | stop | status"

	// if received any kind of command, do it
	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		case "hello":
			return "Hello there", nil
		default:
			return usage, nil
		}
	}

	// Do something, call your goroutines, etc

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	// Set up listener for defined host and port
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return "Possibly was a problem with the port binding", err
	}

	// set up channel on which to send accepted connections
	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	// loop work cycle with accept connections or interrupt
	// by system signal
	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			stdlog.Println("Got signal:", killSignal)
			stdlog.Println("Stoping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				return "Daemon was interruped by system signal", nil
			}
			return "Daemon was killed", nil
		}
	}

	// never happen, but need to complete code
	return usage, nil
}

// Accept a client connection and collect it in a channel
func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		client.Write(buf)
	}
}

func init() {
	stdlog = log.New(os.Stdout, "", log.Ldate|log.Ltime)
	errlog = log.New(os.Stderr, "", log.Ldate|log.Ltime)
}

func main() {
	srv, err := daemon.New(name, description)
	if err != nil {
		errlog.Println("Error: ", err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		errlog.Println(status, "\nError: ", err)
		os.Exit(1)
	}
	fmt.Println(status)
}
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
*/
