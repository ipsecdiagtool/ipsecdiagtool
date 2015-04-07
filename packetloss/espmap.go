package packetloss

import (
	"time"
)

type EspMap struct {
	windowsize uint32
	elements   map[Connection]Packets
}

type Connection struct {
	src, dst string
	SPI      uint32
}

type Packets struct {
	head             uint32
	lostpackets      []LostPacket
	maybelostpackets []uint32
}

type LostPacket struct {
	sequencenumber uint32
	timestamp      time.Time
}

func NewEspMap(windowsize uint32) *EspMap {
	return &EspMap{elements: make(map[Connection]Packets),
		windowsize: windowsize} //, lostpackets: make(map[Connection][]LostPacket)}
}

func (espm EspMap) MakeEntry(key Connection, value uint32) {
	if espm.elements[key].head == 0 {
		espm.elements[key] = Packets{head: value}
	} else {
		packets := espm.elements[key]

		if value > packets.head {

			if value != packets.head+1 {
				for i := packets.head + 1; i < value; i++ {
					packets.maybelostpackets = append(packets.maybelostpackets, i)
				}

			}
			packets.head = value

		} else {
			if (packets.head-espm.windowsize) < value || value < espm.windowsize {
				for i, v := range packets.maybelostpackets {
					if v == value {
						packets.maybelostpackets = append(packets.maybelostpackets[:i], packets.maybelostpackets[i+1:]...)
					}
				}

			}
		}
		espm.elements[key] = CheckForLost(packets, espm.windowsize)
	}
}

func CheckForLost(packets Packets, windowsize uint32) Packets {
	var newMaybelost []uint32
	for _, e := range packets.maybelostpackets {
		if packets.head > windowsize && packets.head-windowsize >= e {
			packets.lostpackets = append(packets.lostpackets, LostPacket{e, time.Now().Local()})
		} else {
			newMaybelost = append(newMaybelost, e)
		}
	}
	packets.maybelostpackets = newMaybelost
	return packets

}
