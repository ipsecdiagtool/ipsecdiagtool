package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"os"
	"time"
	"log"
)

//Debug is mainly used to determine whether to report a log message or not.
var Debug = false

const configFile string = "config.json"
const configVersion int = 9

//Config contains the user configurable values for IPSecDiagTool.
//It can hold multiple MTUConfig's to handle MTU detection for multiple tunnels.
type Config struct {
	ApplicationID int
	Debug         bool

	//MTU specific:
	MTUConfList []MTUConfig

	//Packet loss specific:
	WindowSize    uint32
	InterfaceName string
	AlertTime     int    // Time in Seconds for LostPacketsCheck
	AlertCounter  int    // Packets in LostPacketsCheck
	SyslogServer  string //Path inclusive Port

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
	conf := Config{0, false, mtuList, 32, "any", 60, 10, "localhost:514", configVersion}
	Write(conf)
	conf.ApplicationID = setupAppID(conf.ApplicationID)
	return conf
}

//Read reads an existing config file and returns it as a config object. If you're loading
//the configuration for the first time you should use LoadConfig() instead.
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

//Write accepts a Config object and writes it to the disk.
func Write(conf Config) {
	jsonConfig, err := json.MarshalIndent(conf, "", "    ")
	check(err)

	w, err := os.Create(configFile)
	check(err)

	defer w.Close()
	w.Write(jsonConfig)
}

//LoadConfig reads an existing config file if it exists. If it doesn't
//exist a new config, containing default values, is written.
func LoadConfig() Config {
	var conf Config
	if _, err := os.Stat(configFile); err == nil {
		fmt.Println("Existing config found.")
		conf = Read()
	} else {
		fmt.Println("No config found, running init.")
		conf = initialize()
	}
	Debug = conf.Debug
	if(Debug){
		log.Println("Debug-Mode enabled. Verbose reporting of status & errors.")
	}
	return conf
}

//setupAppID generates a new Application ID if the existing appID equals 0.
//If the existing ID doesn't equal 0, then it will be returned instead.
func setupAppID(applicationID int) int {
	if applicationID == 0 {
		rand.Seed(time.Now().UnixNano())
		applicationID = rand.Intn(100000)

		//Prevent generation of 1337
		if applicationID == 1337 {
			applicationID = setupAppID(0)
		}

		fmt.Println("Application ID generated:", applicationID)
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
