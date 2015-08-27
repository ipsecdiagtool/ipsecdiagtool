package mtu

import (
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"fmt"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
	"net"
	"strconv"
	"strings"
)

//Commands
const cmdDaemonFindMTU string = "START"
const cmdMTU string = "MTU"

func handlePackets(icmpPacketsStage1 chan gopacket.Packet, icmpPacketsStage2 chan gopacket.Packet, appID int) {
	for packet := range icmpPacketsStage1 {
		s := string(packet.NetworkLayer().LayerPayload()[:])
		arr := strings.Split(s, ",")
		if len(arr) > 3 {
			packetAppID, err1 := strconv.Atoi(arr[1])
			if err1 != nil && config.Debug {
				log.Println("Bad AppID")
			}
			icmpLayer := packet.Layer(layers.LayerTypeICMPv4)
			if arr[3] == cmdMTU && icmpLayer.LayerContents()[0] == 8 && appID == packetAppID {
				//log.Println("HELLO", icmpLayer.LayerContents()[0])
				select {
				case icmpPacketsStage2 <- packet: // Put packet in channel unless full
				default:
					if config.Debug {
						log.Println("icmpPacketsStage2 is full or doesn't exist. Dropping OK-Information.")
					}
				}
			} else if arr[3] == cmdDaemonFindMTU {
				go FindAll()
			}
		} else {
			if config.Debug {
				log.Println("mtu.handlePackets: ICMP packet doesn't contain IPSecDiagTool data, dropping packet.")
			}
		}
	}
}

func distributeMtuOkPackets(icmpPacketsStage2 chan gopacket.Packet, mtuOkChannels map[int]chan int, quit chan bool) {
	for {
		select {
		case packet := <-icmpPacketsStage2:
			s := string(packet.NetworkLayer().LayerPayload()[:])
			arr := strings.Split(s, ",")
			chanID, err := strconv.Atoi(arr[2])
			if err == nil {
				select {
				case mtuOkChannels[chanID] <- originalSize(packet): // Put packet in channel unless full
				default:
					if config.Debug {
						log.Println("mtuOkChannel is full or doesn't exist. Dropping OK-Information.")
					}
				}
			}
		case <-quit:
			if config.Debug {
				fmt.Println("Received quit message, stopping Distributor.")
			}
			return
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
