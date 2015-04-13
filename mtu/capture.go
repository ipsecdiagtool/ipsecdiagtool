package mtu

import (
	//Google packages
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	"log"
	"net"
	"strconv"
	"strings"
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
			handlePacket(packet)
		}
	}
}

//handlePacket decides if a packet contains a valid IPSecDiagTool-MTU instruction
//and if the packet is from itself or the neighbouring node. If the packet is
//not from itself it either responds with a OK or sends an internal message
//to the findMTU goroutine that it has received an OK.
func handlePacket(packet gopacket.Packet) {
	s := string(packet.NetworkLayer().LayerPayload()[:])

	//Cutting off the filler material
	arr := strings.Split(s, ",")
	if len(arr) > 2 {
		remoteAppID, err := strconv.Atoi(arr[1])
		if err == nil {
			//Check that packet is not from this application
			if conf.ApplicationID == remoteAppID {
				//log.Println("Packet is from us.. ignoring.", remoteAppID)
			} else if arr[2] == "OK" {
				log.Println("Received OK-packet with length", packet.Metadata().Length, "bytes.")
				mtuOKchan <- 1
			} else if arr[2] == "MTU?" {
				log.Println("Received MTU?-packet with length", packet.Metadata().Length, "bytes.")
				sendOKResponse(packet)
			}
		} else {
			log.Println("ERROR: Cought a packet with an invalid App-ID. ", packet.NetworkLayer().LayerPayload())
		}
	}
}

//TODO: maybe throw error when packet without IP layer..
func getIP(packet gopacket.Packet) net.IP {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	// Get IP data from this layer
	ip, _ := ipLayer.(*layers.IPv4)
	return ip.SrcIP
}

func originalSize(packet gopacket.Packet) int {
	return len(packet.NetworkLayer().LayerPayload()) + len(packet.NetworkLayer().LayerContents())
}
