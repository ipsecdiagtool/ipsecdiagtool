package packetloss

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	"fmt"
	"log"
	"net"
)

var espmap *EspMap

func Detect(windowsize int) {

	espmap = NewEspMap(windowsize)

	fmt.Println("detecting packetloss ...")

	iface, err := net.InterfaceByName("any")
	if err != nil {
		// error handling
	}

	//handle, err := pcap.OpenLive("any", 1500, true, 100)
	handle, err := pcap.OpenOffline("/home/student/TestIpSec_Ping.pcap")
	if err != nil {
		panic(err)
	}

	stop := make(chan struct{})
	go readIPSec(handle, iface, stop)
	defer close(stop)
	for {
		fmt.Scanln()
	}
}

func readIPSec(handle *pcap.Handle, iface *net.Interface, stop chan struct{}) {
	src := gopacket.NewPacketSource(handle, handle.LinkType())

	for packet := range src.Packets() {
		ipsecLayer := packet.Layer(layers.LayerTypeIPSecESP)
		if ipsecLayer != nil {
			netFlow := packet.NetworkLayer().NetworkFlow()
			src, dst := netFlow.Endpoints()
			ipsec := ipsecLayer.(*layers.IPSecESP)
			log.Println("Source: ", src, "Destination: ", dst, "Seqnumber: ", ipsec.Seq)

			espmap.MakeEntry(Connection{src.String(), dst.String(), ipsec.SPI}, ipsec.Seq)
		}
	}

	for k, element := range espmap.elements {
		fmt.Println(k)
		for _, seqnumber := range element {
			fmt.Println(seqnumber)
		}
	}
}
