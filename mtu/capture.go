package mtu

import (
	//Google packages
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
	"log"
	"strings"
	"strconv"
)

//StartCapture captures all packets on any interface for an unlimited duration.
//Packets can be filtered by a BPF filter string. (E.g. tcp port 22)
func StartCapture(bpfFilter string) {
	log.Println("Waiting for MTU-Analyzer packet")
	if handle, err := pcap.OpenLive("any", 3000, true, 100); err != nil {
		panic(err)
		//https://www.wireshark.org/tools/string-cf.html
	} else if err := handle.SetBPFFilter(bpfFilter); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for packet := range packetSource.Packets() {
			log.Println("Received packet with length", packet.Metadata().Length, "bytes.")

			//Check if packet is not from ourselves, then handle response.
			if analyzePayload(packet) {
				composeResponse(packet)
			}
		}
	}
}

//TODO: refactor to test properly
//analyzePayload detects where the packet is a valid IPSecDiagTool
//MTU-Detection packet.
func analyzePayload(packet gopacket.Packet) bool{
	s := string(packet.TransportLayer().LayerPayload()[:])

	//Cutting off the filler material
	arr := strings.Split(s, ",")
	if len(arr) > 1{
		//TODO: add error stream
		remoteApp, _ := strconv.Atoi(arr[0])

		//Check that packet is not from this application
		if ApplicationID == remoteApp {
			log.Println("Packet is from us.. ignoring.")
			return false
		}

		log.Println("Packet comming from AppID:", remoteApp)
		log.Println("Packet contains instructions:", arr[1])

		return true
	} else {
		return false
	}
}

func composeResponse(packet gopacket.Packet) {
	//TODO: determine source automatically
	var source = "127.0.0.1"
	var destination = "127.0.0.1"
	sendPacket(source, destination, 22, 200, "OK")
}

/*
	Proposed structure for instructions in payload
	 1. AppID
	 2. task
	 3 ..
*/
