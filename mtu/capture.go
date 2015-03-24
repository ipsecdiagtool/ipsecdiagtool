package mtu

import (
	//Google packages
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
	"log"
	"strings"
	"strconv"
	"code.google.com/p/gopacket/layers"
	"net"
)

//startCapture captures all packets on any interface for an unlimited duration.
//Packets can be filtered by a BPF filter string. (E.g. tcp port 22)
func startCapture(bpfFilter string) {
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
			handlePacket(packet)
		}
	}
}

func handlePacket(packet gopacket.Packet){
	s := string(packet.TransportLayer().LayerPayload()[:])

	//Cutting off the filler material
	arr := strings.Split(s, ",")
	if len(arr) > 1 {
		remoteAppID, err := strconv.Atoi(arr[0])
		if err != nil {
			panic(err)
		}

		//Check that packet is not from this application
		if appID == remoteAppID {
			log.Println("Packet is from us.. ignoring.")
		} else if arr[1] == "OK" {
			sendIncreasedMTU(packet)
		} else if arr[1] == "MTU?"{
			sendOKResponse(packet)
		}
	}
}

func getIP(packet gopacket.Packet) net.IP {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	// Get IP data from this layer
	ip, _ := ipLayer.(*layers.IPv4)
	return ip.SrcIP
}
