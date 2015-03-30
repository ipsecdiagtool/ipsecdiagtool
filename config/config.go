package config

import (
	"encoding/json"
	"os"
	"fmt"
	"io/ioutil"
)

//Config contains the user configurable values for IPSecDiagTool.
//It contains only primitive datatypes so that it is easily serializable.
type Config struct {
	ApplicationID int

	//MTU specific:
	SourceIP string
	DestinationIP string
	Port int
	IncrementationStep int

	//Packet loss specific:
	//add here..
}

//initalize creates a new config with default values and writes it to disk.
func initalize() Config{
	conf := Config{0, "127.0.0.1", "127.0.0.1", 22, 100}
	Write(conf)
	return conf
}

//Read reads an existing config file and returns it as a config object. If you're loading
//the configuration for the first time you should use LoadConfig() instead.
func Read() Config{
	//TODO: magic constants into one file
	jsonConfig, err := ioutil.ReadFile("config.json")
	check(err)

	var conf Config
	err2 := json.Unmarshal(jsonConfig, &conf)
	check(err2)
	return conf
}

//Write accepts a Config object and writes it to the disk.
func Write(conf Config){
	jsonConfig, err := json.MarshalIndent(conf,"", "    ")
	check(err)

	w, err := os.Create("config.json")
	check(err)

	defer w.Close()
	w.Write(jsonConfig)
}

//LoadConfig reads an existing config file if it exists. If it doesn't
//exist a new config, containing default values, is written.
func LoadConfig() Config{
	if _, err := os.Stat("config.json"); err == nil {
		fmt.Println("Existing config found.")
		return Read()
	} else {
		fmt.Println("No config found, running init.")
		return initalize()
	}
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
