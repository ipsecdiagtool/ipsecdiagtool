package packetloss

import (
    "testing"
)

func TestXYZ(t *testing.T) {
	esp:=NewEspMap(32)
	esp.MakeEntry(Connection{"192.168.0.1", "192.168.0.1", 12345} , 5)
	esp.MakeEntry(Connection{"192.168.0.1", "192.168.0.1", 12345} , 6)
	esp.MakeEntry(Connection{"192.168.0.1", "192.168.0.1", 12345} , 7)
	esp.MakeEntry(Connection{"192.168.0.1", "192.168.0.1", 12345} , 8)
	esp.MakeEntry(Connection{"192.168.0.1", "192.168.0.1", 12345} , 9)
	esp.CheckForLost()
if len(esp.lostpackets) != 0 {
		t.Error("Expected", length, "bytes of payload, instead got", len(payload), "bytes.")
	}
} 