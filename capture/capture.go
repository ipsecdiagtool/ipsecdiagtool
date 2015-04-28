package capture

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
	"log"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"code.google.com/p/gopacket/layers"
)

var quit chan bool

//Start creates a new goroutine that captures data from device "ANY".
//This function blocks until the capturing-goroutine is ready.
func Start(c config.Config, icmpPackets chan gopacket.Packet){
	quit := make(chan bool)
	captureReady := make(chan bool)
	go startCapture("icmp", 3000, icmpPackets, quit, captureReady)
	<- captureReady
	if(c.Debug){
		log.Println("Capture Goroutine Ready")
	}
}

//Sends a quit-Message to the capturing-goroutine to gracefully shutdown.
func Stop() {
	//TODO: check if nil
	quit <- true
}

//startCapture captures all packets on any interface for an unlimited duration.
//Packets can be filtered by a BPF filter string. (E.g. tcp port 22)
func startCapture(bpfFilter string, snaplen int32, icmpPackets chan gopacket.Packet, quit chan bool, captureReady chan bool) {
	log.Println("Waiting for MTU-Analyzer packet")
	handle, err := pcap.OpenLive("any", snaplen, true, 100)
	if err != nil {
		panic(err)
		//https://www.wireshark.org/tools/string-cf.html
	} else if err := handle.SetBPFFilter(bpfFilter); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		captureReady <- true
		for {
			select {
			case packet := <-packetSource.Packets():
				//1. do packetloss stuff here:

				//2. Handle ICMP packets for MTU-Detection if relevant.
				if packet.Layer(layers.LayerTypeICMPv4) != nil{
					//TODO: make overflow
					//2.1 ICMP packets are handled by mtu package
					icmpPackets <- packet
				}
			case <-quit:
				log.Println("Received quit message, stopping Listener.")
				return
			}
		}
	}
}
