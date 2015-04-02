package packetloss

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	esp := NewEspMap(32)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 5)
	fmt.Println(esp.elements[con].head)
	esp.MakeEntry(con, 6)
	esp.MakeEntry(con, 7)
	esp.MakeEntry(con, 8)
	esp.MakeEntry(con, 10)

	for _, e := range esp.elements {
		fmt.Println("Head: ", e.head)
		for _, e := range e.lostpackets {
			fmt.Println("LostPackets: ", e.sequencenumber)
		}
		for _, e := range e.maybelostpackets {
			fmt.Println("MaybeLostPackets: ", e)
		}

		lp := len(esp.elements[con].lostpackets)
		mlp := len(esp.elements[con].maybelostpackets)

		if lp != 0 || mlp != 1 {
			t.Error("Expected lostpackets 0 but it's: ", len(esp.elements[con].lostpackets), "and maybelostpackets 1 but it's:", len(esp.elements[con].maybelostpackets))
		}
	}
}

/*
func Test2(t *testing.T) {
	esp := NewEspMap(32)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 5)
	esp.MakeEntry(con, 6)
	esp.MakeEntry(con, 8)
	esp.MakeEntry(con, 9)
	esp.MakeEntry(con, 13)
	esp.CheckForLost()
	if len(esp.lostpackets[con]) != 2 {
		t.Error("Expected lostpackets 1 but it's: ", esp.lostpackets[con])
	}
}

func Test3(t *testing.T) {
	esp := NewEspMap(32)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 5)
	esp.MakeEntry(con, 6)
	esp.MakeEntry(con, 7)
	esp.MakeEntry(con, 8)
	esp.MakeEntry(con, 9)
	esp.CheckForLost()
	if len(esp.lostpackets[con]) != 0 {
		t.Error("Expected lostpackets 1 but it's: ", esp.lostpackets[con])
	}
}

func Test4(t *testing.T) {
	esp := NewEspMap(32)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 5)
	esp.MakeEntry(con, 6)
	esp.MakeEntry(con, 7)
	esp.MakeEntry(con, 9)
	esp.MakeEntry(con, 8)
	esp.CheckForLost()
	if len(esp.lostpackets[con]) != 0 {
		t.Error("Expected lostpackets 1 but it's: ", esp.lostpackets[con])
	}
}
*/
