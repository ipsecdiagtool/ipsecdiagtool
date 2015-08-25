package capture

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"errors"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
	"time"
)

var ipSecChannel chan gopacket.Packet
var icmpChannel chan gopacket.Packet

//Start creates a new goroutine that captures data from device "ANY".
//It is blocking until the capture-goroutine is ready. Start returns a quit-channel
//that can be used to gracefully shutdown it's capture-goroutine.
func Start(c config.Config, icmpPackets chan gopacket.Packet, ipsecESP chan gopacket.Packet) chan bool {
	initChannels(icmpPackets, ipsecESP)
	quit := make(chan bool)
	captureReady := make(chan bool)
	go capture(c.PcapSnapLen, quit, captureReady, c.PcapFile)
	<-captureReady
	if c.Debug {
		log.Println("Capture Goroutine Ready.")
	}
	return quit
}

//initChannels is needed to initialize this package in the tests
func initChannels(icmpPackets chan gopacket.Packet, ipsecESP chan gopacket.Packet) {
	ipSecChannel = ipsecESP
	icmpChannel = icmpPackets
}

//startCapture captures all packets on any interface for an unlimited duration.
//Packets can be filtered by a BPF filter string. (E.g. tcp port 22)
func capture(snaplen int32, quit chan bool, captureReady chan bool, pcapFile string) error {
	var handle *pcap.Handle
	var err error
	if pcapFile != "" {
		log.Println("Reading packet loss data from pcap-file:", pcapFile)
		handle, err = pcap.OpenOffline(pcapFile)
	} else {
		//https://godoc.org/code.google.com/p/gopacket/pcap
		//This might have been the culprit
		//Alternative to try: 250*time.Millisecond
		handle, err = pcap.OpenLive("any", snaplen, true, 250*time.Millisecond)
	}

	if err != nil {
		log.Println("Error while start capturing packets", err)
		return err
	}

	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	captureReady <- true

	for {
		select {
		case packet := <-packetSource.Packets():
			if packet != nil {
				if packet.Layer(layers.LayerTypeIPSecESP) != nil {
					putChannel(packet, ipSecChannel)
				}
				if packet.Layer(layers.LayerTypeICMPv4) != nil {
					putChannel(packet, icmpChannel)
				}
			}
		case <-quit:
			log.Println("Received quit message, stopping Listener.")
			return nil
		}
	}
}

func putChannel(packet gopacket.Packet, channel chan gopacket.Packet) error {
	select {
	// Put packets in channel unless full
	case channel <- packet:
	default:
		msg := "Channel full, discarding packet."
		if config.Debug {
			log.Println(msg)
		}
		return errors.New(msg)
	}
	return nil
}
