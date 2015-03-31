package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"math/rand"
	"time"
)

//Constants & magic values:
const configFile string = "config.json"
const configVersion int = 1
const StartingMTU int = 500
const MTUIterations int = 3

//Config contains the user configurable values for IPSecDiagTool.
//It contains only primitive datatypes so that it is easily serializable.
type Config struct {
	ApplicationID int
	Debug bool

	//MTU specific:
	SourceIP           string
	DestinationIP      string
	Port               int
	IncrementationStep int

	//Packet loss specific:
	//add here..

	//Used to determine whether configuration needs to be updated.
	CfgVers int
}

//initalize creates a new config with default values and writes it to disk.
func initalize() Config {
	conf := Config{0, false, "127.0.0.1", "127.0.0.1", 22, 100, configVersion}
	Write(conf)
	//TODO: perhaps write AppID to file later?
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
	if(configOutdated(conf)){
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
	if _, err := os.Stat(configFile); err == nil {
		fmt.Println("Existing config found.")
		return Read()
	} else {
		fmt.Println("No config found, running init.")
		return initalize()
	}
}

//setupAppID generates a new Application ID if the existing appID equals 0.
//If the existing ID doesn't equal 0, then it will be returned instead.
func setupAppID(applicationID int) int {
	if applicationID == 0 {
		rand.Seed(time.Now().UnixNano()) //Seed is required otherwise we always get the same number
		id := rand.Intn(100000)
		fmt.Println("Application ID generated:",id)
		return id
	} else {
		return applicationID
	}
}

func configOutdated(c Config) bool{
	if(c.CfgVers < configVersion){
		return true
	}
	return false
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}
