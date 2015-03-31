package config

import (
	"testing"
	"os"
)

func TestInitReadWriteCompare(t *testing.T) {
	initializedConf := initalize()
	readConf := Read()

	initializedConf.ApplicationID = 0;
	readConf.ApplicationID = 0;
	if(initializedConf != readConf){
		t.Error("Initialized configuration and read configuration do not match.")
	}
	os.Remove(configFile)
}

func TestSetupAppID(t *testing.T) {
	id1 := setupAppID(0);
	id2 := setupAppID(0);
	id3 := setupAppID(666);

	if(id1 == id2){
		t.Error("Random App ID not random..")
	}
	if(id3 != 666){
		t.Error("Expected id3 to be 666, not", id3)
	}
}

func TestOutDatedConfigMechanism(t *testing.T){
	initializedConf := initalize()
	initializedConf.CfgVers = 0
	Write(initializedConf)

	updatedConf := LoadConfig()

	if(updatedConf.CfgVers == 0){
		t.Error("Configuration not properly updated.")
	}
	os.Remove(configFile)
}
