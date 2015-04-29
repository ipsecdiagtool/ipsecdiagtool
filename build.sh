#!/usr/bin/env bash

echo "Cleaning workspace"
if [ -d go ]; then
	rm -rf go1.4.2.linux-amd64.tar.gz
	rm -rf go
    rm -rf workspace
fi

echo "Setting up Go"
wget https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz
tar xf go1.4.2.linux-amd64.tar.gz
mkdir workspace

export GOROOT=$(pwd)/go
export GOPATH=$(pwd)/workspace
export PATH="$PATH:$GOROOT/bin:$GOPATH/bin"

cd workspace

echo "Downloading dependencies"
#sudo apt-get install libpcap-0.8-dev
go get code.google.com/p/gopacket
go get code.google.com/p/gopacket/pcap
go get golang.org/x/net/ipv4

echo "Moving to program directory"
cd src
mkdir -p github.com/ipsecdiagtool/
cd github.com/ipsecdiagtool
git clone https://github.com/IPSecDiagTool/IPSecDiagTool.git
mv IPSecDiagTool ipsecdiagtool
cd ipsecdiagtool

echo "Building IPSecDiagTool"
go build