package mtu

import (
	"code.google.com/p/gopacket/examples/util"
	"net"
	"log"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket"
	"golang.org/x/net/ipv4"
)


func Analyze(){
	defer util.Run()()
	log.Println("Analyzing MTU..")
	buildPacket("127.0.0.1","127.0.0.1")


}

func buildPacket(sourceIP string, destinationIP string){
	var srcIP = net.ParseIP(sourceIP)
	var dstIP = net.ParseIP(destinationIP)

	if srcIP == nil || dstIP == nil {
		log.Println("Invalid IP")
	}

	//Convert IP to 4bit representation
	srcIP = srcIP.To4()
	dstIP = dstIP.To4()

	//IP Layer
	ip := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolTCP,
	}

	srcPort := layers.TCPPort(666)
	dstPort := layers.TCPPort(22)

	//TCP Layer
	tcp := layers.TCP{
		SrcPort: srcPort,
		DstPort: dstPort,
		Window:  1505,
		Urgent:  0,
		Seq:     11050,
		Ack:     0,
		ACK:     false,
		SYN:     false,
		FIN:     false,
		RST:     false,
		URG:     false,
		ECE:     false,
		CWR:     false,
		NS:      false,
		PSH:     false,
	}

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	tcp.SetNetworkLayerForChecksum(&ip)

	ipHeaderBuf := gopacket.NewSerializeBuffer()

	err := ip.SerializeTo(ipHeaderBuf, opts)
	if err != nil {
		panic(err)
	}

	//Set Don't Fragment in Header
	ipHeader, err := ipv4.ParseHeader(ipHeaderBuf.Bytes())
	ipHeader.Flags |= ipv4.DontFragment
	if err != nil {
		panic(err)
	}

	tcpPayloadBuf := gopacket.NewSerializeBuffer()

	//Influence the payload size
	payload := gopacket.Payload([]byte("Hello IPSec Hello IPSec Hello IPSec Hello IPSec"))
	err = gopacket.SerializeLayers(tcpPayloadBuf, opts, &tcp, payload)
	if err != nil {
		panic(err)
	}

	//Send packet
	var packetConn net.PacketConn
	var rawConn *ipv4.RawConn
	packetConn, err = net.ListenPacket("ip4:tcp", sourceIP)
	if err != nil {
		panic(err)
	}
	rawConn, err = ipv4.NewRawConn(packetConn)
	if err != nil {
		panic(err)
	}

	err = rawConn.WriteTo(ipHeader, tcpPayloadBuf.Bytes(), nil)

	log.Println("Packet with length", (len(tcpPayloadBuf.Bytes()) + len(ipHeaderBuf.Bytes())), "sent.")
}
