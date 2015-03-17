package capture

import "testing"
import "fmt"
import "strconv"

func TestReadPcapFile(t *testing.T) {
	var numberOfPackets = ReadPcapFile("/home/parallels/Desktop/capture.pcap")
	fmt.Println("Number of Packets :"+strconv.Itoa(numberOfPackets))
	if numberOfPackets != 841 {
		t.Error("Expected 841, got", numberOfPackets)
	}
}
