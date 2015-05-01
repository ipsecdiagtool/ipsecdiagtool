package packetloss

import (
	"github.com/ipsecdiagtool/ipsecdiagtool/logging"
	"strings"
	"time"
	"fmt"
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
	LastAlert        time.Time
	head             uint32
	lostpackets      []LostPacket
	maybelostpackets []uint32
}

type LostPacket struct {
	sequencenumber uint32
	Timestamp      time.Time
	late           bool
}

func NewEspMap(windowsize uint32) *EspMap {
	return &EspMap{elements: make(map[Connection]Packets),
		windowsize: windowsize}
}

//Processing new Packets and adjusting the value of Head.
func (espm EspMap) MakeEntry(key Connection, value uint32) {
	if espm.elements[key].head != 0 {

		packets := espm.elements[key]

		if value > packets.head {
			handleNewpacket(&packets, value)
			checkLost(&packets, espm.windowsize)

		} else {
			handleOldpacket(&packets, value, espm.windowsize)

		}
		
		if time.Now().Local().Sub(packets.LastAlert).Seconds() > 10 && CheckLog(packets.lostpackets) {
			s := []string{"Too much LostPackets in Connection: (SPI: ", string(key.SPI), " SRC: ", key.src, " DST: ", key.dst, ")"}
			logging.AlertLog(strings.Join(s, ""))
			packets.LastAlert = time.Now().Local()
		}
		espm.elements[key] = packets
	} else {
		 t1, e := time.Parse(time.RFC3339,"2012-11-01T22:08:41+00:00")
		 if(e!=nil){
		 	panic(e)
		 }
		 fmt.Println(t1.Format("2006-01-02T15:04:05.999999-07:00"))
		espm.elements[key] = Packets{head: value, LastAlert:  t1}

	}
}

//Checks if MaybeLost values are valid.
//If values are not within the WindowSize they are definitly lost
func checkLost(packets *Packets, windowsize uint32) {
	var newMaybelost []uint32
	for _, e := range packets.maybelostpackets {
		if packets.head > windowsize && packets.head-windowsize >= e {
			packets.lostpackets = append(packets.lostpackets, LostPacket{e, time.Now().Local(), false})
		} else {
			newMaybelost = append(newMaybelost, e)
		}
	}
	packets.maybelostpackets = newMaybelost
}

//Handles a Packet with a Sequencnumber bigger than current Head
func handleNewpacket(packets *Packets, value uint32) {
	if value != packets.head+1 {
		for i := packets.head + 1; i < value; i++ {
			packets.maybelostpackets = append(packets.maybelostpackets, i)
		}
	}
	packets.head = value
}

//Handles a Packet with a Sequencnumber less than current Head
func handleOldpacket(packets *Packets, value uint32, windowsize uint32) {
	if (packets.head-windowsize) < value || value < windowsize {

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

//Checks the current espmap if alert logging is necessary
func CheckLog(lostpackets []LostPacket) bool {
	var counter int
	currenttime := time.Now().Local()
	for _, v := range lostpackets {
		seconds := currenttime.Sub(v.Timestamp).Seconds()
		if seconds < float64(logging.AlertTime()) {
			counter++
		}
	}
	return counter > logging.AlertCounter()
}
