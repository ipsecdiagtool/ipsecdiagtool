package capture

import (
	//Google packages
	"fmt"
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
	"log"
)

//TODO: Code here is no longer required for MTU-analyzer. Code is left as an example
//TODO: for now. But should be removed in the future.

//LiveCapture captures all packets on any interface for an unlimited duration.
//Packets can be filtered by a BPF filter string. (E.g. tcp port 22)
func LiveCapture(bpfFilter string) {
	log.Println("Waiting for packet")
	if handle, err := pcap.OpenLive("any", 1500, true, 100); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter(bpfFilter); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			//Handling packets here
			log.Println(packet.Dump())
			log.Println("Received packet with size", packet.Metadata().Length)
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
