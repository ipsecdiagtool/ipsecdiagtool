package packetloss

import (
	"time"
)

//Datastructure for different ESP Connections
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
	late           bool
}

func NewEspMap(windowsize uint32) *EspMap {
	return &EspMap{elements: make(map[Connection]Packets),
		windowsize: windowsize}
}

//Processing new Packets and adjusting the value of Head.
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
			espm.elements[key] = checkLost(packets, espm.windowsize)
		} else {
			if (packets.head-espm.windowsize) < value || value < espm.windowsize {
				for i, v := range packets.maybelostpackets {
					if v == value {
						packets.maybelostpackets = append(packets.maybelostpackets[:i], packets.maybelostpackets[i+1:]...)
					}
				}
			} else {
				for k, v := range packets.lostpackets {
					if v.sequencenumber == value {
						packets.lostpackets[k].late = true
					}
				}
			}
		}
	}
}

//Checks if MaybeLost values are valid.
//If values are not within the WindowSize they are definitly lost
func checkLost(packets Packets, windowsize uint32) Packets {
	var newMaybelost []uint32
	for _, e := range packets.maybelostpackets {
		if packets.head > windowsize && packets.head-windowsize >= e {
			packets.lostpackets = append(packets.lostpackets, LostPacket{e, time.Now().Local(), false})
		} else {
			newMaybelost = append(newMaybelost, e)
		}
	}
	packets.maybelostpackets = newMaybelost
	return packets
}
