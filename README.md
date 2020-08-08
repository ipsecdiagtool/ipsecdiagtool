# IPSecDiagTool

[![GoDoc](https://godoc.org/github.com/IPSecDiagTool/IPSecDiagTool?status.svg)](https://godoc.org/github.com/IPSecDiagTool/IPSecDiagTool)

**Go Code Quality:** http://goreportcard.com/report/IPSecDiagTool/IPSecDiagTool

## About
IPSecDiagTool is a diagnosis tool for the continuous monitoring of [IPSec](http://en.wikipedia.org/wiki/IPsec) VPN tunnels.
It has two main features. Firstly it's capable of passively detecting packet loss occurring within the IPSec tunnels. If
the packet loss exceeds a specified threshold a Syslog warning is automatically dispatched. Secondly IPSecDiagTool can
actively determine the exact MTU within a tunnel. This is useful when you're dealing with badly configured routers that
block regular ICMP messages outside of the tunnel. IPSecDiagTool is designed to run as a daemon/service, but for testing
purposes it also has a interactive mode.

## Main features
+ Passive detection of IPSec packet loss by capturing arriving ESP packets and gathering their sequence numbers.
+ Active diagnosis of IPSec fragmentation problems and discovery of the ideal MTU (Maximum Transmission Unit).
+ Daemon/service that can bee kept running indefinitely.
+ Machine and human readable JSON configuration.
+ Optimized for minimal performance impact.

## Usage

| Command       | Alt.    | Explanation                                                                                                  |
|---------------|---------|--------------------------------------------------------------------------------------------------------------|
| install       |         | Installs IPSecDiagTool as a service/daemon.                                                                  |
| uninstall     | remove  | Removes IPSecDiagTool service.                                                                               |
| interactive   | demo    | Allows for interactive testing. Results are directly printed to the console. The service/daemon is not used. |
| mtu-discovery | mtu     | Tells the service/daemon to start finding the MTU for all configured tunnels.                                |
| about         | version | General information about IPSecDiagTool.                                                                     |
| help          |         | A list of commands and how to use them.

## Project structure

    YOUR_GO_WORKSPACE
        bin/
            executables                          # command executable
        pkg/
            packages                             # package object
        src/
            github.com/ipsecdiagtool/ipsecdiagtool/
                .git/                            # Git repository metadata
                .gitignore
                main.go                          # main source file
                build.sh
                update.sh
                README.md
                capture/
                    capture.go                   # Captures pcap data
                config/
                    config.go                    # Read/write the JSON config file
                logging/
                    logging.go                   # Send messages to a Syslog server
                mtu/
                    analyze.go                   # Find the MTU for all configured tunnels
                packetloss/
                    detect.go                    # Detect packet loss and report to a Syslog server

## Golang
IPSecDiagTool is being programmed in Golang. You can pull our code into your own application by
running `go get github.com/ipsecdiagtool/ipsecdiagtool/`. Our code is also documented on 
[Godoc](https://godoc.org/github.com/IPSecDiagTool/IPSecDiagTool).

## License
The MIT License (MIT)

Copyright (c) 2015 Theo Winter, Jan Balmer

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
