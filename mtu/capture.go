package mtu

import (
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
func startCapture(bpfFilter string, snaplen int32, appID int) {
	log.Println("Waiting for MTU-Analyzer packet")
	if handle, err := pcap.OpenLive("any", snaplen, true, 100); err != nil {
		panic(err)
		//https://www.wireshark.org/tools/string-cf.html
	} else if err := handle.SetBPFFilter(bpfFilter); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())

		for packet := range packetSource.Packets() {
			handlePacket(packet, appID)
		}
	}
}

//handlePacket decides if a packet contains a valid IPSecDiagTool-MTU instruction
//and if the packet is from itself or the neighbouring node. If the packet is
//not from itself it either responds with a OK or sends an internal message
//to the findMTU goroutine that it has received an OK.
func handlePacket(packet gopacket.Packet, appID int) {
	s := string(packet.NetworkLayer().LayerPayload()[:])

	//Cutting off the filler material
	arr := strings.Split(s, ",")
	if len(arr) > 2 {
		remoteAppID, err := strconv.Atoi(arr[1])
		if err == nil {
			//Check that packet is not from this application
			//1337 is used to disable the id check for unit-tests. It can't be generated
			//in production use.
			if appID == remoteAppID && appID != 1337 {
				//log.Println("Packet is from us.. ignoring.")
			} else if arr[2] == "OK" {
				//log.Println("Received OK-packet with length", packet.Metadata().Length, "bytes.")
				mtuOKchan <- originalSize(packet)
			} else if arr[2] == "MTU?" {
				//log.Println("Received MTU?-packet with length", packet.Metadata().Length, "bytes.")
				sendOKResponse(packet, appID)
			} else {/*
				if(c.Debug){
					log.Println("Discarded packet because neither MTU? nor OK command were included.")
				}*/ //TODO: fix
			}
		} else {
			log.Println("ERROR: Cought a packet with an invalid App-ID. ", packet.NetworkLayer().LayerPayload())
		}
	}
}

//TODO: maybe throw error when packet without IP layer..
//Returns both the source & destination IP.
func getSrcDstIP(packet gopacket.Packet) (net.IP, net.IP) {
	ipLayer := packet.Layer(layers.LayerTypeIPv4)
	// Get IP data from this layer
	ip, _ := ipLayer.(*layers.IPv4)
	return ip.SrcIP, ip.DstIP
}

func originalSize(packet gopacket.Packet) int {
	return len(packet.NetworkLayer().LayerPayload()) + len(packet.NetworkLayer().LayerContents())
}
