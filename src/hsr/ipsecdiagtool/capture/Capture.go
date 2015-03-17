package capture

import (
	//Google packages
	"fmt"
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/pcap"
)

func LiveCapture() {
	//Capturing via gopacket-pcap on eth0
	if handle, err := pcap.OpenLive("eth0", 1600, true, 0); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter("tcp and port 80"); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			//Handling packets here
			fmt.Println(packet.Dump())
		}
	}
}

func ReadPcapFile(filepath string) {
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

			//Filtering out only IPSecESP packets.
			if len(layers) == 3  {
				if layers[2].LayerType().String() == "IPSecESP"{
					//Printing the layers each packet has
					fmt.Println(packet.String())
				}
			}
		}
	}
}
