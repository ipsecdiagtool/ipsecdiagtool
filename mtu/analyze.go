package mtu

import (
	"code.google.com/p/gopacket/examples/util"
	"net"
	"log"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket"
	"golang.org/x/net/ipv4"
)

//Analyze computes the ideal MTU for a conntection between to computers.
func Analyze(){
	defer util.Run()()
	log.Println("Analyzing MTU..")
	sendPacket("127.0.0.1","127.0.0.1", 80)

	//TODO:
	//-Record packet
	//-Loop several times to find ideal MTU

}

//sendPacket generates & sends a packet of arbitrary size to a specific destination.
//The size specified should be larger then 40bytes.
func sendPacket(sourceIP string, destinationIP string, size int){
	var payloadSize int
	if size < 40 {
		log.Println("Unable to create a packet smaller then 40bytes.")
		payloadSize = 0
	} else {
		payloadSize = size-40
	}

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
		Seq:     0x539, //Hex 1337
		Window: 1337,
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
	payload := gopacket.Payload(generatePayload(payloadSize))
	err = gopacket.SerializeLayers(tcpPayloadBuf, opts, &tcp, payload)
	if err != nil {
		panic(err)
	}

	//Send packet
	var packetConn net.PacketConn
	var rawConn *ipv4.RawConn
	packetConn, err = net.ListenPacket("ip4:tcp", destinationIP)
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

//generatePayload generates a payload of the given size (bytes).
//If the payload is longer then 11 bytes the first eleven bytes are used to spell "Hello IPsec".
func generatePayload(size int) []byte {
	var payload []byte
	if size > 11 {
		payload = make([]byte, size-11)
		payload = append([]byte("Hello IPSec"), payload...)
	} else {
		payload = make([]byte, size)
	}
	return payload
}
