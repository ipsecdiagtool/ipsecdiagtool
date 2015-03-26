package packetloss

import (
        //"bytes"
        "fmt"
        "log"
        "net"

        "code.google.com/p/gopacket"
        "code.google.com/p/gopacket/layers"
        "code.google.com/p/gopacket/pcap"
)

//Detect, for Jan.
func Detect() {
	fmt.Println("detecting packetloss ...")
	
	//net.Interfaces()
	iface, err := net.InterfaceByName("any")
	if err != nil {
	// error handling
	}
	
	handle, err := pcap.OpenLive("any", 1500, true, 100)
	if err != nil {
		panic(err)
	}
	
	stop := make(chan struct{})
    go readIPSec(handle, iface, stop)
    defer close(stop)
    for{fmt.Scanln()}
}

func readIPSec(handle *pcap.Handle, iface *net.Interface, stop chan struct{}) {
		fmt.Println("testpunkt1")
        src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
        fmt.Println("testpunkt2")
        in := src.Packets()
        fmt.Println("testpunkt3")
        for {
                var packet gopacket.Packet
                select {
                case <-stop:
                		fmt.Print("*")
                        return
                case packet = <-in:
                		
                		
                		
                		test:=packet.ApplicationLayer()
                		if test != nil{
                			fmt.Println(test.Payload())
                		}
                		
                        arpLayer := packet.Layer(layers.LayerTypeIPSecESP)
                        
                        if arpLayer == nil {
                                continue
                        }
                        arp := arpLayer.(*layers.IPSecESP)
                        //if arp.Operation != layers.ARPReply || bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
                                // This is a packet I sent.
                        //        continue
                        //}
                        // Note:  we might get some packets here that aren't responses to ones we've sent,
                        // if for example someone else sends US an ARP request.  Doesn't much matter, though...
                        // all information is good information :)
                        //log.Printf("IP %v", net.IP(arp.DstIP))
                        log.Println(arp.Seq)
                }
        }
}


//LiveCapture captures all packets on any interface for an unlimited duration.
//Packets can be filtered by a BPF filter string. (E.g. tcp port 22)
/*
func LiveCapture(bpfFilter string) {
	log.Println("Waiting for packet")
	if handle, err := pcap.OpenLive("any", 1500, true, 100); err != nil {
		panic(err)
	} else if err := handle.SetBPFFilter(bpfFilter); err != nil {
		panic(err)
	} else {
		packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
		for packet := range packetSource.Packets() {
			var packet gopacket.Packet
			

		}
	}
}
*/



//ReadPcapFile iterates over all packets in a .pcap-file and counts them.
//Returns the total number  of packets and outputs the layers of all IPSecESP-Type packets.
func ReadPcapFile(filepath string) int {
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
			if len(layers) == 3 {
				if layers[2].LayerType().String() == "IPSecESP" {
					//Printing the layers each packet has
					fmt.Println(packet.String())
				}
			}
		}
	}
	return counter
}
