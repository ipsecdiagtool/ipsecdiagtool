package packetloss

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket/pcap"
	"fmt"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
	"net"
)

var espmap *EspMap

func Detect(configuration config.Config) error {

	espmap = NewEspMap(configuration.WindowSize)

	fmt.Println("detecting packetloss ...")

	if configuration.InterfaceName != "any" {
		iface, err := net.InterfaceByName(configuration.InterfaceName)
		if err != nil {
			return err
		}
		fmt.Println(iface)
	}

	//handle, err := pcap.OpenLive("any", 1500, true, 100)
	//handle, err := pcap.OpenOffline("/home/student/TestIpSec_Ping.pcap")
	handle, err := pcap.OpenOffline("/home/student/test.pcap")
	if err != nil {
		return err
	}

	readIPSec(handle)
	return nil
}

func readIPSec(handle *pcap.Handle) {
	src := gopacket.NewPacketSource(handle, handle.LinkType())
	var counter uint32
	for packet := range src.Packets() {
		ipsecLayer := packet.Layer(layers.LayerTypeIPSecESP)
		if ipsecLayer != nil {
			netFlow := packet.NetworkLayer().NetworkFlow()
			src, dst := netFlow.Endpoints()
			ipsec := ipsecLayer.(*layers.IPSecESP)
			log.Println("Source: ", src, "Destination: ", dst, "Seqnumber: ", ipsec.Seq)

			espmap.MakeEntry(Connection{src.String(), dst.String(), ipsec.SPI}, ipsec.Seq)

			//espmap.CheckForLost()

		}
		counter++
	}
	for _, element := range espmap.elements {
		fmt.Println(element.lostpackets)
		fmt.Println(element.maybelostpackets)

	}
	fmt.Println("Packets: ", counter)

}
