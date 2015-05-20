package packetloss

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

//Writes a csv file.
//The file provides information about lost packets
func WriteLostFile(con Connection, lostpackets []LostPacket) {
	s := []string{strconv.FormatUint(uint64(con.SPI), 10), ".csv"}
	f, err := os.Create(strings.Join(s, ""))
	check(err)

	defer f.Close()
	s = []string{"SPI:", strconv.FormatUint(uint64(con.SPI), 10), " SRC: ", con.src, " DST: ", con.dst, "\n"}
	f.WriteString(strings.Join(s, ""))
	f.WriteString("Sequencenumber;Timestamp;ReceivedLater\n")
	//f.WriteString("
	for _, v := range lostpackets {
		fmt.Println(v.sequencenumber)
		s = []string{
			strconv.FormatUint(uint64(v.sequencenumber), 10), ";",
			string(v.Timestamp.Format(time.RFC850)), ";",
			strconv.FormatBool(v.late), "\n"}
		f.WriteString(strings.Join(s, ""))
	}

	f.Sync()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
