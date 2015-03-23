package mtu

import (
	//Google packages
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
	"log"
	"strings"
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
			//Handling packets here
			//log.Println(packet.Dump())
			analyzePayload(packet)
			log.Println("Received packet with size", packet.Metadata().Length)
		}
	}
}

//analyzePayload detects where the packet is a valid IPSecDiagTool
//MTU-Detection packet.
func analyzePayload(packet gopacket.Packet) bool{
	s := string(packet.TransportLayer().LayerPayload()[:])

	//Cutting off the filler material
	arr := strings.Split(s, "/END")
	log.Println("Packet contains instructions:",arr[0])
	if len(arr) > 1{
		return true
	} else {
		return false
	}
}
