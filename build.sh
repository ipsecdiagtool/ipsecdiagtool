#!/usr/bin/env bash

#-----------------------------------------------------------#
#  Title:  build.sh                                         #
#  URL:    https://github.com/IPSecDiagTool/IPSecDiagTool   #
#  Author: Jan Balmer, Theo Winter                          #
#                                                           #
#  This script can be used to build IPSecDiagTool in a      #
#  Linux environment. Tested on Ubuntu 14.04.1 LTS.         #
#                                                           #
#  Dependencies:                                            #
#   - libpcap0.8-dev                                        #
#-----------------------------------------------------------#

echo "Cleaning workspace"
if [ -d go ]; then
	rm -rf go1.4.2.linux-amd64.tar.gz
	rm -rf go
    rm -rf workspace
fi

echo "Debug Infos"
go version
go env

echo "Setting up Go"
wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
tar xf go1.4.2.linux-amd64.tar.gz
mkdir workspace

export GOROOT=$(pwd)/go
export GOPATH=$(pwd)/workspace
export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"

cd workspace

#Can be deactivated if libpcap is already installed.
echo "Downloading dependencies"
sudo apt-get install libpcap0.8-dev

echo "Getting and building ipsecdiagtool"
go get github.com/ipsecdiagtool/ipsecdiagtool