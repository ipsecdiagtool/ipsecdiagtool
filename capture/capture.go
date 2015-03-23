package capture

import (
	//Google packages
	"fmt"
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
	"log"
)

//LiveCapture captures all tcp & port 80 packets on eth0.
func LiveCapture() {
	log.Println("Waiting for packet")
	if handle, err := pcap.OpenLive("any", 1600, true, 100); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			//Handling packets here
			log.Println(packet.Dump())
		}
	}
}

//ReadPcapFile iterates over all packets in a .pcap-file and counts them.
//Returns the total number  of packets and outputs the layers of all IPSecESP-Type packets.
func ReadPcapFile(filepath string) int{
	var counter = 0
	if handle, err := pcap.OpenOffline(filepath); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			//Handling packets here
			//fmt.Println(packet.Dump())
			//fmt.Println(packet.String())
			var layers []gopacket.Layer
			layers = append(layers, packet.Layers()...) //Three dots to signify that we're combing two a slices.

			counter++

			//Filtering out only IPSecESP packets.
			if len(layers) == 3  {
				if layers[2].LayerType().String() == "IPSecESP"{
					//Printing the layers each packet has
					fmt.Println(packet.String())
				}
			}
		}
	}
	return counter
}
