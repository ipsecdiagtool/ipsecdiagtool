package mtu

import (
	"code.google.com/p/gopacket"
	"code.google.com/p/gopacket/layers"
	"golang.org/x/net/ipv4"
	"net"
	"strconv"
)

func sendOKResponse(packet gopacket.Packet, appID int, chanID int) {
	srcIP, dstIP := getSrcDstIP(packet)
	sendPacket(dstIP.String(), srcIP.String(), originalSize(packet), "OK", appID, chanID)
}

//sendPacket generates & sends a packet of arbitrary size to a specific destination.
//The size specified should be larger then 40bytes.
func sendPacket(sourceIP string, destinationIP string, size int, message string, appID int, chanID int) []byte {

	var payloadSize int
	if size < 28 {
		//log.Println("Unable to create a packet smaller then 28bytes.")
		payloadSize = 0
	} else {
		payloadSize = size - 28
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
		Protocol: layers.IPProtocolICMPv4,
	}

	icmp := layers.ICMPv4{}

	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}

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

	payloadBuf := gopacket.NewSerializeBuffer()

	//Influence the payload size
	payload := gopacket.Payload(generatePayload(payloadSize, ","+strconv.Itoa(appID)+","+strconv.Itoa(chanID)+","+message+","))
	err = gopacket.SerializeLayers(payloadBuf, opts, &icmp, payload)
	if err != nil {
		panic(err)
	}

	//Send packet
	var packetConn net.PacketConn
	var rawConn *ipv4.RawConn

	packetConn, err = net.ListenPacket("ip4:icmp", srcIP.String())
	if err != nil {
		panic(err)
	}
	rawConn, err = ipv4.NewRawConn(packetConn)
	if err != nil {
		panic(err)
	}

	err = rawConn.WriteTo(ipHeader, payloadBuf.Bytes(), nil)

	//log.Println("Packet with length", (len(payloadBuf.Bytes()) + len(ipHeaderBuf.Bytes())), "sent.")
	return append(ipHeaderBuf.Bytes(), payloadBuf.Bytes()...)
}

//generatePayload generates a payload of the given size (bytes) unless the wanted size is smaller then the message
//in that case it simply returns the message as payload. We're assuming the message is more important than the exact size.
func generatePayload(size int, message string) []byte {
	var payload []byte
	if size > len(message) {
		payload = make([]byte, size-len(message))
		payload = append([]byte(message), payload...)
	} else {
		payload = []byte(message)
	}
	return payload
}
