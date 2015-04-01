package packetloss

import (
	"time"
)

type EspMap struct {
	windowsize  int
	lostpackets []LostPacket
	elements    map[Connection][]uint32
}

type Connection struct {
	src, dst string
	SPI      uint32
}

type LostPacket struct {
	conn      Connection
	lost      uint32
	timestamp time.Time
}

func NewEspMap(windowsize int) *EspMap {
	return &EspMap{elements: make(map[Connection][]uint32),
		windowsize: windowsize}
}

func NewLostPacket(lost uint32, conn Connection) *LostPacket {
	return &LostPacket{lost: lost, timestamp: time.Now().Local(), conn: conn}
}

func (espm EspMap) MakeEntry(key Connection, value uint32) {
	if espm.elements[key] == nil {
		espm.elements[key] = []uint32{value}
	} else {
		espm.elements[key] = append(espm.elements[key], value)
	}
}

func (espm EspMap) CheckForLost() {

	for k, element := range espm.elements {
		for key, seqnumber := range element {
			l := seqnumber - element[key-1]
			if key != 0 && l != 1 {
				espm.lostpackets = append(espm.lostpackets, LostPacket{k, l, time.Now().Local()})
			}
		}
	}

}
