package capture

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
)

var ipSecChannel chan gopacket.Packet
var icmpChannel chan gopacket.Packet

//Start creates a new goroutine that captures data from device "ANY".
//It is blocking until the capture-goroutine is ready. Start returns a quit-channel
//that can be used to gracefully shutdown it's capture-goroutine.
func Start(c config.Config, icmpPackets chan gopacket.Packet, ipSecESP chan gopacket.Packet) chan bool {
	ipSecChannel = ipSecESP
	icmpChannel = icmpPackets

	quit := make(chan bool)
	captureReady := make(chan bool)
	go startCapture(3000, quit, captureReady)
	<-captureReady
	if c.Debug {
		log.Println("Capture Goroutine Ready")
	}
	return quit
}

//startCapture captures all packets on any interface for an unlimited duration.
//Packets can be filtered by a BPF filter string. (E.g. tcp port 22)
func startCapture(snaplen int32, quit chan bool, captureReady chan bool) {
	log.Println("Waiting for MTU-Analyzer packet")
	//handle, err := pcap.OpenLive("any", snaplen, true, 100)
	handle, err := pcap.OpenOffline("/home/student/test.pcap")
	if err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		captureReady <- true

		for {
			select {
			case packet := <-packetSource.Packets():
				if packet != nil {
					//1. do packetloss stuff here:
					if packet.Layer(layers.LayerTypeIPSecESP) != nil {
						putChannel(packet,ipSecChannel)
					}

					//2. Handle ICMP packets for MTU-Detection if relevant.
					if packet.Layer(layers.LayerTypeICMPv4) != nil {
						//2.1 ICMP packets are handled by mtu package
						putChannel(packet,icmpChannel)
					}
				}
			case <-quit:
				log.Println("Received quit message, stopping Listener.")
				return
			}
		}
	}
}

func putChannel(packet gopacket.Packet, channel chan gopacket.Packet) {
	select {
	case channel <- packet: // Put packet in channel unless full
	default:
		if config.Debug {
			log.Println("Channel full, discarding packet.")
		}
	}
}
