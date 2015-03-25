package mtu

import (
	"log"
	"net"
	"code.google.com/p/gopacket/layers"
	"code.google.com/p/gopacket"
	"golang.org/x/net/ipv4"
	"strconv"
	"time"
)

func sendIncreasedMTU(packet gopacket.Packet) {
	//TODO: slow down the response speed.
	time.Sleep(1000 * time.Millisecond)
	currentMTU += incStep
	sendPacket(srcIP, destIP, destPort, currentMTU, "MTU?")
}

func sendOKResponse(packet gopacket.Packet) {
	sendPacket(srcIP, getIP(packet), 22, 200, "OK")
}

//sendPacket generates & sends a packet of arbitrary size to a specific destination.
//The size specified should be larger then 40bytes.
func sendPacket(srcIP net.IP, dstIP net.IP, destinationPort int, size int, message string) []byte {

	var payloadSize int
	if size < 40 {
		log.Println("Unable to create a packet smaller then 40bytes.")
		payloadSize = 0
	} else {
		payloadSize = size-40
	}

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
	dstPort := layers.TCPPort(destinationPort)

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
	payload := gopacket.Payload(generatePayload(payloadSize, strconv.Itoa(appID)+","+message+","))
	err = gopacket.SerializeLayers(tcpPayloadBuf, opts, &tcp, payload)
	if err != nil {
		panic(err)
	}

	//Send packet
	var packetConn net.PacketConn
	var rawConn *ipv4.RawConn

	packetConn, err = net.ListenPacket("ip4:tcp", dstIP.String())
	if err != nil {
		panic(err)
	}
	rawConn, err = ipv4.NewRawConn(packetConn)
	if err != nil {
		panic(err)
	}

	err = rawConn.WriteTo(ipHeader, tcpPayloadBuf.Bytes(), nil)

	log.Println("Packet with length", (len(tcpPayloadBuf.Bytes())+len(ipHeaderBuf.Bytes())), "sent.")
	return append(ipHeaderBuf.Bytes(),tcpPayloadBuf.Bytes()...)
}

//generatePayload generates a payload of the given size (bytes).
//If the payload is longer then 11 bytes the first eleven bytes are used to spell "Hello IPsec".
func generatePayload(size int, message string) []byte {
	var payload []byte
	if size > len(message) {
		payload = make([]byte, size-len(message))
		payload = append([]byte(message), payload...)
	} else {
		//TODO: Case is probably not relevant. Remove.
		payload = make([]byte, size)
	}
	return payload
}
