# IPSecDiagTool

[![GoDoc](https://godoc.org/github.com/IPSecDiagTool/IPSecDiagTool?status.svg)](https://godoc.org/github.com/IPSecDiagTool/IPSecDiagTool)
[![Build Status](http://152.96.56.53:40000/job/IPSecDiagTool%20-%20Application/badge/icon)](http://152.96.56.53:40000/job/IPSecDiagTool%20-%20Application/)

**Go Code Quality:** http://goreportcard.com/report/IPSecDiagTool/IPSecDiagTool

**YouTrack-Server:** http://sinv-56053.edu.hsr.ch

##About
IPSecDiagTool is a diagnosis tool for IPSec VPN connections. It is being developed as a
semester/bachelor thesis at [HSR](http://www.hsr.ch).

##Main features
+ Passive detection of IPSec packet loss by capturing arriving ESP packets and gathering their sequence numbers.
+ Active diagnosis of IPSec fragmentation problems and discovery of the ideal MTU (Maximum Transmission Unit).

##Project structure

    YOUR_GO_WORKSPACE
        bin/
            executables                          # command executable
        pkg/
            packages                             # package object
        src/
            github.com/ipsecdiagtool/ipsecdiagtool/
                .git/                            # Git repository metadata
                main.go                          # main source file
                mtu/
                    analyze.go                   # MTU Analyzer Code
                packetloss/
                    detect.go                    # Packetloss Detector Code


##Golang
IPSecDiagTool is being programmed in Golang. You can pull our code into your own application by
running `go get github.com/ipsecdiagtool/ipsecdiagtool/`. Our code is also documented on 
[Godoc](https://godoc.org/github.com/IPSecDiagTool/IPSecDiagTool).


##License
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
