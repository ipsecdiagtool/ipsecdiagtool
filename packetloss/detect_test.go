package packetloss

import (
	"fmt"
	"github.com/ipsecdiagtool/ipsecdiagtool/config"
	"testing"
)

func TestInterfaceName(t *testing.T) {
	fmt.Println("**********TestInterfaceName**********")
	configuration := config.LoadConfig()
	configuration.InterfaceName = "test"
	if err := Detect(configuration); err != nil {
		// handle err
		fmt.Println("err: ", err)
		fmt.Println("**********TestInterfaceName OK**********")
	} else {
		t.Error("Wrong Interfacename not detected!")
	}
}
