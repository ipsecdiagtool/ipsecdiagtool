package mtu

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
	"net"
	"strconv"
	"strings"
)

func handlePackets(icmpPacketsStage1 chan gopacket.Packet, icmpPacketsStage2 chan gopacket.Packet, appID int) {
	for packet := range icmpPacketsStage1 {
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
					select {
					case icmpPacketsStage2 <- packet: // Put packet in channel unless full
					default:
						if config.Debug {
							log.Println("icmpPacketsStage2 is full or doesn't exist. Dropping OK-Information.")
						}
					}
				} else if arr[3] == "MTU?" {
					//log.Println("Received MTU?-packet with length", packet.Metadata().Length, "bytes.")
					sendOKResponse(packet, appID, chanID)
				} else {
					//log.Println("Discarded packet because neither MTU? nor OK command were included.")
				}
			} else {
				if config.Debug {
					log.Println("ERROR: Cought a packet with an invalid app- or chan-ID. ", packet.NetworkLayer().LayerPayload())
				}
			}
		}
	}
}

func distributeMtuOkPackets(icmpPacketsStage2 chan gopacket.Packet, mtuOkChannels map[int]chan int) {
	for packet := range icmpPacketsStage2 {
		s := string(packet.NetworkLayer().LayerPayload()[:])
		arr := strings.Split(s, ",")
		chanID, err := strconv.Atoi(arr[2])
		if err == nil {
			select {
			case mtuOkChannels[chanID] <- originalSize(packet): // Put packet in channel unless full
			default:
				if config.Debug {
					log.Println("mtuOkChannels is full or doesn't exist. Dropping OK-Information.")
				}
			}
		}
	}
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
