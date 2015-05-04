package packetloss

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Writes a csv file.
//The file provides information about lost packets
func WriteLostFile(con Connection, lostpackets []LostPacket) {
	f, err := os.Create("/home/student/Desktop/test.csv")
	check(err)

	defer f.Close()
	s := []string{"SPI:", strconv.FormatUint(uint64(con.SPI), 10), " SRC: ", con.src, " DST: ", con.dst, "\n"}
	f.WriteString(strings.Join(s, ""))
	f.WriteString("Sequencenumber;Timestamp;ReceivedLater\n")
	//f.WriteString("
	for _, v := range lostpackets {
		fmt.Println(v.sequencenumber)
		s = []string{
			strconv.FormatUint(uint64(v.sequencenumber), 10), ";",
			string(v.Timestamp.Format("2006-01-02T15:04:05.999999-07:00")), ";",
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

/*
func main() {

	// To start, here's how to dump a string (or just
	// bytes) into a file.
	d1 := []byte("hello\ngo\n")
	err := ioutil.WriteFile("/tmp/dat1", d1, 0644)
	check(err)

	// For more granular writes, open a file for writing.
	f, err := os.Create("/tmp/dat2")
	check(err)

	// It's idiomatic to defer a `Close` immediately
	// after opening a file.
	defer f.Close()

	// You can `Write` byte slices as you'd expect.
	d2 := []byte{115, 111, 109, 101, 10}
	n2, err := f.Write(d2)
	check(err)
	fmt.Printf("wrote %d bytes\n", n2)

	// A `WriteString` is also available.
	n3, err := f.WriteString("writes\n")
	fmt.Printf("wrote %d bytes\n", n3)

	// Issue a `Sync` to flush writes to stable storage.
	f.Sync()

	// `bufio` provides buffered writers in addition
	// to the buffered readers we saw earlier.
	w := bufio.NewWriter(f)
	n4, err := w.WriteString("buffered\n")
	fmt.Printf("wrote %d bytes\n", n4)

	// Use `Flush` to ensure all buffered operations have
	// been applied to the underlying writer.
	w.Flush()
}*/
