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
	go startCapture(3000, quit, captureReady, c.PcapFile)
	<-captureReady
	if c.Debug {
		log.Println("Capture Goroutine Ready")
	}
	return quit
}

//startCapture captures all packets on any interface for an unlimited duration.
//Packets can be filtered by a BPF filter string. (E.g. tcp port 22)
func startCapture(snaplen int32, quit chan bool, captureReady chan bool, pcapFile string) {
	log.Println("Waiting for MTU-Analyzer packet")
	var handle *pcap.Handle
	var err error
	if(pcapFile != ""){
		log.Println("Reading packet loss data from pcap-file:", pcapFile)
		handle, err = pcap.OpenOffline(pcapFile) //Path: /home/student/test.pcap
	} else {
		handle, err = pcap.OpenLive("any", snaplen, true, 100)
	}

	if err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		captureReady <- true

		for {
			select {
			case packet := <-packetSource.Packets():
				if packet != nil {
					//Handling packet loss
					if packet.Layer(layers.LayerTypeIPSecESP) != nil {
						putChannel(packet,ipSecChannel)
					}

					//Handling mtu detection
					if packet.Layer(layers.LayerTypeICMPv4) != nil {
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
	// Put packets in channel unless full
	case channel <- packet:
	default:
		if config.Debug {
			log.Println("Channel full, discarding packet.")
		}
	}
}
