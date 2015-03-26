package packetloss

import (
		"strings"
        "fmt"
        "log"
        "net"
        "code.google.com/p/gopacket"
        "code.google.com/p/gopacket/layers"
        "code.google.com/p/gopacket/pcap"
)

var elements map[string][]uint32

func Detect() {
	
	elements = make(map[string][]uint32)
	
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
    for{fmt.Scanln()}
}

func readIPSec(handle *pcap.Handle, iface *net.Interface, stop chan struct{}) {
        src := gopacket.NewPacketSource(handle, handle.LinkType())
     
        for packet := range src.Packets(){
        	ipsecLayer := packet.Layer(layers.LayerTypeIPSecESP)
            if ipsecLayer != nil{
            	netFlow := packet.NetworkLayer().NetworkFlow()
				src, dst := netFlow.Endpoints()
                ipsec := ipsecLayer.(*layers.IPSecESP)		
				log.Println("Source: ",src, "Destination: ",dst,"Seqnumber: ",ipsec.Seq)
				makeEntry(strings.Join([]string{src.String(), dst.String()},","), ipsec.Seq)
             }
        }
        
        for k,element := range elements{
        	fmt.Println(k);
        	for _,seqnumber := range element{
        		fmt.Println(seqnumber)
        	}
        }
}

func makeEntry(key string, value uint32){
	if(elements[key] == nil){
		elements[key] = []uint32{value}
	}else{
		elements[key] = append(elements[key],value)
	}
}
