package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"time"
)

//Debug is mainly used to determine whether to report a log message or not.
var Debug = false

const configFile string = "config.json"
const configVersion int = 11

//Config contains the user configurable values for IPSecDiagTool.
//It can hold multiple MTUConfig's to handle MTU detection for multiple tunnels.
type Config struct {
	ApplicationID int
	Debug         bool
	SyslogServer  string //IP:Port
	PcapSnapLen   int32

	//MTU specific:
	MTUConfList []MTUConfig

	//Packet loss specific:
	WindowSize    uint32
	InterfaceName string
	AlertTime     int // Time in Seconds for LostPacketsCheck
	AlertCounter  int // Packets in LostPacketsCheck
	PcapFile      string

	//Used to determine whether configuration needs to be updated.
	CfgVers int
}

//MTUConfig contains all the necessary settings to detect the MTU of one tunnel.
type MTUConfig struct {
	SourceIP      string
	DestinationIP string
	Timeout       time.Duration
	MTURangeStart int
	MTURangeEnd   int
}

//initialize creates a new config with default values and writes it to disk.
func initialize() Config {
	mtuSample := MTUConfig{"127.0.0.1", "127.0.0.1", 10, 0, 2000}
	mtuList := []MTUConfig{mtuSample, mtuSample}
	conf := Config{0, false, "localhost:514", 3000, mtuList, 32, "any", 60, 10, "", configVersion}
	Write(conf)
	conf.ApplicationID = setupAppID(conf.ApplicationID)
	return conf
}

//Read an existing config file and return it.
func Read() Config {
	jsonConfig, err := ioutil.ReadFile(configFile)
	check(err)

	var conf Config
	err2 := json.Unmarshal(jsonConfig, &conf)
	check(err2)

	//Update config file if outdated
	if configOutdated(conf) {
		fmt.Println("Outdated configuration found, updating it now.")
		conf.CfgVers = configVersion
		Write(conf)
	}

	conf.ApplicationID = setupAppID(conf.ApplicationID)
	return conf
}

//Write a config to the disk
func Write(conf Config) {
	jsonConfig, err := json.MarshalIndent(conf, "", "    ")
	check(err)

	w, err := os.Create(configFile)
	check(err)

	defer w.Close()
	w.Write(jsonConfig)
}

//LoadConfig tries to read an existing config from the program directory first and in the users current working
//directory second. If neither folder contains a config it will initialize a new config.
func LoadConfig(location string) Config {
	var conf Config
	if _, err := os.Stat(location + "/" + configFile); err == nil {
		log.Println("Existing config found in program directory.")
		conf = Read()
	} else if _, err := os.Stat(configFile); err == nil {
		log.Println("Existing config found in current working directory.")
		conf = Read()
	} else {
		log.Println("No config found, running init.")
		conf = initialize()
	}
	Debug = conf.Debug
	if Debug {
		log.Println("Debug-Mode enabled. Verbose reporting of status & errors.")
	}
	return conf
}

//setupAppID generates a new Application ID if the existing appID equals 0.
//If the existing ID doesn't equal 0, then it will be returned instead.
func setupAppID(applicationID int) int {
	if applicationID == 0 {
		rand.Seed(time.Now().UnixNano())
		applicationID = rand.Intn(100000)

		//Prevent generation of 1337
		if applicationID == 1337 {
			applicationID = setupAppID(0)
		}
		if Debug {
			log.Println("Application ID generated:", applicationID)
		}
	}
	return applicationID
}

func configOutdated(c Config) bool {
	if c.CfgVers < configVersion {
		return true
	}
	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
