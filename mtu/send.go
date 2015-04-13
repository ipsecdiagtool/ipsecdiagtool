package mtu

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"strconv"
)

func sendOKResponse(packet gopacket.Packet) {
	sendPacket(conf.SourceIP, getIP(packet).String(), conf.Port, originalSize(packet), "OK")
}

//sendPacket generates & sends a packet of arbitrary size to a specific destination.
//The size specified should be larger then 40bytes.
func sendPacket(sourceIP string, destinationIP string, destinationPort int, size int, message string) []byte {

	var payloadSize int
	if size < 40 {
		log.Println("Unable to create a packet smaller then 40bytes.")
		payloadSize = 0
	} else {
		payloadSize = size - 40
	}

	//Convert IP to 4bit representation
	srcIP := net.ParseIP(sourceIP).To4()
	dstIP := net.ParseIP(destinationIP).To4()

	//IP Layer
	ip := layers.IPv4{
		SrcIP:    srcIP,
		DstIP:    dstIP,
		Version:  4,
		TTL:      64,
		Protocol: layers.IPProtocolUDP,
	}

	srcPort := layers.UDPPort(666)
	dstPort := layers.UDPPort(destinationPort)

	//UDP Layer
	udp := layers.UDP{
		SrcPort: srcPort,
		DstPort: dstPort,
	}

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

	udp.SetNetworkLayerForChecksum(&ip)

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

	udpPayloadBuf := gopacket.NewSerializeBuffer()

	//Influence the payload size
	payload := gopacket.Payload(generatePayload(payloadSize, strconv.Itoa(conf.ApplicationID)+","+message+","))
	err = gopacket.SerializeLayers(udpPayloadBuf, opts, &udp, payload)
	if err != nil {
		panic(err)
	}

	//Send packet
	var packetConn net.PacketConn
	var rawConn *ipv4.RawConn

	packetConn, err = net.ListenPacket("ip4:udp", srcIP.String())
	if err != nil {
		panic(err)
	}
	rawConn, err = ipv4.NewRawConn(packetConn)
	if err != nil {
		panic(err)
	}

	err = rawConn.WriteTo(ipHeader, udpPayloadBuf.Bytes(), nil)

	log.Println("Packet with length", (len(udpPayloadBuf.Bytes()) + len(ipHeaderBuf.Bytes())), "sent.")
	return append(ipHeaderBuf.Bytes(), udpPayloadBuf.Bytes()...)
}

//generatePayload generates a payload of the given size (bytes).
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
