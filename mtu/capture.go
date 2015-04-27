package mtu

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"net"
	"strconv"
	"strings"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
)

func handlePackets(icmpPackets chan gopacket.Packet, appID int, mtuOkChannels map[int]chan int) {
	for packet := range icmpPackets {
		handlePacketForChannel(packet, appID, mtuOkChannels)
	}
}

//handlePacket decides if a packet contains a valid IPSecDiagTool-MTU instruction
//and if the packet is from itself or the neighbouring node. If the packet is
//not from itself it either responds with a OK or sends an internal message
//to the findMTU goroutine that it has received an OK.
func handlePacketForChannel(packet gopacket.Packet, appID int, mtuOkChannels map[int]chan int) bool {
	s := string(packet.NetworkLayer().LayerPayload()[:])

	arr := strings.Split(s, ",")
	if len(arr) > 3 {
		remoteAppID, err1 := strconv.Atoi(arr[1])
		chanID, err2 := strconv.Atoi(arr[2])
		if (err1 == nil) && (err2 == nil) {
			//Check that packet is not from this application
			//1337 is used to disable the id check for unit-tests. It can't be generated in production use.
			if appID == remoteAppID && appID != 1337 {
				//log.Println("Packet is from us.. ignoring.")
			} else if arr[3] == "OK" {
				//log.Println("Received OK-packet with length", packet.Metadata().Length, "bytes.")
				select {
				case mtuOkChannels[chanID] <- originalSize(packet): // Put packet in channel unless full
				default:
					if(config.Debug){
						log.Println("mtuOkChannels is full or doesn't exist. Dropping OK-Information.")
					}
				}
			} else if arr[3] == "MTU?" {
				//log.Println("Received MTU?-packet with length", packet.Metadata().Length, "bytes.")
				sendOKResponse(packet, appID, chanID)
			} else {
				//log.Println("Discarded packet because neither MTU? nor OK command were included.")
			}
		} else {
			if(config.Debug){
				log.Println("ERROR: Cought a packet with an invalid app- or chan-ID. ", packet.NetworkLayer().LayerPayload())
			}
		}
	}
	return false
}

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
