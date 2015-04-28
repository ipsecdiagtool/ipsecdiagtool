package capture

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
	"log"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"code.google.com/p/gopacket/layers"
)

//Start creates a new goroutine that captures data from device "ANY".
//It is blocking until the capture-goroutine is ready. Start returns a quit-channel
//that can be used to gracefully shutdown it's capture-goroutine.
func Start(c config.Config, icmpPackets chan gopacket.Packet) chan bool{
	quit := make(chan bool)
	captureReady := make(chan bool)
	go startCapture(3000, icmpPackets, quit, captureReady)
	<- captureReady
	if(c.Debug){
		log.Println("Capture Goroutine Ready")
	}
	return quit
}

//startCapture captures all packets on any interface for an unlimited duration.
//Packets can be filtered by a BPF filter string. (E.g. tcp port 22)
func startCapture(snaplen int32, icmpPackets chan gopacket.Packet, quit chan bool, captureReady chan bool) {
	log.Println("Waiting for MTU-Analyzer packet")
	handle, err := pcap.OpenLive("any", snaplen, true, 100)
	if err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		captureReady <- true
		for {
			select {
			case packet := <-packetSource.Packets():
				//1. do packetloss stuff here:
				//probably best to send it relevant packets to separate packetloss goroutine, similar to how I've done it below..

				//2. Handle ICMP packets for MTU-Detection if relevant.
				if packet.Layer(layers.LayerTypeICMPv4) != nil{

					//2.1 ICMP packets are handled by mtu package
					select {
					case icmpPackets <- packet: // Put packet in channel unless full
					default:
						//TODO: only log when debug-mode is enabled.
						log.Println("Channel full, discarding ICMP packet.")
					}
				}
			case <-quit:
				log.Println("Received quit message, stopping Listener.")
				return
			}
		}
	}
}
