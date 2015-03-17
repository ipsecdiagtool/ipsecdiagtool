package capture

import "testing"
import "fmt"
import "strconv"

/*
	Golang Testing HowTo:
	+ Name of Test: Something_test.go
	+ always import "testing"
	+ Put test in the same folder where the original .go-file is located
 */

func TestReadPcapFile(t *testing.T) {
	var numberOfPackets = ReadPcapFile("/home/parallels/Desktop/capture.pcap")
	fmt.Println("Number of Packets :"+strconv.Itoa(numberOfPackets))
	if numberOfPackets != 841 {
		t.Error("Expected 841, got", numberOfPackets)
	}
}
