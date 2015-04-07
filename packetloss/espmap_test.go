package packetloss

import (
	"fmt"
	"testing"
)

func Test1(t *testing.T) {
	fmt.Println("**********Test1**********")
	esp := NewEspMap(32)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 5)
	esp.MakeEntry(con, 6)
	esp.MakeEntry(con, 7)
	esp.MakeEntry(con, 8)
	esp.MakeEntry(con, 10)

	lp := len(esp.elements[con].lostpackets)
	mlp := len(esp.elements[con].maybelostpackets)

	if lp != 0 || mlp != 1 {
		t.Error("Expected lostpackets 0 but it's: ", len(esp.elements[con].lostpackets), "and maybelostpackets 1 but it's:", len(esp.elements[con].maybelostpackets))
	}
}

func Test2(t *testing.T) {
	fmt.Println("**********Test2**********")
	esp := NewEspMap(2)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 5)
	esp.MakeEntry(con, 6)
	esp.MakeEntry(con, 8)
	esp.MakeEntry(con, 9)
	esp.MakeEntry(con, 13)

	lp := len(esp.elements[con].lostpackets)
	mlp := len(esp.elements[con].maybelostpackets)

	if lp != 3 || mlp != 1 {
		t.Error("Expected lostpackets 3 but it's: ", len(esp.elements[con].lostpackets), "and maybelostpackets 1 but it's:", len(esp.elements[con].maybelostpackets))
	}

}

func Test3(t *testing.T) {
	fmt.Println("**********Test3**********")
	esp := NewEspMap(1)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 5)
	esp.MakeEntry(con, 6)
	esp.MakeEntry(con, 7)
	esp.MakeEntry(con, 8)
	esp.MakeEntry(con, 9)

	lp := len(esp.elements[con].lostpackets)
	mlp := len(esp.elements[con].maybelostpackets)

	if lp != 0 || mlp != 0 {
		t.Error("Expected lostpackets 0 but it's: ", len(esp.elements[con].lostpackets), "and maybelostpackets 0 but it's:", len(esp.elements[con].maybelostpackets))
	}
}

func Test4(t *testing.T) {
	fmt.Println("**********Test4**********")
	esp := NewEspMap(3)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 5)
	esp.MakeEntry(con, 7)
	esp.MakeEntry(con, 6)
	esp.MakeEntry(con, 8)
	esp.MakeEntry(con, 9)
	esp.MakeEntry(con, 10)

	lp := len(esp.elements[con].lostpackets)
	mlp := len(esp.elements[con].maybelostpackets)

	if lp != 0 || mlp != 0 {
		t.Error("Expected lostpackets 0 but it's: ", len(esp.elements[con].lostpackets), "and maybelostpackets 0 but it's:", len(esp.elements[con].maybelostpackets))
	}
}

func Test5(t *testing.T) {
	fmt.Println("**********Test5**********")
	esp := NewEspMap(3)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 1)
	esp.MakeEntry(con, 20)

	lp := len(esp.elements[con].lostpackets)
	mlp := len(esp.elements[con].maybelostpackets)

	if lp != 16 || mlp != 2 {
		t.Error("Expected lostpackets 16 but it's: ", len(esp.elements[con].lostpackets), "and maybelostpackets 2 but it's:", len(esp.elements[con].maybelostpackets))
	}
}

func Test6(t *testing.T) {
	fmt.Println("**********Test6**********")
	esp := NewEspMap(3)
	con := Connection{"192.168.0.1", "192.168.0.1", 12345}
	esp.MakeEntry(con, 1)
	esp.MakeEntry(con, 20)
	esp.MakeEntry(con, 50)

	lp := len(esp.elements[con].lostpackets)
	mlp := len(esp.elements[con].maybelostpackets)

	if lp != 45 || mlp != 2 {
		t.Error("Expected lostpackets 3 but it's: ", len(esp.elements[con].lostpackets), "and maybelostpackets 1 but it's:", len(esp.elements[con].maybelostpackets))
	}
}
