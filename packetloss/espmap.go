package packetloss

import (
	"time"
)

type EspMap struct {
	windowsize  int
	lostpackets map[Connection][]LostPacket
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
		windowsize: windowsize, lostpackets: make(map[Connection][]LostPacket)}
}

func (espm EspMap) MakeEntry(key Connection, value uint32) {
		espm.elements[key] = append(espm.elements[key], value)
}

func (espm EspMap) CheckForLost() {

	for k, element := range espm.elements {
		for key, seqnumber := range element {

			if key != 0 {
				l := seqnumber - element[key-1]
				if l != 1 {
					espm.lostpackets[k] = append(espm.lostpackets[k], LostPacket{k, l-1, time.Now().Local()})			
				}
			}
		}
	}

}
