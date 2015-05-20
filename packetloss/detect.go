package packetloss

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"log"
)

var espmap *EspMap

//Detect starts the packet loss detection processs.
func Detect(c config.Config, ipSecESP chan gopacket.Packet, output bool) {
	espmap = NewEspMap(c.WindowSize)
	log.Println("Packet loss detection started..")

	for packet := range ipSecESP {
		ipsecLayer := packet.Layer(layers.LayerTypeIPSecESP)
		if ipsecLayer != nil {
			netFlow := packet.NetworkLayer().NetworkFlow()
			src, dst := netFlow.Endpoints()
			ipsec := ipsecLayer.(*layers.IPSecESP)

			if output {
				log.Println("Source: ", src, "Destination: ", dst, "Seqnumber: ", ipsec.Seq)
			}
			espmap.MakeEntry(Connection{src.String(), dst.String(), ipsec.SPI}, ipsec.Seq)
		}
	}
}
