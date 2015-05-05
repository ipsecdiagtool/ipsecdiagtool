package config

import (
	"os"
	"testing"
)

//Init a new config, write it to a file with a specific value,
//read it again and check value.
func TestReadWrite(t *testing.T){
	conf := initialize()
	conf.ApplicationID = 1337
	conf.Debug = true
	Write(conf)
	readConf := Read()

	if readConf.ApplicationID != 1337 {
		t.Error("Wrote a config with 1337 as AppID, read the file and got", readConf.ApplicationID)
	} else if readConf.Debug != true {
		t.Error("Wrote a config with Debug=true, read the file and got", readConf.Debug)
	}
	os.Remove(configFile)
}

//Check that a random AppID is generated if AppID is 0
func TestSetupAppID(t *testing.T) {
	id1 := setupAppID(0)
	id2 := setupAppID(0)
	id3 := setupAppID(666)

	if id1 == id2 {
		t.Error("Random App ID not random..")
	}
	if id3 != 666 {
		t.Error("Expected id3 to be 666, not", id3)
	}
}

//Check that updating an outdated config works
func TestOutDatedConfigMechanism(t *testing.T) {
	initializedConf := initialize()
	initializedConf.CfgVers = 0
	Write(initializedConf)

	updatedConf := LoadConfig(os.Args[0])

	if updatedConf.CfgVers == 0 {
		t.Error("Configuration not properly updated.")
	}
	os.Remove(configFile)
}
