package mtu

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"net"
	"strconv"
	"strings"
)

func handlePackets(icmpPackets chan gopacket.Packet, appID int, mtuOkChannels map[int]chan int){
	for packet := range icmpPackets{
		handlePacket(packet, appID, mtuOkChannels)
	}
}

//handlePacket decides if a packet contains a valid IPSecDiagTool-MTU instruction
//and if the packet is from itself or the neighbouring node. If the packet is
//not from itself it either responds with a OK or sends an internal message
//to the findMTU goroutine that it has received an OK.
func handlePacket(packet gopacket.Packet, appID int, mtuOkChannels map[int]chan int) bool {
	s := string(packet.NetworkLayer().LayerPayload()[:])

	//Cutting off the filler material
	arr := strings.Split(s, ",")
	if len(arr) > 3 {
		//TODO: clean error handling
		remoteAppID, err := strconv.Atoi(arr[1])
		chanID, err := strconv.Atoi(arr[2])
		if err == nil {
			//Check that packet is not from this application
			//1337 is used to disable the id check for unit-tests. It can't be generated
			//in production use.
			if appID == remoteAppID && appID != 1337 {
				//log.Println("Packet is from us.. ignoring.")
			} else if arr[3] == "OK" {
				//log.Println("Received OK-packet with length", packet.Metadata().Length, "bytes.")
				//TODO: overflow if chan doesn't exist !!!
				mtuOkChannels[chanID] <- originalSize(packet)
			} else if arr[3] == "MTU?" {
				//log.Println("Received MTU?-packet with length", packet.Metadata().Length, "bytes.")
				sendOKResponse(packet, appID, chanID)
			} else {
				//log.Println("Discarded packet because neither MTU? nor OK command were included.")
			}
		} else {
			//log.Println("ERROR: Cought a packet with an invalid App-ID. ", packet.NetworkLayer().LayerPayload())
		}
	}
	return false
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
